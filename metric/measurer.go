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

import "github.com/mtneug/spate/api/types"

// Measurer measures a metric for a given service.
type Measurer interface {
	Measure() (float64, error)
}

// CPUMeasurer measures the CPU utilization.
type CPUMeasurer struct {
	ServiceID string
	Metric    types.Metric
}

// Measure the CPU utilization.
func (m *CPUMeasurer) Measure() (float64, error) {
	panic("not implemented")
}

// MemoryMeasurer measures the memory utilization.
type MemoryMeasurer struct {
	ServiceID string
	Metric    types.Metric
}

// Measure the memory utilization.
func (m *MemoryMeasurer) Measure() (float64, error) {
	panic("not implemented")
}

// PrometheusMeasurer measures the Prometheus metric.
type PrometheusMeasurer struct {
	ServiceID string
	Metric    types.Metric
}

// Measure the Prometheus metric.
func (m *PrometheusMeasurer) Measure() (float64, error) {
	panic("not implemented")
}
