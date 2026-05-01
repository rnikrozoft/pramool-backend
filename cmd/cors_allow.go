package cmd

import (
	"net"
	"net/url"
	"strings"
)

// corsAllowDevLAN returns true for browser origins safe for local/LAN dev:
// localhost, 127.0.0.1, and RFC1918 / loopback / link-local IPs (e.g. 192.168.x.x:3000).
// Fiber passes Origin already lowercased.
func corsAllowDevLAN(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return false
	}
	switch u.Scheme {
	case "http", "https":
	default:
		return false
	}
	host := strings.TrimSpace(u.Hostname())
	if host == "localhost" || host == "127.0.0.1" {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast()
}
