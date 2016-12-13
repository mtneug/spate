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

package types

// MetricType represents some category of metrics.
type MetricType string

const (
	// MetricTypeCPU is a CPUMetric
	MetricTypeCPU MetricType = "cpu"

	// MetricTypeMemory is a MemoryMetric
	MetricTypeMemory MetricType = "memory"

	// MetricTypePrometheus is a PrometheusMetric
	MetricTypePrometheus MetricType = "prometheus"
)

// Metric represents a service metric.
type Metric struct {
	// ID of the metric.
	ID string
	// Name of the metric.
	Name string
	// Type of the metric.
	Type MetricType
}
