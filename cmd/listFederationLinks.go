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

var listFederationLinksCmd = &cobra.Command{
	Use:   "federation-links",
	Short: "Lists status for all federation links. Requires the rabbitmq_federation_management plugin to be enabled.",
	Long:  `Lists status for all federation links. Requires the rabbitmq_federation_management plugin to be enabled.`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetFederationLinks(), nil
	}),
}

var listVhostFederationLinksCmd = &cobra.Command{
	Use:   "vhost-federation-links",
	Short: "Lists status for vhost federation links. Requires the rabbitmq_federation_management plugin to be enabled.",
	Long:  `Lists status for vhost federation links. Requires the rabbitmq_federation_management plugin to be enabled.`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetFederationLinksForVhost(api.Config.VHost), nil
	}),
}

func init() {
	listCmd.AddCommand(listFederationLinksCmd)
	listCmd.AddCommand(listVhostFederationLinksCmd)
}
