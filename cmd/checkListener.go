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
	"strings"

	"github.com/LogiqsAgro/rmq/api"
	"github.com/spf13/cobra"
)

// checkListenerCmd represents the checkListener command
var checkListenerCmd = &cobra.Command{
	Use:   "listener",
	Short: "Checks for an active listener on port or protocol.",
	Long:  `Responds a 200 OK if there is an active listener on the given port or protocol, otherwise responds with a 503 Service Unavailable.`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		if cmd.Flags().Changed("port") {
			if checkListenerPort < 1 || checkListenerPort > 65535 {
				return nil, fmt.Errorf("invalid port value: %d", checkListenerPort)
			}
			return api.GetHealthChecksPortListener(checkListenerPort), nil
		} else if cmd.Flags().Changed("protocol") {
			return api.GetHealthChecksProtocolListener(checkListenerProtocol), nil

		} else {
			return nil, fmt.Errorf("no port or protocol specified")
		}
	}),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		count := 0
		names := []string{"port", "protocol"}
		for _, name := range names {
			if cmd.Flags().Changed(name) {
				count++
			}
		}
		if count != 1 {
			return fmt.Errorf("exactly one of the --%s flags must be specified", strings.Join(names, " or --"))
		}
		return nil
	},
}

var checkListenerPort int
var checkListenerProtocol string

func init() {
	checkCmd.AddCommand(checkListenerCmd)
	checkListenerCmd.PersistentFlags().IntVarP(&checkListenerPort, "port", "", 0, "The RabbitMQ listener port (0-65535), cannot be used together with --protocol")
	checkListenerCmd.PersistentFlags().StringVarP(&checkListenerProtocol, "protocol", "", "", "Some valid protocol names are: amqp091, amqp10, mqtt, stomp, web-mqtt, web-stomp, cannot be used together with --port")
}
