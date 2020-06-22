package exporter

import (
	"bufio"
	"io"
	"time"

	"github.com/coreos/go-systemd/v22/sdjournal"
	"golang.org/x/sync/errgroup"
)

type SystemdLog struct {
	pr     io.ReadCloser
	br     *bufio.Reader
	eg     errgroup.Group
	stopCh chan time.Time
}

func NewSystemdLog() (*SystemdLog, error) {
	pr, pw := io.Pipe()

	l := &SystemdLog{
		pr:     pr,
		br:     bufio.NewReader(pr),
		stopCh: make(chan time.Time),
	}

	jr, err := newSystemdJournalReader(sdjournal.JournalReaderConfig{
		Since:   1 * time.Nanosecond,
		Matches: []sdjournal.Match{{Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT, Value: "ssh.service"}},
	})
	if err != nil {
		pw.Close()
		pr.Close()
		return nil, err
	}
	l.eg.Go(func() error {
		defer pw.Close()
		defer jr.Close()

		err := jr.Follow(l.stopCh, pw)
		if err == sdjournal.ErrExpired {
			return nil
		}
		return err
	})

	return l, nil
}

func (l *SystemdLog) Close() error {
	close(l.stopCh)
	l.eg.Go(l.pr.Close)
	return l.eg.Wait()
}

func (l *SystemdLog) Read() (string, error) {
	return l.br.ReadString('\n')
}

type systemdJournalReader interface {
	io.Closer
	Follow(<-chan time.Time, io.Writer) error
}

var newSystemdJournalReader = func(config sdjournal.JournalReaderConfig) (systemdJournalReader, error) {
	return sdjournal.NewJournalReader(config)
}
