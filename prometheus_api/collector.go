package prometheus_api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/schmiddim/kibana-alert-exporter/helper"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"
)

type HealthWrapper struct {
	alertRule          kibana_api.AlertRule
	descNewAlert       *prometheus.Desc
	descActiveAlert    *prometheus.Desc
	descIgnoredAlert   *prometheus.Desc
	descRecoveredAlert *prometheus.Desc
}

type KibanaCollector struct {
	kClient     kibana_api.KclientInterface
	versionInfo *prometheus.Desc
}

func NewKibanaCollector(kclient kibana_api.KclientInterface) *KibanaCollector {
	return &KibanaCollector{
		kClient:     kclient,
		versionInfo: prometheus.NewDesc("exporter_info", "Build Information about the Exporter", []string{"code_version"}, nil),
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
	ch <- collector.versionInfo

}
func (collector *KibanaCollector) Collect(ch chan<- prometheus.Metric) {
	hws := collector.getHealthWrappers()
	for _, h := range hws {
		m1, err := prometheus.NewConstMetric(h.descNewAlert, prometheus.GaugeValue, h.alertRule.LastRun.AlertsCount.Active, h.alertRule.LabelValues...)
		if err != nil {
			log.Fatal().Err(err)
		}

		m2, err := prometheus.NewConstMetric(h.descActiveAlert, prometheus.GaugeValue, h.alertRule.LastRun.AlertsCount.Active, h.alertRule.LabelValues...)
		if err != nil {
			log.Fatal().Err(err)
		}

		m3, err := prometheus.NewConstMetric(h.descIgnoredAlert, prometheus.GaugeValue, h.alertRule.LastRun.AlertsCount.Active, h.alertRule.LabelValues...)
		if err != nil {
			log.Fatal().Err(err)
		}
		m4, err := prometheus.NewConstMetric(h.descRecoveredAlert, prometheus.GaugeValue, h.alertRule.LastRun.AlertsCount.Active, h.alertRule.LabelValues...)
		if err != nil {
			log.Fatal().Err(err)
		}

		ch <- m1
		ch <- m2
		ch <- m3
		ch <- m4

	}
	m5, err := prometheus.NewConstMetric(collector.versionInfo, prometheus.GaugeValue, 1, helper.GitCommit)
	if err != nil {
		log.Fatal().Err(err)
	}
	ch <- m5
}
