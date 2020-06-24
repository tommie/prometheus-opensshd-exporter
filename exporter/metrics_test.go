package exporter

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRunLogLines(t *testing.T) {
	want := []string{
		ts + "Failed password for root from 177.129.191.142 port 36571 ssh2",
	}
	var got []string
	err := runLogLines(&fakeLog{Lines: want}, func(line string) error {
		got = append(got, line)
		return nil
	})
	if err != nil {
		t.Fatalf("runLogLines failed: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("runLogLines: got %+v, want %+v", got, want)
	}
}

func TestHandleLogLine(t *testing.T) {
	authResults.Reset()

	err := handleLogLine(ts + "Failed password for root from 177.129.191.142 port 36571 ssh2 ")
	if err != nil {
		t.Fatalf("handleLogLine failed: %v", err)
	}

	want := []string{
		`# HELP opensshd_auth_results_total OpenSSHd authentication results`,
		`# TYPE opensshd_auth_results_total counter`,
		`opensshd_auth_results_total{method="password",result="failed",user="root",valid_user="1"} 1`,
	}
	if err := testutil.CollectAndCompare(authResults, strings.NewReader(strings.Join(want, "\n")+"\n")); err != nil {
		t.Errorf("handleLogLine metrics: %v", err)
	}
}

func TestHandleLogLine_HandlesUnknownLine(t *testing.T) {
	authResults.Reset()

	err := handleLogLine(ts + "something uninteresting")
	if err != nil {
		t.Fatalf("handleLogLine failed: %v", err)
	}

	if err := testutil.CollectAndCompare(authResults, strings.NewReader("")); err != nil {
		t.Errorf("handleLogLine metrics: %v", err)
	}
}

func TestPromLint(t *testing.T) {
	lints, err := testutil.GatherAndLint(prometheus.DefaultGatherer)
	if err != nil {
		t.Fatal(err)
	}
	for _, lint := range lints {
		t.Errorf("Prometheus lint for %q: %v", lint.Metric, lint.Text)
	}
}

type fakeLog struct {
	Lines []string
	i     int
}

func (l *fakeLog) Read() (string, error) {
	if l.i == len(l.Lines) {
		return "", io.EOF
	}
	l.i++
	return l.Lines[l.i-1], nil
}
