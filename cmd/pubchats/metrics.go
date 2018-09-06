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
	}, []string{"chat", "source"})
)

func init() {
	prometheus.MustRegister(messagesCounter)
}

func startMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
