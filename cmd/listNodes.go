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

// listNodesCmd represents the listNodes command
var listNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Lists all cluster nodes",
	Long:  `Lists all cluster nodes`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetNodes(), nil
	}),
}

// listNodeCmd represents the listNode command
var listNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Lists node details",
	Long:  `Lists node details`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetNode(listNodeName).Memory(listNodeMemory).Binary(listNodeBinary), nil
	}),
}

var listNodeName string
var listNodeMemory bool
var listNodeBinary bool

func init() {
	listCmd.AddCommand(listNodesCmd)
	listCmd.AddCommand(listNodeCmd)
	listNodeCmd.PersistentFlags().StringVarP(&listNodeName, "name", "n", "", "The node name")
	addMemoryAndBinaryFlags(listNodeCmd)
}

func addMemoryAndBinaryFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&listNodeMemory, "memory", "m", false, "Add memory statistics to node, can cause performance degradation, use with care")
	cmd.PersistentFlags().BoolVarP(&listNodeBinary, "binary", "b", false, "Add binary statistics to node, can cause performance degradation, use with care")
}
