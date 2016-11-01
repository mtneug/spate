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

package cmd

import (
	"fmt"

	"github.com/mtneug/spate/version"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print information about spate",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("spate:\n")
		fmt.Printf("  Version: %v\n", version.Spate)
		fmt.Printf("  GitCommit: %v\n", version.Spate.GitCommit)
		fmt.Printf("  GitTreeState: %v\n", version.Spate.GitTreeState)
		fmt.Printf("  BuildDate: %v\n", version.Spate.BuildDate)
		fmt.Printf("  GoVersion: %v\n", version.Spate.GoVersion)
		fmt.Printf("  Compiler: %v\n", version.Spate.Compiler)
		fmt.Printf("  Platform: %v\n", version.Spate.Platform)
	},
}
