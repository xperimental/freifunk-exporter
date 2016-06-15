package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_model/go"
)

func TestUpdateLinkMetric(t *testing.T) {
	nodes := &Nodes{
		Links: []Link{
			Link{},
			Link{},
			Link{},
		},
	}

	updateMetrics(nodes)

	ch := make(chan prometheus.Metric)
	go func() {
		linksCount.Collect(ch)
		close(ch)
	}()

	for m := range ch {
		var collected io_prometheus_client.Metric
		m.Write(&collected)

		g := collected.GetGauge()
		if g == nil {
			t.Errorf("got nil gauge")
		}

		var expected float64 = 3
		if *g.Value != expected {
			t.Errorf("got %f, wanted %f", *g.Value, expected)
		}
	}
}

func TestUpdateClientsMetric(t *testing.T) {
	nodes := &Nodes{
		Nodes: []Node{
			Node{
				ID:      "one",
				Clients: 10,
			},
			Node{
				ID:      "two",
				Clients: 20,
			},
		},
	}

	updateMetrics(nodes)

	ch := make(chan prometheus.Metric)
	go func() {
		clientCount.Collect(ch)
		close(ch)
	}()

	clients := 0
	for m := range ch {
		var collected io_prometheus_client.Metric
		m.Write(&collected)

		g := collected.GetGauge()
		if g == nil {
			t.Errorf("got nil gauge")
		}

		if g.Value == nil {
			t.Errorf("got nil value")
		}

		clients += int(*g.Value)
	}

	expected := 30
	if clients != expected {
		t.Errorf("got %d, wanted %d", clients, expected)
	}
}
