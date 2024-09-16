package prometheus_api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schmiddim/kibana-alert-exporter/helper"
	"github.com/schmiddim/kibana-alert-exporter/kibana_api"
	log "github.com/sirupsen/logrus"
)

type HealthWrapper struct {
	alertRule          kibana_api.AlertRule
	descNewAlert       *prometheus.Desc
	descActiveAlert    *prometheus.Desc
	descIgnoredAlert   *prometheus.Desc
	descRecoveredAlert *prometheus.Desc
}

type KibanaCollector struct {
	kClient        kibana_api.KclientInterface
	versionInfo    *prometheus.Desc
	labelsToExport []string
}

func NewKibanaCollector(kclient kibana_api.KclientInterface, labelsToExport []string) *KibanaCollector {
	return &KibanaCollector{
		kClient:        kclient,
		versionInfo:    prometheus.NewDesc("exporter_info", "Build Information about the Exporter", []string{"code_version"}, nil),
		labelsToExport: labelsToExport,
	}
}
func (collector *KibanaCollector) getHealthWrappers() []HealthWrapper {
	var hws []HealthWrapper

	rules, _ := collector.kClient.GetRules()
	for _, rule := range rules {

		hw := HealthWrapper{
			alertRule: *rule,
		}
		labelNames, _ := rule.GetLabels(collector.labelsToExport)
		hw.descNewAlert = prometheus.NewDesc("kibana_new_alerts",
			"New Alerts in Kibana",
			labelNames, nil,
		)
		hw.descActiveAlert = prometheus.NewDesc("kibana_active_alerts",
			"Active Alerts in Kibana",
			labelNames, nil,
		)

		hw.descIgnoredAlert = prometheus.NewDesc("kibana_ignored_alerts",
			"Ignored Alerts in Kibana",
			labelNames, nil,
		)
		hw.descRecoveredAlert = prometheus.NewDesc("kibana_recovered_alerts",
			"Recovered Alerts in Kibana",
			labelNames, nil,
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
		_, labelValues := h.alertRule.GetLabels(collector.labelsToExport)
		m1, err := prometheus.NewConstMetric(h.descNewAlert, prometheus.GaugeValue, h.alertRule.LastRun.AlertsCount.Active, labelValues...)
		if err != nil {
			log.Fatal(err)
		}

		m2, err := prometheus.NewConstMetric(h.descActiveAlert, prometheus.GaugeValue, h.alertRule.LastRun.AlertsCount.Active, labelValues...)
		if err != nil {
			log.Fatal(err)
		}

		m3, err := prometheus.NewConstMetric(h.descIgnoredAlert, prometheus.GaugeValue, h.alertRule.LastRun.AlertsCount.Active, labelValues...)
		if err != nil {
			log.Fatal(err)
		}
		m4, err := prometheus.NewConstMetric(h.descRecoveredAlert, prometheus.GaugeValue, h.alertRule.LastRun.AlertsCount.Active, labelValues...)
		if err != nil {
			log.Fatal(err)
		}

		ch <- m1
		ch <- m2
		ch <- m3
		ch <- m4

	}
	m5, err := prometheus.NewConstMetric(collector.versionInfo, prometheus.GaugeValue, 1, helper.GitCommit)
	if err != nil {
		log.Fatal(err)
	}
	ch <- m5
}
