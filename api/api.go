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

package api

import (
	"context"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

// Config for an API server.
type Config struct {
	Addr string
}

// Server implements an API server.
type Server struct {
	config *Config
}

// New creates a new API server.
func New(c *Config) (*Server, error) {
	s := &Server{
		config: c,
	}
	return s, nil
}

// Start the API server and listens for requests.
func (s *Server) Start(context.Context) error {
	http.Handle("/metrics", prometheus.Handler())

	log.Info("Starting API server")
	go http.ListenAndServe(s.config.Addr, nil)
	log.Info("API server started")

	return nil
}
