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

// listDefinitionsCmd represents the listDefinitions command
var listDefinitionsCmd = &cobra.Command{
	Use:   "definitions",
	Short: "Lists all definitions (queues, exchanges ,etc... ) for all vhosts",
	Long:  `Lists all definitions (queues, exchanges ,etc... ) for all vhosts`,
	Run: func(cmd *cobra.Command, args []string) {
		json, err := api.GetDefinitionsJson()
		api.Print(json, err)
	},
}

// listVhostDefinitionsCmd represents the listVhostDefinitions command
var listVHostDefinitionsCmd = &cobra.Command{
	Use:   "vhost-definitions",
	Short: "Lists all definitions (queues, exchanges ,etc... ) in the vhost.",
	Long:  `Lists all definitions (queues, exchanges ,etc... ) in the vhost.`,
	Run: func(cmd *cobra.Command, args []string) {
		json, err := api.GetVHostDefinitionsJson(api.Config.VHost)
		api.Print(json, err)
	},
}

func init() {
	listCmd.AddCommand(listDefinitionsCmd)
	listCmd.AddCommand(listVHostDefinitionsCmd)
}
