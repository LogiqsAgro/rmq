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

// listConnectionsCmd represents the listConnections command
var listConnectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "Lists all connections",
	Long:  `Lists all connections`,
	RunE: func(cmd *cobra.Command, args []string) error {
		req := api.GetConnections()
		api.ApplyConfig(req)
		resp, err := api.Do(req)
		return api.Print(resp, err)
	},
}

// listConnectionsCmd represents the listConnections command
var listVHostConnectionsCmd = &cobra.Command{
	Use:   "vhost-connections",
	Short: "Lists the connections for a vhost",
	Long:  `Lists the connections for a vhost`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetConnectionsForVhost(api.Config.VHost), nil
	}),
}

// listConnectionsCmd represents the listConnections command
var listConnectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "Lists details for the connection with the given --name",
	Long:  `Lists details for the connection with the given --name`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetConnection(connectionName), nil
	}),
}

func init() {
	listCmd.AddCommand(listConnectionsCmd)
	api.AddPagingFlags(listConnectionsCmd)

	listCmd.AddCommand(listVHostConnectionsCmd)
	api.AddPagingFlags(listVHostConnectionsCmd)

	listCmd.AddCommand(listConnectionCmd)
	listConnectionCmd.PersistentFlags().StringVarP(&connectionName, "name", "n", "", "The connection name")
}

var (
	connectionName string
)
