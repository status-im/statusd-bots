package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	messagesCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "shh",
		Name:      "messages_total",
		Help:      "Received messages counter.",
	}, []string{"chat"})
	uniqueCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "shh",
		Name:      "unique_authors_total",
		Help:      "Unique authoers of the messages.",
	}, []string{"chat"})
)

func init() {
	prometheus.MustRegister(messagesCounter)
	prometheus.MustRegister(uniqueCounter)
}

func startMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
