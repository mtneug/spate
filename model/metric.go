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

package model

import "net/url"

// MetricType represents some category of metrics.
type MetricType string

const (
	// MetricTypeCPU is a CPU metric.
	MetricTypeCPU MetricType = "cpu"

	// MetricTypeMemory is a memory metric.
	MetricTypeMemory MetricType = "memory"

	// MetricTypePrometheus is a Prometheus metric.
	MetricTypePrometheus MetricType = "prometheus"
)

// MetricKind represents some kind of metric.
type MetricKind string

const (
	// MetricKindReplica is a replica metric.
	MetricKindReplica MetricKind = "replica"

	// MetricKindSystem is a system metric.
	MetricKindSystem MetricKind = "system"
)

// Metric represents a service metric.
type Metric struct {
	// ID of the metric.
	ID string
	// Name of the metric.
	Name string
	// Type of the metric.
	Type MetricType
	// Kind of the metric.
	Kind MetricKind
	// Prometheus spec.
	Prometheus PrometheusSpec
}

// PrometheusSpec specifies a Prometheus metric.
type PrometheusSpec struct {
	// Endpoint of the Prometheus metrics.
	Endpoint url.URL
	// Name of the Prometheus metrics.
	Name string
}
