package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	BlockFetchInterval = 5 * time.Second
	latestBlockHeight  = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "latest_block_height",
		Help:        "latest bitcoind block height",
		ConstLabels: prometheus.Labels{"location": "chengdu"},
	})
)

func init() {
	prometheus.MustRegister(latestBlockHeight)
}

func main() {
	go func() {
		for {
			<-time.After(BlockFetchInterval)
			latestBlockHeight.Set(float64(rand.Intn(1000)))
		}
	}()

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":9200", nil)
}
