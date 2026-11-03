package contracts

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// UDPBytes decodes a bounded user-supplied datagram or response assertion.
func UDPBytes(value, encoding string) ([]byte, error) {
	if len(value) > 16384 {
		return nil, fmt.Errorf("UDP data is too large")
	}
	switch encoding {
	case "", "text":
		return []byte(value), nil
	case "hex":
		return hex.DecodeString(strings.Join(strings.Fields(value), ""))
	default:
		return nil, fmt.Errorf("UDP encoding must be text or hex")
	}
}
