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
	"sync"
)

// Map maps string keys to StartStopper.
type Map interface {
	sync.Locker

	// AddAndStart adds a StartStopper to the store and starts it.
	AddAndStart(ctx context.Context, key string, ss StartStopper) (changed bool, err error)

	// UpdateAndRestart stops the old version of the StartStopper, replaces it
	// with the started new one.
	UpdateAndRestart(ctx context.Context, key string, ss StartStopper) (changed bool, err error)

	// DeleteAndStop stops the StartStopper and deletes it. If it is not in the
	// store no change is made.
	DeleteAndStop(ctx context.Context, key string) (changed bool, err error)

	// Get the StartStopper of given key.
	Get(key string) (ss StartStopper, ok bool)

	// ForEach executes the given function for each StartStopper passing it and
	// its key as arguments.
	ForEach(func(key string, ss StartStopper))

	// Len returns how many StartStopper are stored.
	Len() int
}

type inMemoryMap struct {
	sync.RWMutex
	store map[string]StartStopper
}

// NewInMemoryMap creates a new Map that stores
// StartStopper in memory.
func NewInMemoryMap() Map {
	return &inMemoryMap{
		store: make(map[string]StartStopper),
	}
}

func (i *inMemoryMap) Len() int {
	return len(i.store)
}

func (i *inMemoryMap) Get(key string) (ss StartStopper, ok bool) {
	ss, ok = i.store[key]
	return
}

func (i *inMemoryMap) AddAndStart(ctx context.Context, key string, ss StartStopper) (changed bool, err error) {
	_, ok := i.store[key]

	if !ok {
		err = ss.Start(ctx)
		if err != nil {
			return
		}
		i.store[key] = ss
		changed = true
	}

	return
}

func (i *inMemoryMap) UpdateAndRestart(ctx context.Context, key string, ss StartStopper) (changed bool, err error) {
	changed, err = i.DeleteAndStop(ctx, key)
	if err != nil {
		return
	}

	changed2, err := i.AddAndStart(ctx, key, ss)
	changed = changed || changed2

	return
}

func (i *inMemoryMap) DeleteAndStop(ctx context.Context, key string) (changed bool, err error) {
	ss, ok := i.store[key]

	if ok {
		err = ss.Stop(ctx)
		if err != nil {
			return
		}
		delete(i.store, key)
		changed = true
	}

	return
}

func (i *inMemoryMap) ForEach(f func(key string, ss StartStopper)) {
	for key, ss := range i.store {
		f(key, ss)
	}
}
