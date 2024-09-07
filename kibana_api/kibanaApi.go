package kibana_api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"

	"net/http"
	"strings"
)

type KclientInterface interface {
	GetRules() []*AlertRule
	GetAlertingHealth() AlertingHealthResponse
}

type Config struct {
	KibanaBaseUri string
	ApiKey        string
}

type Kclient struct {
	Config *Config
	client *http.Client
}

func (c *Kclient) doHttpRequest(req *http.Request) []byte {
	res, err := c.client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("Status Code is %d", res.StatusCode)
	}
	if err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal("unable to close body")
			}
		}(res.Body)
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return body
}
func (c *Kclient) newRequest(endpoint string, method string) *http.Request {
	url := c.Config.KibanaBaseUri + endpoint
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("User-Agent", "kibana-alert-exporter")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("kbn-xsrf", "true")
	req.Header.Set("Authorization", fmt.Sprintf("ApiKey %s", c.Config.ApiKey))
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func NewKibanaClient(baseUri string, apiKey string, client http.Client) KclientInterface {

	if strings.HasSuffix(baseUri, "/") == false {
		baseUri += "/"
	}
	c := Config{
		KibanaBaseUri: baseUri,
		ApiKey:        apiKey,
	}

	return &Kclient{Config: &c, client: &client}
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
func (c *Kclient) GetRules() []*AlertRule {
	resultCount := 1

	var idStrings []string
	var alertRules []*AlertRule
	page := 1
	for resultCount > 0 {
		endpoint := fmt.Sprintf("api/alerting/rules/_find?per_page=100")
		if page > 1 {
			endpoint = fmt.Sprintf("api/alerting/rules/_find?per_page=100&page=%d", page)
		}

		response := alertRulesFindResponse{}
		req := c.newRequest(endpoint, http.MethodGet)
		body := c.doHttpRequest(req)
		jsonErr := json.Unmarshal(body, &response)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		resultCount = len(response.AlertRules)
		alertRules = append(alertRules, response.AlertRules...)
		page = response.Page + 1

	}
	for _, rule := range alertRules {
		idStrings = append(idStrings, rule.Id)
	}
	trimmed := removeDuplicateStr(idStrings)
	if len(trimmed) != len(alertRules) {
		log.Errorf("trimmed %d not equal to R %d", len(trimmed), len(alertRules))
		fmt.Println("All")
		for _, l := range idStrings {
			log.Error(l)

		}
		log.Error("============")
		log.Error("Trimmed")
		for _, l := range trimmed {
			log.Error(l)

		}
	}
	return alertRules
}

func (c *Kclient) GetAlertingHealth() AlertingHealthResponse {
	// https://www.elastic.co/docs/api/doc/kibana/v8/operation/operation-getalertinghealth
	endpoint := "api/alerting/_health"

	req := c.newRequest(endpoint, http.MethodGet)
	body := c.doHttpRequest(req)

	response := AlertingHealthResponse{}
	jsonErr := json.Unmarshal(body, &response)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return response
}
