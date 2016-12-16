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

package labels

import (
	"math"
	"strconv"

	"github.com/mtneug/spate/api/types"
)

var (
	// DefaultTargetValueCPU if no value given.
	DefaultTargetValueCPU = 0.7

	// DefaultTargetValueMemory if no value given.
	DefaultTargetValueMemory = 0.7

	// TODO: Default lower/upper deviation for CPU/memory
)

// ParseTarget parses the labels and sets the corresponding values for given
// target.
func ParseTarget(target *types.Target, labels map[string]string) error {
	// value
	var (
		ok          bool
		err         error
		value       float64
		valueParsed = false
	)

	valueStr, ok := labels[MetricTargetSuffix]
	if !ok {
		// default value apply?
		var metricTypeStr string
		metricTypeStr, ok = labels[MetricTypeSuffix]
		if !ok {
			return ErrNoValue
		}

		switch types.MetricType(metricTypeStr) {
		case types.MetricTypeCPU:
			value = DefaultTargetValueCPU
			valueParsed = true
		case types.MetricTypeMemory:
			value = DefaultTargetValueMemory
			valueParsed = true
		default:
			return ErrNoValue
		}
	}
	if ok && !valueParsed {
		value, err = strconv.ParseFloat(valueStr, 64)
		if err != nil || math.IsNaN(value) {
			return ErrInvalidFloat
		}
	}
	target.Value = value

	// lower deviation
	deviationLowerStr, ok := labels[MetricTargetDeviationLowerSuffix]
	if ok {
		deviationLower, err := strconv.ParseFloat(deviationLowerStr, 64)
		if err != nil || math.IsNaN(deviationLower) || deviationLower < 0 {
			return ErrInvalidDeviation
		}
		target.LowerDeviation = deviationLower
	}

	// upper deviation
	deviationUpperStr, ok := labels[MetricTargetDeviationUpperSuffix]
	if ok {
		deviationUpper, err := strconv.ParseFloat(deviationUpperStr, 64)
		if err != nil || math.IsNaN(deviationUpper) || deviationUpper < 0 {
			return ErrInvalidDeviation
		}
		target.UpperDeviation = deviationUpper
	}

	return nil
}
