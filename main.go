package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
)

var (
	addr          string
	sourceURL     string
	cacheInterval time.Duration
)

func init() {
	pflag.StringVar(&addr, "addr", ":9295", "Address to listen on.")
	pflag.StringVar(&sourceURL, "source-url", "", "URL to Meshviewer JSON file.")
	pflag.DurationVar(&cacheInterval, "cache-interval", time.Minute*3, "Interval for local caching of Meshviewer data.")
}

func main() {
	pflag.Parse()

	if len(sourceURL) == 0 {
		log.Println("Need to provide source URL.")
		return
	}

	collector := newCollector(sourceURL, cacheInterval)
	prometheus.MustRegister(collector)

	http.Handle("/", http.RedirectHandler("/metrics", http.StatusFound))
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on %s...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
