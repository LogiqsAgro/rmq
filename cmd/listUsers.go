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

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Lists all users",
	Long:  `Lists all users`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetUsers(), nil
	}),
}

func init() {
	listCmd.AddCommand(listUsersCmd)
}
