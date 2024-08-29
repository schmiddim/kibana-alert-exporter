package cmd

import (
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kibana-alert-exporter",
	Short: "A brief description of your application",
	Long:  ``,
}
var kibanaUrl string
var kibanaAuthToken string
var insecureTLS = false

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().BoolVarP(&insecureTLS, "insecure", "k", false, "skip verification of tls certificates")
	viper.SetDefault("KIBANA_URL", "http://localhost:5601")
	viper.SetDefault("KIBANA_AUTH_TOKEN", "tooSecret")
	viper.AutomaticEnv()
	kibanaUrl = viper.GetString("KIBANA_URL")
	kibanaAuthToken = viper.GetString("KIBANA_AUTH_TOKEN")

}
