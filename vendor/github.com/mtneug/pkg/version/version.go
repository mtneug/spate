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

package version

import (
	"fmt"
	"io"
	"text/template"
)

var printTemplate = `{{.Name}}:
  Version:        {{.Major}}.{{.Minor}}.{{.Patch}}+{{.GitCommit}}
  Git Commit:     {{.GitCommit}}
  Git Tree State: {{.GitTreeState}}
  Build Date:     {{.BuildDate}}
  Go Version:     {{.GoVersion}}
  Compiler:       {{.Compiler}}
  Platform:       {{.Platform}}
`

// Info holds version and build information. The fields are largly the same as
// in the `k8s.io/kubernetes/pkg/version` package of the Kubernetes project.
type Info struct {
	// Name of the versioned object.
	Name string `json:"name"`
	// Major version number.
	Major string `json:"major"`
	// Minor version number.
	Minor string `json:"minor"`
	// Patch version number.
	Patch string `json:"patch"`
	// GitCommit SHA.
	GitCommit string `json:"gitCommit"`
	// GitTreeState is either "clean" or "dirty".
	GitTreeState string `json:"gitTreeState"`
	// BuildDate of the binary.
	BuildDate string `json:"buildDate"`
	// GoVersion of the binary.
	GoVersion string `json:"goVersion"`
	// Compiler used for the binary.
	Compiler string `json:"compiler"`
	// Platform the binary is compiled for.
	Platform string `json:"platform"`
}

// String returns a formated version string.
func (i Info) String() string {
	return fmt.Sprintf("%v.%v.%v+%v", i.Major, i.Minor, i.Patch, i.GitCommit)
}

// PrintFull writes the version
func (i Info) PrintFull(w io.Writer) (err error) {
	t, err := template.New("info").Parse(printTemplate)
	if err != nil {
		return
	}
	err = t.Execute(w, i)
	return
}
