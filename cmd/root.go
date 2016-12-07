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

	log "github.com/Sirupsen/logrus"
	"github.com/mtneug/spate/api"
	"github.com/mtneug/spate/autoscaler"
	"github.com/mtneug/spate/docker"
	"github.com/mtneug/spate/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "spate",
	Short:         "Horizontal service autoscaler for Docker Swarm mode",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		flag, err := cmd.Flags().GetString("log-level")
		if err != nil {
			log.Fatal(err)
		}
		level, err := log.ParseLevel(flag)
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(level)

		i, err := cmd.Flags().GetBool("info")
		if err != nil {
			log.Fatal(err)
		}
		if i {
			_ = version.Spate.PrintFull(os.Stdout)
			if docker.Err == nil {
				_ = docker.PrintInfo(context.Background(), os.Stdout)
			} else {
				fmt.Println("docker: could not connect")
			}
			os.Exit(0)
		}

		v, err := cmd.Flags().GetBool("version")
		if err != nil {
			log.Fatal(err)
		}
		if v {
			fmt.Println(version.Spate)
			os.Exit(0)
		}
	},
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

		if err = a.Start(ctx); err != nil {
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
			_ = s.Stop(ctx)
		}()

		return s.Err(ctx)
	},
}

func init() {
	rootCmd.Flags().Bool("info", false, "Print spate environment information and exit")
	rootCmd.Flags().String("listen-address", ":8080", "Interface to bind to")
	rootCmd.Flags().String("log-level", "info", `Log level ("debug", "info", "warn", "error", "fatal", "panic")`)
	rootCmd.Flags().BoolP("version", "v", false, "Print the version and exit")
}

// Execute invoces the top-level command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("An error occurred")
	}
}
