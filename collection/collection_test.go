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
	"database/sql"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	var tests = []EventsStatementsSummaryByDigest{
		{COUNT_STAR: 0, SUM_TIMER_WAIT: 0}, // default
		{COUNT_STAR: 100, SUM_TIMER_WAIT: 0},
		{COUNT_STAR: 0, SUM_TIMER_WAIT: 100},
		{COUNT_STAR: 100, SUM_TIMER_WAIT: 100},
	}
	var expected = []bool{
		true,
		false,
		false,
		false,
	}

	for i := range tests {
		result := tests[i].IsEmpty()

		if result != expected[i] {
			t.Errorf("e.IsEmpty(%v): expected %v, actual %v", tests[i], expected[i], result)
		}
	}
}

func TestNullString(t *testing.T) {
	var tests = []sql.NullString{
		{Valid: false},
		{Valid: true, String: "hello"},
	}
	var expected = []string{
		"",
		"hello",
	}
	for i := range tests {
		result := nullStringToString(tests[i])

		if result != expected[i] {
			t.Errorf("e.nullStringToString(%v): expected %v, actual %v", tests[i], expected[i], result)
		}
	}
}
