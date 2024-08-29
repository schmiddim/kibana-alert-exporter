package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"

	"github.com/spf13/cobra"
)

var configTestCmd = &cobra.Command{
	Use:   "configTest",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		kibanaClient := kibana_api.NewKibanaClient(kibanaUrl, kibanaAuthToken, insecureTLS)
		result := kibanaClient.GetAlertingHealth()
		out, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
	},
}

func init() {
	rootCmd.AddCommand(configTestCmd)

}
