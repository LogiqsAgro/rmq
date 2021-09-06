/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/LogiqsAgro/rmq/api"
	"github.com/spf13/cobra"
)

// checkAlarmsCmd represents the checkAlarms command
var checkAlarmsCmd = &cobra.Command{
	Use:   "alarms",
	Short: "Check if there are no alarms in effect in the cluster.",
	Long:  `Check if there are no alarms in effect in the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		json, err := api.GetHealthChecksAlarmsJson()
		api.Print(json, err)
	},
}

// checkAlarmsCmd represents the checkAlarms command
var checkLocalAlarmsCmd = &cobra.Command{
	Use:   "alarms",
	Short: "Check if there are no local alarms in effect on the target node.",
	Long:  `Check if there are no local alarms in effect on the target node.`,
	Run: func(cmd *cobra.Command, args []string) {
		json, err := api.GetHealthChecksLocalAlarmsJson()
		api.Print(json, err)
	},
}

func init() {
	checkCmd.AddCommand(checkAlarmsCmd)
	checkCmd.AddCommand(checkLocalAlarmsCmd)
}
