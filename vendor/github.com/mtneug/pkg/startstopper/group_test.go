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
	"strconv"
	"testing"
	"time"

	"github.com/mtneug/pkg/startstopper"
	"github.com/mtneug/pkg/startstopper/testutils"
	"github.com/stretchr/testify/require"
)

func TestGroup(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	sss := []startstopper.StartStopper{
		&testutils.MockStartStopper{},
		&testutils.MockStartStopper{},
		&testutils.MockStartStopper{},
		&testutils.MockStartStopper{},
	}
	for _, ss := range sss {
		ss := ss.(*testutils.MockStartStopper)
		ss.On("Start", ctx).Return(nil).Once()
		ss.On("Stop", ctx).Return(nil).Once()
		ss.On("Err", ctx).Return(nil).Once()
	}

	group := startstopper.NewGroup(sss)
	require.NotNil(t, group)

	err := group.Start(ctx)
	require.NoError(t, err)

	err = group.Stop(ctx)
	require.NoError(t, err)

	select {
	case <-group.Done():
	case <-time.After(time.Second):
		t.Fatal("Runner not stopped")
	}

	err = group.Err(ctx)
	require.NoError(t, err)

	for _, ss := range sss {
		ss := ss.(*testutils.MockStartStopper)
		ss.AssertExpectations(t)
	}
}

func TestGroupStartErr(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	sss := []startstopper.StartStopper{
		&testutils.MockStartStopper{},
		&testutils.MockStartStopper{},
		&testutils.MockStartStopper{},
		&testutils.MockStartStopper{},
	}
	for i, ss := range sss {
		ss := ss.(*testutils.MockStartStopper)
		if i%2 == 0 {
			err := errors.New("test_error " + strconv.Itoa(i))
			ss.On("Start", ctx).Return(err).Once()
		} else {
			ss.On("Start", ctx).Return(nil).Once()
		}
	}

	group := startstopper.NewGroup(sss)
	require.NotNil(t, group)

	err := group.Start(ctx)
	require.NoError(t, err)

	select {
	case <-group.Done():
	case <-time.After(time.Second):
		t.Fatal("Runner not stopped")
	}

	err = group.Err(ctx)
	require.EqualError(t, err, "test_error 0, test_error 2")

	for _, ss := range sss {
		ss := ss.(*testutils.MockStartStopper)
		ss.AssertExpectations(t)
	}
}

func TestGroupStopErr(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	sss := []startstopper.StartStopper{
		&testutils.MockStartStopper{},
		&testutils.MockStartStopper{},
		&testutils.MockStartStopper{},
		&testutils.MockStartStopper{},
	}
	for i, ss := range sss {
		ss := ss.(*testutils.MockStartStopper)
		ss.On("Start", ctx).Return(nil).Once()
		if i%2 == 0 {
			err := errors.New("test_error " + strconv.Itoa(i))
			ss.On("Stop", ctx).Return(err).Once()
		} else {
			ss.On("Stop", ctx).Return(nil).Once()
			ss.On("Err", ctx).Return(nil).Once()
		}
	}

	group := startstopper.NewGroup(sss)
	require.NotNil(t, group)

	err := group.Start(ctx)
	require.NoError(t, err)

	err = group.Stop(ctx)
	require.NoError(t, err)

	select {
	case <-group.Done():
	case <-time.After(time.Second):
		t.Fatal("Runner not stopped")
	}

	err = group.Err(ctx)
	require.EqualError(t, err, "test_error 0, test_error 2")

	for _, ss := range sss {
		ss := ss.(*testutils.MockStartStopper)
		ss.AssertExpectations(t)
	}
}
