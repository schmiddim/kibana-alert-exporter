package cmd

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"
	"github.com/schmiddim/kibana-alert-exporter/prometheus_api"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"time"
)

var port = 9101

var startTime = time.Now()
var waitReadinessTime = time.Duration(10 * time.Second)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		kibanaClient := kibana_api.NewKibanaClient(kibanaUrl, kibanaAuthToken, insecureTLS)

		foo := prometheus_api.NewKibanaCollector(*kibanaClient)
		prometheus.MustRegister(foo)

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
