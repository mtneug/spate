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

package label

import "strings"

// ExtractSpateLabels from service labels.
func ExtractSpateLabels(labels map[string]string, metricLabels map[string]map[string]string,
	srvLabels map[string]string) error {
	for label, value := range srvLabels {
		if strings.HasPrefix(label, Namespace+".") {
			label = label[len(Namespace)+1:]
			labels[label] = value

			if strings.HasPrefix(label, MetricSuffix+".") {
				// metric labels must have at least three parts:
				//   * "metric"
				//   *  name
				//   *  metric label
				parts := strings.SplitN(label, ".", 3)
				if len(parts) < 3 {
					return ErrInvalidMetricLabel
				}

				metricName, metricLabelSuffix := parts[1], parts[2]
				metricLabels[metricName] = initMapAndAdd(
					metricLabels[metricName],
					metricLabelSuffix,
					value,
				)
			}
		}
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
