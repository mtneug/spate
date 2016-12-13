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
	// ErrNoType indicates that the metric could not be created because no type
	// was specified.
	ErrNoType = errors.New("metric: no type specified")

	// ErrUnknownType indicates that the metric could not be created because the
	// type is unknown.
	ErrUnknownType = errors.New("metric: type unknown")

	// ErrNoKind indicates that the metric could not be created because no kind
	// was specified.
	ErrNoKind = errors.New("metric: wrong kind")

	// ErrUnknownKind indicates that the metric could not be created because the
	// kind is unknown.
	ErrUnknownKind = errors.New("metric: kind unknown")

	// ErrWrongKind indicates that a wrong kind was used for the metric. For some
	// metric types this can be automatically corrected so that the metric still
	// can be created.
	ErrWrongKind = errors.New("metric: wrong kind")

	emptyMetric = types.Metric{}

	stdCPUMetric = types.Metric{
		Type: types.MetricTypeCPU,
		Kind: types.MetricKindReplica,
	}

	stdMemoryMetric = types.Metric{
		Type: types.MetricTypeMemory,
		Kind: types.MetricKindReplica,
	}

	stdPrometheusMetric = types.Metric{
		Type: types.MetricTypePrometheus,
	}
)

// NewByLabels creates a new Metric based on labels.
func NewByLabels(name string, label map[string]string) (types.Metric, error) {
	tStr, ok := label[consts.LabelMetricTypeSuffix]
	if !ok {
		return emptyMetric, ErrNoType
	}

	switch types.MetricType(tStr) {
	case types.MetricTypeCPU:
		return newCPUMetricByLabels(name, label)
	case types.MetricTypeMemory:
		return newMemoryMetricByLabels(name, label)
	case types.MetricTypePrometheus:
		return newPrometheusMetricByLabels(name, label)
	}

	return emptyMetric, ErrUnknownType
}

func newCPUMetricByLabels(name string, label map[string]string) (m types.Metric, err error) {
	m = stdCPUMetric
	m.ID = ulid.New().String()
	m.Name = name

	k, ok := label[consts.LabelMetricKindSuffix]
	if ok && k != string(types.MetricKindReplica) {
		err = ErrWrongKind
	}

	return
}

func newMemoryMetricByLabels(name string, label map[string]string) (m types.Metric, err error) {
	m = stdMemoryMetric
	m.ID = ulid.New().String()
	m.Name = name

	k, ok := label[consts.LabelMetricKindSuffix]
	if ok && k != string(types.MetricKindReplica) {
		err = ErrWrongKind
	}

	return
}

func newPrometheusMetricByLabels(name string, label map[string]string) (m types.Metric, err error) {
	m = stdPrometheusMetric
	m.ID = ulid.New().String()
	m.Name = name

	k, ok := label[consts.LabelMetricKindSuffix]
	if !ok {
		return emptyMetric, ErrNoKind
	}
	switch types.MetricKind(k) {
	case types.MetricKindSystem:
		m.Kind = types.MetricKindSystem
	case types.MetricKindReplica:
		m.Kind = types.MetricKindReplica
	default:
		return emptyMetric, ErrUnknownKind
	}

	return
}
