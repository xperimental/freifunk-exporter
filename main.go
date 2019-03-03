package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/pflag"
	"github.com/xperimental/freifunk-exporter/info"
)

var addr string
var nodesURL string
var interval time.Duration

func init() {
	pflag.StringVar(&addr, "addr", ":9295", "Address to listen on.")
	pflag.StringVar(&nodesURL, "source", "", "URL of nodes.json file.")
	pflag.DurationVar(&interval, "interval", time.Minute*3, "Interval to use for getting updates.")
}

func main() {
	pflag.Parse()

	if len(nodesURL) == 0 {
		log.Println("Need to provide source URL.")
		return
	}

	infoReader := func() (*info.Nodes, error) {
		return info.GetNodes(nodesURL)
	}
	collector := newCollector(infoReader)
	prometheus.MustRegister(collector)

	http.Handle("/", http.RedirectHandler("/metrics", http.StatusFound))
	http.Handle("/metrics", prometheus.Handler())

	log.Printf("Listening on %s...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
