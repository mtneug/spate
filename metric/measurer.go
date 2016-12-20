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

package metric

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/mtneug/pkg/reducer"
	"github.com/mtneug/spate/api/types"
	"github.com/mtneug/spate/docker"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
)

// CriticalFailurePercentage indicates for a replica metric how many
// measurements are allowed to fail.
const CriticalFailurePercentage = 0.5

var (
	// ErrUnknownType indicates that the type is unknown.
	ErrUnknownType = errors.New("metric: unknown type")

	// ErrContainerNotFound indicates that the container/s was/were not found.
	ErrContainerNotFound = errors.New("metric: Container not found")

	// ErrMetricNotFound indicates that the metric was not found.
	ErrMetricNotFound = errors.New("metric: Metric not found")

	// ErrTooManyFailedMeasurements indicates that too many measurements failed.
	ErrTooManyFailedMeasurements = errors.New("metric: Too many failed measurements")
)

// Measurer measures a metric for a given service.
type Measurer interface {
	Measure(ctx context.Context) (float64, error)
}

// NewMeasurer creates the right measurer for given metric.
func NewMeasurer(serviceID string, metric types.Metric) (measurer Measurer, err error) {
	switch metric.Type {
	case types.MetricTypeCPU:
		measurer = &CPUMeasurer{ServiceID: serviceID, Metric: metric}
	case types.MetricTypeMemory:
		measurer = &MemoryMeasurer{ServiceID: serviceID, Metric: metric}
	case types.MetricTypePrometheus:
		measurer = &PrometheusMeasurer{ServiceID: serviceID, Metric: metric}
	default:
		err = ErrUnknownType
	}
	return
}

// CPUMeasurer measures the CPU utilization.
type CPUMeasurer struct {
	ServiceID string
	Metric    types.Metric
}

// Measure the CPU utilization.
func (m *CPUMeasurer) Measure(ctx context.Context) (float64, error) {
	// TODO: implement
	return 0, errors.New("not implemented")
}

// MemoryMeasurer measures the memory utilization.
type MemoryMeasurer struct {
	ServiceID string
	Metric    types.Metric
}

// Measure the memory utilization.
func (m *MemoryMeasurer) Measure(ctx context.Context) (float64, error) {
	// TODO: implement
	return 0, errors.New("not implemented")
}

// PrometheusMeasurer measures the Prometheus metric.
type PrometheusMeasurer struct {
	ServiceID string
	Metric    types.Metric
	client    http.Client
}

// Measure the Prometheus metric.
func (m *PrometheusMeasurer) Measure(ctx context.Context) (float64, error) {
	args := filters.NewArgs()
	args.Add("label", "com.docker.swarm.service.id="+m.ServiceID)

	opts := dockerTypes.ContainerListOptions{Filter: args}
	if m.Metric.Kind == types.MetricKindSystem {
		opts.Limit = 1
	}

	cs, err := docker.C.ContainerList(ctx, opts)
	if err != nil {
		return 0, err
	}
	if len(cs) == 0 {
		// TODO: system metrics might be able to continue if 'localhost' is not used.
		return 0, ErrContainerNotFound
	}

	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	measures := make([]float64, 0, len(cs))

	measure := func(c dockerTypes.Container) {
		defer wg.Done()
		for _, n := range c.NetworkSettings.Networks {
			url := m.Metric.Prometheus.Endpoint
			url.Host = strings.Replace(url.Host, "localhost", n.IPAddress, 1)

			req, err2 := http.NewRequest("GET", url.String(), nil)
			if err2 != nil {
				continue
			}

			resp, err2 := m.client.Do(req.WithContext(ctx))
			if err2 != nil {
				continue
			}
			defer drainAndCloseReader(resp.Body)

			sample, err2 := decodeAndFindPrometheusSample(resp.Body, m.Metric.Prometheus.Name)
			if err2 != nil {
				continue
			}

			mutex.Lock()
			measures = append(measures, float64(sample.Value))
			mutex.Unlock()
			break
		}
	}

	for _, c := range cs {
		if c.NetworkSettings != nil {
			wg.Add(1)
			go measure(c)
		}
	}

	wg.Wait()

	if len(measures) < len(cs) {
		if float64(len(measures)) < float64(len(cs))*CriticalFailurePercentage {
			return 0, ErrTooManyFailedMeasurements
		}

		var avg float64
		skipped := float64(len(cs) - len(measures))
		avg, err = reducer.Avg().Reduce(measures)
		if err != nil {
			return 0, err
		}
		measures = append(measures, skipped*avg)
	}

	sum, err := reducer.Sum().Reduce(measures)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func decodeAndFindPrometheusSample(r io.Reader, metricName string) (*model.Sample, error) {
	dec := expfmt.SampleDecoder{
		Dec:  expfmt.NewDecoder(r, expfmt.FmtText),
		Opts: &expfmt.DecodeOptions{},
	}

	for {
		var vec model.Vector

		err := dec.Decode(&vec)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			continue
		}

		if len(vec) == 0 {
			continue
		}

		sample := vec[0]
		name := sample.Metric[model.MetricNameLabel]
		if name == model.LabelValue(metricName) {
			return sample, nil
		}
	}

	return nil, ErrMetricNotFound
}

func drainAndCloseReader(r io.ReadCloser) {
	_, _ = io.CopyN(ioutil.Discard, r, 1024)
	_ = r.Close()
}
