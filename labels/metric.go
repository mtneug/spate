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
	"net/url"

	"github.com/mtneug/spate/model"
)

// ParseMetric parses the labels and sets the corresponding values for given
// metric.
func ParseMetric(metric *model.Metric, labels map[string]string) error {
	typeStr, ok := labels[MetricTypeSuffix]
	if !ok {
		return ErrNoType
	}

	switch model.MetricType(typeStr) {
	case model.MetricTypeCPU:
		return parseCPUMetric(metric, labels)
	case model.MetricTypeMemory:
		return parseMemoryMetric(metric, labels)
	case model.MetricTypePrometheus:
		return parsePrometheusMetric(metric, labels)
	}

	return ErrUnknownType
}

func parseCPUMetric(metric *model.Metric, labels map[string]string) error {
	metric.Type = model.MetricTypeCPU

	kindStr, ok := labels[MetricKindSuffix]
	if ok && kindStr != string(model.MetricKindReplica) {
		return ErrWrongKind
	}
	metric.Kind = model.MetricKindReplica

	return nil
}

func parseMemoryMetric(metric *model.Metric, labels map[string]string) error {
	metric.Type = model.MetricTypeMemory

	kindStr, ok := labels[MetricKindSuffix]
	if ok && kindStr != string(model.MetricKindReplica) {
		return ErrWrongKind
	}
	metric.Kind = model.MetricKindReplica

	return nil
}

func parsePrometheusMetric(metric *model.Metric, labels map[string]string) error {
	metric.Type = model.MetricTypePrometheus

	// Kind
	kindStr, ok := labels[MetricKindSuffix]
	if !ok {
		return ErrNoKind
	}

	kind := model.MetricKind(kindStr)
	ok = validMetricKind(kind)
	if !ok {
		return ErrUnknownKind
	}

	metric.Kind = kind

	// Prometheus endpoint
	endpointStr, ok := labels[MetricPrometheusEndpointSuffix]
	if !ok {
		return ErrNoPrometheusEndpoint
	}
	endpoint, err := url.Parse(endpointStr)
	if err != nil || endpoint.Scheme != "http" {
		return ErrInvalidHTTPUrl
	}
	metric.Prometheus.Endpoint = *endpoint

	// Prometheus metric name
	prometheusName, ok := labels[MetricPrometheusNameSuffix]
	if !ok {
		return ErrNoPrometheusMetricName
	}
	metric.Prometheus.Name = prometheusName

	return nil
}

func validMetricKind(kind model.MetricKind) (ok bool) {
	return kind == model.MetricKindSystem ||
		kind == model.MetricKindReplica
}
