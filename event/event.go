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

package event

import (
	"docker.io/go-docker/api/types/swarm"
	"github.com/mtneug/pkg/ulid"
)

// Type represents some category of events.
type Type string

const (
	// TypeServiceCreated indicates a service creation event.
	TypeServiceCreated Type = "service_created"
	// TypeServiceUpdated indicates a service update event.
	TypeServiceUpdated Type = "service_updated"
	// TypeServiceDeleted indicates a service deletion event.
	TypeServiceDeleted Type = "service_deleted"
	// TypeServiceScaledUp indicates a service scale up event.
	TypeServiceScaledUp Type = "service_scaled_up"
	// TypeServiceScaledDown indicates a service scale down event.
	TypeServiceScaledDown Type = "service_scaled_down"
)

// Event represents some incident.
type Event struct {
	// ID of the event.
	ID string
	// Type of the event.
	Type Type
	// Service relevant to the event.
	Service swarm.Service
}

// New creates a new event.
func New(t Type, srv swarm.Service) Event {
	return Event{
		ID:      ulid.New().String(),
		Type:    t,
		Service: srv,
	}
}
