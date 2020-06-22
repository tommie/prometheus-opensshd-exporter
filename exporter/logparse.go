package exporter

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrMalformedLine = errors.New("malformed log line")
	ErrUnknownLine   = errors.New("unknown log line")

	timestampRE = regexp.MustCompile(`^[-\d :.+]+ \S+ (.*)`)
)

var authResultRE = regexp.MustCompile(`(.+?) (\S+) for (invalid user )?(.+?) from (\S+) port (\S+) (.+?)(?:: (.+))?$`)

type AuthResult struct {
	AuthMsg  string
	Method   string
	IsValid  bool
	Username string
	Addr     string
	Protocol string
	Extra    string
}

func parseLogLine(s string) (interface{}, error) {
	m := timestampRE.FindStringSubmatch(s)
	if m == nil {
		return nil, ErrMalformedLine
	}

	s = m[1]
	if m := authResultRE.FindStringSubmatch(s); m != nil {
		return &AuthResult{
			AuthMsg:  strings.ToLower(m[1]),
			Method:   strings.ToLower(m[2]),
			IsValid:  m[3] == "",
			Username: m[4],
			Addr:     "tcp:" + m[5] + ":" + m[6],
			Protocol: m[7],
			Extra:    m[8],
		}, nil
	}

	return nil, ErrUnknownLine
}
