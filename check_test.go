package contracts

import (
	"strings"
	"testing"
)

func TestUDPValidation(t *testing.T) {
	for _, tc := range []struct {
		name, encoding, payload, expected string
		port                              int
		valid                             bool
	}{
		{"text", "text", "ping", "pong", 9000, true},
		{"binary", "hex", "00 ff 12", "ff", 9000, true},
		{"default encoding", "", "ping", "", 9000, true},
		{"empty probe", "text", "", "", 9000, false},
		{"bad port", "text", "ping", "", 0, false},
		{"bad hex", "hex", "0x12", "", 9000, false},
		{"bad assertion", "hex", "12", "odd", 9000, false},
		{"oversized", "text", strings.Repeat("x", 1401), "", 9000, false},
		{"unknown encoding", "base64", "aA==", "", 9000, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := Check{Type: "udp", Target: "example.com", Port: tc.port, UDPEncoding: tc.encoding, UDPPayload: tc.payload, UDPExpected: tc.expected}
			if err := c.Validate(); (err == nil) != tc.valid {
				t.Fatalf("valid=%v, error=%v", tc.valid, err)
			}
		})
	}
}
