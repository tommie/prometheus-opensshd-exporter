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

// preauthEndRE matches the last line of a session that never
// authenticated. With password authentication disabled, OpenSSH logs no
// "Failed" line for these. The username can be empty.
var preauthEndRE = regexp.MustCompile(`^(?:Connection closed by|Disconnected from) (authenticating|invalid) user (.*?) (\S+) port (\S+) \[preauth\]$`)

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

	if m := preauthEndRE.FindStringSubmatch(s); m != nil {
		// The log doesn't say which method was tried.
		return &AuthResult{
			AuthMsg:  "closed",
			Method:   "unknown",
			IsValid:  m[1] == "authenticating",
			Username: m[2],
			Addr:     "tcp:" + m[3] + ":" + m[4],
		}, nil
	}

	return nil, ErrUnknownLine
}
