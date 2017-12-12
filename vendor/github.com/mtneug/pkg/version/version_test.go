// Copyright (c) 2016 Matthias Neugebauer
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package version_test

import (
	"bytes"
	"testing"

	"github.com/mtneug/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t.Parallel()

	i := version.Info{
		Major:     "1",
		Minor:     "2",
		Patch:     "3",
		GitCommit: "abcdefg",
	}
	assert.Equal(t, "1.2.3+abcdefg", i.String())
}

func TestPrintFull(t *testing.T) {
	t.Parallel()

	i := version.Info{
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
	assert.NoError(t, err)

	want := `myapp:
  Version:        1.2.3+abcdefg
  Git Commit:     abcdefg
  Git Tree State: clean
  Build Date:     2016-12-05 19:52:06 UTC
  Go Version:     go1.7.4
  Compiler:       gc
  Platform:       linux/amd64
`
	assert.Equal(t, want, b.String())
}
