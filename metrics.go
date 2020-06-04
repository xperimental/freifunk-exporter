package main

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xperimental/freifunk-exporter/info"
)

var (
	prefix = "freifunk_"

	collectorUpDesc = prometheus.NewDesc(
		prefix+"_collector_up",
		"Is set to 1 when the collector is able to get information.",
		[]string{}, nil)
	collectorTimestampDesc = prometheus.NewDesc(
		prefix+"info_timestamp",
		"Contains the timestamp of the currently cached information.",
		[]string{}, nil)
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
	sourceURL     string
	cacheDuration time.Duration
	lastInfo      *info.Meshinfo
	timestamp     time.Time
}

func newCollector(sourceURL string, cacheDuration time.Duration) *collector {
	return &collector{
		sourceURL:     sourceURL,
		cacheDuration: cacheDuration,
		lastInfo:      nil,
		timestamp:     time.Unix(0, 0),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- clientCountDesc
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	now := time.Now()
	if now.Sub(c.timestamp) > c.cacheDuration {
		if err := c.updateCache(now); err != nil {
			log.Printf("Error updating data: %s", err)
		}
	}

	c.sendCollectorMetrics(ch)
	if c.lastInfo == nil {
		return
	}

	c.updateNodes(ch, c.lastInfo.Nodes)
}

func (c *collector) updateCache(now time.Time) error {
	info, err := info.GetInfo(c.sourceURL)
	if err != nil {
		return err
	}

	c.timestamp = now
	c.lastInfo = info
	return nil
}

func (c *collector) sendCollectorMetrics(ch chan<- prometheus.Metric) {
	collectorUp := 0.0
	if c.lastInfo != nil {
		collectorUp = 1
	}

	sendMetric(ch, collectorUpDesc, collectorUp, []string{})
	sendMetric(ch, collectorTimestampDesc, float64(c.timestamp.Unix()), []string{})
}

func (c *collector) updateNodes(ch chan<- prometheus.Metric, nodes []info.Node) {
	for _, node := range nodes {
		metaLabels := []string{node.ID, node.Hostname, node.Model, node.Firmware.Release, node.SiteCode}
		sendMetric(ch, metaDesc, 1.0, metaLabels)

		idLabel := []string{node.ID}
		sendMetric(ch, clientCountDesc, float64(node.Clients), idLabel)

		sendMetric(ch, loadAvgDesc, node.LoadAvg, idLabel)
		sendMetric(ch, rootFsUsageDesc, node.RootfsUsage, idLabel)

		if node.MemoryUsage > 0 {
			sendMetric(ch, memoryUsageDesc, node.MemoryUsage, idLabel)
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
