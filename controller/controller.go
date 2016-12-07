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
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
)

var (
	// ErrNotStarted indicates that the controller has not been started yet.
	ErrNotStarted = errors.New("controller: not started")

	// ErrStarted indicates that the controller was already started once. Note
	// that it does not indicate whether it is still running.
	ErrStarted = errors.New("controller: already started once")
)

// Config for a controller.
type Config struct {
}

// Controller monitors Docker Swarm services and scales them if needed.
type Controller struct {
	config *Config
	client *client.Client
	err    error

	startOnce sync.Once
	stopOnce  sync.Once

	startChan chan struct{}
	stopChan  chan struct{}
	doneChan  chan struct{}
}

// New creates a new controller.
func New(c *Config) (*Controller, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	a := &Controller{
		config:    c,
		client:    cli,
		startChan: make(chan struct{}),
		stopChan:  make(chan struct{}),
		doneChan:  make(chan struct{}),
	}

	return a, nil
}

// Start the controller in the background.
func (c *Controller) Start(ctx context.Context) error {
	err := ErrStarted

	c.startOnce.Do(func() {
		log.Info("Starting controller")
		go func() {
			close(c.startChan)
			log.Info("Controller started")
			c.err = c.run(ctx)
			close(c.doneChan)
			log.Info("Controller stopped")
		}()
		err = nil
	})

	return err
}

// Stop the controller. After it is stopped, this controller cannot be started
// again.
func (c *Controller) Stop(ctx context.Context) error {
	select {
	case <-c.startChan:
	default:
		return ErrNotStarted
	}

	c.stopOnce.Do(func() {
		log.Info("Stopping controller")
		close(c.stopChan)
	})

	select {
	case <-c.doneChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Err returns an error object after the controller has stopped.
func (c *Controller) Err(ctx context.Context) error {
	select {
	case <-c.doneChan:
		return c.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Controller) run(ctx context.Context) (err error) {
	for {
		select {
		case <-c.stopChan:
			// TODO: clean up
			return nil
		}
	}
}
