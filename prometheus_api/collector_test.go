package prometheus_api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

	client := kibana_api.NewKibanaClient(server.URL, "SuperSecret", false)
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

	if len(metrics) != 4 {
		t.Errorf("expected 4 metrics, got %d", len(metrics))
	}

}
