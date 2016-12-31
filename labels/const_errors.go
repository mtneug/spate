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

import "errors"

var (
	// ErrInvalidMetricLabel indicates that the parsing failed because the metric
	// label is invalid.
	ErrInvalidMetricLabel = errors.New("labels: metric label invalid")

	// ErrDuplicateMetric indicates that at least two metrics are the same.
	ErrDuplicateMetric = errors.New("labels: duplicate metric")

	// ErrNoType indicates that the parsing failed because no type was specified.
	ErrNoType = errors.New("labels: no type specified")

	// ErrUnknownType indicates that the parsing failed because the type is
	// unknown.
	ErrUnknownType = errors.New("labels: unknown type")

	// ErrNoKind indicates that the parsing failed because no kind was specified.
	ErrNoKind = errors.New("labels: no kind specified")

	// ErrUnknownKind indicates that the parsing failed because the kind is
	// unknown.
	ErrUnknownKind = errors.New("labels: unknown kind")

	// ErrWrongKind indicates that the parsing failed because a wrong kind was
	// used for the metric.
	ErrWrongKind = errors.New("labels: wrong kind")

	// ErrNoPrometheusEndpoint indicates that the parsing failed because no
	// Prometheus endpoint was specified.
	ErrNoPrometheusEndpoint = errors.New("labels: no Prometheus endpoint specified")

	// ErrInvalidHTTPUrl indicates that the parsing failed because the specified
	// URL is an invalid HTTP URL.
	ErrInvalidHTTPUrl = errors.New("labels: invalid HTTP URL")

	// ErrNoPrometheusMetricName indicates that the parsing failed because no
	// Prometheus metric name was specified.
	ErrNoPrometheusMetricName = errors.New("labels: no Prometheus metric name specified")

	// ErrUnknownAggregationMethod indicates that the parsing failed because the
	// aggregation method is unknown.
	ErrUnknownAggregationMethod = errors.New("labels: aggregation method is unknown")

	// ErrNoValue indicates that the parsing failed because no value was
	// specified.
	ErrNoValue = errors.New("labels: no value specified")

	// ErrInvalidFloat indicates that the parsing failed because the float is
	// invalid.
	ErrInvalidFloat = errors.New("labels: invalid float")

	// ErrInvalidDeviation indicates that the parsing failed because the deviation
	// is invalid.
	ErrInvalidDeviation = errors.New("labels: invalid deviation")

	// ErrInvalidDuration indicates that the parsing failed because the duration
	// is invalid.
	ErrInvalidDuration = errors.New("labels: invalid duration")

	// ErrInvalidUint indicates that the parsing failed because the uint is
	// invalid.
	ErrInvalidUint = errors.New("labels: invalid uint")
)
