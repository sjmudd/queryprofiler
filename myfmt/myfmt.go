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
	"fmt"
	"strconv"
)

// return rounded version of x with prec precision.
func Round(x float64, prec int) string {
	frep := strconv.FormatFloat(x, 'g', prec, 64)
	f, _ := strconv.ParseFloat(frep, 64)
	format := "%." + fmt.Sprintf("%d", prec) + "f"
	return fmt.Sprintf(format, f)
}

// input in whole seconds
func SecToTime(t uint64) string {
	var h, m, s uint64
	s = t % 60
	m = ((t - s) / 60) % 60
	h = (t - s - 60*m) % 3600
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// input in P_S picoseconds
func FormatTime(picoseconds uint64) string {
	if picoseconds >= 60000000000000 {
		return SecToTime(picoseconds / 1000000000000)
	}
	if picoseconds >= 1000000000000 {
		return Round(float64(picoseconds/1000000000000), 2) + " s"
	}
	if picoseconds >= 1000000000 {
		return Round(float64(picoseconds/1000000000), 2) + " ms"
	}
	if picoseconds >= 1000000 {
		return Round(float64(picoseconds/1000000), 2) + " µs"
	}
	if picoseconds >= 1000 {
		return Round(float64(picoseconds/1000), 2) + " ns"
	}
	return strconv.Itoa(int(picoseconds)) + " ps"
}

// FloatTime formats a time to a nice to read format
// Input in seconds (float64)
func FloatTime(t float64) string {
	if t >= 1 { // 999
		return fmt.Sprintf("%.0f s", t)
	}
	if t >= 0.1 { // 0.1 ---> 10.00 ms
		return fmt.Sprintf("%.1f ms", t*1000)
	}
	if t >= 0.01 { // 0.01 ---> 10.00 ms
		return fmt.Sprintf("%.2f ms", t*1000)
	}
	if t >= 0.001 { // 0.001 ---> 1.000 ms
		return fmt.Sprintf("%.3f ms", t*1000)
	}
	if t >= 0.0001 { // 0.001 ---> 100.0 µs
		return fmt.Sprintf("%.1f µs", t*1000000)
	}
	if t >= 0.00001 { // 0.0001 ---> 1.00 µs
		return fmt.Sprintf("%.2f µs", t*1000000)
	}
	if t >= 0.000001 { // 0.00001 ---> 1.00 µs
		return fmt.Sprintf("%.3f µs", t*1000000)
	}
	if t >= 0.0000001 { // 0.00001 ---> 1.00 µs
		return fmt.Sprintf("%.1f ns", t*1000000000)
	}
	if t >= 0.00000001 { // 0.00001 ---> 1.00 µs
		return fmt.Sprintf("%.2f ns", t*1000000000)
	}
	if t >= 0.000000001 { // 0.00001 ---> 1.00 ns
		return fmt.Sprintf("%.3f ns", t*1000000000)
	}
	return fmt.Sprintf("%f s", t)
}

// FloatTime formats a time to a nice to read format
// Input in seconds (float64)
func FloatNumber(t float64) string {
	var v string

	switch {
	case t == 0: // 0
		v = fmt.Sprintf("%.0f", t)
	case t >= 1000: // 9999
		v = fmt.Sprintf("%.2f k", t/1000)
	case t >= 100: // 999
		v = fmt.Sprintf("%.0f", t)
	case t >= 10: // 99.9
		v = fmt.Sprintf("%.1f", t)
	case t >= 1: // 9.99
		v = fmt.Sprintf("%.2f", t)
	case t >= 0.1: // 0.1 ---> 0.999
		v = fmt.Sprintf("%.3f", t)
	case t >= 0.01: // 0.01 ---> 0.0123
		v = fmt.Sprintf("%.4f", t)
	case true:
		v = fmt.Sprintf("%f", t)
	}

	return v
}
