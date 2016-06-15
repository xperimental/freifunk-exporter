package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var addr string
var nodesURL string
var interval time.Duration

func init() {
	flag.StringVar(&addr, "addr", ":8080", "Address to listen on.")
	flag.StringVar(&nodesURL, "source", "", "URL of nodes.json file.")
	flag.DurationVar(&interval, "interval", time.Minute*3, "Interval to use for getting updates.")
}

func main() {
	flag.Parse()

	if len(nodesURL) == 0 {
		log.Println("Need to provide source URL.")
		return
	}

	tick := make(chan struct{})
	gen := generator(nodesURL, tick)

	log.Println("Starting update loop.")
	go runLoop(tick, gen)

	log.Println("Trigger first update.")
	tick <- struct{}{}

	log.Printf("Listening on %s...", addr)
	http.Handle("/metrics", prometheus.Handler())
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func runLoop(tick chan<- struct{}, gen <-chan *Nodes) {
	for {
		select {
		case <-time.After(interval):
			tick <- struct{}{}
		case nodes := <-gen:
			updateMetrics(nodes)
		}
	}
}
