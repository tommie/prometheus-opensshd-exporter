package exporter

import (
	"reflect"
	"testing"
)

const ts = "2006-01-02 15:04:05.000000 -0700 MST "

func TestParseLogLine(t *testing.T) {
	tsts := []struct {
		Line    string
		Want    interface{}
		WantErr error
	}{
		{
			Line:    "",
			Want:    nil,
			WantErr: ErrMalformedLine,
		},
		{
			Line:    ts + "Disconnected from invalid user web 128.199.91.26 port 48966 [preauth]",
			Want:    nil,
			WantErr: ErrUnknownLine,
		},
		{
			Line: ts + "Failed password for root from 177.129.191.142 port 36571 ssh2",
			Want: &AuthResult{
				AuthMsg:  "failed",
				Method:   "password",
				IsValid:  true,
				Username: "root",
				Addr:     "tcp:177.129.191.142:36571",
				Protocol: "ssh2",
				Extra:    "",
			},
			WantErr: nil,
		},
		{
			Line: ts + "Failed password for invalid user admin from 35.220.210.160 port 35048 ssh2",
			Want: &AuthResult{
				AuthMsg:  "failed",
				Method:   "password",
				IsValid:  false,
				Username: "admin",
				Addr:     "tcp:35.220.210.160:35048",
				Protocol: "ssh2",
				Extra:    "",
			},
			WantErr: nil,
		},
		{
			Line: ts + "Accepted publickey for root from 131.210.128.50 port 7258 ssh2: RSA SHA256:asdlij3rnfk4kfnknf4tl949g4ittwetwerlke3423Q",
			Want: &AuthResult{
				AuthMsg:  "accepted",
				Method:   "publickey",
				IsValid:  true,
				Username: "root",
				Addr:     "tcp:131.210.128.50:7258",
				Protocol: "ssh2",
				Extra:    "RSA SHA256:asdlij3rnfk4kfnknf4tl949g4ittwetwerlke3423Q",
			},
			WantErr: nil,
		},
	}

	for _, tst := range tsts {
		t.Run(tst.Line, func(t *testing.T) {
			got, err := parseLogLine(tst.Line)
			if tst.WantErr != nil {
				if !reflect.DeepEqual(err, tst.WantErr) {
					t.Fatalf("ParseLogLine(%q): got error %v, want %v", tst.Line, err, tst.WantErr)
				}
			} else if err != nil {
				t.Fatalf("ParseLogLine(%q) failed: %v", tst.Line, err)
			}
			if !reflect.DeepEqual(tst.Want, got) {
				t.Fatalf("ParseLogLine(%q): got %+v, want %+v", tst.Line, got, tst.Want)
			}
		})
	}
}
