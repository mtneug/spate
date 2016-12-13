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

	"github.com/mtneug/spate/api/types"
	"github.com/mtneug/spate/consts"
	"github.com/mtneug/spate/metric"
	"github.com/stretchr/testify/require"
)

func TestNewByLabelsNoType(t *testing.T) {
	t.Parallel()

	_, err := metric.NewByLabels("test", nil)
	require.EqualError(t, err, metric.ErrNoType.Error())
}

func TestNewByLabelsUnknownType(t *testing.T) {
	t.Parallel()

	l := map[string]string{
		consts.LabelMetricTypeSuffix: "unknown",
	}

	_, err := metric.NewByLabels("test", l)
	require.EqualError(t, err, metric.ErrUnknownType.Error())
}

func TestNewByLabelsMetricTypeCPU(t *testing.T) {
	t.Parallel()

	l := map[string]string{
		consts.LabelMetricTypeSuffix: string(types.MetricTypeCPU),
	}

	m, err := metric.NewByLabels("test", l)
	require.NoError(t, err)
	require.Equal(t, "test", m.Name)
	require.Equal(t, types.MetricTypeCPU, m.Type)
}

func TestNewByLabelsMetricTypeMemory(t *testing.T) {
	t.Parallel()

	l := map[string]string{
		consts.LabelMetricTypeSuffix: string(types.MetricTypeMemory),
	}

	m, err := metric.NewByLabels("test", l)
	require.NoError(t, err)
	require.Equal(t, "test", m.Name)
	require.Equal(t, types.MetricTypeMemory, m.Type)
}

func TestNewByLabelsMetricTypePrometheus(t *testing.T) {
	t.Parallel()

	l := map[string]string{
		consts.LabelMetricTypeSuffix: string(types.MetricTypePrometheus),
	}

	m, err := metric.NewByLabels("test", l)
	require.NoError(t, err)
	require.Equal(t, "test", m.Name)
	require.Equal(t, types.MetricTypePrometheus, m.Type)
}
