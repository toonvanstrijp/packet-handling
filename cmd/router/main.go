// Copyright 2022 Stichting ThingsIX Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/ThingsIXFoundation/packet-handling/router"
	"github.com/ThingsIXFoundation/packet-handling/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "router",
	Short:   "run the router service",
	Args:    cobra.RangeArgs(0, 1),
	Run:     router.Run,
	Version: utils.Version(),
}

func init() {
	keyCmd.AddCommand(genKeyCmd)
	rootCmd.AddCommand(keyCmd)

	rootCmd.PersistentFlags().String("config", "", "configuration file")
	err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		logrus.WithError(err).Fatal("could not find viper flag")
	}
}
