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

package autoscaler

import (
	"context"
	"errors"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
)

var (
	// ErrNotStarted indicates that the autoscaler has not been started yet.
	ErrNotStarted = errors.New("autoscaler: not started")

	// ErrStarted indicates that the autoscaler was already started once. Note
	// that it does not indicate whether it is still running.
	ErrStarted = errors.New("autoscaler: already started once")
)

// Config for a autoscaler.
type Config struct {
}

// Autoscaler monitors Docker Swarm services and scales them if needed.
type Autoscaler struct {
	config *Config
	client *client.Client
	err    error

	startOnce sync.Once
	stopOnce  sync.Once

	startChan chan struct{}
	stopChan  chan struct{}
	doneChan  chan struct{}
}

// New creates a new autoscaler.
func New(c *Config) (*Autoscaler, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	a := &Autoscaler{
		config:    c,
		client:    cli,
		startChan: make(chan struct{}),
		stopChan:  make(chan struct{}),
		doneChan:  make(chan struct{}),
	}

	return a, nil
}

// Start the autoscaler in the background.
func (a *Autoscaler) Start(ctx context.Context) error {
	err := ErrStarted

	a.startOnce.Do(func() {
		log.Info("Starting autostaler")
		close(a.startChan)
		go a.run(ctx)
		err = nil
	})

	return err
}

// Stop the autoscaler. After it is stopped, this autoscaler cannot be started
// again.
func (a *Autoscaler) Stop(ctx context.Context) error {
	select {
	case <-a.startChan:
	default:
		return ErrNotStarted
	}

	a.stopOnce.Do(func() {
		log.Info("Stopping autostaler")
		close(a.stopChan)
	})

	select {
	case <-a.doneChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Err returns an error object after the autoscaler has stopped.
func (a *Autoscaler) Err(ctx context.Context) error {
	select {
	case <-a.doneChan:
		return a.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *Autoscaler) run(ctx context.Context) (err error) {
	defer func() { a.err = err }()
	defer func() { log.Info("Autoscaler stopped") }()

	log.Info("Autoscaler started")
	for {
		select {
		case <-a.stopChan:
			// TODO: clean up
			close(a.doneChan)
			return nil
		}
	}
}
