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
	"fmt"

	"github.com/LogiqsAgro/rmq/api"
	"github.com/spf13/cobra"
)

// checkCertificateExpirationCmd represents the checkCertificateExpiration command
var checkCertificateExpirationCmd = &cobra.Command{
	Use:   "certificate-expiration",
	Short: "Checks the expiration date on the certificates for every listener configured to use TLS.",
	Long:  `Checks the expiration date on the certificates for every listener configured to use TLS.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, unit := range api.AllUnits() {
			name := withinFlagName(unit)
			if cmd.Flags().Changed(name) {
				within, err := cmd.PersistentFlags().GetInt(name)
				if err == nil && within > 0 {
					json, err := api.GetHealthChecksCertificateExpirationJson(within, unit)
					api.Print(json, err)
					return
				}
			}
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		count := 0
		names := []interface{}{}
		for _, unit := range api.AllUnits() {
			name := withinFlagName(unit)
			names = append(names, "--"+name)
			if cmd.Flags().Changed(name) {
				count++
			}
		}
		if count != 1 {
			return fmt.Errorf("exactly one of the %v parameters must be specified", names)
		}
		return nil
	},
}

func withinFlagName(name string) string {
	return "within-" + name
}

func init() {
	checkCmd.AddCommand(checkCertificateExpirationCmd)
	pflags := checkCertificateExpirationCmd.PersistentFlags()
	for _, unit := range api.AllUnits() {
		pflags.IntP(withinFlagName(unit), unit[0:1], 0, "The number of "+unit+" within which the certificate expires")
	}
}
