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
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

// Config for a spate API server.
type Config struct {
	Addr string
}

// Server implements a spate API server.
type Server struct {
	config   *Config
	server   *http.Server
	err      error
	doneChan chan struct{}
}

// New creates a new server.
func New(c *Config) (*Server, error) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", prometheus.Handler())
	// TODO: connect ErrorLog to logrus
	srv := &Server{
		config:   c,
		server:   &http.Server{Handler: mux},
		doneChan: make(chan struct{}),
	}

	return srv, nil
}

// Start the API server and listens for requests.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return err
	}

	log.Info("Starting API server")
	go func() {
		s.err = s.server.Serve(ln)
		close(s.doneChan)
		log.Info("API server stopped")
	}()
	log.Info("API server started")

	return nil
}

// Err returns an error object after the server has stopped.
func (s *Server) Err(ctx context.Context) error {
	select {
	case <-s.doneChan:
		return s.err
	case <-ctx.Done():
		return ctx.Err()
	}
}
