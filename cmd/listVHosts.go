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

// listBindingsCmd represents the listBindings command
var listVHostsCmd = &cobra.Command{
	Use:   "vhosts",
	Short: "Lists all vhosts",
	Long:  `Lists all vhosts`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetVhosts(), nil
	}),
}

// listLimitsCmd represents the listLimits command
var listLimitsCmd = &cobra.Command{
	Use:   "limits",
	Short: "Lists limits for all vhosts",
	Long:  `Lists limits for all vhosts`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetVhostLimits(), nil
	}),
}

// listVHostLimitsCmd represents the listVHostLimits command
var listVHostLimitsCmd = &cobra.Command{
	Use:   "vhost-limits",
	Short: "Lists limits for the vhost",
	Long:  `Lists limits for the vhost`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetVhostLimitsForVhost(api.Config.VHost), nil
	}),
}

func init() {
	listCmd.AddCommand(listVHostsCmd)
	listCmd.AddCommand(listLimitsCmd)
	listCmd.AddCommand(listVHostLimitsCmd)
}
