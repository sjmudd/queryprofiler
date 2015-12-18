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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	//	"github.com/sjmudd/ps-top/connector"
	"github.com/sjmudd/queryprofiler/collection"
	"github.com/sjmudd/queryprofiler/connection"
//	"github.com/sjmudd/queryprofiler/querycache"
	"github.com/sjmudd/queryprofiler/querykey"
	"github.com/sjmudd/queryprofiler/sample"
)

const (
	driver = "mysql"

	defaultIgnorePerformanceSchema = false
	defaultInterval                = 1
	defaultIterations              = 10
	defaultQueryFilter             = ""
	defaultTopN                    = 5
)

// put this somewhere else or use time or whatever
func toSeconds(t float64) string {
	if t == 0 {
		return "0s"
	}
	if t > 1 {
		return fmt.Sprintf("%.1f%s", t)
	}
	if t < 0.001 {
		return fmt.Sprintf(".1%f us", t*1000000)
	}
	return fmt.Sprintf(".1%f ms", t*1000)
}

// connect to the databases using the given DSNs
// - return a slice of handles to each one.
func connect(dsns []string) connection.Connections {
	var connections connection.Connections
	var Connections string = "Connection" // UPPER CASE C
	var servers string = "connection"

	if len(dsns) != 1 {
		servers += "s"
		Connections += "s"
	}
	log.Println(fmt.Sprintf("Connecting to %d servers...", len(dsns)))
	for i, d := range dsns {
		conn := connection.NewConnection(d)
		conn.Open() // should always work, or be fatal

		log.Println(fmt.Sprintf("%d: %s connected. Version: %s", i, conn.Name(), conn.Version()))

		connections = append(connections, conn)
	}

	return connections
}

// collect collects data from the given connections
func collect(connections connection.Connections, iterations int, intervalSeconds int, ignorePerformanceSchema bool, queryFilter string) []collection.Collections {
	var wg sync.WaitGroup
	var interval = time.Second * time.Duration(intervalSeconds)

	// preallocate the empty collection settings
	collections := make([]collection.Collections, len(connections))

	log.Println(fmt.Sprintf("Collecting data %d times...", iterations))
	for iteration := 1; iteration <= iterations; iteration++ {

		for i, conn := range connections {
			i := i       // avoid go concurrency issues
			conn := conn // for concurrency avoid overwriting stuff

			wg.Add(1)

			// launch collector in parallel
			go func() {
				defer wg.Done()

				collection, err := collection.Collect(conn, ignorePerformanceSchema, queryFilter)
				if err != nil {
					log.Fatal("Failed to collect rows:", err)
				}
				collections[i] = append(collections[i], collection)
			}()
		}
		wg.Wait()

		var lengths []string
		for i := range collections {
			collectedLen := collections[i][len(collections[i])-1].Len()
			lengths = append(lengths, fmt.Sprintf("%d: %d", i, collectedLen))
		}
		log.Println(fmt.Sprintf("  %d: Collected rows: %s", iteration, strings.Join(lengths, ", ")))

		time.Sleep(interval) // wait required time
	}

	log.Println("Collecting finished.")
	return collections
}

// close closes database connections
func close(connections connection.Connections) {
	log.Println("Closing database connections.")
	for _, dbh := range connections {
		dbh.Close()
	}
}

func main() {
	var (
		collections             []collection.Collections
		connections             connection.Connections
		dsns                    []string
		ignorePerformanceSchema bool
		interval                int
		iterations              int
		queryFilter             string
		samples                 sample.SamplesSlice
		topN                    int
	)

	flag.StringVar(&queryFilter, "query-filter", defaultQueryFilter, "Filter queries matching this regexp")
	flag.IntVar(&interval, "interval", defaultInterval, "interval to collect samples")
	flag.IntVar(&iterations, "iterations", defaultIterations, "number of iterations to collect from each connection")
	flag.BoolVar(&ignorePerformanceSchema, "ignore-performance-schema", defaultIgnorePerformanceSchema, "do we ignore performance_schema queries")
	flag.IntVar(&topN, "top-n", defaultTopN, "the top-n queries to show results for")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: ", os.Args[0], "[<options>] <dsn1> [<dsn2> ...]")
		fmt.Println("")
		fmt.Println("Dsn is the Go dsn format such as: user:pass@tcp(host:3306)/performance_schema")
		fmt.Println("Flags:")
		fmt.Println(fmt.Sprintf("--ignore-performance-schema   ignores queries using performance_schema. Default: %v", defaultIgnorePerformanceSchema))
		fmt.Println(fmt.Sprintf("--interval=X                  interval to collect data. Default: %v", defaultInterval))
		fmt.Println(fmt.Sprintf("--iterations=X                number of collections to make. Should be a minimum of 2. Default: %v", defaultIterations))
		fmt.Println(fmt.Sprintf("--query-filter=X              regexp to use when filtering queries. Default: '%v'", defaultQueryFilter))
		fmt.Println(fmt.Sprintf("--top-n=X                     show the top N queries by query elapsed time. Default: %v", defaultTopN))
		os.Exit(1)
	}

	log.Println(fmt.Sprintf("Configuration: Showing top %d queries (by total elapsed time)", topN))
	log.Println(fmt.Sprintf("Configuration: Collecting %d metrics every %d second(s)", iterations, interval))
	log.Println(fmt.Sprintf("Configuration: Query Filter: %s", queryFilter))
	dsns = flag.Args()

	// now start to do stuff
	log.Println("=== Phase I ===")
	connections = connect(dsns)
	log.Println("=== Phase II ===")
	collections = collect(connections, iterations, interval, ignorePerformanceSchema, queryFilter)
	connections.Close()

	// show collection amounts
	log.Println("Collection sizes:")
	for i := range collections {
		log.Println(fmt.Sprintf("%d: Have %d entries, first entry: %d rows, last entry: %d rows",
			i,
			len(collections[i]),
			collections[i][0].Len(),
			collections[i][len(collections[i])-1].Len()))
	}

	log.Println("Converting collection information to samples...")
	for i := range collections {
		s := sample.SamplesFromCollections(collections[i])
		log.Println(fmt.Sprintf("%d: Generated %d sample(s)", i, len(s)))
		samples = append(samples, s)
	}

	/********************************************************************
	** figure out which are the busiest queries                         *
	*  - find the first connection's top n queries based on the sample  *
	*    between the first and last collections.                        *
	*********************************************************************/

	log.Println("=== Phase III ===")
	log.Println("Determining the longest running queries over collection period...")

	var topKeys []querykey.QueryKey
	first := collections[0][0]
	last := collections[0][len(collections[0])-1]
	log.Println(fmt.Sprintf("3.1 Making a single sample from the first/last collections of connection: %s", connections[0].Name()))
	sample := sample.NewSample(
		first.Rows,
		last.Rows,
		first.CollectedTime,
		last.CollectedTime.Sub(first.CollectedTime))

	log.Println("3.2 Sorting results")
	sort.Sort(sample) /* decreasing duration */

	var j int
	for i := 0; i < topN; i++ {
		if len(sample.Rows) > 1 {
			j++

			r := sample.Rows[i]
			topKeys = append(topKeys, r.Key)
		}
	}
	if j == 0 {
		log.Println("Samples from", connections[0].Name(), "have no sample information.")
		log.Println("This usually means that no queries have run in the collection period.")
		log.Println("Unable to continue as this server appears idle.")
		return
	}

	/********************************************************************
	** We now have the top queries (in topKeys).                        *
	*  Collect metrics from the samples for these specific queries and  *
	*  then report them to the user.                                    *
	*********************************************************************/

	log.Println("=== Phase IV ===")
	log.Println("Printing metrics from calculated samples")
	samples.CompareMetrics(connections, topKeys)
}
