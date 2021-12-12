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

// checkNodeIsMirrorSyncCriticalCmd represents the checkNodeIsMirrorSyncCritical command
var checkNodeIsMirrorSyncCriticalCmd = &cobra.Command{
	Use:   "node-is-mirror-sync-critical",
	Short: "Checks if there are classic mirrored queues without synchronised mirrors online (queues that would potentially lose data if the target node is shut down).",
	Long:  `Checks if there are classic mirrored queues without synchronised mirrors online (queues that would potentially lose data if the target node is shut down).`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetHealthChecksNodeIsMirrorSyncCritical(), nil
	}),
}

// checkVHostAlivenessCmd represents the checkVHostAliveness command
var checkNodeIsQuorumCriticalCmd = &cobra.Command{
	Use:   "node-is-quorum-critical",
	Short: "Checks if there are quorum queues with minimum online quorum (queues that would lose their quorum and availability if the target node is shut down).",
	Long:  `Checks if there are quorum queues with minimum online quorum (queues that would lose their quorum and availability if the target node is shut down).`,
	RunE: RunE(func(cmd *cobra.Command, args []string) (api.Builder, error) {
		return api.GetHealthChecksNodeIsQuorumCritical(), nil
	}),
}

func init() {
	checkCmd.AddCommand(checkNodeIsMirrorSyncCriticalCmd)
	checkCmd.AddCommand(checkNodeIsQuorumCriticalCmd)
}
