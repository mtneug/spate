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
	"fmt"
	"runtime"

	"github.com/mtneug/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	major        = "0"
	minor        = "0"
	patch        = "1"
	gitCommit    string // set by Makefile
	gitTreeState string // set by Makefile
	buildDate    string // set by Makefile
	goVersion    = runtime.Version()
	compiler     = runtime.Compiler
	platform     = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	// Spate exposes the version number.
	Spate = version.Info{
		Name:         "spate",
		Major:        major,
		Minor:        minor,
		Patch:        patch,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    goVersion,
		Compiler:     compiler,
		Platform:     platform,
	}
)

func init() {
	spateInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "spate_info",
			Help: "A metric with a constant '1' value labeled by major, minor, " +
				"patch, git commit, git tree state, build date, Go version, " +
				"compiler, and platform.",
		},
		[]string{
			"major",
			"minor",
			"patch",
			"gitCommit",
			"gitTreeState",
			"buildDate",
			"goVersion",
			"compiler",
			"platform",
		},
	)
	spateInfo.WithLabelValues(
		Spate.Major,
		Spate.Minor,
		Spate.Patch,
		Spate.GitCommit,
		Spate.GitTreeState,
		Spate.BuildDate,
		Spate.GoVersion,
		Spate.Compiler,
		Spate.Platform,
	).Set(1)

	prometheus.MustRegister(spateInfo)
}
