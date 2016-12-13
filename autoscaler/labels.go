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

package autoscaler

import (
	"errors"
	"strings"

	"github.com/mtneug/spate/api/types"
	"github.com/mtneug/spate/consts"
	"github.com/mtneug/spate/metric"
)

var (
	// ErrInvalidMetricLabel indicates that the metric label is invalid.
	ErrInvalidMetricLabel = errors.New("autoscaler: invalid metric label")

	// ErrDuplMetric indicates that at least two metrics are duplicates.
	ErrDuplMetric = errors.New("autoscaler: duplicate metrics defined")

	// ErrNoMetrics indicates that no metrics were defined.
	ErrNoMetrics = errors.New("autoscaler: no metrics defined")
)

func (a *Autoscaler) processLabels() error {
	metricsLabels := make(map[string]map[string]string)

	for label, value := range a.Service.Spec.Labels {
		if strings.HasPrefix(label, consts.LabelSpate+".") {
			label = label[len(consts.LabelSpate)+1:]

			labelSuffixIndex := strings.Index(label, ".")
			labelSuffix := label[:labelSuffixIndex]

			switch labelSuffix {
			case consts.LabelMetricSuffix:
				// metric labels must have at least three parts: "metric", name, suffix
				parts := strings.SplitN(label, ".", 3)
				if len(parts) < 3 {
					return ErrInvalidMetricLabel
				}

				metricName := parts[1]
				metricLabelSuffix := parts[2]
				metricsLabels[metricName] = initMapAndAdd(
					metricsLabels[metricName],
					metricLabelSuffix,
					value,
				)
			}
		}
	}

	metrics, err := createMetrics(metricsLabels)
	if err != nil {
		return err
	}
	if len(metrics) == 0 {
		return ErrNoMetrics
	}

	return nil
}

func initMapAndAdd(m map[string]string, key, value string) map[string]string {
	if m == nil {
		m = make(map[string]string)
	}
	m[key] = value
	return m
}

func createMetrics(metricsLabels map[string]map[string]string) (map[string]types.Metric, error) {
	metrics := make(map[string]types.Metric, len(metricsLabels))
	seen := make(map[types.Metric]bool)

	for metricName, labels := range metricsLabels {
		m, err := metric.NewByLabels(metricName, labels)
		if err != nil {
			return nil, err
		}

		// To look for duplicate metrics we have to use equal values for the metric
		// ID and Name
		normM := m
		normM.ID = ""
		normM.Name = ""
		if seen[normM] {
			return nil, ErrDuplMetric
		}
		seen[normM] = true

		metrics[metricName] = m
	}

	return metrics, nil
}
