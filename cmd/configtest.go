package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"
	promClient "github.com/travelaudience/go-promhttp"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var configTestCmd = &cobra.Command{
	Use:   "configTest",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		promClient := &promClient.Client{
			Client: &http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureTLS}}, Timeout: 2 * time.Second},
			Registerer: prometheus.DefaultRegisterer,
		}
		httpClient, _ := promClient.ForRecipient("kibanaApi")
		kibanaClient := kibana_api.NewKibanaClient(kibanaUrl, kibanaAuthToken, *httpClient)
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
