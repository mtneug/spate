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
	"errors"
	"testing"
	"time"

	"github.com/mtneug/pkg/startstopper"
	"github.com/mtneug/pkg/startstopper/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestNew(t *testing.T) {
	t.Parallel()

	p := time.Second
	m := startstopper.NewInMemoryMap()

	ctrl, err := New(p, m)
	require.NotNil(t, ctrl)
	require.NoError(t, err)

	require.IsType(t, &eventLoop{}, ctrl.eventLoop)

	require.IsType(t, &changeLoop{}, ctrl.changeLoop)
	cl := ctrl.changeLoop.(*changeLoop)
	require.Equal(t, p, cl.period)
}

type ControllerTestSuite struct {
	suite.Suite

	ctx        context.Context
	eventLoop  *testutils.MockStartStopper
	changeLoop *testutils.MockStartStopper
	ctrl       *Controller
}

func (s *ControllerTestSuite) SetupTest() {
	s.ctx = context.Background()

	s.changeLoop = &testutils.MockStartStopper{}
	s.eventLoop = &testutils.MockStartStopper{}
	s.ctrl, _ = New(time.Second, startstopper.NewInMemoryMap())
	s.ctrl.changeLoop = s.changeLoop
	s.ctrl.eventLoop = s.eventLoop
}

func (s *ControllerTestSuite) TestRun() {
	s.changeLoop.On("Start", s.ctx).Return(nil).Once()
	s.changeLoop.On("Stop", s.ctx).Return(nil).Once()
	s.eventLoop.On("Start", s.ctx).Return(nil).Once()
	s.eventLoop.On("Stop", s.ctx).Return(nil).Once()

	err := s.ctrl.Start(s.ctx)
	require.NoError(s.T(), err)

	err = s.ctrl.Stop(s.ctx)
	require.NoError(s.T(), err)

	s.changeLoop.AssertExpectations(s.T())
	s.eventLoop.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestRunCLStartErr() {
	s.changeLoop.On("Start", s.ctx).Return(errors.New("cl start")).Once()

	err := s.ctrl.Start(s.ctx)
	require.NoError(s.T(), err)

	err = s.ctrl.Stop(s.ctx)
	require.NoError(s.T(), err)

	select {
	case <-s.ctrl.Done():
	case <-time.After(time.Second):
		s.T().Fatal("Did not stop after 1s")
	}

	err = s.ctrl.Err(s.ctx)
	require.EqualError(s.T(), err, "cl start")

	s.changeLoop.AssertExpectations(s.T())
	s.eventLoop.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestRunELStartErr() {
	s.changeLoop.On("Start", s.ctx).Return(nil).Once()
	s.changeLoop.On("Stop", s.ctx).Return(nil).Once()
	s.eventLoop.On("Start", s.ctx).Return(errors.New("el start")).Once()

	err := s.ctrl.Start(s.ctx)
	require.NoError(s.T(), err)

	err = s.ctrl.Stop(s.ctx)
	require.NoError(s.T(), err)

	select {
	case <-s.ctrl.Done():
	case <-time.After(time.Second):
		s.T().Fatal("Did not stop after 1s")
	}

	err = s.ctrl.Err(s.ctx)
	require.EqualError(s.T(), err, "el start")

	s.changeLoop.AssertExpectations(s.T())
	s.eventLoop.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestRunCLStopErr() {
	s.changeLoop.On("Start", s.ctx).Return(nil).Once()
	s.changeLoop.On("Stop", s.ctx).Return(errors.New("cl stop")).Once()
	s.eventLoop.On("Start", s.ctx).Return(nil).Once()
	s.eventLoop.On("Stop", s.ctx).Return(nil).Once()

	err := s.ctrl.Start(s.ctx)
	require.NoError(s.T(), err)

	err = s.ctrl.Stop(s.ctx)
	require.NoError(s.T(), err)

	select {
	case <-s.ctrl.Done():
	case <-time.After(time.Second):
		s.T().Fatal("Did not stop after 1s")
	}

	err = s.ctrl.Err(s.ctx)
	require.EqualError(s.T(), err, "cl stop")

	s.changeLoop.AssertExpectations(s.T())
	s.eventLoop.AssertExpectations(s.T())
}

func (s *ControllerTestSuite) TestRunELStopErr() {
	s.changeLoop.On("Start", s.ctx).Return(nil).Once()
	s.changeLoop.On("Stop", s.ctx).Return(nil).Once()
	s.eventLoop.On("Start", s.ctx).Return(nil).Once()
	s.eventLoop.On("Stop", s.ctx).Return(errors.New("el stop")).Once()

	err := s.ctrl.Start(s.ctx)
	require.NoError(s.T(), err)

	err = s.ctrl.Stop(s.ctx)
	require.NoError(s.T(), err)

	select {
	case <-s.ctrl.Done():
	case <-time.After(time.Second):
		s.T().Fatal("Did not stop after 1s")
	}

	err = s.ctrl.Err(s.ctx)
	require.EqualError(s.T(), err, "el stop")

	s.changeLoop.AssertExpectations(s.T())
	s.eventLoop.AssertExpectations(s.T())
}

func TestController(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ControllerTestSuite{})
}
