package main

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xperimental/freifunk-exporter/info"
)

var (
	prefix = "freifunk_"

	clientCountDesc = prometheus.NewDesc(
		prefix+"router_client_count_total",
		"Number of connected clients",
		[]string{"id", "name", "hardware", "firmware", "community"}, nil)
)

type collector struct {
	readerFunc func() (*info.Nodes, error)
}

func newCollector(reader func() (*info.Nodes, error)) *collector {
	return &collector{
		readerFunc: reader,
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- clientCountDesc
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	nodes, err := c.readerFunc()
	if err != nil {
		log.Printf("Error collecting node information: %s", err)
		return
	}

	c.updateNodes(ch, nodes)
}

func (c *collector) updateNodes(ch chan<- prometheus.Metric, nodes *info.Nodes) {
	for _, node := range nodes.List {
		info := node.Nodeinfo
		if info == nil {
			continue
		}

		stats := node.Statistics
		if stats == nil {
			continue
		}

		values := []string{info.NodeID, info.Hostname, info.Hardware.Model, info.Software.Firmware.Release, info.System.SiteCode}

		m, err := prometheus.NewConstMetric(clientCountDesc, prometheus.GaugeValue, float64(stats.Clients), values...)
		if err != nil {
			log.Printf("Error creating metric: %s", err)
			continue
		}

		ch <- m
	}

}
