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
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// initConfig reads in config file and ENV variables if set.
func initializeConfig(cmd *cobra.Command) error {

	v := viper.GetViper()

	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		// Search config in home directory with name ".rmq" (without extension).
		v.AddConfigPath(home)
		v.SetConfigType("yaml")
		v.SetConfigName(".rmq")
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	v.ReadInConfig()

	//if err := v.ReadInConfig(); err == nil {
	// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	//}

	bindFlags(cmd, v, envPrefix)
	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper, envPrefix string) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			strVal := fmt.Sprintf("%v", val)
			cmd.Flags().Set(f.Name, strVal)
			// fmt.Println("Set flag " + f.Name + " to '" + strVal + "' from config file")
		}
	})
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long:  ` `,
	Run:   nil,
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Save a setting in the config file",
	Long:  ` `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("specify a name and a value to store")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if strings.EqualFold(args[0], "password") {
			fmt.Println("You really shouldn't do that, now you will have a cleartext password stored in " + viper.ConfigFileUsed())
			fmt.Println("You can undo your mistake with the command 'rmq unset " + args[0] + "'")
		}

		viper.Set(args[0], args[1])

		err := viper.SafeWriteConfig()
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
			err := viper.WriteConfig()
			panicIf(err)
		}
	},
}

var configUnsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Remove a setting from the config file",
	Long:  ` `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("specify one or more config setting names")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		configFile := viper.ConfigFileUsed()
		toDelete := map[string]string{}
		for _, key := range args {
			lowerKey := strings.ToLower(key)
			toDelete[lowerKey] = key
		}

		vOld := viper.New()
		vOld.SetConfigFile(configFile)
		err := vOld.ReadInConfig()
		if err != nil {
			panicIf(err)
		}

		vNew := viper.New()
		vNew.SetConfigFile(configFile)

		for _, key := range vOld.AllKeys() {
			lowerKey := strings.ToLower(key)
			if _, ok := toDelete[lowerKey]; ok {
				if strings.EqualFold(key, "password") {
					fmt.Println("I'm really happy you reconsidered storing cleartext passwords on disk!" +
						" Removing it now from " + viper.ConfigFileUsed())
				}
				continue
			}
			vNew.Set(key, vOld.Get(key))
		}

		err = vNew.WriteConfig()
		panicIf(err)
	},
}

var configCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Created the config file if it doesn't exist yet",
	Long:  ` `,

	Run: func(cmd *cobra.Command, args []string) {
		err := viper.SafeWriteConfig()
		if err == nil {
			fmt.Println("Created config file at", viper.ConfigFileUsed())
		} else {
			if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
				fmt.Println("Config file already exists at", viper.ConfigFileUsed())
			} else {
				writeError(err)
			}
		}
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current settings loaded from config file, environment and commandline parameters",
	Long:  ` `,

	Run: func(cmd *cobra.Command, args []string) {

		if fileOnly, err := cmd.PersistentFlags().GetBool(fileOnlyFlag); err == nil && fileOnly {
			v := viper.New()
			v.SetConfigFile(viper.ConfigFileUsed())
			err = v.ReadInConfig()
			if writeError(err) {
				return
			}
			keys := v.AllKeys()
			if len(keys) == 0 {
				fmt.Println("No settings configured in config file " + viper.ConfigFileUsed())
			} else {
				fmt.Printf("The current settings loaded from config file '%s'\n", viper.ConfigFileUsed())
			}
			printViperContents(v)

		} else if envOnly, err := cmd.PersistentFlags().GetBool(envOnlyFlag); err == nil && envOnly {

			v := viper.New()
			v.SetEnvPrefix(envPrefix)
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
			v.AutomaticEnv() // read in environment variables that match

			for _, key := range viper.AllKeys() {
				// the environment values are loaded on demand
				// so we check if we find a value, and then
				// explicitly set it so the keys with settings
				// from the environment will appear in v.AllKeys()
				var val = v.Get(key)
				if val != nil {
					v.Set(key, val)
				}
			}

			keys := v.AllKeys()
			if len(keys) == 0 {
				fmt.Println("No settings configured in environment")
			} else {
				fmt.Printf("The current settings loaded from environment variables with prefix '%s_'\n", envPrefix)
			}
			printViperContents(v)

		} else {

			v := viper.GetViper()
			printViperContents(v)

		}
	},
}

func printViperContents(v *viper.Viper) {
	keys := v.AllKeys()

	l := 0
	for i := 0; i < len(keys); i++ {
		if l < len(keys[i]) {
			l = len(keys[i])
		}
	}

	sort.Strings(keys)
	format := fmt.Sprintf("%%-%ds : %%v\n", l)
	for _, key := range keys {
		val := v.Get(key)
		if val == nil {
			continue
		}
		fmt.Printf(format, key, val)
	}
}

const (
	fileOnlyFlag string = "file-only"
	envOnlyFlag  string = "env-only"
	envPrefix    string = "RMQ"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
	configCmd.AddCommand(configCreateCmd)
	configCmd.AddCommand(configShowCmd)
	configShowCmd.PersistentFlags().Bool(fileOnlyFlag, false, "Only print the values defined in the configuration file")
	configShowCmd.PersistentFlags().Bool(envOnlyFlag, false, "Only print the values defined in environment variables")
}
