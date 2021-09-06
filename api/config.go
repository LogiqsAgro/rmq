package api

import (
	"net/url"

	"github.com/spf13/cobra"
)

var (
	// Config contains the default global config settings
	Config *cfg = new(cfg)
	// Page contains the default global pagination settings
	Page *page = newPage()

	defaults = map[string]interface{}{
		"scheme":       "http",
		"host":         "localhost",
		"port":         15672,
		"user":         "guest",
		"password":     "guest",
		"vhost":        "/",
		"debug":        false,
		"pretty-print": false,
	}
)

// AddConfigFlags adds commandline parameters to the command for the RabbitMQ api
func AddConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&Config.Scheme, "scheme", defaults["scheme"].(string), "RabbitMQ api protocol")
	cmd.PersistentFlags().StringVar(&Config.Host, "host", defaults["host"].(string), "RabbitMQ machine name")
	cmd.PersistentFlags().IntVar(&Config.Port, "port", defaults["port"].(int), "RabbitMQ management api port")
	cmd.PersistentFlags().StringVar(&Config.User, "user", defaults["user"].(string), "RabbitMQ user name")
	cmd.PersistentFlags().StringVar(&Config.Password, "password", defaults["password"].(string), "RabbitMQ user name")
	cmd.PersistentFlags().StringVar(&Config.VHost, "vhost", defaults["vhost"].(string), "RabbitMQ virtual host")
	cmd.PersistentFlags().BoolVar(&Config.Debug, "debug", defaults["debug"].(bool), "Enable http request and response details logging")
	cmd.PersistentFlags().BoolVar(&Config.IndentJson, "pretty-print", defaults["pretty-print"].(bool), "Enable formatting of the json responses")
}

// AddPagingFlags adds commandline parameters to the command for paging parameters in the RabbitMQ api
func AddPagingFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVarP(&Page.Page, "page", "p", 1, "The results page number (one-based)")
	cmd.PersistentFlags().IntVarP(&Page.Size, "page-size", "s", 100, "The results page size")
	cmd.PersistentFlags().StringVarP(&Page.Name, "name", "n", "", "The name to filter for")
	cmd.PersistentFlags().BoolVarP(&Page.UseRegex, "regex", "r", false, "Enables regular expressions for the --name filter")
}

type (
	cfg struct {
		Scheme     string
		Host       string
		Port       int
		VHost      string
		User       string
		Password   string
		Debug      bool
		IndentJson bool
	}
)

func (cfg *cfg) escapedVHost() string {
	return url.PathEscape(cfg.VHost)
}
