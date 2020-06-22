package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tommie/prometheus-opensshd-exporter/exporter"
	"golang.org/x/sync/errgroup"
)

var (
	listenAddr = flag.String("web.listen-address", "tcp:localhost:9100", "Address to listen for HTTP requests on.")
)

func main() {
	flag.Parse()

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", http.RedirectHandler("/metrics", http.StatusMovedPermanently))

	l, err := exporter.NewSystemdLog()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	var eg errgroup.Group
	eg.Go(func() error { return httpListenAndServe(*listenAddr) })
	eg.Go(func() error { return exporter.RunMetricUpdates(l) })

	log.Fatal(eg.Wait())
}

func httpListenAndServe(addr string) error {
	ss := strings.SplitN(addr, ":", 2)

	if len(ss) != 2 {
		return fmt.Errorf("expected '<type>:<address>': %s", addr)
	}

	l, err := net.Listen(ss[0], ss[1])
	if err != nil {
		return err
	}
	defer l.Close()

	log.Printf("Listening on %s...", l.Addr())

	s := http.Server{Addr: ss[1]}
	return s.Serve(l)
}
