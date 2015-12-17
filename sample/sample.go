/*
Copyright (c) 2015, Simon J Mudd
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package sample

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sjmudd/queryprofiler/collection"
	"github.com/sjmudd/queryprofiler/connection"
	"github.com/sjmudd/queryprofiler/event"
	"github.com/sjmudd/queryprofiler/metric"
	"github.com/sjmudd/queryprofiler/myfmt"
	"github.com/sjmudd/queryprofiler/querycache"
	"github.com/sjmudd/queryprofiler/querykey"
)

// collection data is not sorted so to extract similar information we need to key by Key

// Row has the values of a collection between 2 points in time.
// COUNT_* AND SUM_* fields are values below are differences from the next value compared to the collected time.
type Row struct {
	Key            querykey.QueryKey
	SUM_TIMER_WAIT uint64
	COUNT_STAR     uint64
	SUM_ERRORS     uint64
	SUM_WARNINGS   uint64

	//	SUM_LOCK_TIME               uint64
	//	SUM_ROWS_AFFECTED           uint64
	//	SUM_ROWS_SENT               uint64
	//	SUM_ROWS_EXAMINED           uint64
	//	SUM_CREATED_TMP_DISK_TABLES uint64
	//	SUM_CREATED_TMP_TABLES      uint64
	//	SUM_SELECT_FULL_JOIN        uint64
	//	SUM_SELECT_FULL_RANGE_JOIN  uint64
	//	SUM_SELECT_RANGE            uint64
	//	SUM_SELECT_RANGE_CHECK      uint64
	//	SUM_SELECT_SCAN             uint64
	//	SUM_SORT_MERGE_PASSES       uint64
	//	SUM_SORT_RANGE              uint64
	//	SUM_SORT_ROWS               uint64
	//	SUM_SORT_SCAN               uint64
	//	SUM_NO_INDEX_USED           uint64
	//	SUM_NO_GOOD_INDEX_USED      uint64
}

// Sample has the rows between 2 collections plus the period they were collected over
type Sample struct {
	Rows      []Row
	StartTime time.Time
	Duration  time.Duration
}

// Samples is a slice of Sample
type Samples []Sample

// SamplesSlice is for a slice of samples
type SamplesSlice []Samples

func (s Sample) Print() {
	log.Println(fmt.Sprintf("StartTime:%+v, Duration:%+v", s.StartTime, s.Duration))
	for i := range s.Rows {
		log.Println(fmt.Sprintf("[%d]: Key: %+v, SUM_TIMER_WAIT: %+v, COUNT_STAR: %+v, SUM_ERRORS: %+v, SUM_WARNINGS: %+v",
			i,
			s.Rows[i].Key,
			s.Rows[i].SUM_TIMER_WAIT,
			s.Rows[i].COUNT_STAR,
			s.Rows[i].SUM_ERRORS,
			s.Rows[i].SUM_WARNINGS))
	}
}

// in queries / second
func (r Row) AvgQueryLatency() float64 {
	if r.COUNT_STAR != 0 {
		return time.Duration(r.SUM_TIMER_WAIT/1000).Seconds() / float64(r.COUNT_STAR)
	}
	return 0
}

// in queries / second
func (s Sample) AvgQueryLatency(i int) float64 {
	return s.Rows[i].AvgQueryLatency()
}

// for sorting
func (sample Sample) Len() int      { return len(sample.Rows) }
func (sample Sample) Swap(i, j int) { sample.Rows[i], sample.Rows[j] = sample.Rows[j], sample.Rows[i] }
func (sample Sample) Less(i, j int) bool {
	return (sample.Rows[i].SUM_TIMER_WAIT > sample.Rows[j].SUM_TIMER_WAIT) ||
		((sample.Rows[i].SUM_TIMER_WAIT == sample.Rows[j].SUM_TIMER_WAIT) &&
			(sample.Rows[i].Key.String() < sample.Rows[j].Key.String())) // this is a useless column to sort on but at least we stay stable
}

type key2eventMap map[string]event.Event

// filledEventMap() makes a lookup map based on the key for each valid row
func filledEventMap(e event.Events) key2eventMap {
	h := make(key2eventMap)

	for i := range e {
		if !e[i].IsEmpty() {
			h[e[i].Key.String()] = e[i]
		}
	}

	return h
}

// StatementRowDiff calculate the difference of 2 different rows which must have the SAME KEY
func StatementRowDiff(e1, e2 event.Event, duration time.Duration) (Row, error) {
	if e1.Key != e2.Key {
		log.Fatal("StatementRowDiff() not comparing rows with the same Key")
	}

	// this should not happen
	if duration == 0 {
		return Row{}, errors.New("sample.StatementRowDiff() duration should not be 0, please ignore")
	}

	// We expect e2 to be after e1, so check this by checking COUNT_STAR is larger and also SUM_TIMER_WAIT is higher
	if (e1.SUM_TIMER_WAIT > e2.SUM_TIMER_WAIT) || (e1.COUNT_STAR > e2.COUNT_STAR) {
		return Row{}, errors.New("sample.StatementRowDiff() e1 is ahead of e2, please ignore")
	}

	//	log.Println("StatementRowDiff()")
	//	log.Println("  e1:", e1)
	//	log.Println("  e2:", e2)

	s := Row{
		Key:            e1.Key,
		COUNT_STAR:     e2.COUNT_STAR - e1.COUNT_STAR,
		SUM_TIMER_WAIT: e2.SUM_TIMER_WAIT - e1.SUM_TIMER_WAIT,
		SUM_ERRORS:     e2.SUM_ERRORS - e1.SUM_ERRORS,
		SUM_WARNINGS:   e2.SUM_WARNINGS - e1.SUM_WARNINGS,

		//		SUM_LOCK_TIME:               e2.SUM_LOCK_TIME - e1.SUM_LOCK_TIME,
		//		SUM_ROWS_AFFECTED:           e2.SUM_ROWS_AFFECTED - e1.SUM_ROWS_AFFECTED,
		//		SUM_ROWS_SENT:               e2.SUM_ROWS_SENT - e1.SUM_ROWS_SENT,
		//		SUM_ROWS_EXAMINED:           e2.SUM_ROWS_EXAMINED - e1.SUM_ROWS_EXAMINED,
		//		SUM_CREATED_TMP_DISK_TABLES: e2.SUM_CREATED_TMP_DISK_TABLES - e1.SUM_CREATED_TMP_DISK_TABLES,
		//		SUM_CREATED_TMP_TABLES:      e2.SUM_CREATED_TMP_TABLES - e1.SUM_CREATED_TMP_TABLES,
		//		SUM_SELECT_FULL_JOIN:        e2.SUM_SELECT_FULL_JOIN - e1.SUM_SELECT_FULL_JOIN,
		//		SUM_SELECT_FULL_RANGE_JOIN:  e2.SUM_SELECT_FULL_RANGE_JOIN - e1.SUM_SELECT_FULL_RANGE_JOIN,
		//		SUM_SELECT_RANGE:            e2.SUM_SELECT_RANGE - e1.SUM_SELECT_RANGE,
		//		SUM_SELECT_RANGE_CHECK:      e2.SUM_SELECT_RANGE_CHECK - e1.SUM_SELECT_RANGE_CHECK,
		//		SUM_SELECT_SCAN:             e2.SUM_SELECT_SCAN - e1.SUM_SELECT_SCAN,
		//		SUM_SORT_MERGE_PASSES:       e2.SUM_SORT_MERGE_PASSES - e1.SUM_SORT_MERGE_PASSES,
		//		SUM_SORT_RANGE:              e2.SUM_SORT_RANGE - e1.SUM_SORT_RANGE,
		//		SUM_SORT_ROWS:               e2.SUM_SORT_ROWS - e1.SUM_SORT_ROWS,
		//		SUM_SORT_SCAN:               e2.SUM_SORT_SCAN - e1.SUM_SORT_SCAN,
		//		SUM_NO_INDEX_USED:           e2.SUM_NO_INDEX_USED - e1.SUM_NO_INDEX_USED,
		//		SUM_NO_GOOD_INDEX_USED:      e2.SUM_NO_GOOD_INDEX_USED - e1.SUM_NO_GOOD_INDEX_USED,
	}
	//	log.Println("  s:", s)

	return s, nil
}

// NewSample returns a Sample based on the difference in two collections.
// We calculate the diff of all rows and return the diff
// values in e2 not in e1 are ignored
// if there's not a value in e2 that matches e1 then we ingore it
func NewSample(e1, e2 event.Events, start time.Time, duration time.Duration) Sample {
	sample := Sample{StartTime: start, Duration: duration}

	e2lookup := filledEventMap(e2)

	for i := range e1 {
		if e2value, found := e2lookup[e1[i].Key.String()]; found {
			// ignoring errors if the counts do not match
			if row, err := StatementRowDiff(e1[i], e2value, duration); err == nil {
				// log.Println("  sample:", row)

				// Only take samples with COUNT_STAR > 0 and SUM_TIMER_WAIT > 0
				// as null samples are not useful.
				if row.COUNT_STAR > 0 && row.SUM_TIMER_WAIT > 0 {
					sample.Rows = append(sample.Rows, row)
				}
			}
		}
	}

	return sample
}

// SamplesFromCollections takes the collection and looks at each interval
// generating sample values for those that match
func SamplesFromCollections(c collection.Collections) Samples {
	var samples Samples

	// log.Println(fmt.Sprintf("SamplesFromCollections() %d collections", len(c)))

	for i := range c[0 : len(c)-1] {
		sample := NewSample(
			c[i].Rows,
			c[i+1].Rows,
			c[i].CollectedTime,
			c[i+1].CollectedTime.Sub(c[i].CollectedTime))

		samples = append(samples, sample)
	}

	// log.Println(fmt.Sprintf("- return %d samples", len(samples)))

	return samples
}

// RowByKey returns the row in a Sample that has the given key (if it exists)
func (s Sample) RowByKey(key querykey.QueryKey) (Row, bool) {
	// log.Println(fmt.Sprintf("RowByKey(%v)", key))
	for i := range s.Rows {
		if s.Rows[i].Key == key {
			return s.Rows[i], true
		}
	}
	return Row{}, false
}

// FindMetricsByDigest returns a metric for the digested query from the given samples.
func (s Samples) MetricsByKey(key querykey.QueryKey) metric.NamedMetrics {
	// log.Println(fmt.Sprintf("MetricsByKey(%v) len() = %d", &s, len(s)))

	named := make(metric.NamedMetrics)

	named["queries"] = nil
	named["QPS"] = nil
	named["Latency"] = nil
	named["errors"] = nil
	named["warnings"] = nil

	// iterate over each sample looking for the metric values of the given key
	for i := range s {
		if r, ok := s[i].RowByKey(key); ok {
			// Ensure we have metrics for this query.
			if r.COUNT_STAR > 0 {
				named["QPS"] = append(named["QPS"], float64(r.COUNT_STAR)/s[i].Duration.Seconds())
				named["Latency"] = append(named["Latency"], r.AvgQueryLatency())
				named["queries"] = append(named["queries"], float64(r.COUNT_STAR))
				named["errors"] = append(named["errors"], float64(r.SUM_ERRORS))
				named["warnings"] = append(named["warnings"], float64(r.SUM_WARNINGS))
			}
		}
	}

	return named
}

// PrintMetrics for each sample for the topKeys provided.
func (s Samples) PrintMetrics(conn *connection.Connection, topKeys []querykey.QueryKey) {
	log.Println("------------------------")
	log.Println(fmt.Sprintf("Host: %s / %s. Have %d entries. Looking at top %d queries", conn.Name(), conn.Version(), len(s), len(topKeys)))
	for i := range topKeys {
		key := topKeys[i]

		// query, _ := querycache.Get(key)
		m := s.MetricsByKey(key)

		log.Println(fmt.Sprintf("Query %d:", i+1))

		print(m)
	}
}

// short name removing the middle...
func shortName(s string) string {
	if len(s) <= 40 /* don't hard code */ {
		return s
	}

	// 40 = 2n + 3, => 37 / 2
	i := 18 // int(37 / 2 )

	return s[0:i] + "..." + s[len(s)-1-i:]
}

// CompareMetrics prints out the metrics for s SampleSlices. That is
// compare the values of the first one with the other servers.
func (s SamplesSlice) CompareMetrics(connections connection.Connections, topKeys []querykey.QueryKey) {
	if s == nil {
		log.Println("No data provided")
		return
	}

	log.Println(fmt.Sprintf("Looking at top %d queries", len(topKeys)))

	// generate server name headings
	names := []string{}
	for i := range connections {
		cutDown := shortName(connections[i].Name())
		names = append(names, fmt.Sprintf("%-40s", cutDown))
	}
	serverNames := "Server:            " + strings.Join(names, " | ")

	for i := range topKeys {
		key := topKeys[i]
		query, _ := querycache.Get(key)

		log.Println(fmt.Sprintf("Query %d: %s", i+1, query))
		log.Println(serverNames)

		var m []metric.NamedMetrics

		// get the metrics
		for j := range s {
			m = append(m, s[j].MetricsByKey(key))
		}

		// collect formatted output
		headings := formattedHeadings(m[0])
		formattedOutput := [][]string{}

		for j := range s {
			formattedOutput = append(formattedOutput, formattedValues(m[j]))
		}

		// print in the right order
		for row := range formattedOutput[0] {
			output := []string{}
			for server := range formattedOutput {
				output = append(output, formattedOutput[server][row])
			}
			log.Println(headings[row], strings.Join(output, " | "))
		}
	}
}

// print a metric out
func print(n metric.NamedMetrics) {
	names := []string{"queries", "QPS", "Latency", "errors", "warnings"}

	if n == nil {
		return
	}

	for i := range names {
		v := names[i]

		// this should NOT be done this way (hard-coded but provide some function pointer to do this !
		switch {
		case v == "queries":
			log.Println(fmt.Sprintf("  Avg %s: %s, σ %s",
				v,
				myfmt.FloatNumber(n[v].Avg()),
				myfmt.FloatNumber(n[v].StdDev())))
		case v == "QPS":
			log.Println(fmt.Sprintf("  Avg %s: %s qps, σ %s",
				v,
				myfmt.FloatNumber(n[v].Avg()),
				myfmt.FloatNumber(n[v].StdDev())))
		case v == "Latency":
			log.Println(fmt.Sprintf("  Avg %s: %s, σ %s",
				v,
				myfmt.FloatTime(n[v].Avg()),
				myfmt.FloatTime(n[v].StdDev())))
		case true:
			log.Println(fmt.Sprintf("  Avg %v: %v, σ %v",
				v,
				n[v].Avg(),
				n[v].StdDev()))
		}
	}
}

func formattedHeadings(n metric.NamedMetrics) []string {
	if n == nil {
		return nil
	}

	s := []string{}

	names := []string{"queries", "QPS", "Latency", "errors", "warnings"}

	for i := range names {
		value := fmt.Sprintf("  Avg %s:", names[i])
		s = append(s, fmt.Sprintf("%-18s", value))
	}

	return s
}

func formattedValues(n metric.NamedMetrics) []string {
	var value string

	if n == nil {
		return nil
	}

	s := []string{}

	names := []string{"queries", "QPS", "Latency", "errors", "warnings"}

	// should sort maybe ?
	for i := range names {
		v := names[i]

		// this should NOT be done this way (hard-coded but provide some function pointer to do this !
		switch {
		case v == "queries":
			value = fmt.Sprintf("%s, σ %s", myfmt.FloatNumber(n[v].Avg()), myfmt.FloatNumber(n[v].StdDev()))
		case v == "QPS":
			value = fmt.Sprintf("%s qps, σ %s", myfmt.FloatNumber(n[v].Avg()), myfmt.FloatNumber(n[v].StdDev()))
		case v == "Latency":
			value = fmt.Sprintf("%s, σ %s", myfmt.FloatTime(n[v].Avg()), myfmt.FloatTime(n[v].StdDev()))
		case true:
			value = fmt.Sprintf("%v, σ %v", n[v].Avg(), n[v].StdDev())
		}
		s = append(s, fmt.Sprintf("%-40s", value)) // hard coded BAD
	}

	return s
}
