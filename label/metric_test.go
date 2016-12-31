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

package label_test

import (
	"net/url"
	"testing"

	"github.com/mtneug/spate/label"
	"github.com/mtneug/spate/metric"
	"github.com/stretchr/testify/require"
)

func TestParseMetric(t *testing.T) {
	t.Parallel()

	testUrlStr := "http://localhost:8080/metrics?test=1"
	testUrl, _ := url.Parse(testUrlStr)

	testCases := []struct {
		labels map[string]string
		err    error
		metric metric.Metric
	}{
		// errors
		{
			labels: nil,
			err:    label.ErrNoType,
		},
		{
			labels: map[string]string{
				"type": "unknown",
			},
			err: label.ErrUnknownType,
		},
		{
			labels: map[string]string{
				"type": "prometheus",
			},
			err: label.ErrNoKind,
		},
		{
			labels: map[string]string{
				"type": "prometheus",
				"kind": "unknown",
			},
			err: label.ErrUnknownKind,
		},
		{
			labels: map[string]string{
				"type": "cpu",
				"kind": "system",
			},
			err: label.ErrWrongKind,
		},
		{
			labels: map[string]string{
				"type": "memory",
				"kind": "system",
			},
			err: label.ErrWrongKind,
		},
		{
			labels: map[string]string{
				"type": "prometheus",
				"kind": "system",
			},
			err: label.ErrNoPrometheusEndpoint,
		},
		{
			labels: map[string]string{
				"type":                "prometheus",
				"kind":                "system",
				"prometheus.endpoint": "ftp://why",
			},
			err: label.ErrInvalidHTTPUrl,
		},
		{
			labels: map[string]string{
				"type":                "prometheus",
				"kind":                "system",
				"prometheus.endpoint": testUrlStr,
			},
			err: label.ErrNoPrometheusMetricName,
		},

		// cpu
		{
			labels: map[string]string{
				"type": "cpu",
			},
			err: nil,
			metric: metric.Metric{
				Type: metric.TypeCPU,
				Kind: metric.KindReplica,
			},
		},
		{
			labels: map[string]string{
				"type": "cpu",
				"kind": "replica",
			},
			err: nil,
			metric: metric.Metric{
				Type: metric.TypeCPU,
				Kind: metric.KindReplica,
			},
		},

		// memory
		{
			labels: map[string]string{
				"type": "memory",
			},
			err: nil,
			metric: metric.Metric{
				Type: metric.TypeMemory,
				Kind: metric.KindReplica,
			},
		},
		{
			labels: map[string]string{
				"type": "memory",
				"kind": "replica",
			},
			err: nil,
			metric: metric.Metric{
				Type: metric.TypeMemory,
				Kind: metric.KindReplica,
			},
		},

		// prometheus
		{
			labels: map[string]string{
				"type":                "prometheus",
				"kind":                "replica",
				"prometheus.endpoint": testUrlStr,
				"prometheus.name":     "test_metric",
			},
			err: nil,
			metric: metric.Metric{
				Type: metric.TypePrometheus,
				Kind: metric.KindReplica,
				Prometheus: metric.PrometheusSpec{
					Endpoint: *testUrl,
					Name:     "test_metric",
				},
			},
		},
		{
			labels: map[string]string{
				"type":                "prometheus",
				"kind":                "system",
				"prometheus.endpoint": testUrlStr,
				"prometheus.name":     "test_metric",
			},
			err: nil,
			metric: metric.Metric{
				Type: metric.TypePrometheus,
				Kind: metric.KindSystem,
				Prometheus: metric.PrometheusSpec{
					Endpoint: *testUrl,
					Name:     "test_metric",
				},
			},
		},
	}

	for _, c := range testCases {
		m := metric.Metric{}
		err := label.ParseMetric(&m, c.labels)
		require.Equal(t, c.err, err)
		if c.err == nil {
			require.Equal(t, c.metric, m)
		}
	}
}
