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

package label

import (
	"math"
	"strconv"
	"time"

	"github.com/mtneug/spate/autoscaler"
)

var (
	// DefaultAutoscalerPeriod for autoscaler.
	DefaultAutoscalerPeriod = 30 * time.Second

	// DefaultCooldownScaledUp for autoscaler.
	DefaultCooldownScaledUp = 3 * time.Minute

	// DefaultCooldownScaledDown for autoscaler.
	DefaultCooldownScaledDown = 5 * time.Minute

	// DefaultCooldownServiceAdded for autoscaler.
	DefaultCooldownServiceAdded = 0 * time.Second

	// DefaultCooldownServiceUpdated for autoscaler.
	DefaultCooldownServiceUpdated = 0 * time.Second
)

// ParseAutoscaler parses the labels and sets the corresponding values for given
// autoscaler.
func ParseAutoscaler(a *autoscaler.Autoscaler, labels map[string]string) error {
	// period
	autoscalerPeriodStr, ok := labels[AutoscalerPeriod]
	if !ok {
		a.Period = DefaultAutoscalerPeriod
	} else {
		autoscalerPeriod, err := time.ParseDuration(autoscalerPeriodStr)
		if err != nil {
			return ErrInvalidDuration
		}
		a.Period = autoscalerPeriod
	}

	// cooldown ScaledUp
	autoscalerCooldownScaledUpStr, ok := labels[AutoscalerCooldownScaledUp]
	if !ok {
		a.CooldownServiceScaledUp = DefaultCooldownScaledUp
	} else {
		autoscalerCooldownScaledUp, err := time.ParseDuration(autoscalerCooldownScaledUpStr)
		if err != nil {
			return ErrInvalidDuration
		}
		a.CooldownServiceScaledUp = autoscalerCooldownScaledUp
	}

	// cooldown ScaledDown
	autoscalerCooldownScaledDownStr, ok := labels[AutoscalerCooldownScaledDown]
	if !ok {
		a.CooldownServiceScaledDown = DefaultCooldownScaledDown
	} else {
		autoscalerCooldownScaledDown, err := time.ParseDuration(autoscalerCooldownScaledDownStr)
		if err != nil {
			return ErrInvalidDuration
		}
		a.CooldownServiceScaledDown = autoscalerCooldownScaledDown
	}

	// cooldown ServiceAdded
	autoscalerCooldownServiceAddedStr, ok := labels[AutoscalerCooldownServiceAdded]
	if !ok {
		a.CooldownServiceCreated = DefaultCooldownServiceAdded
	} else {
		autoscalerCooldownServiceAdded, err := time.ParseDuration(autoscalerCooldownServiceAddedStr)
		if err != nil {
			return ErrInvalidDuration
		}
		a.CooldownServiceCreated = autoscalerCooldownServiceAdded
	}

	// cooldown ServiceUpdated
	autoscalerCooldownServiceUpdatedStr, ok := labels[AutoscalerCooldownServiceUpdated]
	if !ok {
		a.CooldownServiceUpdated = DefaultCooldownServiceUpdated
	} else {
		autoscalerCooldownServiceUpdated, err := time.ParseDuration(autoscalerCooldownServiceUpdatedStr)
		if err != nil {
			return ErrInvalidDuration
		}
		a.CooldownServiceUpdated = autoscalerCooldownServiceUpdated
	}

	// min replicas
	replicasMinStr, ok := labels[ReplicasMin]
	if !ok {
		a.MinReplicas = 1
	} else {
		replicasMin, err := strconv.ParseUint(replicasMinStr, 10, 64)
		if err != nil {
			return ErrInvalidUint
		}
		a.MinReplicas = replicasMin
	}

	// max replicas
	replicasMaxStr, ok := labels[ReplicasMax]
	if !ok {
		a.MaxReplicas = math.MaxUint64
	} else {
		replicasMax, err := strconv.ParseUint(replicasMaxStr, 10, 64)
		if err != nil {
			return ErrInvalidUint
		}
		a.MaxReplicas = replicasMax
	}

	return nil
}
