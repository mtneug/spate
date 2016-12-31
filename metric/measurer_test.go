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

package metric_test

import (
	"testing"

	"github.com/mtneug/spate/metric"
	"github.com/mtneug/spate/model"
	"github.com/stretchr/testify/require"
)

func TestNewMeasurer(t *testing.T) {
	t.Parallel()

	serviceID := "bv24k4vkkdnch70zpt81ia4kz"
	serviceName := "test_service_name"

	m := model.Metric{}
	measurer, err := metric.NewMeasurer(serviceID, serviceName, m)
	require.EqualError(t, metric.ErrUnknownType, err.Error())
	require.Nil(t, measurer)

	m = model.Metric{Type: model.MetricTypeCPU}
	measurer, err = metric.NewMeasurer(serviceID, serviceName, m)
	require.NoError(t, err)
	require.IsType(t, &metric.CPUMeasurer{}, measurer)

	m = model.Metric{Type: model.MetricTypeMemory}
	measurer, err = metric.NewMeasurer(serviceID, serviceName, m)
	require.NoError(t, err)
	require.IsType(t, &metric.MemoryMeasurer{}, measurer)

	m = model.Metric{Type: model.MetricTypePrometheus}
	measurer, err = metric.NewMeasurer(serviceID, serviceName, m)
	require.NoError(t, err)
	require.IsType(t, &metric.PrometheusMeasurer{}, measurer)
}
