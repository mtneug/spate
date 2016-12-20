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
	"context"
	"testing"

	"github.com/mtneug/spate/metric"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockMeasure struct {
	mock.Mock
}

func (m *MockMeasure) Measure(ctx context.Context) (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

type MockReducer struct {
	mock.Mock
}

func (m *MockReducer) Reduce(data []float64) (float64, error) {
	args := m.Called(data)
	return args.Get(0).(float64), args.Error(1)
}

func TestNewObserver(t *testing.T) {
	t.Parallel()

	m := &MockMeasure{}
	r := &MockReducer{}

	o := metric.NewObserver(m, r)
	require.Equal(t, m, o.Measurer)
	require.Equal(t, r, o.Reducer)

	m.AssertExpectations(t)
	r.AssertExpectations(t)
}
