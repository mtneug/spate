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

package controller

import (
	"context"
	"testing"
	"time"

	"github.com/mtneug/spate/event"
	"github.com/stretchr/testify/require"
)

func TestChangeLoopRun(t *testing.T) {
	t.Parallel()
	t.Skip("Refactor so that the Docker Client can be mocked")

	eq := make(chan event.Event)
	cl := newChangeLoop(time.Second, eq, nil)

	// stopChan
	stopChan := make(chan struct{})
	close(stopChan)
	err := cl.run(context.Background(), stopChan)
	require.NoError(t, err)
	stopChan = make(chan struct{})

	// ctx
	ctx, cancle := context.WithCancel(context.Background())
	cancle()
	err = cl.run(ctx, stopChan)
	require.EqualError(t, err, "context canceled")
}
