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
	"bytes"
	"context"
	"sync"
)

// Group StartStopper together and run them as one.
type Group struct {
	StartStopper

	sss []StartStopper
}

// GroupError is the concrete error type returned by Group structs. It bundles
// the errors that could be returned by the StartStoppers.
type GroupError struct {
	Errors []error
}

// Error implements the error interface.
func (g GroupError) Error() string {
	var buffer bytes.Buffer

	for _, err := range g.Errors {
		if err != nil {
			buffer.WriteString(err.Error())
			buffer.WriteString(", ")
		}
	}

	str := buffer.String()
	return str[:len(str)-2]
}

// NewGroup creates a new group.
func NewGroup(sss []StartStopper) *Group {
	group := &Group{sss: sss}
	group.StartStopper = NewGo(RunnerFunc(group.run))
	return group
}

func (g *Group) run(ctx context.Context, stopChan <-chan struct{}) error {
	var once sync.Once
	errorOccured := false
	setErrorOccured := func() {
		once.Do(func() {
			errorOccured = true
		})
	}

	groupErr := GroupError{
		Errors: make([]error, len(g.sss)),
	}

	for i, ss := range g.sss {
		groupErr.Errors[i] = ss.Start(ctx)
		if groupErr.Errors[i] != nil {
			setErrorOccured()
		}
	}

	if errorOccured {
		return groupErr
	}

	<-stopChan

	var wg sync.WaitGroup
	wg.Add(len(g.sss))

	for i, ss := range g.sss {
		_i := i
		_ss := ss
		go func() {
			groupErr.Errors[_i] = _ss.Stop(ctx)
			if groupErr.Errors[_i] != nil {
				setErrorOccured()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	if errorOccured {
		return groupErr
	}
	return nil
}
