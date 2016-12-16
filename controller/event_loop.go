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

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mtneug/pkg/startstopper"
	"github.com/mtneug/spate/api/types"
	"github.com/mtneug/spate/autoscaler"
)

type eventLoop struct {
	startstopper.StartStopper

	eventQueue  <-chan types.Event
	autoscalers startstopper.Map
}

func newEventLoop(eq <-chan types.Event, m startstopper.Map) *eventLoop {
	el := &eventLoop{
		eventQueue:  eq,
		autoscalers: m,
	}
	el.StartStopper = startstopper.NewGo(startstopper.RunnerFunc(el.run))
	return el
}

func (el *eventLoop) run(ctx context.Context, stopChan <-chan struct{}) error {
	log.Debug("Event loop started")
	defer log.Debug("Event loop stopped")

	for {
		select {
		case e := <-el.eventQueue:
			el.handleEvent(ctx, e)
		case <-stopChan:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (el *eventLoop) handleEvent(ctx context.Context, e types.Event) {
	log.Debugf("Received %s event", e.Type)

	srv, ok := e.Object.(swarm.Service)
	if !ok {
		log.
			WithError(errors.New("controller: type assertion failed")).
			Error("Failed to get service")
		return
	}

	var changed bool
	var err error

	switch e.Type {
	case types.EventTypeServiceCreated:
		changed, err = el.addAutoscaler(ctx, srv)
		if err != nil {
			log.WithError(err).Error("Could not add autoscaler")
		}
	case types.EventTypeServiceUpdated:
		changed, err = el.updateAutoscaler(ctx, srv)
		if err != nil {
			log.WithError(err).Error("Could not update autoscaler")
		}
	case types.EventTypeServiceDeleted:
		changed, err = el.deleteAutoscaler(ctx, srv)
		if err != nil {
			log.WithError(err).Error("Could not delete autoscaler")
		}
	}

	if !changed {
		log.Debug("State unchanged")
	}
}

func (el *eventLoop) addAutoscaler(ctx context.Context, srv swarm.Service) (bool, error) {
	a := autoscaler.New(srv, nil)
	return el.autoscalers.AddAndStart(ctx, srv.ID, a)
}

func (el *eventLoop) updateAutoscaler(ctx context.Context, srv swarm.Service) (bool, error) {
	a := autoscaler.New(srv, nil)
	a.Update = true

	return el.autoscalers.UpdateAndRestart(ctx, srv.ID, a)
}

func (el *eventLoop) deleteAutoscaler(ctx context.Context, srv swarm.Service) (bool, error) {
	return el.autoscalers.DeleteAndStop(ctx, srv.ID)
}
