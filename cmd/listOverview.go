/*
Copyright © 2021 Remco Schoeman <remco.schoeman@logiqs.nl>

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

// listOverviewCmd represents the listOverview command
var listOverviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Lists information about the cluster, node and runtime.",
	Long:  `Lists information about the cluster, node and runtime.`,
	Run: func(cmd *cobra.Command, args []string) {
		json, err := api.GetOverviewJson()
		api.Print(json, err)
	},
}

func init() {
	listCmd.AddCommand(listOverviewCmd)
}
