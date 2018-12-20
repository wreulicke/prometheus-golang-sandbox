package main

import (
	"time"
	"log"
	"net/http"
	"fmt"
	"math/rand"
	
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", prometheus.InstrumentHandlerFuncWithOpts(
		prometheus.SummaryOpts{
			Subsystem:   "http",
			ConstLabels: prometheus.Labels{"handler": "hello_world"},
			Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001},
		},
		hello))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hello(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	time.Sleep(time.Duration(r.Int63n(100)) * time.Millisecond)
	fmt.Fprintln(w, "Hello World")
}
