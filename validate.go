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
	switch c.Type {
	case "http":
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
	case "tcp", "ping", "dns", "mtr", "traceroute":
		if !ValidHost(c.Target) {
			return fmt.Errorf("请输入有效域名或 IP 地址")
		}
		if c.Type == "tcp" && (c.Port < 1 || c.Port > 65535) {
			return fmt.Errorf("端口必须在 1—65535 之间")
		}
		if c.Type == "dns" && !strings.Contains("|A|AAAA|CNAME|MX|TXT|NS|SOA|", "|"+c.DNSRecord+"|") {
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
