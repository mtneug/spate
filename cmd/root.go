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
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/mtneug/spate/api"
	"github.com/mtneug/spate/autoscaler"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spate",
	Short: "Horizontal service autoscaler for Docker Swarm mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		flags := cmd.Flags()

		addr, err := flags.GetString("listen-address")
		if err != nil {
			return err
		}

		// API server
		a, err := api.New(&api.Config{
			Addr: addr,
		})
		if err != nil {
			return err
		}

		if err := a.Start(ctx); err != nil {
			return err
		}

		// Scaler
		s, err := autoscaler.New(&autoscaler.Config{})
		if err != nil {
			return err
		}

		if err := s.Start(ctx); err != nil {
			return err
		}

		// Handle termination
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		go func() {
			<-sig
			s.Stop(ctx)
		}()

		return s.Err(ctx)
	},
}

func init() {
	rootCmd.Flags().String("listen-address", ":8080", "Interface to bind to")

	rootCmd.AddCommand(
		infoCmd,
		versionCmd,
	)
}

// Execute invoces the top-level command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
