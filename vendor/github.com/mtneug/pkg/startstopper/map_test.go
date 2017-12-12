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

	"github.com/mtneug/pkg/startstopper"
	"github.com/mtneug/pkg/startstopper/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type InMemoryMapTestSuite struct {
	suite.Suite

	ctx context.Context
	err error
	ss1 *testutils.MockStartStopper
	ss2 *testutils.MockStartStopper
	m   startstopper.Map
}

func (s *InMemoryMapTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.err = errors.New("test error")
	s.ss1 = &testutils.MockStartStopper{}
	s.ss2 = &testutils.MockStartStopper{}
	s.m = startstopper.NewInMemoryMap()
}

func (s *InMemoryMapTestSuite) TestLen() {
	s.ss1.On("Start", s.ctx).Return(nil).Once()
	s.ss2.On("Start", s.ctx).Return(nil).Once()

	require.Equal(s.T(), 0, s.m.Len())
	_, _ = s.m.AddAndStart(s.ctx, "test1", s.ss1)
	require.Equal(s.T(), 1, s.m.Len())
	_, _ = s.m.AddAndStart(s.ctx, "test2", s.ss2)
	require.Equal(s.T(), 2, s.m.Len())

	s.ss1.AssertExpectations(s.T())
	s.ss2.AssertExpectations(s.T())
}

func (s *InMemoryMapTestSuite) TestGetNotInserted() {
	ss, ok := s.m.Get("test")
	require.Nil(s.T(), ss)
	require.False(s.T(), ok)
}

func (s *InMemoryMapTestSuite) TestAddAndStart() {
	s.ss1.On("Start", s.ctx).Return(nil).Once()

	changed, err := s.m.AddAndStart(s.ctx, "test", s.ss1)
	require.NoError(s.T(), err)
	require.True(s.T(), changed)

	ss, ok := s.m.Get("test")
	require.Equal(s.T(), s.ss1, ss)
	require.True(s.T(), ok)

	s.ss1.AssertExpectations(s.T())
}

func (s *InMemoryMapTestSuite) TestAddAndStartErrStart() {
	s.ss1.On("Start", s.ctx).Return(s.err).Once()

	changed, err := s.m.AddAndStart(s.ctx, "test", s.ss1)
	require.EqualError(s.T(), err, s.err.Error())
	require.False(s.T(), changed)

	ss, ok := s.m.Get("test")
	require.Nil(s.T(), ss)
	require.False(s.T(), ok)

	s.ss1.AssertExpectations(s.T())
}

func (s *InMemoryMapTestSuite) TestDeleteAndStopNotInserted() {
	changed, err := s.m.DeleteAndStop(s.ctx, "test")
	require.NoError(s.T(), err)
	require.False(s.T(), changed)

	s.ss1.AssertExpectations(s.T())
	s.ss1.AssertNotCalled(s.T(), "Stop")
}

func (s *InMemoryMapTestSuite) TestDeleteAndStopInserted() {
	s.ss1.On("Start", s.ctx).Return(nil).Once()
	s.ss1.On("Stop", s.ctx).Return(nil).Once()
	_, _ = s.m.AddAndStart(s.ctx, "test", s.ss1)

	changed, err := s.m.DeleteAndStop(s.ctx, "test")
	require.NoError(s.T(), err)
	require.True(s.T(), changed)

	ss, ok := s.m.Get("test")
	require.Nil(s.T(), ss)
	require.False(s.T(), ok)

	s.ss1.AssertExpectations(s.T())
}

func (s *InMemoryMapTestSuite) TestDeleteAndStopInsertedErrStart() {
	s.ss1.On("Start", s.ctx).Return(nil).Once()
	s.ss1.On("Stop", s.ctx).Return(s.err).Once()
	_, _ = s.m.AddAndStart(s.ctx, "test", s.ss1)

	changed, err := s.m.DeleteAndStop(s.ctx, "test")
	require.EqualError(s.T(), err, s.err.Error())
	require.False(s.T(), changed)

	ss, ok := s.m.Get("test")
	require.Equal(s.T(), s.ss1, ss)
	require.True(s.T(), ok)

	s.ss1.AssertExpectations(s.T())
}

func (s *InMemoryMapTestSuite) TestUpdateAndRestartNotInserted() {
	s.ss1.On("Start", s.ctx).Return(nil).Once()

	changed, err := s.m.UpdateAndRestart(s.ctx, "test", s.ss1)
	require.NoError(s.T(), err)
	require.True(s.T(), changed)

	ss, ok := s.m.Get("test")
	require.Equal(s.T(), s.ss1, ss)
	require.True(s.T(), ok)

	s.ss1.AssertExpectations(s.T())
}

func (s *InMemoryMapTestSuite) TestUpdateAndRestartNotInsertedErrStart() {
	s.ss1.On("Start", s.ctx).Return(s.err).Once()

	changed, err := s.m.UpdateAndRestart(s.ctx, "test", s.ss1)
	require.EqualError(s.T(), err, s.err.Error())
	require.False(s.T(), changed)

	ss, ok := s.m.Get("test")
	require.Nil(s.T(), ss)
	require.False(s.T(), ok)

	s.ss1.AssertExpectations(s.T())
}

func (s *InMemoryMapTestSuite) TestUpdateAndRestartInserted() {
	s.ss1.On("Start", s.ctx).Return(nil).Once()
	s.ss1.On("Stop", s.ctx).Return(nil).Once()
	s.ss2.On("Start", s.ctx).Return(nil).Once()
	_, _ = s.m.AddAndStart(s.ctx, "test", s.ss1)

	changed, err := s.m.UpdateAndRestart(s.ctx, "test", s.ss2)
	require.NoError(s.T(), err)
	require.True(s.T(), changed)

	ss, ok := s.m.Get("test")
	require.Equal(s.T(), s.ss2, ss)
	require.True(s.T(), ok)

	s.ss1.AssertExpectations(s.T())
	s.ss2.AssertExpectations(s.T())
}

func (s *InMemoryMapTestSuite) TestUpdateAndRestartInsertedErrStop() {
	s.ss1.On("Start", s.ctx).Return(nil).Once()
	s.ss1.On("Stop", s.ctx).Return(s.err).Once()
	_, _ = s.m.AddAndStart(s.ctx, "test", s.ss1)

	changed, err := s.m.UpdateAndRestart(s.ctx, "test", s.ss2)
	require.EqualError(s.T(), s.err, err.Error())
	require.False(s.T(), changed)

	ss, ok := s.m.Get("test")
	require.Equal(s.T(), s.ss1, ss)
	require.True(s.T(), ok)

	s.ss1.AssertExpectations(s.T())
	s.ss2.AssertExpectations(s.T())
}

func (s *InMemoryMapTestSuite) TestUpdateAndRestartInsertedErrStart() {
	s.ss1.On("Start", s.ctx).Return(nil).Once()
	s.ss1.On("Stop", s.ctx).Return(nil).Once()
	s.ss2.On("Start", s.ctx).Return(s.err).Once()
	_, _ = s.m.AddAndStart(s.ctx, "test", s.ss1)

	changed, err := s.m.UpdateAndRestart(s.ctx, "test", s.ss2)
	require.EqualError(s.T(), s.err, err.Error())
	require.True(s.T(), changed)

	ss, ok := s.m.Get("test")
	require.Nil(s.T(), ss)
	require.False(s.T(), ok)

	s.ss1.AssertExpectations(s.T())
	s.ss2.AssertExpectations(s.T())
}

func (s *InMemoryMapTestSuite) TestForEach() {
	s.ss1.On("Start", s.ctx).Return(nil).Once()
	s.ss2.On("Start", s.ctx).Return(nil).Once()
	_, _ = s.m.AddAndStart(s.ctx, "test1", s.ss1)
	_, _ = s.m.AddAndStart(s.ctx, "test2", s.ss2)

	called := 0
	notSeen := map[string]*testutils.MockStartStopper{
		"test1": s.ss1,
		"test2": s.ss2,
	}

	s.m.ForEach(func(key string, ss startstopper.StartStopper) {
		called++
		require.Equal(s.T(), notSeen[key], ss)
		delete(notSeen, key)
	})

	require.Equal(s.T(), 2, called)
	require.Empty(s.T(), notSeen)

	s.ss1.AssertExpectations(s.T())
	s.ss2.AssertExpectations(s.T())
}

func TestInMemoryMap(t *testing.T) {
	t.Parallel()
	suite.Run(t, &InMemoryMapTestSuite{})
}
