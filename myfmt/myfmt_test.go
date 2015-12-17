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

package myfmt

import (
	"testing"
)

type s1 struct {
	input  uint64
	output string
}
type s2 struct {
	input  float64
	output string
}

func TestFormatInt(t *testing.T) {
	var tests = []s1{
		{1, "1 ps"},
		{10, "10 ps"},
		{100, "100 ps"},
		{1000, "1.00 ns"},
		{10000, "10.00 ns"},
		{100000, "100.00 ns"},
		{1000000, "1.00 µs"},
		{10000000, "10.00 µs"},
		{100000000, "100.00 µs"},
		{1000000000, "1.00 ms"},
		{10000000000, "10.00 ms"},
		{100000000000, "100.00 ms"},
		{1000000000000, "1.00 s"},
		{10000000000000, "10.00 s"},
		{100000000000000, "00:01:40"},
	}

	for i := range tests {
		result := FormatTime(tests[i].input)

		if result != tests[i].output {
			t.Errorf("FormatTime(%v): expected %v, actual %v", tests[i].input, tests[i].output, result)
		}
	}
}

func TestFormatFloat(t *testing.T) {
	var tests = []s2{
		{1000, "1000 s"},
		{100, "100 s"},
		{10, "10 s"},
		{1, "1 s"},
		{0.1, "100.0 ms"},
		{0.01, "10.00 ms"},
		{0.001, "1.000 ms"},
		{0.0001, "100.0 µs"},
		{0.00001, "10.00 µs"},
		{0.000001, "1.000 µs"},
		{0.0000001, "100.0 ns"},
		{0.00000001, "10.00 ns"},
		{0.000000001, "1.000 ns"},
	}

	for i := range tests {
		result := FloatTime(tests[i].input)

		if result != tests[i].output {
			t.Errorf("FormatTime(%v): expected %v, actual %v", tests[i].input, tests[i].output, result)
		}
	}
}
