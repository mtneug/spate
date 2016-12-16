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
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/mtneug/pkg/startstopper"
	"github.com/mtneug/spate/api/types"
	"github.com/mtneug/spate/metric"
)

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

	Period                 time.Duration
	CooldownScaledUp       time.Duration
	CooldownScaledDown     time.Duration
	CooldownServiceAdded   time.Duration
	CooldownServiceUpdated time.Duration

	MaxReplicas uint64
	MinReplicas uint64
}

// New creates an autoscaler for the given service.
func New(srv swarm.Service, goal []Goal) *Autoscaler {
	a := &Autoscaler{
		Service: srv,
		Goals:   goal,
	}
	a.StartStopper = startstopper.NewGo(startstopper.RunnerFunc(a.run))
	return a
}

func (a *Autoscaler) run(ctx context.Context, stopChan <-chan struct{}) error {
	// TODO: implement
	return nil
}
