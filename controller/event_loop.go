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

	"docker.io/go-docker/api/types/swarm"
	log "github.com/Sirupsen/logrus"
	"github.com/mtneug/pkg/startstopper"
	"github.com/mtneug/spate/event"
	"github.com/mtneug/spate/label"
)

type eventLoop struct {
	startstopper.StartStopper

	eventQueue  <-chan event.Event
	autoscalers startstopper.Map
}

func newEventLoop(eq <-chan event.Event, m startstopper.Map) *eventLoop {
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

func (el *eventLoop) handleEvent(ctx context.Context, e event.Event) {
	log.Debugf("Received %s event", e.Type)

	var changed bool
	var err error

	switch e.Type {
	case event.TypeServiceCreated:
		changed, err = el.addAutoscaler(ctx, e.Service)
		if err != nil {
			log.WithError(err).Error("Could not add autoscaler")
		} else {
			log.Info("Autoscaler added")
		}
	case event.TypeServiceUpdated:
		changed, err = el.updateAutoscaler(ctx, e.Service)
		if err != nil {
			log.WithError(err).Error("Could not update autoscaler")
		} else {
			log.Info("Autoscaler updated")
		}
	case event.TypeServiceDeleted:
		changed, err = el.deleteAutoscaler(ctx, e.Service)
		if err != nil {
			log.WithError(err).Error("Could not delete autoscaler")
		} else {
			log.Info("Autoscaler deleted")
		}
	}

	if !changed {
		log.Debug("State unchanged")
	}
}

func (el *eventLoop) addAutoscaler(ctx context.Context, srv swarm.Service) (bool, error) {
	a, err := label.ConstructAutoscaler(srv)
	if err != nil {
		return false, err
	}

	return el.autoscalers.AddAndStart(ctx, srv.ID, a)
}

func (el *eventLoop) updateAutoscaler(ctx context.Context, srv swarm.Service) (bool, error) {
	a, err := label.ConstructAutoscaler(srv)
	if err != nil {
		return false, err
	}
	a.Update = true

	return el.autoscalers.UpdateAndRestart(ctx, srv.ID, a)
}

func (el *eventLoop) deleteAutoscaler(ctx context.Context, srv swarm.Service) (bool, error) {
	return el.autoscalers.DeleteAndStop(ctx, srv.ID)
}
