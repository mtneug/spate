// Copyright © 2016 Matthias Neugebauer <mtneug@mailbox.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"context"
	"testing"
	"time"

	"github.com/mtneug/pkg/startstopper/testutils"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	p := time.Second

	ctrl, err := New(p)
	require.NotNil(t, ctrl)
	require.NoError(t, err)

	require.IsType(t, &eventLoop{}, ctrl.eventLoop)

	require.IsType(t, &changeLoop{}, ctrl.changeLoop)
	cl := ctrl.changeLoop.(*changeLoop)
	require.Equal(t, p, cl.period)
}

func TestController(t *testing.T) {
	ctx := context.Background()

	changeLoop := &testutils.MockStartStopper{}
	changeLoop.On("Start", ctx).Return(nil).Once()
	changeLoop.On("Stop", ctx).Return(nil).Once()

	eventLoop := &testutils.MockStartStopper{}
	eventLoop.On("Start", ctx).Return(nil).Once()
	eventLoop.On("Stop", ctx).Return(nil).Once()

	ctrl, _ := New(time.Second)
	ctrl.changeLoop = changeLoop
	ctrl.eventLoop = eventLoop

	err := ctrl.Start(ctx)
	require.NoError(t, err)

	err = ctrl.Stop(ctx)
	require.NoError(t, err)

	select {
	case <-ctrl.Done():
	case <-time.After(time.Second):
		t.Fatal("Runner not stopped")
	}

	err = ctrl.Err(ctx)
	require.NoError(t, err)

	changeLoop.AssertExpectations(t)
	eventLoop.AssertExpectations(t)
}
