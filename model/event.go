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

package model

// EventType represents some category of events.
type EventType string

const (
	// EventTypeServiceCreated indicates a service creation event.
	EventTypeServiceCreated EventType = "service_created"
	// EventTypeServiceUpdated indicates a service update event.
	EventTypeServiceUpdated EventType = "service_updated"
	// EventTypeServiceDeleted indicates a service deletion event.
	EventTypeServiceDeleted EventType = "service_deleted"
	// EventTypeServiceScaledUp indicates a service scale up event.
	EventTypeServiceScaledUp EventType = "service_scaled_up"
	// EventTypeServiceScaledDown indicates a service scale down event.
	EventTypeServiceScaledDown EventType = "service_scaled_down"
)

// Event represents some incident.
type Event struct {
	// ID of the event.
	ID string
	// Type of the event.
	Type EventType
	// Object relevant to the event.
	Object interface{}
}
