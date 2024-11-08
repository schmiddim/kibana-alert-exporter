package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	log "github.com/sirupsen/logrus"
	"io"
)

type ActiveAlert struct {
	InstanceId string
	RuleUUID   string
}

type EsWrapper struct {
	esClient *elasticsearch.Client
}

func NewActiveAlerts(esClient *elasticsearch.Client) *EsWrapper {

	return &EsWrapper{
		esClient: esClient,
	}
}

func (a *EsWrapper) GetInstanceIdsForActiveAlerts() []ActiveAlert {

	var buf bytes.Buffer
	res, err := a.esClient.Search(
		a.esClient.Search.WithContext(context.Background()),
		a.esClient.Search.WithIndex(".internal.alerts-*"),
		a.esClient.Search.WithBody(&buf),
		a.esClient.Search.WithTrackTotalHits(true),
		a.esClient.Search.WithPretty(),
		a.esClient.Search.WithQuery("kibana.alert.status: active"),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)
	if res.IsError() {
		log.Fatal(err)
	}
	var (
		r map[string]interface{}
	)

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	var instanceIds []ActiveAlert
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		instanceId := hit.(map[string]interface{})["_source"].(map[string]interface{})["kibana.alert.instance.id"]
		uuid := hit.(map[string]interface{})["_source"].(map[string]interface{})["kibana.alert.rule.uuid"]
		foo := ActiveAlert{
			InstanceId: fmt.Sprintf("%s", instanceId),
			RuleUUID:   fmt.Sprintf("%s", uuid),
		}
		instanceIds = append(instanceIds, foo)
	}
	return instanceIds
}

func (a *EsWrapper) Info() (*esapi.Response, error) {
	return a.esClient.Info()
}
