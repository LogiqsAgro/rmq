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
package api

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Config contains the default global config settings
	Config *cfg = new(cfg)
	// Page contains the default global pagination settings
	Page *pageFilter = new(pageFilter)

	defaults = map[string]interface{}{
		"scheme":       "http",
		"host":         "localhost",
		"api-port":     15672,
		"user":         "guest",
		"password":     "guest",
		"vhost":        "/",
		"debug":        false,
		"pretty-print": false,
		"columns":      []string{},
		"sort":         "",
		"sort-reverse": false,
	}
)

// AddConfigFlags adds commandline parameters to the command for the RabbitMQ api
func AddConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&Config.Scheme, "scheme", defaults["scheme"].(string), "RabbitMQ api protocol")
	cmd.PersistentFlags().StringVar(&Config.Host, "host", defaults["host"].(string), "RabbitMQ machine name")
	cmd.PersistentFlags().IntVar(&Config.ApiPort, "api-port", defaults["api-port"].(int), "RabbitMQ management api port")
	cmd.PersistentFlags().StringVar(&Config.User, "user", defaults["user"].(string), "RabbitMQ user name")
	cmd.PersistentFlags().StringVar(&Config.Password, "password", defaults["password"].(string), "RabbitMQ user name")
	cmd.PersistentFlags().StringVar(&Config.VHost, "vhost", defaults["vhost"].(string), "RabbitMQ virtual host")
	cmd.PersistentFlags().BoolVar(&Config.Debug, "debug", defaults["debug"].(bool), "Enable http request and response details logging")
	cmd.PersistentFlags().BoolVar(&Config.IndentJson, "pretty-print", defaults["pretty-print"].(bool), "Enable formatting of the json responses")
}

// AddListFlags adds parameters to the command that change the shape and sort order of returned data from the RabbitMQ api.
func AddListFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringArrayVar(&Config.Columns, "columns", defaults["columns"].([]string), "Fields to include in list responses, use commas to separate fields, use dots to include sub-fields like: field.subfield")
	cmd.PersistentFlags().StringVar(&Config.Sort, "sort", defaults["sort"].(string), "Field to sort list responses by, use dots to specify a sub-field like: message_stats.deliver_details.rate, You cannot specify multiple sort fields, only 1 field is supported")
	cmd.PersistentFlags().BoolVar(&Config.SortReverse, "sort-reverse", defaults["sort-reverse"].(bool), "Reverses the sort order")
}

// AddPagingFlags adds parameters to the command for paging parameters in the RabbitMQ api
func AddPagingFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVarP(&Page.Page, "page", "p", 0, "The results page number (one-based)")
	cmd.PersistentFlags().IntVarP(&Page.PageSize, "page-size", "s", 0, "The results page size")
	cmd.PersistentFlags().StringVarP(&Page.Name, "name", "n", "", "The name to filter for")
	cmd.PersistentFlags().BoolVarP(&Page.UseRegex, "regex", "r", false, "Enables regular expressions for the --name filter")
}

func ApplyConfig(b Builder) {
	Config.Apply(b)
	Page.Apply(b)
}

type (
	cfg struct {
		Scheme      string
		Host        string
		ApiPort     int
		VHost       string
		User        string
		Password    string
		Debug       bool
		IndentJson  bool
		Columns     []string
		Sort        string
		SortReverse bool
	}
)

func (cfg *cfg) Apply(b Builder) Builder {
	return b.
		BaseUrl(fmt.Sprintf("%s://%s:%d", cfg.Scheme, cfg.Host, cfg.ApiPort)).
		BasicAuth(cfg.User, cfg.Password).
		Sort(cfg.Sort, cfg.SortReverse).
		Columns(cfg.Columns...)
}
