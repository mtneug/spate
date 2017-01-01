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
	"net/url"

	"github.com/mtneug/spate/metric"
)

// ParseMetric parses the labels and sets the corresponding values for given
// metric.
func ParseMetric(m *metric.Metric, labels map[string]string) error {
	typeStr, ok := labels[MetricTypeSuffix]
	if !ok {
		return ErrNoType
	}

	switch metric.Type(typeStr) {
	case metric.TypeCPU:
		return parseCPUMetric(m, labels)
	case metric.TypeMemory:
		return parseMemoryMetric(m, labels)
	case metric.TypePrometheus:
		return parsePrometheusMetric(m, labels)
	}

	return ErrUnknownType
}

func parseCPUMetric(m *metric.Metric, labels map[string]string) error {
	m.Type = metric.TypeCPU

	kindStr, ok := labels[MetricKindSuffix]
	if ok && kindStr != string(metric.KindReplica) {
		return ErrWrongKind
	}
	m.Kind = metric.KindReplica

	return nil
}

func parseMemoryMetric(m *metric.Metric, labels map[string]string) error {
	m.Type = metric.TypeMemory

	kindStr, ok := labels[MetricKindSuffix]
	if ok && kindStr != string(metric.KindReplica) {
		return ErrWrongKind
	}
	m.Kind = metric.KindReplica

	return nil
}

func parsePrometheusMetric(m *metric.Metric, labels map[string]string) error {
	m.Type = metric.TypePrometheus

	// Kind
	kindStr, ok := labels[MetricKindSuffix]
	if !ok {
		return ErrNoKind
	}
	kind := metric.Kind(kindStr)
	ok = validMetricKind(kind)
	if !ok {
		return ErrUnknownKind
	}
	m.Kind = kind

	// Prometheus endpoint
	endpointStr, ok := labels[MetricPrometheusEndpointSuffix]
	if !ok {
		return ErrNoPrometheusEndpoint
	}
	endpoint, err := url.Parse(endpointStr)
	if err != nil || endpoint.Scheme != "http" {
		return ErrInvalidHTTPUrl
	}
	m.Prometheus.Endpoint = *endpoint

	// Prometheus metric name
	prometheusName, ok := labels[MetricPrometheusNameSuffix]
	if !ok {
		return ErrNoPrometheusMetricName
	}
	m.Prometheus.Name = prometheusName

	return nil
}

func validMetricKind(kind metric.Kind) (ok bool) {
	return kind == metric.KindSystem ||
		kind == metric.KindReplica
}
