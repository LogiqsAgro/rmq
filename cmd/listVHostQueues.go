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

// listVHostQueuesCmd represents the listQueues command
var listVHostQueuesCmd = &cobra.Command{
	Use:   "vhost-queues",
	Short: "Lists the queues defined in a vhost",
	Long:  `Lists all the queues defined in a vhost`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetQueuesForVhost(api.Config.VHost), nil
	}),
}

// listVHostQueuesCmd represents the listQueues command
var listVHostQueueCmd = &cobra.Command{
	Use:   "vhost-queue",
	Short: "Lists the named queue defined in the vhost",
	Long:  `Lists the named queue defined in the vhost`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetQueueForVhost(api.Config.VHost, listVHostQueueName), nil
	}),
}

var listVHostQueueName string = ""

func init() {
	listCmd.AddCommand(listVHostQueuesCmd)
	listCmd.AddCommand(listVHostQueueCmd)
	listVHostQueueCmd.PersistentFlags().StringVarP(&listVHostQueueName, "name", "n", "NAME", "The name of the queue to list")
}
