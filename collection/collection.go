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

package collection

import (
	"fmt"
	"log"
	"time"

	"github.com/sjmudd/queryprofiler/connection"
	"github.com/sjmudd/queryprofiler/event"
)

type Collection struct {
	Rows          event.Events
	CollectedTime time.Time
}

type Collections []Collection

// return the size of the collection rows (saves going indirectly)
func (c Collection) Len() int {
	return len(c.Rows)
}

// Return the rows obtained from the query as a slice.
// See: http://bugs.mysql.com/bug.php?id=79533
// While this table is supposed to have a row per DIGEST / SCHEMA_NAME it does not
// so to ensure ths we need to collect all rows and then merge in the changes to generate
// a single row. Consequently store into a hash, and if you find duplicates adjust
// the final row, but finally return clean slice
func Collect(conn *connection.Connection, ignorePerformanceSchema bool, queryFilter string) (Collection, error) {
	var (
		err        error
		collection Collection
	)

	digestRows, err := event.CollectEvents(conn, ignorePerformanceSchema, queryFilter)
	if err != nil {
		log.Fatal(err)
	}

	collection.Rows = event.RemoveDups(digestRows) // remove possible duplicates see: http://bugs.mysql.com/bug.php?id=79533
	collection.CollectedTime = time.Now()

	// log.Println(fmt.Sprintf("Collected() collected %d rows", collection.Len()))

	return collection, err
}

func (c Collection) String() string {
	s := c.CollectedTime.String() + "\n"
	for i := range c.Rows {
		s += fmt.Sprintf("%+v\n", c.Rows[i])
	}
	return s
}
