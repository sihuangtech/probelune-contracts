package contracts

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Resolve 同时检查所有解析地址，连接时使用已检查的 IP，防止 DNS 重绑定。
func (p Policy) Resolve(ctx context.Context, host, family string) ([]net.IP, error) {
	if !ValidHost(host) {
		return nil, fmt.Errorf("无效目标")
	}
	records, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	var ips []net.IP
	for _, record := range records {
		ip := record.IP
		if family == "ipv4" && ip.To4() == nil || family == "ipv6" && ip.To4() != nil {
			continue
		}
		if !p.Allows(ip) {
			return nil, fmt.Errorf("访问策略拒绝目标地址 %s", ip)
		}
		ips = append(ips, ip)
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("目标没有所选协议族的地址")
	}
	return ips, nil
}

func (p Policy) Allows(ip net.IP) bool {
	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return false
	}
	addr = addr.Unmap()
	for _, raw := range p.AllowCIDRs {
		if prefix, e := netip.ParsePrefix(raw); e == nil && prefix.Contains(addr) {
			return true
		}
	}
	if !addr.IsGlobalUnicast() || addr.IsPrivate() || addr.IsLoopback() || addr.IsLinkLocalUnicast() {
		return false
	}
	for _, raw := range []string{"100.64.0.0/10", "192.0.0.0/24", "192.0.2.0/24", "198.18.0.0/15", "198.51.100.0/24", "203.0.113.0/24", "240.0.0.0/4", "2001:db8::/32"} {
		if netip.MustParsePrefix(raw).Contains(addr) {
			return false
		}
	}
	return true
}

func (p Policy) Validate() error {
	if len(p.AllowCIDRs) > 32 {
		return fmt.Errorf("最多允许 32 个 CIDR")
	}
	for _, raw := range p.AllowCIDRs {
		if _, err := netip.ParsePrefix(raw); err != nil {
			return fmt.Errorf("无效 CIDR: %s", raw)
		}
	}
	return nil
}

func (p Policy) Dial(family string) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}
		ips, err := p.Resolve(ctx, host, family)
		if err != nil {
			return nil, err
		}
		var last error
		for _, ip := range ips {
			c, e := (&net.Dialer{Timeout: 10 * time.Second}).DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
			if e == nil {
				return c, nil
			}
			last = e
		}
		return nil, last
	}
}

func Redact(raw string) string {
	raw = regexp.MustCompile(`https?://[^\s]+`).ReplaceAllStringFunc(raw, func(value string) string {
		parsed, err := url.Parse(value)
		if err != nil {
			return "[URL]"
		}
		parsed.User = nil
		parsed.RawQuery = ""
		parsed.Fragment = ""
		return parsed.String()
	})
	if len(raw) > 2048 {
		raw = raw[:2048]
	}
	return strings.ReplaceAll(raw, "\x00", "")
}
