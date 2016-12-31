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
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mtneug/pkg/startstopper"
	"github.com/mtneug/spate/event"
)

// Controller monitors Docker Swarm services and scales them if needed.
type Controller struct {
	startstopper.StartStopper

	autoscalers startstopper.Map
	eventQueue  chan event.Event
	eventLoop   startstopper.StartStopper
	changeLoop  startstopper.StartStopper
}

// New creates a new controller.
func New(p time.Duration) (*Controller, error) {
	m := startstopper.NewInMemoryMap()
	eq := make(chan event.Event, 20)
	ctrl := &Controller{
		autoscalers: m,
		eventQueue:  eq,
		eventLoop:   newEventLoop(eq, m),
		changeLoop:  newChangeLoop(p, eq, m),
	}
	ctrl.StartStopper = startstopper.NewGo(startstopper.RunnerFunc(ctrl.run))

	return ctrl, nil
}

func (c *Controller) run(ctx context.Context, stopChan <-chan struct{}) error {
	log.Debug("Controller loop started")
	defer log.Debug("Controller loop stopped")

	group := startstopper.NewGroup([]startstopper.StartStopper{
		c.changeLoop,
		c.eventLoop,
	})

	_ = group.Start(ctx)

	select {
	case <-stopChan:
	case <-ctx.Done():
	}

	_ = group.Stop(ctx)
	err := group.Err(ctx)

	c.autoscalers.ForEach(func(key string, autoscaler startstopper.StartStopper) {
		_ = autoscaler.Stop(ctx)
	})

	return err
}
