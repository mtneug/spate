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
	"testing"

	"github.com/mtneug/spate/labels"
	"github.com/stretchr/testify/require"
)

func TestParseReducer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		labels map[string]string
		err    error
	}{
		{
			labels: map[string]string{"aggregation.method": "unknown"},
			err:    labels.ErrUnknownAggregationMethod,
		},
		{
			labels: map[string]string{},
			err:    nil,
		},
		{
			labels: map[string]string{"aggregation.method": "max"},
			err:    nil,
		},
		{
			labels: map[string]string{"aggregation.method": "min"},
			err:    nil,
		},
		{
			labels: map[string]string{"aggregation.method": "avg"},
			err:    nil,
		},
	}

	for _, c := range testCases {
		_, err := labels.ParseReducer(c.labels)
		require.Equal(t, c.err, err)
	}
}
