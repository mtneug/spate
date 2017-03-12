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

package startstopper

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrNotStarted indicates that the StartStopper has not been started yet.
	ErrNotStarted = errors.New("startstopper: not started")

	// ErrStarted indicates that the StartStopper was already started once. Note
	// that it does not indicate whether it is still running.
	ErrStarted = errors.New("startstopper: already started once")
)

// StartStopper of an object that can be started, stopped and waited for.
type StartStopper interface {
	// Start the StartStopper in the background. After the first invocation
	// ErrStarted is returned as error.  A stopped StartStopper cannot be started
	// again.
	Start(ctx context.Context) error

	// Stop the StartStopper. This function will block until the StartStopper is
	// done. If the StartStopper has not been running ErrNotStarted is returned as
	// error. If already stopped no error is returned.
	Stop(ctx context.Context) error

	// Done returnes a chanel which will be closed once the StartStopper is done.
	Done() <-chan struct{}

	// Err returnes the error of the stopped StartStopper. This function blocks if
	// it is still running. If the context is done earlier its error is returned
	// instead.
	Err(ctx context.Context) error
}

// Runner implement a Run function that stops
type Runner interface {
	// Run is called by a StartStopper. The function ought to return once stopChan
	// is closed or the context is done.
	Run(ctx context.Context, stopChan <-chan struct{}) error
}

// RunnerFunc is a Runner.
type RunnerFunc func(context.Context, <-chan struct{}) error

// Run implements Runner by calling itself.
func (f RunnerFunc) Run(ctx context.Context, stopChan <-chan struct{}) error {
	return f(ctx, stopChan)
}

// NewGo creates a new GoStartStopper for given Runner. The Runner is executed
// in a Goroutine.
//
// 	package main
//
// 	import (
// 		"context"
// 		"fmt"
// 		"time"
//
// 		"github.com/mtneug/pkg/startstopper"
// 	)
//
// 	func main() {
// 		ctx := context.Background()
//
// 		ss := startstopper.NewGo(startstopper.RunnerFunc(
// 			func(ctx context.Context, stopChan <-chan struct{}) error {
// 				for {
// 					select {
// 					case <-time.After(time.Second):
// 						fmt.Println("Hello World")
// 					case <-stopChan:
// 						return nil
// 					case <-ctx.Done():
// 						return ctx.Err()
// 					}
// 				}
// 			}))
//
// 		ss.Start(ctx)
// 		time.AfterFunc(10*time.Second, func() {
// 			// Stop also blocks
// 			ss.Stop(ctx)
// 		})
// 		<-ss.Done()
// 	}
func NewGo(r Runner) StartStopper {
	return &goStartStopper{
		Runner:    r,
		startChan: make(chan struct{}),
		stopChan:  make(chan struct{}),
		doneChan:  make(chan struct{}),
	}
}

type goStartStopper struct {
	Runner    Runner
	err       error
	startOnce sync.Once
	stopOnce  sync.Once
	startChan chan struct{}
	stopChan  chan struct{}
	doneChan  chan struct{}
}

func (ss *goStartStopper) Start(ctx context.Context) error {
	err := ErrStarted

	ss.startOnce.Do(func() {
		close(ss.startChan)
		go func() {
			ss.err = ss.Runner.Run(ctx, ss.stopChan)
			close(ss.doneChan)
		}()
		err = nil
	})

	return err
}

func (ss *goStartStopper) Stop(ctx context.Context) error {
	select {
	case <-ss.startChan:
	default:
		return ErrNotStarted
	}

	ss.stopOnce.Do(func() {
		close(ss.stopChan)
	})

	select {
	case <-ss.Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (ss *goStartStopper) Done() <-chan struct{} {
	return ss.doneChan
}

func (ss *goStartStopper) Err(ctx context.Context) error {
	select {
	case <-ss.Done():
		return ss.err
	case <-ctx.Done():
		return ctx.Err()
	}
}
