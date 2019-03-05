package main

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xperimental/freifunk-exporter/info"
)

var (
	prefix = "freifunk_"

	metaDesc = prometheus.NewDesc(
		prefix+"router_meta",
		"Contains labels with metadata about a router. Value is fixed to 1.",
		[]string{"id", "name", "hardware", "firmware", "community"}, nil)
	clientCountDesc = prometheus.NewDesc(
		prefix+"router_client_count_total",
		"Number of connected clients",
		[]string{"id"}, nil)
	loadAvgDesc = prometheus.NewDesc(
		prefix+"router_load_avg_5m",
		"Contains the five minutes average load for a router.",
		[]string{"id"}, nil)
	memoryUsageDesc = prometheus.NewDesc(
		prefix+"router_memory_usage_total",
		"Router memory usage as a fraction of the total.",
		[]string{"id"}, nil)
	rootFsUsageDesc = prometheus.NewDesc(
		prefix+"router_rootfs_usage_total",
		"Router root filesystem usage as a fraction of the total.",
		[]string{"id"}, nil)
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

		metaLabels := []string{info.NodeID, info.Hostname, info.Hardware.Model, info.Software.Firmware.Release, info.System.SiteCode}
		sendMetric(ch, metaDesc, 1.0, metaLabels)

		idLabel := []string{info.NodeID}
		sendMetric(ch, clientCountDesc, float64(stats.Clients), idLabel)

		sendMetric(ch, loadAvgDesc, stats.LoadAverage, idLabel)
		sendMetric(ch, rootFsUsageDesc, stats.RootFsUsage, idLabel)

		if stats.MemoryUsage != nil {
			sendMetric(ch, memoryUsageDesc, *stats.MemoryUsage, idLabel)
		}
	}
}

func sendMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, value float64, labels []string) {
	m, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, value, labels...)
	if err != nil {
		log.Printf("Error creating metric %q: %s", desc, err)
		return
	}

	ch <- m
}
