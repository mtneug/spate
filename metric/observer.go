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

package metric

import (
	"context"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mtneug/pkg/reducer"
	"github.com/mtneug/pkg/startstopper"
)

// Observer observes one metric and aggregate measurements.
type Observer struct {
	startstopper.StartStopper

	Measurer          Measurer
	Reducer           reducer.Reducer
	Period            time.Duration
	AggregationAmount uint8
}

// NewObserver creates a new Observer for given measurer and reducer.
func NewObserver(m Measurer, r reducer.Reducer) *Observer {
	o := &Observer{
		Measurer: m,
		Reducer:  r,
	}
	o.StartStopper = startstopper.NewGo(startstopper.RunnerFunc(o.run))

	return o
}

func (o *Observer) run(ctx context.Context, stopChan <-chan struct{}) error {
	log.Debug("Observer started")
	defer log.Debug("Observer stopped")

	for {
		select {
		case <-time.After(o.Period):
			o.tick(ctx)
		case <-stopChan:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (o *Observer) tick(ctx context.Context) {
	// TODO: implement
}
