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

package reducer_test

import (
	"errors"
	"math"
	"testing"

	"github.com/mtneug/pkg/reducer"
	"github.com/stretchr/testify/require"
)

func TestFunc(t *testing.T) {
	t.Parallel()

	var called bool
	testVal := 42.0
	testErr := errors.New("test")

	r := reducer.Func(func(data []float64) (float64, error) {
		called = true
		return testVal, testErr
	})

	called = false
	_, err := r.Reduce(nil)
	require.False(t, called)
	require.EqualError(t, err, reducer.ErrSliceEmpty.Error())

	called = false
	val, err := r.Reduce([]float64{42})
	require.True(t, called)
	require.Equal(t, testErr, err)
	require.Equal(t, testVal, val)
}

func TestMax(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		slice []float64
		val   float64
	}{
		{
			slice: []float64{1, 2, 3, 4},
			val:   4,
		},
		{
			slice: []float64{-1, -2, -3, -4},
			val:   -1,
		},
		{
			slice: []float64{math.Inf(-1), 2, 3, math.Inf(1)},
			val:   math.Inf(1),
		},
	}

	r := reducer.Max()

	for _, c := range testCases {
		val, err := r.Reduce(c.slice)
		require.NoError(t, err)
		require.Equal(t, c.val, val)
	}
}

func TestMin(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		slice []float64
		val   float64
	}{
		{
			slice: []float64{1, 2, 3, 4},
			val:   1,
		},
		{
			slice: []float64{-1, -2, -3, -4},
			val:   -4,
		},
		{
			slice: []float64{math.Inf(-1), 2, 3, math.Inf(1)},
			val:   math.Inf(-1),
		},
	}

	r := reducer.Min()

	for _, c := range testCases {
		val, err := r.Reduce(c.slice)
		require.NoError(t, err)
		require.Equal(t, c.val, val)
	}
}

func TestAvg(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		slice []float64
		val   float64
	}{
		{
			slice: []float64{1, 2, 3, 4},
			val:   2.5,
		},
		{
			slice: []float64{-1, -2, -3, -4},
			val:   -2.5,
		},
	}

	r := reducer.Avg()

	for _, c := range testCases {
		val, err := r.Reduce(c.slice)
		require.NoError(t, err)
		require.Equal(t, c.val, val)
	}
}

func TestSum(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		slice []float64
		val   float64
	}{
		{
			slice: []float64{1, 2, 3, 4},
			val:   10,
		},
		{
			slice: []float64{-1, -2, -3, -4},
			val:   -10,
		},
	}

	r := reducer.Sum()

	for _, c := range testCases {
		val, err := r.Reduce(c.slice)
		require.NoError(t, err)
		require.Equal(t, c.val, val)
	}
}
