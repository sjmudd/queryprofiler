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

// package stats provides some statistics
package metric

import (
	"testing"
)

// results of what we are testing
type results struct {
	min    float64
	max    float64
	avg    float64
	stddev float64
}

// TestStats tests the various stats
func TestStats(t *testing.T) {
	tests := []Metric{
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		{5, 5, 5, 5, 5, 5, 5, 5, 5, 5},
	}
	expected := []results{
		results{avg: 5.5, min: 1.0, max: 10.0, stddev: 2.8722813232690143},
		results{avg: 5.5, min: 1.0, max: 10.0, stddev: 2.8722813232690143},
		results{avg: 5, min: 5.0, max: 5, stddev: 0},
	}

	for i := range tests {
		min := tests[i].Min()
		max := tests[i].Max()
		avg := tests[i].Avg()
		stddev := tests[i].StdDev()

		if min != expected[i].min {
			t.Errorf("Metric(%+v): Min() expected: %v, actual: %v", tests[i], expected[i].min, min)
		}
		if max != expected[i].max {
			t.Errorf("Metric(%+v): max() expected: %v, actual: %v", tests[i], expected[i].max, max)
		}
		if avg != expected[i].avg {
			t.Errorf("Metric(%+v): avg() expected: %v, actual: %v", tests[i], expected[i].avg, avg)
		}
		if stddev != expected[i].stddev {
			t.Errorf("Metric(%+v): stddev() expected: %v, actual: %v", tests[i], expected[i].stddev, stddev)
		}
	}
}
