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

package label_test

import (
	"testing"

	"github.com/mtneug/spate/label"
	"github.com/mtneug/spate/metric"
	"github.com/stretchr/testify/require"
)

func TestParseTarget(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		labels map[string]string
		err    error
		target metric.Target
	}{
		// error
		{
			labels: nil,
			err:    label.ErrNoValue,
		},
		{
			labels: map[string]string{"type": "prometheus"},
			err:    label.ErrNoValue,
		},
		{
			labels: map[string]string{"target": "NaN"},
			err:    label.ErrInvalidFloat,
		},
		{
			labels: map[string]string{
				"target":                 "42",
				"target.deviation.lower": "NaN",
			},
			err: label.ErrInvalidDeviation,
		},
		{
			labels: map[string]string{
				"target":                 "42",
				"target.deviation.lower": "-1",
			},
			err: label.ErrInvalidDeviation,
		},
		{
			labels: map[string]string{
				"target":                 "42",
				"target.deviation.upper": "NaN",
			},
			err: label.ErrInvalidDeviation,
		},
		{
			labels: map[string]string{
				"target":                 "42",
				"target.deviation.upper": "-1",
			},
			err: label.ErrInvalidDeviation,
		},

		// ok
		{
			labels: map[string]string{"type": "cpu"},
			err:    nil,
			target: metric.Target{
				Value:          0.8,
				LowerDeviation: 0,
				UpperDeviation: 0,
			},
		},
		{
			labels: map[string]string{
				"type": "cpu",
				"target.deviation.lower": "1%",
				"target.deviation.upper": "2%",
			},
			err: nil,
			target: metric.Target{
				Value:          0.8,
				LowerDeviation: 0.008,
				UpperDeviation: 0.016,
			},
		},
		{
			labels: map[string]string{"type": "cpu", "target": "0.4"},
			err:    nil,
			target: metric.Target{
				Value:          0.4,
				LowerDeviation: 0,
				UpperDeviation: 0,
			},
		},
		{
			labels: map[string]string{"type": "memory"},
			err:    nil,
			target: metric.Target{
				Value:          0.8,
				LowerDeviation: 0,
				UpperDeviation: 0,
			},
		},
		{
			labels: map[string]string{"type": "memory", "target": "0.4"},
			err:    nil,
			target: metric.Target{
				Value:          0.4,
				LowerDeviation: 0,
				UpperDeviation: 0,
			},
		},
		{
			labels: map[string]string{
				"type": "memory",
				"target.deviation.lower": "1%",
				"target.deviation.upper": "2%",
			},
			err: nil,
			target: metric.Target{
				Value:          0.8,
				LowerDeviation: 0.008,
				UpperDeviation: 0.016,
			},
		},
		{
			labels: map[string]string{
				"target":                 "42",
				"target.deviation.lower": "1",
				"target.deviation.upper": "2",
			},
			err: nil,
			target: metric.Target{
				Value:          42,
				LowerDeviation: 1,
				UpperDeviation: 2,
			},
		},
		{
			labels: map[string]string{
				"target":                 "-42.1234",
				"target.deviation.lower": "1.42",
				"target.deviation.upper": "2.42",
			},
			err: nil,
			target: metric.Target{
				Value:          -42.1234,
				LowerDeviation: 1.42,
				UpperDeviation: 2.42,
			},
		},
		{
			labels: map[string]string{
				"target":                 "42",
				"target.deviation.lower": "50%",
				"target.deviation.upper": "25%",
			},
			err: nil,
			target: metric.Target{
				Value:          42,
				LowerDeviation: 21,
				UpperDeviation: 10.5,
			},
		},
		{
			labels: map[string]string{
				"target":                 "-42.1234",
				"target.deviation.lower": "50%",
				"target.deviation.upper": "25%",
			},
			err: nil,
			target: metric.Target{
				Value:          -42.1234,
				LowerDeviation: 21.0617,
				UpperDeviation: 10.53085,
			},
		},
	}

	for _, c := range testCases {
		target := metric.Target{}
		err := label.ParseTarget(&target, c.labels)

		require.Equal(t, c.err, err)
		if c.err == nil {
			require.Equal(t, c.target, target)
		}
	}
}
