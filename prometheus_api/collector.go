package prometheus_api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"
	"log"
	"time"
)

type HealthWrapper struct {
	alertRule          kibana_api.AlertRule
	descNewAlert       *prometheus.Desc
	descActiveAlert    *prometheus.Desc
	descIgnoredAlert   *prometheus.Desc
	descRecoveredAlert *prometheus.Desc
}

type KibanaCollector struct {
	kClient kibana_api.KclientInterface
}

func NewKibanaCollector(kclient kibana_api.KclientInterface) *KibanaCollector {
	return &KibanaCollector{
		kClient: kclient,
	}
}
func (collector *KibanaCollector) getHealthWrappers() []HealthWrapper {
	var hws []HealthWrapper

	for _, rule := range collector.kClient.GetRules() {

		hw := HealthWrapper{
			alertRule: *rule,
		}
		hw.descNewAlert = prometheus.NewDesc("new_alerts",
			"New Alerts in Kibana",
			rule.LabelNames, nil,
		)
		hw.descActiveAlert = prometheus.NewDesc("active_alerts",
			"Active Alerts in Kibana",
			rule.LabelNames, nil,
		)

		hw.descIgnoredAlert = prometheus.NewDesc("ignored_alerts",
			"Ignored Alerts in Kibana",
			rule.LabelNames, nil,
		)
		hw.descRecoveredAlert = prometheus.NewDesc("recovered_alerts",
			"Recovered Alerts in Kibana",
			rule.LabelNames, nil,
		)
		hws = append(hws, hw)
	}

	return hws
}
func (collector *KibanaCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, hw := range collector.getHealthWrappers() {
		ch <- hw.descNewAlert
		ch <- hw.descActiveAlert
		ch <- hw.descIgnoredAlert
		ch <- hw.descRecoveredAlert

	}

}
func (collector *KibanaCollector) Collect(ch chan<- prometheus.Metric) {
	hws := collector.getHealthWrappers()
	for _, h := range hws {
		m1, err := prometheus.NewConstMetricWithCreatedTimestamp(h.descNewAlert, prometheus.CounterValue, h.alertRule.LastRun.AlertsCount.Active, time.Now(), h.alertRule.LabelValues...)
		if err != nil {
			log.Fatal(err)
		}

		m2, err := prometheus.NewConstMetricWithCreatedTimestamp(h.descActiveAlert, prometheus.CounterValue, h.alertRule.LastRun.AlertsCount.Active, time.Now(), h.alertRule.LabelValues...)
		if err != nil {
			log.Fatal(err)
		}

		m3, err := prometheus.NewConstMetricWithCreatedTimestamp(h.descIgnoredAlert, prometheus.CounterValue, h.alertRule.LastRun.AlertsCount.Active, time.Now(), h.alertRule.LabelValues...)
		if err != nil {
			log.Fatal(err)
		}
		m4, err := prometheus.NewConstMetricWithCreatedTimestamp(h.descRecoveredAlert, prometheus.CounterValue, h.alertRule.LastRun.AlertsCount.Active, time.Now(), h.alertRule.LabelValues...)
		ch <- m1
		ch <- m2
		ch <- m3
		ch <- m4

	}
}
