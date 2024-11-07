package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/prometheus/client_golang/prometheus"
	es "github.com/schmiddim/kibana-alert-exporter/elasticsearch"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"
	log "github.com/sirupsen/logrus"
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
		fmt.Println("Test Kibana API")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))

		if !queryEsForAlerts {
			fmt.Println("query-alerts-in-es flag not set - skip elasticsearch test")
		} else {
			cfg := elasticsearch.Config{
				Addresses: []string{
					elasticSearchUrl,
				},
				Username: elasticSearchUsername,
				Password: elasticSearchPassword,
			}

			fmt.Println("Test ElasticSearch API")
			esClient, err := elasticsearch.NewClient(cfg)

			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			aa := es.NewActiveAlerts(esClient)

			resp, err := aa.Info()
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			fmt.Println(resp)

		}
	},
}

func init() {
	rootCmd.AddCommand(configTestCmd)

}
