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

import "github.com/mtneug/pkg/ulid"

import "net/url"

// Type represents some category of metrics.
type Type string

const (
	// TypeCPU is a CPU metric.
	TypeCPU Type = "cpu"

	// TypeMemory is a memory metric.
	TypeMemory Type = "memory"

	// TypePrometheus is a Prometheus metric.
	TypePrometheus Type = "prometheus"
)

// Kind represents some kind of metric.
type Kind string

const (
	// KindReplica is a replica metric.
	KindReplica Kind = "replica"

	// KindSystem is a system metric.
	KindSystem Kind = "system"
)

// Metric represents a service metric.
type Metric struct {
	// ID of the metric.
	ID string
	// Name of the metric.
	Name string
	// Type of the metric.
	Type Type
	// Kind of the metric.
	Kind Kind
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

// New creates a new metric.
func New(name string) Metric {
	m := Metric{
		ID:   ulid.New().String(),
		Name: name,
	}
	return m
}
