package kibana_api

import (
	"time"
)

type alertRulesFindResponse struct {
	Page       int          `json:"page"`
	Total      int          `json:"total"`
	PerPage    int          `json:"per_page"`
	AlertRules []*AlertRule `json:"data"`
}

type AlertRule struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Tags        []string `json:"tags"`
	Enabled     bool     `json:"enabled"`
	Running     bool     `json:"running"`
	LabelValues []string
	LabelNames  []string
	LastRun     struct {
		Outcome     string `json:"outcome"`
		AlertsCount struct {
			New       float64 `json:"new"`
			Active    float64 `json:"active"`
			Ignored   float64 `json:"ignored"`
			Recovered float64 `json:"recovered"`
		} `json:"alerts_count"`
	} `json:"last_run"`

	Params struct {
		Description string `json:"description"`
	} `json:"params"`
}

func (r *AlertRule) ParseLabels() {
	var names []string
	var values []string

	values = append(values, r.Id)
	values = append(values, r.Name)
	values = append(values, r.LastRun.Outcome)

	names = append(names, "id")
	names = append(names, "name")
	names = append(names, "last_run_outcome")

	//	for _, t := range r.Tags {
	//		splits := strings.Split(t, "=")
	//		if len(splits) == 2 {
	//			names = append(names, splits[0])
	//			values = append(values, splits[1])
	//		}
	//	}
	r.LabelNames = names
	r.LabelValues = values

}

type AlertingHealthResponse struct {
	FrameWorkHealth struct {
		ReadHealth struct {
			Status    string    `json:"status"`
			TimeStamp time.Time `json:"timestamp"`
		} `json:"read_health"`
		ExecutionHealth struct {
			Status    string    `json:"status"`
			TimeStamp time.Time `json:"timestamp"`
		} `json:"execution_health"`
		DecryptionHealth struct {
			Status    string    `json:"status"`
			TimeStamp time.Time `json:"timestamp"`
		} `json:"decryption_health"`
	} `json:"alerting_framework_health"`
	HasPermanentEncryptionKey bool `json:"has_permanent_encryption_key"`
	IsSufficientlySecure      bool `json:"is_sufficiently_secure"`
}
