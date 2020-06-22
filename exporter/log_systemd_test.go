package exporter

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/coreos/go-systemd/v22/sdjournal"
)

func TestSystemdLogRead(t *testing.T) {
	oldNew := newSystemdJournalReader
	defer func() {
		newSystemdJournalReader = oldNew
	}()
	newSystemdJournalReader = func(config sdjournal.JournalReaderConfig) (systemdJournalReader, error) {
		return &fakeSystemdJournalReader{Lines: []string{"abc"}}, nil
	}

	l, err := NewSystemdLog()
	if err != nil {
		t.Fatalf("NewSystemdLog failed: %v", err)
	}
	defer func() {
		if err := l.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	got, err := l.Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if want := "abc\n"; got != want {
		t.Errorf("Read: got %q, want %q", got, want)
	}
}

type fakeSystemdJournalReader struct {
	Lines []string
}

func (r *fakeSystemdJournalReader) Close() error {
	return nil
}

func (r *fakeSystemdJournalReader) Follow(until <-chan time.Time, w io.Writer) error {
	for _, line := range r.Lines {
		select {
		case <-until:
			return sdjournal.ErrExpired
		default:
			// continue
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}

	return nil
}
