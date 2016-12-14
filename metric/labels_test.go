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
	"net/url"
	"testing"

	"github.com/mtneug/spate/api/types"
	"github.com/mtneug/spate/metric"
	"github.com/stretchr/testify/require"
)

func TestNewByLabels(t *testing.T) {
	t.Parallel()

	testUrlStr := "http://localhost:8080/metrics?test=1"
	testUrl, _ := url.Parse(testUrlStr)

	testCases := []struct {
		name   string
		label  map[string]string
		err    error
		metric types.Metric
	}{
		{
			name:   "test",
			label:  nil,
			err:    metric.ErrNoType,
			metric: types.Metric{},
		},
		{
			name: "test",
			label: map[string]string{
				"type": "unknown",
			},
			err:    metric.ErrUnknownType,
			metric: types.Metric{},
		},

		// CPU
		{
			name: "test",
			label: map[string]string{
				"type": "cpu",
			},
			err: nil,
			metric: types.Metric{
				Name: "test",
				Type: types.MetricTypeCPU,
				Kind: types.MetricKindReplica,
			},
		},
		{
			name: "test",
			label: map[string]string{
				"type": "cpu",
				"kind": "system",
			},
			err: metric.ErrWrongKind,
			metric: types.Metric{
				Name: "test",
				Type: types.MetricTypeCPU,
				Kind: types.MetricKindReplica,
			},
		},

		// memory
		{
			name: "test",
			label: map[string]string{
				"type": "memory",
			},
			err: nil,
			metric: types.Metric{
				Name: "test",
				Type: types.MetricTypeMemory,
				Kind: types.MetricKindReplica,
			},
		},
		{
			name: "test",
			label: map[string]string{
				"type": "memory",
				"kind": "system",
			},
			err: metric.ErrWrongKind,
			metric: types.Metric{
				Name: "test",
				Type: types.MetricTypeMemory,
				Kind: types.MetricKindReplica,
			},
		},

		// Prometheus
		{
			name: "test",
			label: map[string]string{
				"type": "prometheus",
			},
			err:    metric.ErrNoKind,
			metric: types.Metric{},
		},
		{
			name: "test",
			label: map[string]string{
				"type": "prometheus",
				"kind": "unknown",
			},
			err:    metric.ErrUnknownKind,
			metric: types.Metric{},
		},
		{
			name: "test",
			label: map[string]string{
				"type": "prometheus",
				"kind": "system",
			},
			err:    metric.ErrNoPrometheusEndpoint,
			metric: types.Metric{},
		},
		{
			name: "test",
			label: map[string]string{
				"type":                "prometheus",
				"kind":                "system",
				"prometheus.endpoint": "ftp://why",
			},
			err:    metric.ErrPrometheusEndpointNotHTTPUrl,
			metric: types.Metric{},
		},
		{
			name: "test",
			label: map[string]string{
				"type":                "prometheus",
				"kind":                "system",
				"prometheus.endpoint": testUrlStr,
			},
			err:    metric.ErrNoPrometheusMetricName,
			metric: types.Metric{},
		},
		{
			name: "test",
			label: map[string]string{
				"type":                "prometheus",
				"kind":                "system",
				"prometheus.endpoint": testUrlStr,
				"prometheus.name":     "test_metric",
			},
			err: nil,
			metric: types.Metric{
				Name: "test",
				Type: types.MetricTypePrometheus,
				Kind: types.MetricKindSystem,
				Prometheus: types.PrometheusSpec{
					Endpoint: *testUrl,
					Name:     "test_metric",
				},
			},
		},
		{
			name: "test",
			label: map[string]string{
				"type":                "prometheus",
				"kind":                "replica",
				"prometheus.endpoint": testUrlStr,
				"prometheus.name":     "test_metric",
			},
			err: nil,
			metric: types.Metric{
				Name: "test",
				Type: types.MetricTypePrometheus,
				Kind: types.MetricKindReplica,
				Prometheus: types.PrometheusSpec{
					Endpoint: *testUrl,
					Name:     "test_metric",
				},
			},
		},
	}

	for _, c := range testCases {
		m, err := metric.NewByLabels(c.name, c.label)
		m.ID = ""
		require.Equal(t, c.err, err)
		require.Equal(t, c.metric, m)
	}
}
