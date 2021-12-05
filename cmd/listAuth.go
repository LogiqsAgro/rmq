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
	"fmt"

	"github.com/LogiqsAgro/rmq/api"
	"github.com/spf13/cobra"
)

// listAuthCmd represents the listAuth command
var listAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Lists details about the OAuth2 configuration",
	Long:  `Lists details about the OAuth2 configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		json, err := api.GetAuthJson()
		api.Print(json, err)
	},
}

// listAuthCmd represents the listAuth command
var listAuthAttempts = &cobra.Command{
	Use:   "auth-attempts",
	Short: "Lists authentication attempts on the specified node.",
	Long:  `Lists authentication attempts on the specified node.`,
	Run: func(cmd *cobra.Command, args []string) {
		json, err := api.GetAuthAttemptsJson(listAuthAttemptsNode)
		api.Print(json, err)
	},
	PreRunE: validateListAuthAttemptsNode,
}

// listAuthCmd represents the listAuth command
var listAuthAttemptsBySource = &cobra.Command{
	Use:   "auth-attempts-by-source",
	Short: "Lists authentication attempts by source on the specified node. 'track_auth_attempt_source' must be enabled in the RabbitMQ advanced config: see https://blog.rabbitmq.com/posts/2021/03/auth-attempts-metrics/",
	Long:  `Lists authentication attempts by source on the specified node. 'track_auth_attempt_source' must be enabled in the RabbitMQ advanced config: see https://blog.rabbitmq.com/posts/2021/03/auth-attempts-metrics/`,
	Run: func(cmd *cobra.Command, args []string) {
		json, err := api.GetAuthAttemptsBySourceJson(listAuthAttemptsNode)
		api.Print(json, err)
	},
	PreRunE: validateListAuthAttemptsNode,
}

var listAuthAttemptsNode string

func validateListAuthAttemptsNode(cmd *cobra.Command, args []string) error {

	if len(listAuthAttemptsNode) == 0 {
		return fmt.Errorf("--node ( or -n ) is a required parameter, use 'rmq list nodes --columns name' to get valid node names")
	}
	return nil
}

func init() {
	listCmd.AddCommand(listAuthCmd)
	listCmd.AddCommand(listAuthAttempts)
	listCmd.AddCommand(listAuthAttemptsBySource)
	description := "The node from where to get the attempts list, use 'rmq list nodes --columns name' to get the node names"
	listAuthAttempts.PersistentFlags().StringVarP(&listAuthAttemptsNode, "node", "n", "", description)
	listAuthAttemptsBySource.PersistentFlags().StringVarP(&listAuthAttemptsNode, "node", "n", "", description)
}
