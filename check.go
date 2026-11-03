package contracts

const (
	TypeHTTP                    = "http"
	TypeTCP                     = "tcp"
	TypeUDP                     = "udp"
	TypePing                    = "ping"
	TypeDNS                     = "dns"
	TypeMTR                     = "mtr"
	TypeTraceroute              = "traceroute"
	CapabilityCertificateExpiry = "certificate_expiry"
)

// CheckDefinition is the shared catalog used by the API and configuration UI.
type CheckDefinition struct {
	Type       string `json:"type"`
	Continuous bool   `json:"continuous"`
	Group      string `json:"group"`
}

var checkCatalog = []CheckDefinition{
	{TypeHTTP, true, "application"},
	{TypeTCP, true, "network"},
	{TypeUDP, true, "network"},
	{TypePing, true, "network"},
	{TypeDNS, true, "application"},
	{TypeMTR, false, "route"},
	{TypeTraceroute, false, "route"},
}

func CheckDefinitions() []CheckDefinition { return append([]CheckDefinition(nil), checkCatalog...) }

func IsContinuous(kind string) bool {
	for _, definition := range checkCatalog {
		if definition.Type == kind {
			return definition.Continuous
		}
	}
	return false
}

type Check struct {
	Type               string            `json:"type"`
	Target             string            `json:"target"`
	Port               int               `json:"port"`
	UDPPayload         string            `json:"udp_payload,omitempty"`
	UDPEncoding        string            `json:"udp_encoding,omitempty"`
	UDPExpected        string            `json:"udp_expected,omitempty"`
	Timeout            int               `json:"timeout"`
	Method             string            `json:"method"`
	Headers            map[string]string `json:"headers,omitempty"`
	ExpectedStatus     int               `json:"expected_status"`
	Contains           string            `json:"contains,omitempty"`
	FollowRedirects    bool              `json:"follow_redirects"`
	InsecureTLS        bool              `json:"insecure_tls"`
	CertificateMinDays int               `json:"certificate_min_days"`
	DNSRecord          string            `json:"dns_record"`
	Resolver           string            `json:"resolver,omitempty"`
	ExpectedDNS        string            `json:"expected_dns,omitempty"`
	Count              int               `json:"count"`
	Size               int               `json:"size"`
	MaxLoss            float64           `json:"max_loss"`
	Family             string            `json:"family"`
	Retries            int               `json:"retries"`
}

// RequiredCapabilities keeps scheduling and manual diagnostics consistent.
func (c Check) RequiredCapabilities() []string {
	required := []string{c.Type}
	if c.Type == TypeHTTP && c.CertificateMinDays > 0 {
		required = append(required, CapabilityCertificateExpiry)
	}
	return required
}

func (c Check) SupportedBy(capabilities []string) bool {
	for _, required := range c.RequiredCapabilities() {
		found := false
		for _, available := range capabilities {
			if available == required {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
