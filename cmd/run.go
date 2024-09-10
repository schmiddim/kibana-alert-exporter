package cmd

import (
	"crypto/tls"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/schmiddim/kibana-alert-exporter/helper"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"
	"github.com/schmiddim/kibana-alert-exporter/prometheus_api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	promClient "github.com/travelaudience/go-promhttp"
	"net/http"
	"time"
)

var port = 9101
var labelsToExport []string

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

		log.Info("labels to export:", labelsToExport)
		kibanaClient := kibana_api.NewKibanaClient(kibanaUrl, kibanaAuthToken, *httpClient)

		collector := prometheus_api.NewKibanaCollector(kibanaClient, labelsToExport)
		prometheus.MustRegister(collector)

		log.Infof("http://localhost:%d/metrics", port)
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
	runCmd.PersistentFlags().StringArrayVarP(&labelsToExport, "export-labels", "l", []string{}, "add tags in the form key=val to kibana alerts to add them as labelsToExport to the metric")

	rootCmd.AddCommand(runCmd)
	helper.LoggerInit()
}
