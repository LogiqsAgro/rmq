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

// listGlobalParametersCmd represents the listGlobalParameters command
var listGlobalParametersCmd = &cobra.Command{
	Use:   "global-parameters",
	Short: "Lists all global parameters",
	Long:  `Lists all global parameters`,
	Run: func(cmd *cobra.Command, args []string) {
		json, err := api.GetGlobalParametersJson()
		api.Print(json, err)
	},
}

// listGlobalParametersCmd represents the listGlobalParameters command
var listGlobalParameterCmd = &cobra.Command{
	Use:   "global-parameter",
	Short: "Lists a global parameter",
	Long:  `Lists a global parameter`,
	Run: func(cmd *cobra.Command, args []string) {
		json, err := api.GetGlobalParameterJson(parameterName)
		api.Print(json, err)
	},
}

func init() {
	listCmd.AddCommand(listGlobalParametersCmd)

	listCmd.AddCommand(listGlobalParameterCmd)
	listGlobalParameterCmd.PersistentFlags().StringVarP(&parameterName, "name", "n", "", "the parameter name")
}

var (
	parameterName string
)
