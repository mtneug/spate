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
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mtneug/pkg/startstopper"
	"github.com/mtneug/spate/api/types"
	"github.com/mtneug/spate/metric"
)

// ErrNoGoals indicates that no goals were specified.
var ErrNoGoals = errors.New("autoscaler: no goals")

// Goal consist of an observer and a target.
type Goal struct {
	Observer *metric.Observer
	Target   types.Target
}

// Autoscaler observes one Docker Swarm service and automatically scales it
// depending on defined metrics.
type Autoscaler struct {
	startstopper.StartStopper

	Service swarm.Service
	Update  bool
	Goals   []Goal

	Period                    time.Duration
	CooldownServiceCreated    time.Duration
	CooldownServiceUpdated    time.Duration
	CooldownServiceScaledUp   time.Duration
	CooldownServiceScaledDown time.Duration

	MaxReplicas uint64
	MinReplicas uint64
}

// New creates an autoscaler for the given service.
func New(srv swarm.Service, goals []Goal) (*Autoscaler, error) {
	if len(goals) == 0 {
		return nil, ErrNoGoals
	}

	a := &Autoscaler{
		Service: srv,
		Goals:   goals,
	}
	a.StartStopper = startstopper.NewGo(startstopper.RunnerFunc(a.run))
	return a, nil
}

func (a *Autoscaler) run(ctx context.Context, stopChan <-chan struct{}) error {
	log.Debug("Autoscaler started")
	defer log.Debug("Autoscaler stopped")

	var err error

	// TODO: Change data structues so that this allocation is unnecessary
	observer := make([]startstopper.StartStopper, len(a.Goals))
	for i, goal := range a.Goals {
		observer[i] = goal.Observer
	}
	observerGroup := startstopper.NewGroup(observer)

	err = observerGroup.Start(ctx)
	if err != nil {
		return err
	}

	if a.Update {
		a.cooldown(ctx, types.EventTypeServiceUpdated)
	}

loop:
	for {
		select {
		case <-time.After(a.Period):
			a.tick(ctx)
		case <-stopChan:
			break loop
		case <-ctx.Done():
			break loop
		}
	}

	err = observerGroup.Stop(ctx)
	if err != nil {
		return err
	}

	return err
}

func (a *Autoscaler) cooldown(ctx context.Context, et types.EventType) {
	// TODO: refactor to use map
	var d time.Duration
	switch et {
	case types.EventTypeServiceCreated:
		d = a.CooldownServiceCreated
	case types.EventTypeServiceUpdated:
		d = a.CooldownServiceUpdated
	case types.EventTypeServiceScaledUp:
		d = a.CooldownServiceScaledUp
	case types.EventTypeServiceScaledDown:
		d = a.CooldownServiceScaledDown
	}

	if d > 0 {
		log.Debug("Autoscaler cooldown after '" + et + "' started")
		select {
		case <-time.After(d):
		case <-ctx.Done():
		}
		log.Debug("Autoscaler cooldown after '" + et + "' stopped")
	}
}

func (a *Autoscaler) tick(ctx context.Context) {
	// TODO: implement
}
