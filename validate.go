package contracts

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

var hostPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9.\-]*$`)

func ValidHost(host string) bool {
	return len(host) <= 253 && (net.ParseIP(host) != nil || hostPattern.MatchString(host))
}

// Validate 在任何执行器入口再次检查；客户端校验不能代替服务端校验。
func (c *Check) Validate() error {
	if c.Timeout == 0 {
		c.Timeout = 10
	}
	if c.Count == 0 {
		c.Count = 4
	}
	if c.Size == 0 {
		c.Size = 32
	}
	if c.Method == "" {
		c.Method = "GET"
	}
	if c.DNSRecord == "" {
		c.DNSRecord = "A"
	}
	if c.Family == "" {
		c.Family = "auto"
	}
	if c.Timeout < 1 || c.Timeout > 60 || c.Count < 1 || c.Count > 20 || c.Size < 16 || c.Size > 1400 || c.Retries < 0 || c.Retries > 2 || c.MaxLoss < 0 || c.MaxLoss > 100 {
		return fmt.Errorf("检测参数超出允许范围")
	}
	if c.Family != "auto" && c.Family != "ipv4" && c.Family != "ipv6" {
		return fmt.Errorf("不支持的 IP 协议族")
	}
	if c.CertificateMinDays < 0 || c.CertificateMinDays > 365 {
		return fmt.Errorf("证书剩余天数阈值应为 0—365；0 表示关闭")
	}
	switch c.Type {
	case TypeHTTP:
		u, err := url.Parse(c.Target)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || !ValidHost(u.Hostname()) || u.Fragment != "" {
			return fmt.Errorf("请输入有效的 HTTP/HTTPS URL（不含用户凭据和片段）")
		}
		if c.Method != "GET" && c.Method != "HEAD" {
			return fmt.Errorf("仅允许 GET/HEAD")
		}
		if c.ExpectedStatus != 0 && (c.ExpectedStatus < 100 || c.ExpectedStatus > 599) {
			return fmt.Errorf("无效期望状态码")
		}
		if len(c.Headers) > 32 || len(c.Contains) > 4096 {
			return fmt.Errorf("HTTP 断言或请求头过大")
		}
		for k, v := range c.Headers {
			if len(k)+len(v) > 4096 || strings.ContainsAny(k+v, "\r\n") {
				return fmt.Errorf("无效请求头")
			}
		}
	case TypeTCP, TypeUDP, TypePing, TypeDNS, TypeMTR, TypeTraceroute:
		if !ValidHost(c.Target) {
			return fmt.Errorf("请输入有效域名或 IP 地址")
		}
		if (c.Type == TypeTCP || c.Type == TypeUDP) && (c.Port < 1 || c.Port > 65535) {
			return fmt.Errorf("端口必须在 1—65535 之间")
		}
		if c.Type == TypeUDP {
			if c.UDPEncoding == "" {
				c.UDPEncoding = "text"
			}
			payload, err := UDPBytes(c.UDPPayload, c.UDPEncoding)
			if err != nil || len(payload) == 0 || len(payload) > 1400 {
				return fmt.Errorf("UDP 请求必须为 1—1400 字节的有效文本或十六进制数据")
			}
			expected, err := UDPBytes(c.UDPExpected, c.UDPEncoding)
			if err != nil || len(expected) > 4096 {
				return fmt.Errorf("UDP 期望响应必须使用相同编码且不超过 4096 字节")
			}
		}
		if c.Type == TypeDNS && !strings.Contains("|A|AAAA|CNAME|MX|TXT|NS|SOA|", "|"+c.DNSRecord+"|") {
			return fmt.Errorf("不支持的 DNS 类型")
		}
	default:
		return fmt.Errorf("不支持的检测类型")
	}
	return nil
}

func RunID(monitor, node string, version int, at int64) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d:%d", monitor, node, version, at)))
	return hex.EncodeToString(sum[:])
}

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
