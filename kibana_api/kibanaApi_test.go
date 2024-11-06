package kibana_api

import (
	"crypto/tls"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	promClient "github.com/travelaudience/go-promhttp"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
func TestKibanaMuteAll(t *testing.T) {
	rCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fileName := "../fixtures/example-response-muted.json"
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
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, Timeout: 2 * time.Second},
		Registerer: prometheus.DefaultRegisterer,
	}
	httpClient, _ := pClient.ForRecipient("kibanaApi")
	apm := NewKibanaClient(server.URL, "SuperSecret", *httpClient)
	kclient := apm.(*Kclient) // Type assert to *Kclient
	kclient.client = server.Client()

	rules, _ := apm.GetRules()
	got := rules[0].MuteAll

	if got != true {
		t.Errorf("got %t, want %t", got, true)
	}
}

func TestKibanaResponse(t *testing.T) {
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
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, Timeout: 2 * time.Second},
		Registerer: prometheus.DefaultRegisterer,
	}
	httpClient, _ := pClient.ForRecipient("kibanaApi")
	apm := NewKibanaClient(server.URL, "SuperSecret", *httpClient)
	kclient := apm.(*Kclient) // Type assert to *Kclient
	kclient.client = server.Client()

	rules, _ := apm.GetRules()
	want := 1
	got := len(rules)
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	wantName := "security_rule"
	gotName := rules[0].Name
	if got != want {
		t.Errorf("got %s, want %s", wantName, gotName)
	}

	wantAlertCounts := []float64{1, 2, 3, 4}
	gotCount := rules[0].LastRun.AlertsCount.New

	if int(gotCount) != int(wantAlertCounts[0]) {
		t.Errorf("got %d, want %d", int(gotCount), int(wantAlertCounts[0]))
	}
	gotCount = rules[0].LastRun.AlertsCount.Active
	if int(gotCount) != int(wantAlertCounts[1]) {
		t.Errorf("got %d, want %d", int(gotCount), int(wantAlertCounts[1]))
	}
	gotCount = rules[0].LastRun.AlertsCount.Ignored
	if int(gotCount) != int(wantAlertCounts[2]) {
		t.Errorf("got %d, want %d", int(gotCount), int(wantAlertCounts[2]))
	}
	gotCount = rules[0].LastRun.AlertsCount.Recovered
	if int(gotCount) != int(wantAlertCounts[3]) {
		t.Errorf("got %d, want %d", int(gotCount), int(wantAlertCounts[3]))
	}
}

func TestJsonResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/api/alerting/rules/_find?per_page=100&sort_field=created_at")
		// Send response to be tested
		_, err := rw.Write([]byte(`{"oki":"doki"}`))
		if err != nil {
			return
		}
	}))
	// Close the server when test finishes
	defer server.Close()
	pClient := &promClient.Client{
		Client: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, Timeout: 2 * time.Second},
		Registerer: prometheus.DefaultRegisterer,
	}
	httpClient, _ := pClient.ForRecipient("kibanaApi")
	apm := NewKibanaClient(server.URL, "SuperSecret", *httpClient)
	kclient := apm.(*Kclient) // Type assert to *Kclient
	kclient.client = server.Client()
	kclient.GetRules()

}
func TestHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		headerTableTest := []struct {
			field string
			want  string
		}{
			{field: "Content-Type", want: "application/json;charset=UTF-8"},
			{field: "Authorization", want: "ApiKey SuperSecret"},
		}
		for _, test := range headerTableTest {
			got := req.Header.Get(test.field)
			if test.want != got {
				t.Errorf("got %s, want %s", got, test.want)
			}
		}
		_, err := rw.Write([]byte(`{"oki":"doki"}`))
		if err != nil {
			return
		}
	}))
	// Close the server when test finishes
	defer server.Close()
	pClient := &promClient.Client{
		Client: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false}}, Timeout: 2 * time.Second},
		Registerer: prometheus.DefaultRegisterer,
	}
	httpClient, _ := pClient.ForRecipient("kibanaApi")

	apm := NewKibanaClient(server.URL, "SuperSecret", *httpClient)

	kclient := apm.(*Kclient) // Type assert to *Kclient
	kclient.client = server.Client()

	kclient.GetRules()

}
