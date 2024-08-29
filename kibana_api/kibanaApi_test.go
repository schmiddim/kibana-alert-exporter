package kibana_api

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Testing Snippets https://github.com/benbjohnson/testing
// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestKibanaResponse(t *testing.T) {
	r_count := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fileName := "../fixtures/example-response.json"
		if r_count > 0 {
			fileName = "../fixtures/example-response-empty.json"
		}
		b, err := os.ReadFile(fileName) // just pass the file name
		if err != nil {
			log.Fatal("text fixture not found")
		}

		r_count += 1
		rw.Write(b)
	}))

	defer server.Close()

	apm := NewKibanaClient(server.URL, "SuperSecret", false)
	apm.client = server.Client()

	rules := apm.GetRules()
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
		equals(t, req.URL.String(), "/api/alerting/rules/_find?per_page=100")
		// Send response to be tested
		rw.Write([]byte(`{"oki":"doki"}`))
	}))
	// Close the server when test finishes
	defer server.Close()

	apm := NewKibanaClient(server.URL, "SuperSecret", false)
	apm.client = server.Client()

	apm.GetRules()

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
		rw.Write([]byte(`{"oki":"doki"}`))
	}))
	// Close the server when test finishes
	defer server.Close()

	apm := NewKibanaClient(server.URL, "SuperSecret", false)
	apm.client = server.Client()

	apm.GetRules()

}
