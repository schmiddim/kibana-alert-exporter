package cmd

import (
	"crypto/tls"
	"fmt"
	prometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"
	"github.com/schmiddim/kibana-alert-exporter/prometheus_api"
	"github.com/spf13/cobra"
	promClient "github.com/travelaudience/go-promhttp"
	"log"
	"net/http"
	"time"
)

var port = 9101

var startTime = time.Now()
var waitReadinessTime = 5 * time.Second

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "start the exporter",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		promClient := &promClient.Client{
			Client: &http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureTLS}}, Timeout: 2 * time.Second},
			Registerer: prometheus.DefaultRegisterer,
		}
		httpClient, _ := promClient.ForRecipient("kibanaApi")

		kibanaClient := kibana_api.NewKibanaClient(kibanaUrl, kibanaAuthToken, *httpClient)

		collector := prometheus_api.NewKibanaCollector(kibanaClient)
		prometheus.MustRegister(collector)

		fmt.Println(fmt.Sprintf("http://localhost:%d/metrics", port))
		http.Handle("/metrics", promhttp.Handler())

		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			if time.Since(startTime) > waitReadinessTime {
				w.WriteHeader(200)
			} else {
				w.WriteHeader(503)
			}
		})
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

	},
}

func init() {
	runCmd.PersistentFlags().IntVarP(&port, "port", "p", 9101, "port to use")
	rootCmd.AddCommand(runCmd)
}
