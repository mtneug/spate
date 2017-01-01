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

package label

import (
	"strconv"
	"time"

	"github.com/mtneug/spate/metric"
)

var (
	// DefaultObserverPeriod for new observers.
	DefaultObserverPeriod = 30 * time.Second

	// DefaultObserverAggregationAmount for new observers.
	DefaultObserverAggregationAmount uint8 = 5
)

// ParseObserver parses the labels and sets the corresponding values for given
// observer.
func ParseObserver(observer *metric.Observer, labels map[string]string) error {
	observerPeriodStr, ok := labels[MetricObserverPeriodSuffix]
	if !ok {
		observer.Period = DefaultObserverPeriod
	} else {
		observerPeriod, err := time.ParseDuration(observerPeriodStr)
		if err != nil {
			return ErrInvalidDuration
		}
		observer.Period = observerPeriod
	}

	aggregationAmountStr, ok := labels[MetricAggregationAmountSuffix]
	if !ok {
		observer.AggregationAmount = DefaultObserverAggregationAmount
	} else {
		aggregationAmount, err := strconv.ParseUint(aggregationAmountStr, 10, 8)
		if err != nil {
			return ErrInvalidUint
		}
		observer.AggregationAmount = uint8(aggregationAmount)
	}

	return nil
}
