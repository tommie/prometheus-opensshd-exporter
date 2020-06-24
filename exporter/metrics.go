package exporter

import (
	"io"
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type Log interface {
	Read() (string, error)
}

func RunMetricUpdates(l Log) error {
	return runLogLines(l, handleLogLine)
}

func runLogLines(l Log, handle func(string) error) error {
	for {
		line, err := l.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		if err := handle(line); err != nil {
			log.Printf("Could not handle log record (ignored: %v): %#v", err, line)
		}
	}
}

func handleLogLine(line string) error {
	line = strings.TrimSpace(line)
	rec, err := parseLogLine(line)
	if err == ErrUnknownLine {
		return nil
	} else if err != nil {
		return err
	}

	return updateMetrics(rec)
}

func updateMetrics(rec interface{}) error {
	switch r := rec.(type) {
	case *AuthResult:
		validUser := "0"
		var user string
		if r.IsValid {
			validUser = "1"
			user = r.Username
		}
		authResults.With(prometheus.Labels{
			"method":     r.Method,
			"result":     r.AuthMsg,
			"user":       user,
			"valid_user": validUser,
		}).Inc()
	}

	return nil
}

var (
	authResults = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "opensshd_auth_results_total",
		Help: "OpenSSHd authentication results",
	}, []string{"method", "result", "user", "valid_user"})
)

func init() {
	prometheus.MustRegister(authResults)
}
