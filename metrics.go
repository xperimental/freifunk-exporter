package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xperimental/freifunk-exporter/info"
)

var (
	namespace   = "freifunk"
	clientCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "router",
		Name:      "client_count_total",
		Help:      "Number of connected clients",
	}, []string{"id", "name", "hardware", "firmware", "community"})
	linksCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "link_count_total",
		Help:      "Number of links between nodes",
	})
)

func init() {
	prometheus.MustRegister(clientCount)
	prometheus.MustRegister(linksCount)
}

func updateMetrics(nodes *info.Nodes) {
	for _, node := range nodes.Nodes {
		values := []string{node.ID, node.Name, node.Hardware, node.Firmware, node.Community}
		clientCount.WithLabelValues(values...).Set(float64(node.Clients))
	}

	linksCount.Set(float64(len(nodes.Links)))
}
