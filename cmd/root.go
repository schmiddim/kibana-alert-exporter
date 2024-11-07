package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "kibana-alert-exporter",
	Short: "",
	Long:  ``,
}
var kibanaUrl string
var kibanaAuthToken string
var insecureTLS = false
var queryEsForAlerts = false
var elasticSearchUsername string
var elasticSearchPassword string
var elasticSearchUrl string

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().BoolVarP(&insecureTLS, "insecure", "k", false, "skip verification of tls certificates")
	rootCmd.PersistentFlags().BoolVarP(&queryEsForAlerts, "query-alerts-in-es", "q", false, "query ElasticSearch for muted Alerts")
	viper.SetDefault("KIBANA_URL", "http://localhost:5601")
	viper.SetDefault("KIBANA_AUTH_TOKEN", "tooSecret")
	viper.SetDefault("ELASTIC_SEARCH_USERNAME", "elastic")
	viper.SetDefault("ELASTIC_SEARCH_PASSWORD", "elastic")
	viper.SetDefault("ELASTIC_SEARCH_URL", "http://localhost:9200")
	viper.AutomaticEnv()
	kibanaUrl = viper.GetString("KIBANA_URL")
	kibanaAuthToken = viper.GetString("KIBANA_AUTH_TOKEN")
	elasticSearchUsername = viper.GetString("ELASTIC_SEARCH_USERNAME")
	elasticSearchPassword = viper.GetString("ELASTIC_SEARCH_PASSWORD")
	elasticSearchUrl = viper.GetString("ELASTIC_SEARCH_URL")
}
