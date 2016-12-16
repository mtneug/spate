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

package labels

const (
	//
	// Spate
	//

	// Namespace for spate
	Namespace = "de.mtneug.spate"

	//
	// Autoscaler
	//

	// AutoscalerPeriod is the suffix (after spate namespace) for the autoscaler
	// period.
	AutoscalerPeriod = "autoscaler.period"

	// AutoscalerCooldownScaledUp is the suffix (after spate namespace) for
	// autoscaler cooldown time after a ScaledUp.
	AutoscalerCooldownScaledUp = "autoscaler.cooldown.scaled_up"

	// AutoscalerCooldownScaledDown is the suffix (after spate namespace) for
	// autoscaler cooldown time after a ScaledDown.
	AutoscalerCooldownScaledDown = "autoscaler.cooldown.scaled_down"

	// AutoscalerCooldownServiceAdded is the suffix (after spate namespace) for
	// autoscaler cooldown time after a ServiceAdded.
	AutoscalerCooldownServiceAdded = "autoscaler.cooldown.service_added"

	// AutoscalerCooldownServiceUpdated is the suffix (after spate namespace) for
	// autoscaler cooldown time after a ServiceUpdated.
	AutoscalerCooldownServiceUpdated = "autoscaler.cooldown.service_updated"

	// ReplicasMin is the suffix (after spate namespace) for the minimum number of
	// replicas.
	ReplicasMin = "replicas.min"

	// ReplicasMax is the suffix (after spate namespace) for the maximum number of
	// replicas.
	ReplicasMax = "replicas.max"

	//
	// Metric
	//

	// AutoscalerReplicasMin is the suffix (after spate namespace) for metrics.
	AutoscalerReplicasMin = "autoscaler.replicas.min"

	// AutoscalerReplicasMax is the suffix (after spate namespace) for metrics.
	AutoscalerReplicasMax = "autoscaler.replicas.max"

	// MetricSuffix is the suffix (after spate namespace) for metrics.
	MetricSuffix = "metric"

	// MetricTypeSuffix is the suffix (after metric name) for the type.
	MetricTypeSuffix = "type"

	// MetricKindSuffix is the suffix (after metric name) for the kind.
	MetricKindSuffix = "kind"

	// MetricPrometheusEndpointSuffix is the suffix (after metric name) for the
	// Prometheus endpoint.
	MetricPrometheusEndpointSuffix = "prometheus.endpoint"

	// MetricPrometheusNameSuffix is the suffix (after metric name) for the
	// Prometheus metric name.
	MetricPrometheusNameSuffix = "prometheus.name"

	//
	// Aggregation
	//

	// MetricAggregationMethodSuffix is the suffix (after metric name) for the
	// aggregation method.
	MetricAggregationMethodSuffix = "aggregation.method"

	// MetricAggregationMethodMax in string form.
	MetricAggregationMethodMax = "max"

	// MetricAggregationMethodMin in string form.
	MetricAggregationMethodMin = "min"

	// MetricAggregationMethodAvg in string form.
	MetricAggregationMethodAvg = "avg"

	//
	// Observer
	//

	// MetricAggregationAmountSuffix is the suffix (after metric name) for the
	// aggregation amount.
	MetricAggregationAmountSuffix = "aggregation.amount"

	// MetricObserverPeriodSuffix is the suffix (after metric name) for the
	// observer period.
	MetricObserverPeriodSuffix = "observer.period"

	//
	// Target
	//

	// MetricTargetSuffix is the suffix (after metric name) for the target.
	MetricTargetSuffix = "target"

	// MetricTargetDeviationLowerSuffix is the suffix (after metric name) for the
	// lower deviation.
	MetricTargetDeviationLowerSuffix = "target.deviation.lower"

	// MetricTargetDeviationUpperSuffix is the suffix (after metric name) for the
	// lower deviation.
	MetricTargetDeviationUpperSuffix = "target.deviation.upper"
)
