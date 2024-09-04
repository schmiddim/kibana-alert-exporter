package prometheus_api

import (
	"crypto/tls"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"
	promClient "github.com/travelaudience/go-promhttp"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestKibanaCollector(t *testing.T) {
	rCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fileName := "../fixtures/example-response.json"
		if rCount > 0 {
			fileName = "../fixtures/example-response-empty.json"
		}
		b, err := os.ReadFile(fileName) // just pass the file name
		if err != nil {
			log.Fatal("text fixture not found")
		}

		rCount += 1
		_, err = rw.Write(b)
		if err != nil {
			return
		}
	}))

	defer server.Close()
	pClient := &promClient.Client{
		Client: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false}}, Timeout: 2 * time.Second},
		Registerer: prometheus.DefaultRegisterer,
	}
	httpClient, _ := pClient.ForRecipient("kibanaApi")
	client := kibana_api.NewKibanaClient(server.URL, "SuperSecret", *httpClient)
	collector := NewKibanaCollector(client)

	ch := make(chan prometheus.Metric)
	go func() {
		collector.Collect(ch)
		close(ch)
	}()

	var metrics []prometheus.Metric
	for metric := range ch {
		metrics = append(metrics, metric)
	}

	if len(metrics) != 5 {
		t.Errorf("expected 5 metrics, got %d", len(metrics))
	}

}
