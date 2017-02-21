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
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/mtneug/pkg/startstopper"
	"github.com/mtneug/spate/autoscaler"
	"github.com/mtneug/spate/docker"
	"github.com/mtneug/spate/event"
)

const labelSpate = "de.mtneug.spate"

var serviceListOptions types.ServiceListOptions

func init() {
	f := filters.NewArgs()
	f.Add("label", fmt.Sprintf("%s=%s", labelSpate, "enable"))
	serviceListOptions = types.ServiceListOptions{Filters: f}
}

type serviceEventPublisher struct {
	startstopper.StartStopper

	period         time.Duration
	eventQueue     chan<- event.Event
	autoscalersMap startstopper.Map

	// stored so that it doesn't need to be reallocated
	seen map[string]bool
}

func newServiceEventPublisher(p time.Duration, eq chan<- event.Event, m startstopper.Map) *serviceEventPublisher {
	sep := &serviceEventPublisher{
		period:         p,
		eventQueue:     eq,
		autoscalersMap: m,
		seen:           make(map[string]bool),
	}
	sep.StartStopper = startstopper.NewGo(startstopper.RunnerFunc(sep.run))
	return sep
}

func (sep *serviceEventPublisher) run(ctx context.Context, stopChan <-chan struct{}) error {
	log.Debug("Service event publisher started")
	defer log.Debug("Service event publisher stopped")

	sep.tick(ctx)
	for {
		select {
		case <-time.After(sep.period):
			sep.tick(ctx)
		case <-stopChan:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (sep *serviceEventPublisher) tick(ctx context.Context) {
	services, err := docker.C.ServiceList(ctx, serviceListOptions)
	if err != nil {
		log.WithError(err).Error("Failed to get list of services")
		return
	}

	for _, srv := range services {
		sep.seen[srv.ID] = true

		ss, present := sep.autoscalersMap.Get(srv.ID)
		if !present {
			// Add
			sep.eventQueue <- event.New(event.TypeServiceCreated, srv)
		} else {
			a := ss.(*autoscaler.Autoscaler)
			a.RLock()
			if a.Service.Version.Index < srv.Version.Index {
				// Update
				sep.eventQueue <- event.New(event.TypeServiceUpdated, srv)
			}
			a.RUnlock()
		}
	}

	sep.autoscalersMap.ForEach(func(id string, ss startstopper.StartStopper) {
		if !sep.seen[id] {
			// Delete
			a := ss.(*autoscaler.Autoscaler)
			a.RLock()
			sep.eventQueue <- event.New(event.TypeServiceDeleted, a.Service)
			a.RUnlock()
		}
		delete(sep.seen, id)
	})
}
