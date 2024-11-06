package kibana_api

import (
	"testing"
)

func TestNoUnknowns(t *testing.T) {

	labelsToExport := []string{"severity", "owner"}

	r := AlertRule{Tags: []string{"severity=page", "owner=fooSqad", "logMetrics", "foo=bar"}}
	names, values := r.GetLabels(labelsToExport)
	if len(names) != len(labelsToExport)+4 {
		t.Errorf("not enough Labels!")
	}
	ctr := 0
	for _, v := range values {
		if v == "" {
			ctr += 1
		}
	}

	if ctr-3 != 0 {

		t.Errorf("expected 0 got %d", ctr-3)

	}
}

func TestUnknown(t *testing.T) {
	labelsToExport := []string{"notfound", "owner"}

	r := AlertRule{Tags: []string{"severity=page", "owner=fooSqad", "logMetrics", "foo=bar"}}
	_, values := r.GetLabels(labelsToExport)

	ctr := 0
	for _, v := range values {
		if v == "" {
			ctr += 1
		}
	}

	if ctr-3 != 1 { // there are three labels that always be set
		t.Errorf("got %d want %d", ctr, 1)
	}
}
