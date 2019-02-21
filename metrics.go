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
	linksCount prometheus.Gauge
}

func newCollector(reader func() (*info.Nodes, error)) *collector {
	return &collector{
		readerFunc: reader,
		linksCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "link_count_total",
			Help: "Number of links between nodes",
		}),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	c.linksCount.Describe(ch)
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	nodes, err := c.readerFunc()
	if err != nil {
		log.Printf("Error collecting node information: %s", err)
		return
	}

	c.linksCount.Set(float64(len(nodes.Links)))
	c.linksCount.Collect(ch)

	c.updateNodes(ch, nodes)
}

func (c *collector) updateNodes(ch chan<- prometheus.Metric, nodes *info.Nodes) {
	for _, node := range nodes.Nodes {
		values := []string{node.ID, node.Name, node.Hardware, node.Firmware, node.Community}

		m, err := prometheus.NewConstMetric(clientCountDesc, prometheus.GaugeValue, float64(node.Clients), values...)
		if err != nil {
			log.Printf("Error creating metric: %s", err)
			continue
		}

		ch <- m
	}

}
