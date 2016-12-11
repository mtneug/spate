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

// Package testutils contains helper code for testing startstopper.
package testutils

import (
	"context"
	"sync"

	"github.com/mtneug/pkg/startstopper"
	"github.com/stretchr/testify/mock"
)

// MockStartStopper is a mocked StartStopper.
type MockStartStopper struct {
	mock.Mock
}

// Start implements interface.
func (m *MockStartStopper) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Stop implements interface.
func (m *MockStartStopper) Stop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Done implements interface.
func (m *MockStartStopper) Done() <-chan struct{} {
	args := m.Called()
	return args.Get(0).(<-chan struct{})
}

// Err implements interface.
func (m *MockStartStopper) Err(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockRunner is a mocked Runner.
type MockRunner struct {
	mock.Mock
}

// Run implements interface.
func (m *MockRunner) Run(ctx context.Context, stopChan <-chan struct{}) error {
	args := m.Called(ctx)
	<-stopChan
	return args.Error(0)
}

// MockMap is a mocked Map.
type MockMap struct {
	mock.Mock
	sync.RWMutex
}

// AddAndStart implements interface.
func (m *MockMap) AddAndStart(ctx context.Context, key string, ss startstopper.StartStopper) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

// UpdateAndRestart implements interface.
func (m *MockMap) UpdateAndRestart(ctx context.Context, key string, ss startstopper.StartStopper) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

// DeleteAndStop implements interface.
func (m *MockMap) DeleteAndStop(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

// Get implements interface.
func (m *MockMap) Get(key string) (startstopper.StartStopper, bool) {
	args := m.Called(key)
	return args.Get(0).(startstopper.StartStopper), args.Bool(1)
}

// ForEach implements interface.
func (m *MockMap) ForEach(f func(key string, ss startstopper.StartStopper)) {
	m.Called()
}

// Len implements interface.
func (m *MockMap) Len() int {
	args := m.Called()
	return args.Int(0)
}
