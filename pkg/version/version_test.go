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

package version

import (
	"bytes"
	"testing"
)

func TestString(t *testing.T) {
	i := Info{
		Major:     "1",
		Minor:     "2",
		Patch:     "3",
		GitCommit: "abcdefg",
	}

	got := i.String()
	want := "1.2.3+abcdefg"
	if got != want {
		t.Errorf("Expected '%s', got '%s'", want, got)
	}
}

func TestPrintFull(t *testing.T) {
	i := Info{
		Name:         "myapp",
		Major:        "1",
		Minor:        "2",
		Patch:        "3",
		GitCommit:    "abcdefg",
		GitTreeState: "clean",
		BuildDate:    "2016-12-05 19:52:06 UTC",
		GoVersion:    "go1.7.4",
		Compiler:     "gc",
		Platform:     "linux/amd64",
	}

	var b bytes.Buffer
	err := i.PrintFull(&b)
	if err != nil {
		t.Fatalf("Expected err to be nil, got '%s'", err)
	}

	got := b.String()
	want := `myapp:
  Version:        1.2.3+abcdefg
  Git Commit:     abcdefg
  Git Tree State: clean
  Build Date:     2016-12-05 19:52:06 UTC
  Go Version:     go1.7.4
  Compiler:       gc
  Platform:       linux/amd64
`

	if got != want {
		t.Fatalf("Expected '%s', got '%s'", want, got)
	}
}
