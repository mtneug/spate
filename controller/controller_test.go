// Copyright (c) 2016 Matthias Neugebauer <mtneug@mailbox.org>
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

	ctrl := New(p)
	require.NotNil(t, ctrl)

	require.IsType(t, &eventLoop{}, ctrl.eventLoop)

	require.IsType(t, &serviceEventPublisher{}, ctrl.serviceEventPublisher)
	cl := ctrl.serviceEventPublisher.(*serviceEventPublisher)
	require.Equal(t, p, cl.period)
}

func TestController(t *testing.T) {
	ctx := context.Background()

	sep := &testutils.MockStartStopper{}
	sep.On("Start", ctx).Return(nil).Once()
	sep.On("Stop", ctx).Return(nil).Once()

	el := &testutils.MockStartStopper{}
	el.On("Start", ctx).Return(nil).Once()
	el.On("Stop", ctx).Return(nil).Once()

	ctrl := New(time.Second)
	ctrl.serviceEventPublisher = sep
	ctrl.eventLoop = el

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

	sep.AssertExpectations(t)
	el.AssertExpectations(t)
}
