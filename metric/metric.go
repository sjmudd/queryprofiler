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

package metric

// we probably can do this with an outside routine but for now this gives me some numbers

import (
	"math"
)

// we want to get a series of Metricss and then do some number crunching on them.
//

type Values []float64
type NamedMetrics map[string]Values

func (m Values) Min() float64 {
	if len(m) == 0 {
		return 0
	}
	min := float64(999999999)
	for i := range m {
		if m[i] < min {
			min = m[i]
		}
	}
	return min
}

func (m Values) Max() float64 {
	if len(m) == 0 {
		return 0
	}
	max := float64(-999999999)
	for i := range m {
		if m[i] > max {
			max = m[i]
		}
	}
	return max
}

func (m Values) Avg() float64 {
	if len(m) == 0 {
		return 0
	}
	var avg, sum float64
	count := 0

	for i := range m {
		sum += m[i]
		count++
	}

	if count > 0 {
		avg = sum / float64(count)
	}

	return avg
}

// sqrt( sum (x - mean)^ 2 / n )
func (m Values) StdDev() float64 {
	if len(m) == 0 {
		return 0
	}
	var sum, mean float64

	if len(m) == 0 {
		return 0
	}
	mean = m.Avg()
	for i := range m {
		sum += math.Pow(m[i]-mean, 2)
	}
	sum /= float64(len(m))
	return math.Sqrt(sum)
}
