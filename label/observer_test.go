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
	"time"

	"github.com/mtneug/spate/label"
	"github.com/mtneug/spate/metric"
	"github.com/stretchr/testify/require"
)

func TestParseObserver(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		label    map[string]string
		err      error
		observer metric.Observer
	}{

		{
			label: map[string]string{"observer.period": "abc"},
			err:   label.ErrInvalidDuration,
		},
		{
			label: map[string]string{"aggregation.amount": "abc"},
			err:   label.ErrInvalidUint,
		},
		{
			label: map[string]string{},
			err:   nil,
			observer: metric.Observer{
				Period:            30 * time.Second,
				AggregationAmount: 5,
			},
		},
		{
			label: map[string]string{
				"observer.period":    "42s",
				"aggregation.amount": "42",
			},
			err: nil,
			observer: metric.Observer{
				Period:            42 * time.Second,
				AggregationAmount: 42,
			},
		},
	}

	for _, c := range testCases {
		observer := metric.Observer{}
		err := label.ParseObserver(&observer, c.label)
		require.Equal(t, c.err, err)
		if c.err == nil {
			require.Equal(t, c.observer, observer)
		}
	}
}
