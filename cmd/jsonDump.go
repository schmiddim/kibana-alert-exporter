package cmd

import (
	"crypto/tls"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"
	"github.com/spf13/cobra"
	promClient "github.com/travelaudience/go-promhttp"
	"net/http"
	"time"
)

// jsonDumpCmd represents the jsonDump command
var jsonDumpCmd = &cobra.Command{
	Use:   "json-dump",
	Short: "Write Kibana Api Response to stdout",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		promClient := &promClient.Client{
			Client: &http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureTLS}}, Timeout: 2 * time.Second},
			Registerer: prometheus.DefaultRegisterer,
		}
		httpClient, _ := promClient.ForRecipient("kibanaApi")

		kibanaClient := kibana_api.NewKibanaClient(kibanaUrl,
			kibanaAuthToken,
			*httpClient)

		_, responses := kibanaClient.GetRules()
		str := " ["
		for _, response := range responses {
			str += fmt.Sprintf(string(response)) + ", "
		}
		str += "]"
		fmt.Println(str)
	},
}

func init() {
	rootCmd.AddCommand(jsonDumpCmd)

}
