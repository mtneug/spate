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

package metric

import (
	"context"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mtneug/pkg/reducer"
	"github.com/mtneug/pkg/startstopper"
)

// Observer observes one metric and aggregate measurements.
type Observer struct {
	startstopper.StartStopper

	measures []float64
	i        uint8
	mutex    sync.RWMutex

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
	log.Debug("Observer loop started")
	defer log.Debug("Observer loop stopped")

	o.measures = make([]float64, 0, o.AggregationAmount)

	o.tick(ctx)
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

// AggregatedMeasure of current measurements.
func (o *Observer) AggregatedMeasure() (float64, error) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	return o.Reducer.Reduce(o.measures)
}

func (o *Observer) tick(ctx context.Context) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	measure, err := o.Measurer.Measure(ctx)
	if err != nil {
		log.WithError(err).Warn("Measuring failed")
		return
	}

	log.Debugf("Measured %f", measure)

	if len(o.measures) < int(o.AggregationAmount) {
		o.measures = append(o.measures, measure)
	} else {
		o.measures[o.i] = measure
		o.i = (o.i + 1) % o.AggregationAmount
	}
}
