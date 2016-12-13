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
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/mtneug/spate/metric"
	"github.com/stretchr/testify/require"
)

func TestNewSrvNoMetrics(t *testing.T) {
	t.Parallel()

	srv := swarm.Service{}
	a, err := New(srv)
	require.Nil(t, a)
	require.EqualError(t, err, ErrNoMetrics.Error())
}

func TestNewSrvInvalidMetricLabel(t *testing.T) {
	t.Parallel()

	srv := swarm.Service{
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"de.mtneug.spate.metric.": "test",
				},
			},
		},
	}
	a, err := New(srv)
	require.Nil(t, a)
	require.EqualError(t, err, ErrInvalidMetricLabel.Error())
}

func TestNewSrvDuplMetric(t *testing.T) {
	t.Parallel()

	srv := swarm.Service{
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"de.mtneug.spate.metric.test1.type": "cpu",
					"de.mtneug.spate.metric.test2.type": "cpu",
				},
			},
		},
	}
	a, err := New(srv)
	require.Nil(t, a)
	require.EqualError(t, err, ErrDuplMetric.Error())
}

func TestNewMetricNoTyoe(t *testing.T) {
	t.Parallel()

	srv := swarm.Service{
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"de.mtneug.spate.metric.test.test": "test",
				},
			},
		},
	}
	a, err := New(srv)
	require.Nil(t, a)
	require.EqualError(t, err, metric.ErrNoType.Error())
}

func TestNewMetricUnknownType(t *testing.T) {
	t.Parallel()

	srv := swarm.Service{
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"de.mtneug.spate.metric.test.type": "unknown",
				},
			},
		},
	}
	a, err := New(srv)
	require.Nil(t, a)
	require.EqualError(t, err, metric.ErrUnknownType.Error())
}

func TestNew(t *testing.T) {
	t.Parallel()

	srv := swarm.Service{
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"de.mtneug.spate.metric.cpu.type":    "cpu",
					"de.mtneug.spate.metric.memory.type": "memory",
					"de.mtneug.spate.metric.custom.type": "prometheus",
				},
			},
		},
	}
	_, err := New(srv)
	// require.Nil(t, a)
	require.NoError(t, err)
}
