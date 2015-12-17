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

// Package Event performs operations on the performance_schema.events_statements_summary_by_digest table.
package event

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/sjmudd/queryprofiler/connection"
	"github.com/sjmudd/queryprofiler/querycache"
	"github.com/sjmudd/queryprofiler/querykey"
)

// Event represents a row in
// performance_schema.events_statements_summary_by_digest.  Key
// represents the MD5_DIGEST/SCHEMA_NAME combination and is used as
// the DIGEST value is not consisent over different versions of
// MySQL. It is generated during query collection.
type Event struct {
	//	MD5_DIGEST  string // consistent digest of DIGEST_TEXT as different servers may generate a different value
	//	SCHEMA_NAME string
	//	DIGEST                      string
	//	DIGEST_TEXT                 string

	Key            querykey.QueryKey
	COUNT_STAR     uint64
	SUM_TIMER_WAIT uint64

	//	MIN_TIMER_WAIT              uint64
	//	AVG_TIMER_WAIT              uint64
	//	MAX_TIMER_WAIT              uint64
	//	SUM_LOCK_TIME               uint64
	SUM_ERRORS   uint64
	SUM_WARNINGS uint64
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
	//	FIRST_SEEN                  string
	//	LAST_SEEN                   string
}

type Events []Event

// Convert a NULL to "" which should be fine.
func nullStringToString(original sql.NullString) string {
	if original.Valid {
		return original.String
	}
	return ""
}

func (e Event) String() string {
	s := []string{
		fmt.Sprintf("MD5_DIGEST: %s", e.Key.Digest()),
		fmt.Sprintf("SCHEMA: %s", e.Key.Schema()),
		fmt.Sprintf("COUNT_STAR: %+v", e.COUNT_STAR),
		fmt.Sprintf("SUM_TIMER_WAIT: %+v", e.SUM_TIMER_WAIT),
		fmt.Sprintf("SUM_ERRORS: %+v", e.SUM_ERRORS),
		fmt.Sprintf("SUM_WARNINGS: %+v", e.SUM_WARNINGS),
	}

	//	s += fmt.Sprintf("MD5_DIGEST: %+v\n", e.MD5_DIGEST)
	//	s += fmt.Sprintf("SCHEMA_NAME: %+v\n", e.SCHEMA_NAME)
	//	s += fmt.Sprintf("DIGEST: %+v\n", e.DIGEST)
	//	s += fmt.Sprintf("DIGEST_TEXT: %+v\n", e.DIGEST_TEXT)
	//	s += fmt.Sprintf("MIN_TIMER_WAIT: %+v\n", e.MIN_TIMER_WAIT)
	//	s += fmt.Sprintf("AVG_TIMER_WAIT: %+v\n", e.AVG_TIMER_WAIT)
	//	s += fmt.Sprintf("MAX_TIMER_WAIT: %+v\n", e.MAX_TIMER_WAIT)
	//	s += fmt.Sprintf("SUM_LOCK_TIME: %+v\n",  e.SUM_LOCK_TIME)

	//	s += fmt.Sprintf("SUM_ROWS_AFFECTED: %+v\n", e.SUM_ROWS_AFFECTED)
	//	s += fmt.Sprintf("SUM_ROWS_SENT: %+v\n",     e.SUM_ROWS_SENT)
	//	s += fmt.Sprintf("SUM_ROWS_EXAMINED: %+v\n", e.SUM_ROWS_EXAMINED)
	//	s += fmt.Sprintf("SUM_CREATED_TMP_DISK_TABLES: %+v\n", e.SUM_CREATED_TMP_DISK_TABLES)
	//	s += fmt.Sprintf("SUM_CREATED_TMP_TABLES: %+v\n",      e.SUM_CREATED_TMP_TABLES)
	//	s += fmt.Sprintf("SUM_SELECT_FULL_JOIN: %+v\n",        e.SUM_SELECT_FULL_JOIN)
	//	s += fmt.Sprintf("SUM_SELECT_FULL_RANGE_JOIN: %+v\n",  e.SUM_SELECT_FULL_RANGE_JOIN)
	//	s += fmt.Sprintf("SUM_SELECT_RANGE: %+v\n",            e.SUM_SELECT_RANGE)
	//	s += fmt.Sprintf("SUM_SELECT_RANGE_CHECK: %+v\n",      e.SUM_SELECT_RANGE_CHECK)
	//	s += fmt.Sprintf("SUM_SELECT_SCAN: %+v\n",             e.SUM_SELECT_SCAN)
	//	s += fmt.Sprintf("SUM_SORT_MERGE_PASSES: %+v\n",       e.SUM_SORT_MERGE_PASSES)
	//	s += fmt.Sprintf("SUM_SORT_RANGE: %+v\n",              e.SUM_SORT_RANGE)
	//	s += fmt.Sprintf("SUM_SORT_ROWS: %+v\n",               e.SUM_SORT_ROWS)
	//	s += fmt.Sprintf("SUM_SORT_SCAN: %+v\n",               e.SUM_SORT_SCAN)
	//	s += fmt.Sprintf("SUM_NO_INDEX_USED: %+v\n",           e.SUM_NO_INDEX_USED)
	//	s += fmt.Sprintf("SUM_NO_GOOD_INDEX_USED: %+v\n",      e.SUM_NO_GOOD_INDEX_USED)
	//	s += fmt.Sprintf("FIRST_SEEN: %+v\n",                  e.FIRST_SEEN)
	//	s += fmt.Sprintf("LAST_SEEN: %+v\n",                   e.LAST_SEEN)

	return strings.Join(s, ", ")
}

// Return the rows obtained from the query as a slice.
func CollectEvents(conn *connection.Connection, ignorePerformanceSchema bool, queryFilter string) (Events, error) {
	// record when we collected the data
	var (
		events Events
		err    error
		result *sql.Rows
	)

	dbh := conn.Handle()

	query := `
SELECT	MD5(DIGEST_TEXT) AS 'MD5_DIGEST',
	SCHEMA_NAME,
/*	DIGEST, */
	DIGEST_TEXT,
	COUNT_STAR,
	SUM_TIMER_WAIT, /*
	MIN_TIMER_WAIT,
	AVG_TIMER_WAIT,
	MAX_TIMER_WAIT,
	SUM_LOCK_TIME, */
	SUM_ERRORS,
	SUM_WARNINGS /*,
	SUM_ROWS_AFFECTED,
	SUM_ROWS_SENT,
	SUM_ROWS_EXAMINED,
	SUM_CREATED_TMP_DISK_TABLES,
	SUM_CREATED_TMP_TABLES,
	SUM_SELECT_FULL_JOIN,
	SUM_SELECT_FULL_RANGE_JOIN,
	SUM_SELECT_RANGE,
	SUM_SELECT_RANGE_CHECK,
	SUM_SELECT_SCAN,
	SUM_SORT_MERGE_PASSES,
	SUM_SORT_RANGE,
	SUM_SORT_ROWS,
	SUM_SORT_SCAN,
	SUM_NO_INDEX_USED,
	SUM_NO_GOOD_INDEX_USED,
	FIRST_SEEN,
	LAST_SEEN */
FROM	events_statements_summary_by_digest
`

	if queryFilter != "" {
		query += ` WHERE DIGEST_TEXT LIKE ?`
		// log.Println("Query:", query)
		result, err = dbh.Query(query, "%"+queryFilter+"%") // this is ugly having 2 code-paths
	} else {
		// log.Println("Query:", query)
		result, err = dbh.Query(query) // this is ugly having 2 code-paths
	}
	if err != nil {
		log.Fatal(err)
	}

	defer result.Close()

	for result.Next() {
		var row Event
		var nullableMd5Digest, nullableSchemaName, nullableDigestText sql.NullString // groan
		var md5Digest, schemaName, digestText string

		if err := result.Scan(
			&nullableMd5Digest,
			&nullableSchemaName,
			&nullableDigestText,
			&row.COUNT_STAR,
			&row.SUM_TIMER_WAIT,
			&row.SUM_ERRORS,
			&row.SUM_WARNINGS); err != nil {
			/*
				&row.MIN_TIMER_WAIT,
				&row.AVG_TIMER_WAIT,
				&row.MAX_TIMER_WAIT,
				&row.SUM_LOCK_TIME,

				&row.SUM_ROWS_AFFECTED,
				&row.SUM_ROWS_SENT,
				&row.SUM_ROWS_EXAMINED,
				&row.SUM_CREATED_TMP_DISK_TABLES,
				&row.SUM_CREATED_TMP_TABLES,
				&row.SUM_SELECT_FULL_JOIN,
				&row.SUM_SELECT_FULL_RANGE_JOIN,
				&row.SUM_SELECT_RANGE,
				&row.SUM_SELECT_RANGE_CHECK,
				&row.SUM_SELECT_SCAN,
				&row.SUM_SORT_MERGE_PASSES,
				&row.SUM_SORT_RANGE,
				&row.SUM_SORT_ROWS,
				&row.SUM_SORT_SCAN,
				&row.SUM_NO_INDEX_USED,
				&row.SUM_NO_GOOD_INDEX_USED,
				&row.FIRST_SEEN,
				&row.LAST_SEEN  ); err != nil { */
			log.Fatal(err)
		}

		md5Digest = nullStringToString(nullableMd5Digest)
		schemaName = nullStringToString(nullableSchemaName)
		digestText = nullStringToString(nullableDigestText)
		key := querykey.NewQueryKey(md5Digest, schemaName)

		row.Key = key

		querycache.Put(key, digestText)

		if ignorePerformanceSchema && (row.Key.Schema() == "performance_schema") {
			// log.Println("dropping row:", row)
		} else {
			events = append(events, row)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	//	log.Println("CollectEvents() collects", len(events), "row(s)")

	return events, err
}

// IsEmpty tells us if we have any interesting data
func (e Event) IsEmpty() bool {
	return e.COUNT_STAR == 0 && e.SUM_TIMER_WAIT == 0
}

// Merge merges the values of two events IFF they have the same Key.
// If they do not then log.Fatal() will be called.
func Merge(this, that Event) Event {
	if this.Key != that.Key {
		log.Println("Merge(", this, ",", that, ")")
		log.Fatal("Unable to merge 2 events due to differing keys")
	}
	newOne := this

	newOne.COUNT_STAR += that.COUNT_STAR
	newOne.SUM_TIMER_WAIT += that.SUM_TIMER_WAIT
	newOne.SUM_ERRORS += that.SUM_ERRORS
	newOne.SUM_WARNINGS += that.SUM_WARNINGS

	return newOne
}

// RemoveDups removes duplicate events from a slice and returns the cleaned slice.
// See http://bugs.mysql.com/bug.php?id=79533. If there are duplicates their values
// will be merged together (adding values).
func RemoveDups(original Events) Events {
	var dupes int
	var h = make(map[string]Event) // lookup on the string to event
	var newSlice Events

	for i := range original {
		//		var later_text, prev_text, updated_text string
		var description string
		key := original[i].Key.String()

		//		later_text = fmt.Sprintf("NEW: %+v", original[i])

		prev, found := h[key]
		if !found {
			h[key] = original[i] // add to hash
		} else {
			//			prev_text = fmt.Sprintf("OLD: %+v", prev)
			dupes++
			if prev.IsEmpty() {
				// ignore the existing empty row
				h[key] = original[i] // add later value to hash, ignoring the prev empty one
				description = "Overwriting Empty prev value"
			} else {
				if original[i].IsEmpty() {
					description = "Not overwriting. Value is empty and prev value was not"
				} else {
					description = "Merging as prev and later are not empty"
					h[key] = Merge(h[key], original[i])
					//					updated_text = fmt.Sprintf("UPD: %+v", h[key])
				}
			}
		}
		if description != "" {
			//			log.Println(fmt.Sprintf("RemoveDups() Row %d: WARNING [%v] %s", i, key, description))
			//			log.Println(later_text)
			//			log.Println(prev_text)
			//			if updated_text != "" {
			//				log.Println(updated_text)
			//			}
		}
	}

	// now iterate over keys in map to make a new slice
	for _, v := range h {
		newSlice = append(newSlice, v)
	}

	//	log.Println(fmt.Sprintf("RemoveDups() removed %d dupes, leaving %d rows", dupes, len(newSlice)))

	return newSlice
}
