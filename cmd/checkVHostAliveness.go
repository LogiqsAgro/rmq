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

// checkVHostAlivenessCmd represents the checkVHostAliveness command
var checkVHostAlivenessCmd = &cobra.Command{
	Use:   "vhost-aliveness",
	Short: "A very basic vhost health check.",
	Long:  `Declares a test queue on the target node, then publishes and consumes a message. Intended to be used as a very basic health check.`,
	Run: func(cmd *cobra.Command, args []string) {
		json, err := api.GetAlivenessTestJson(api.Config.VHost)
		api.Print(json, err)
	},
}

func init() {
	checkCmd.AddCommand(checkVHostAlivenessCmd)
}
