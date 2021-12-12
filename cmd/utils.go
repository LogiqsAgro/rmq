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
	"os"

	"github.com/LogiqsAgro/rmq/api"
	"github.com/spf13/cobra"
)

func writeError(err error) bool {
	if err != nil {
		os.Stderr.WriteString("ERROR: " + err.Error())
		return true
	}
	return false
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func RunE(request func(*cobra.Command, []string) (api.Builder, error)) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		req, err := request(cmd, args)
		if err != nil {
			return err
		}

		// silence usage message from here on...
		// any errors following here
		// are not due to mis-usage of the command
		cmd.SilenceUsage = true

		api.ApplyConfig(req)
		resp, err := api.Do(req)
		return api.Print(resp, err)
	}
}
