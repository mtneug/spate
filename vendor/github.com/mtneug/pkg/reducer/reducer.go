// Copyright (c) 2016 Matthias Neugebauer
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package reducer

import (
	"errors"
	"math"
)

var (
	// ErrSliceEmpty indicates that the reduction failed because the slice was
	// empty.
	ErrSliceEmpty = errors.New("reducer: slice is empty")
)

// Reducer aggreate float64 slices.
type Reducer interface {
	Reduce([]float64) (float64, error)
}

// Func is a reducers.
type Func func([]float64) (float64, error)

// Reduce implements the Reducer interface. It calls itself ensuring that the
// slice is not empty.
func (f Func) Reduce(data []float64) (float64, error) {
	if len(data) == 0 {
		return 0, ErrSliceEmpty
	}
	return f(data)
}

// Max returns an reducer that reduces the float64 slice to the maximum.
func Max() Reducer {
	return Func(func(data []float64) (float64, error) {
		max := data[0]
		for i := 1; i < len(data); i++ {
			max = math.Max(max, data[i])
		}
		return max, nil
	})
}

// Min returns an reducer that reduces the float64 slice to the minimum.
func Min() Reducer {
	return Func(func(data []float64) (float64, error) {
		min := data[0]
		for i := 1; i < len(data); i++ {
			min = math.Min(min, data[i])
		}
		return min, nil
	})
}

// Avg returns an reducer that reduces the float64 slice to the average.
func Avg() Reducer {
	return Func(func(data []float64) (float64, error) {
		sum, _ := Sum().Reduce(data)
		return sum / float64(len(data)), nil
	})
}

// Sum returns an reducer that reduces the float64 slice to the sum.
func Sum() Reducer {
	return Func(func(data []float64) (float64, error) {
		sum := data[0]
		for i := 1; i < len(data); i++ {
			sum += data[i]
		}
		return sum, nil
	})
}
