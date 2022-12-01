package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/awoodbeck/gnp/ch13/instrumentation/metrics"
)

var (
	metricsAddr = flag.String("metrics", "127.0.0.1:8081", "metrics listen address")
	webAddr    = flag.String("web", "127.0.0.1:8082", "web listen address")
)

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	metrics.Requests.Add(1)
	defer func(start time.Time) {
		metrics.RequestDuration.Observe(time.Since(start).Seconds())
	}(time.Now())

	_, err := w.Write([]byte("Hello!"))
	if err != nil {
		metrics.WriteErrors.Add(1)
	}
}

func newHTTPServer(addr string, mux http.Handler, stateFunc func(net.Conn, http.ConnState)) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		ConnState:         stateFunc,
	}

	go func() {
		log.Fatal(srv.Serve(l))
	}()
	return nil
}

func connStateMetrics(_ net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		metrics.OpenConnections.Add(1)
	case http.StateClosed:
		metrics.OpenConnections.Add(-1)
	}
}

func main()  {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	mux := http.NewServeMux()
	mux.Handle("/metrics/", promhttp.Handler)
	if err := newHTTPServer(*metricsAddr, mux, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Metrics listening on %q ...\n", *metricsAddr)
}