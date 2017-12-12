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

package startstopper_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mtneug/pkg/startstopper"
	"github.com/mtneug/pkg/startstopper/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewGo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	rErr := errors.New("Test error")
	r := &testutils.MockRunner{}
	r.On("Run", ctx).Return(rErr).Twice()

	ss := startstopper.NewGo(r)
	require.NotNil(t, ss)

	err := ss.Stop(ctx)
	require.EqualError(t, startstopper.ErrNotStarted, err.Error())

	ctx2, cancel := context.WithCancel(ctx)
	cancel()
	err = ss.Err(ctx2)
	require.EqualError(t, err, "context canceled")

	err = ss.Start(ctx)
	require.NoError(t, err)

	err = ss.Start(ctx)
	require.EqualError(t, startstopper.ErrStarted, err.Error())

	err = ss.Stop(ctx)
	require.NoError(t, err)

	ss = startstopper.NewGo(r)
	_ = ss.Start(ctx)
	err = ss.Stop(ctx2)
	require.EqualError(t, err, "context canceled")

	select {
	case <-ss.Done():
	case <-time.After(time.Second):
		t.Fatal("Runner not stopped")
	}

	err = ss.Err(ctx)
	require.EqualError(t, err, rErr.Error())

	r.AssertExpectations(t)
}

func TestRunnerFunc(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	stopChan := make(<-chan struct{})

	rErr := errors.New("Test error")
	called := 0
	r := startstopper.RunnerFunc(func(passedCtx context.Context, passedStopChan <-chan struct{}) error {
		require.Equal(t, ctx, passedCtx)
		require.Equal(t, stopChan, passedStopChan)
		called++
		return rErr
	})

	err := r.Run(ctx, stopChan)
	require.EqualError(t, err, rErr.Error())
	require.Equal(t, 1, called)
}
