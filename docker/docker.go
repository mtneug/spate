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

package docker

import (
	"context"
	"io"
	"strconv"
	"text/template"

	"github.com/docker/docker/client"
)

var (
	// C is the default Docker client instance.
	C *client.Client

	// Err is the error that was returned when creating the Docker client.
	Err error

	printTemplate = `docker:
  ID:             {{.ID}}
  API Version:    {{.APIVersion}}
  Server Version: {{.ServerVersion}}
  Swarm:          {{.SwarmLocalNodeState}}
    Cluster ID:   {{.SwarmID}}
    Nodes:        {{.SwarmNodes}}
    Managers:     {{.SwarmManagers}}
`
)

func init() {
	C, Err = client.NewEnvClient()
}

// PrintInfo writes informations about Docker relevant to spate to the given
// Writer.
//
//	// Print to Stdout
//	docker.PrintInfo(context.Background(), os.Stdout)
func PrintInfo(ctx context.Context, w io.Writer) (err error) {
	i, err := C.Info(ctx)
	if err != nil {
		return
	}

	s, err := C.ServerVersion(ctx)
	if err != nil {
		return
	}

	t, err := template.New("info").Parse(printTemplate)
	if err != nil {
		return
	}

	d := struct {
		ID,
		APIVersion,
		ServerVersion,
		SwarmLocalNodeState,
		SwarmID,
		SwarmNodes,
		SwarmManagers string
	}{
		i.ID,
		s.APIVersion,
		i.ServerVersion,
		string(i.Swarm.LocalNodeState),
		i.Swarm.Cluster.ID,
		strconv.Itoa(i.Swarm.Nodes),
		strconv.Itoa(i.Swarm.Managers),
	}

	err = t.Execute(w, d)
	return
}
