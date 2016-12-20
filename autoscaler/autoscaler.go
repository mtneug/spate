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
	"math"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mtneug/pkg/startstopper"
	"github.com/mtneug/spate/api/types"
	"github.com/mtneug/spate/docker"
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
	sync.RWMutex

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

	// start observer
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

	// service created/updated cooldown
	if a.Update {
		a.cooldown(ctx, stopChan, types.EventTypeServiceUpdated)
	} else {
		a.cooldown(ctx, stopChan, types.EventTypeServiceCreated)
	}

	// start autoscaling
loop:
	for {
		select {
		case <-time.After(a.Period):
			a.tick(ctx, stopChan)
		case <-stopChan:
			break loop
		case <-ctx.Done():
			break loop
		}
	}

	// stop observer
	err = observerGroup.Stop(ctx)
	if err != nil {
		return err
	}

	return err
}

func (a *Autoscaler) cooldown(ctx context.Context, stopChan <-chan struct{}, et types.EventType) {
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
		case <-stopChan:
		case <-ctx.Done():
		}
		log.Debug("Autoscaler cooldown after '" + et + "' stopped")
	}
}

func (a *Autoscaler) tick(ctx context.Context, stopChan <-chan struct{}) {
	var (
		err  error
		ag   float64
		once sync.Once
	)

	a.Lock()
	unlock := func() { once.Do(func() { a.Unlock() }) }
	defer unlock()

	a.Service, _, err = docker.C.ServiceInspectWithRaw(ctx, a.Service.ID)
	if err != nil {
		log.WithError(err).Warn("Service inspection failed")
		return
	}

	srvMode := a.Service.Spec.Mode
	if srvMode.Replicated == nil || srvMode.Replicated.Replicas == nil {
		// TODO: check this beforehand
		log.Warn("Service is not a replicated service")
		return
	}

	currentScale := float64(*srvMode.Replicated.Replicas)
	desiredScale := float64(a.MinReplicas)

	for _, goal := range a.Goals {
		ag, err = goal.Observer.AggregatedMeasure()
		if err != nil {
			log.WithError(err).Warn("Measure aggregation failed")
			return
		}

		// deviation acceptable?
		deviation := (ag / currentScale) - goal.Target.Value
		log.Debugf("Deviation from target is %f", deviation)
		if -goal.Target.LowerDeviation <= deviation && deviation <= goal.Target.UpperDeviation {
			break
		}

		// update desired scale
		desiredScale = math.Max(desiredScale, math.Ceil(ag/goal.Target.Value))
	}

	newScale := uint64(math.Min(desiredScale, float64(a.MaxReplicas)))
	srvMode.Replicated.Replicas = &newScale

	log.Debugf("Current #: %d Desired #: %d Constrained #: %d", uint64(currentScale), uint64(desiredScale), newScale)

	if currentScale == float64(newScale) {
		return
	}

	err = docker.C.ServiceUpdate(ctx, a.Service.ID, a.Service.Version, a.Service.Spec, dockerTypes.ServiceUpdateOptions{})
	if err != nil {
		log.WithError(err).Warn("Service scaling failed")
		return
	}

	a.Service, _, err = docker.C.ServiceInspectWithRaw(ctx, a.Service.ID)
	if err != nil {
		log.WithError(err).Warn("Service inspection failed")
	}

	// Autoscaler should not be locked during cooldown times
	unlock()

	if currentScale < float64(newScale) {
		log.Infof("Service scaled up to %d replica(s)", newScale)
		a.cooldown(ctx, stopChan, types.EventTypeServiceScaledUp)
	} else {
		log.Infof("Service scaled down to %d replica(s)", newScale)
		a.cooldown(ctx, stopChan, types.EventTypeServiceScaledDown)
	}
}
