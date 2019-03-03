package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xperimental/freifunk-exporter/info"
)

var addr string
var nodesURL string
var interval time.Duration

func init() {
	flag.StringVar(&addr, "addr", ":9295", "Address to listen on.")
	flag.StringVar(&nodesURL, "source", "", "URL of nodes.json file.")
	flag.DurationVar(&interval, "interval", time.Minute*3, "Interval to use for getting updates.")
}

func main() {
	flag.Parse()

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
