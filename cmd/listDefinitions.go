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

// listDefinitionsCmd represents the listDefinitions command
var listDefinitionsCmd = &cobra.Command{
	Use:   "definitions",
	Short: "Lists all definitions (queues, exchanges ,etc... ) for all vhosts",
	Long:  `Lists all definitions (queues, exchanges ,etc... ) for all vhosts`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetDefinitions(), nil
	}),
}

// listVhostDefinitionsCmd represents the listVhostDefinitions command
var listVHostDefinitionsCmd = &cobra.Command{
	Use:   "vhost-definitions",
	Short: "Lists all definitions (queues, exchanges ,etc... ) in the vhost.",
	Long:  `Lists all definitions (queues, exchanges ,etc... ) in the vhost.`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetDefinitionsForVhost(api.Config.VHost), nil
	}),
}

func init() {
	listCmd.AddCommand(listDefinitionsCmd)
	listCmd.AddCommand(listVHostDefinitionsCmd)
}
