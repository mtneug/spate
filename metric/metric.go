package metric

import (
	"github.com/mtneug/pkg/ulid"
	"github.com/mtneug/spate/api/types"
)

// New creates a new metric.
func New(name string) types.Metric {
	m := types.Metric{
		ID:   ulid.New().String(),
		Name: name,
	}
	return m
}
