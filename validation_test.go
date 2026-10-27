package contracts

import (
	"net"
	"strings"
	"testing"
)

func TestUntrustedInputs(t *testing.T) {
	for _, target := range []string{"example.com;id", "-c10", "$(whoami)", "host\ncommand"} {
		c := Check{Type: "ping", Target: target}
		if c.Validate() == nil {
			t.Fatalf("accepted command-like host %q", target)
		}
	}
	for _, target := range []string{"file:///etc/passwd", "https://user:password@example.com/", "https://example.com/#fragment"} {
		c := Check{Type: "http", Target: target}
		if c.Validate() == nil {
			t.Fatalf("accepted unsafe URL %q", target)
		}
	}
}

func TestPrivateNetworksRequireExplicitGrant(t *testing.T) {
	for _, value := range []string{"127.0.0.1", "10.1.2.3", "169.254.169.254", "100.64.0.1", "::1", "fc00::1", "::ffff:192.168.1.1"} {
		if (Policy{}).Allows(net.ParseIP(value)) {
			t.Fatalf("allowed private address %s", value)
		}
	}
	policy := Policy{AllowCIDRs: []string{"10.0.0.0/8"}}
	if !policy.Allows(net.ParseIP("10.1.2.3")) || policy.Allows(net.ParseIP("192.168.1.1")) {
		t.Fatal("CIDR boundary violated")
	}
}

func TestErrorAndRunIdentity(t *testing.T) {
	value := Redact("GET https://user:password@example.com/path?token=private failed")
	if strings.Contains(value, "password") || strings.Contains(value, "private") {
		t.Fatal("URL credential leaked")
	}
	if RunID("m", "n", 1, 60) == RunID("m", "n", 2, 60) || RunID("m", "n", 1, 60) != RunID("m", "n", 1, 60) {
		t.Fatal("run identity is not deterministic and versioned")
	}
}
