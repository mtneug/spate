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
	"errors"

	"github.com/mtneug/pkg/ulid"
	"github.com/mtneug/spate/api/types"
	"github.com/mtneug/spate/consts"
)

var (
	// ErrNoType indicates that the Metric could not be created because no type
	// was specified.
	ErrNoType = errors.New("metric: no type specified")

	// ErrUnknownType indicates that the Metric could not be created because the
	// type is unknown.
	ErrUnknownType = errors.New("metric: type unknown")
)

// NewByLabels creates a new Metric based on labels.
func NewByLabels(name string, label map[string]string) (types.Metric, error) {
	tStr, ok := label[consts.LabelMetricTypeSuffix]
	if !ok {
		return types.Metric{}, ErrNoType
	}

	switch types.MetricType(tStr) {
	case types.MetricTypeCPU:
		return newCPUMetricByLabels(name, label)
	case types.MetricTypeMemory:
		return newMemoryMetricByLabels(name, label)
	case types.MetricTypePrometheus:
		return newPrometheusMetricByLabels(name, label)
	}

	return types.Metric{}, ErrUnknownType
}

func newCPUMetricByLabels(name string, label map[string]string) (types.Metric, error) {
	m := types.Metric{
		ID:   ulid.New().String(),
		Name: name,
		Type: types.MetricTypeCPU,
	}

	return m, nil
}

func newMemoryMetricByLabels(name string, label map[string]string) (types.Metric, error) {
	m := types.Metric{
		ID:   ulid.New().String(),
		Name: name,
		Type: types.MetricTypeMemory,
	}

	return m, nil
}

func newPrometheusMetricByLabels(name string, label map[string]string) (types.Metric, error) {
	m := types.Metric{
		ID:   ulid.New().String(),
		Name: name,
		Type: types.MetricTypePrometheus,
	}

	return m, nil
}
