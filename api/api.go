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

package api

import (
	"context"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/mtneug/pkg/startstopper"
	"github.com/prometheus/client_golang/prometheus"
)

// Server implements a spate API server.
type Server struct {
	startstopper.StartStopper

	Addr   string
	server *http.Server
}

// New creates a new server.
func New(addr string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", prometheus.Handler())

	loggedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("%s \"%s %s %s\"", r.RemoteAddr, r.Method, r.URL, r.Proto)
		mux.ServeHTTP(w, r)
	})

	s := &Server{
		Addr:   addr,
		server: &http.Server{Handler: loggedMux},
	}
	s.StartStopper = startstopper.NewGo(startstopper.RunnerFunc(s.run))

	return s
}

func (s *Server) run(ctx context.Context, stopChan <-chan struct{}) error {
	log.Debug("API server started")
	defer log.Debug("API server stopped")

	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	doneChan := make(chan struct{})
	go func() {
		err = s.server.Serve(ln)
		close(doneChan)
	}()

	select {
	case <-doneChan:
	case <-stopChan:
	case <-ctx.Done():
	}

	_ = s.server.Shutdown(ctx)

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}
