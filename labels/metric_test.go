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

package labels_test

import (
	"net/url"
	"testing"

	"github.com/mtneug/spate/labels"
	"github.com/mtneug/spate/model"
	"github.com/stretchr/testify/require"
)

func TestParseMetric(t *testing.T) {
	t.Parallel()

	testUrlStr := "http://localhost:8080/metrics?test=1"
	testUrl, _ := url.Parse(testUrlStr)

	testCases := []struct {
		labels map[string]string
		err    error
		metric model.Metric
	}{
		// errors
		{
			labels: nil,
			err:    labels.ErrNoType,
		},
		{
			labels: map[string]string{
				"type": "unknown",
			},
			err: labels.ErrUnknownType,
		},
		{
			labels: map[string]string{
				"type": "prometheus",
			},
			err: labels.ErrNoKind,
		},
		{
			labels: map[string]string{
				"type": "prometheus",
				"kind": "unknown",
			},
			err: labels.ErrUnknownKind,
		},
		{
			labels: map[string]string{
				"type": "cpu",
				"kind": "system",
			},
			err: labels.ErrWrongKind,
		},
		{
			labels: map[string]string{
				"type": "memory",
				"kind": "system",
			},
			err: labels.ErrWrongKind,
		},
		{
			labels: map[string]string{
				"type": "prometheus",
				"kind": "system",
			},
			err: labels.ErrNoPrometheusEndpoint,
		},
		{
			labels: map[string]string{
				"type":                "prometheus",
				"kind":                "system",
				"prometheus.endpoint": "ftp://why",
			},
			err: labels.ErrInvalidHTTPUrl,
		},
		{
			labels: map[string]string{
				"type":                "prometheus",
				"kind":                "system",
				"prometheus.endpoint": testUrlStr,
			},
			err: labels.ErrNoPrometheusMetricName,
		},

		// cpu
		{
			labels: map[string]string{
				"type": "cpu",
			},
			err: nil,
			metric: model.Metric{
				Type: model.MetricTypeCPU,
				Kind: model.MetricKindReplica,
			},
		},
		{
			labels: map[string]string{
				"type": "cpu",
				"kind": "replica",
			},
			err: nil,
			metric: model.Metric{
				Type: model.MetricTypeCPU,
				Kind: model.MetricKindReplica,
			},
		},

		// memory
		{
			labels: map[string]string{
				"type": "memory",
			},
			err: nil,
			metric: model.Metric{
				Type: model.MetricTypeMemory,
				Kind: model.MetricKindReplica,
			},
		},
		{
			labels: map[string]string{
				"type": "memory",
				"kind": "replica",
			},
			err: nil,
			metric: model.Metric{
				Type: model.MetricTypeMemory,
				Kind: model.MetricKindReplica,
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
			metric: model.Metric{
				Type: model.MetricTypePrometheus,
				Kind: model.MetricKindReplica,
				Prometheus: model.PrometheusSpec{
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
			metric: model.Metric{
				Type: model.MetricTypePrometheus,
				Kind: model.MetricKindSystem,
				Prometheus: model.PrometheusSpec{
					Endpoint: *testUrl,
					Name:     "test_metric",
				},
			},
		},
	}

	for _, c := range testCases {
		metric := model.Metric{}
		err := labels.ParseMetric(&metric, c.labels)
		require.Equal(t, c.err, err)
		if c.err == nil {
			require.Equal(t, c.metric, metric)
		}
	}
}
