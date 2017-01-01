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

package cmd

import (
	"time"

	"github.com/mtneug/spate/label"
	flag "github.com/spf13/pflag"
)

func readAndSetDefaults(flags *flag.FlagSet) error {
	// --default-aggregation-amount
	defaultAggregationAmount, err := flags.GetUint8("default-aggregation-amount")
	if err != nil {
		return err
	}
	label.DefaultObserverAggregationAmount = defaultAggregationAmount

	// --default-autoscaler-period
	defaultAutoscalerPeriodStr, err := flags.GetString("default-autoscaler-period")
	if err != nil {
		return err
	}
	defaultAutoscalerPeriod, err := time.ParseDuration(defaultAutoscalerPeriodStr)
	if err != nil {
		return err
	}
	label.DefaultAutoscalerPeriod = defaultAutoscalerPeriod

	// --default-cooldown-scaled_up
	defaultCooldownScaledUpStr, err := flags.GetString("default-cooldown-scaled_up")
	if err != nil {
		return err
	}
	defaultCooldownScaledUp, err := time.ParseDuration(defaultCooldownScaledUpStr)
	if err != nil {
		return err
	}
	label.DefaultCooldownScaledUp = defaultCooldownScaledUp

	// --default-cooldown-scaled_down
	defaultCooldownScaledDownStr, err := flags.GetString("default-cooldown-scaled_down")
	if err != nil {
		return err
	}
	defaultCooldownScaledDown, err := time.ParseDuration(defaultCooldownScaledDownStr)
	if err != nil {
		return err
	}
	label.DefaultCooldownScaledDown = defaultCooldownScaledDown

	// --default-cooldown-service_added
	defaultCooldownServiceAddedStr, err := flags.GetString("default-cooldown-service_added")
	if err != nil {
		return err
	}
	defaultCooldownServiceAdded, err := time.ParseDuration(defaultCooldownServiceAddedStr)
	if err != nil {
		return err
	}
	label.DefaultCooldownServiceAdded = defaultCooldownServiceAdded

	// --default-cooldown-service_updated
	defaultCooldownServiceUpdatedStr, err := flags.GetString("default-cooldown-service_updated")
	if err != nil {
		return err
	}
	defaultCooldownServiceUpdated, err := time.ParseDuration(defaultCooldownServiceUpdatedStr)
	if err != nil {
		return err
	}
	label.DefaultCooldownServiceUpdated = defaultCooldownServiceUpdated

	// --default-observer-period
	defaultObserverPeriodStr, err := flags.GetString("default-observer-period")
	if err != nil {
		return err
	}
	defaultObserverPeriod, err := time.ParseDuration(defaultObserverPeriodStr)
	if err != nil {
		return err
	}
	label.DefaultObserverPeriod = defaultObserverPeriod

	return nil
}

func init() {
	rootCmd.Flags().Uint8("default-aggregation-amount", 5, "Default amount of measurements to take for aggregation")
	rootCmd.Flags().String("default-autoscaler-period", "30s", "Default autoscaler period")
	rootCmd.Flags().String("default-cooldown-scaled_up", "3m", "Default cool down time after a scale up")
	rootCmd.Flags().String("default-cooldown-scaled_down", "5m", "Default cool down time after a scaled down")
	rootCmd.Flags().String("default-cooldown-service_added", "0s", "Default cool down time after a service is added")
	rootCmd.Flags().String("default-cooldown-service_updated", "0s", "Default cool down time after a service is updated")
	rootCmd.Flags().String("default-observer-period", "30s", "Default observer period")
}
