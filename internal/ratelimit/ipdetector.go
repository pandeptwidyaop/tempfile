package ratelimit

import (
	"net"
	"strings"
)

// ipDetector implements the IPDetector interface
type ipDetector struct {
	trustedProxies map[string]bool
	trustedCIDRs   []*net.IPNet
	headerPriority []string
	whitelistIPs   map[string]bool
	whitelistCIDRs []*net.IPNet
}

// NewIPDetector creates a new IP detector with the given configuration
func NewIPDetector(trustedProxies []string, headerPriority []string) IPDetector {
	return NewIPDetectorWithWhitelist(trustedProxies, headerPriority, []string{})
}

// NewIPDetectorWithWhitelist creates a new IP detector with whitelist support
func NewIPDetectorWithWhitelist(trustedProxies []string, headerPriority []string, whitelistIPs []string) IPDetector {
	detector := &ipDetector{
		trustedProxies: make(map[string]bool),
		headerPriority: headerPriority,
		whitelistIPs:   make(map[string]bool),
	}

	// Parse trusted proxies
	for _, proxy := range trustedProxies {
		proxy = strings.TrimSpace(proxy)
		if proxy == "" {
			continue
		}

		// Try to parse as CIDR first
		if _, cidr, err := net.ParseCIDR(proxy); err == nil {
			detector.trustedCIDRs = append(detector.trustedCIDRs, cidr)
		} else if ip := net.ParseIP(proxy); ip != nil {
			// Parse as individual IP
			detector.trustedProxies[proxy] = true
		}
	}

	// Parse whitelist IPs
	for _, ip := range whitelistIPs {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			continue
		}

		// Try to parse as CIDR first
		if _, cidr, err := net.ParseCIDR(ip); err == nil {
			detector.whitelistCIDRs = append(detector.whitelistCIDRs, cidr)
		} else if parsedIP := net.ParseIP(ip); parsedIP != nil {
			// Parse as individual IP
			detector.whitelistIPs[ip] = true
		}
	}

	return detector
}

// GetRealIP extracts the real client IP from HTTP headers
func (d *ipDetector) GetRealIP(headers map[string]string, remoteAddr string) string {
	// Try each header in priority order
	for _, header := range d.headerPriority {
		if value := headers[header]; value != "" {
			if ip := d.parseIPFromHeader(value); ip != "" {
				return ip
			}
		}
	}

	// Fallback to remote address
	return d.extractIPFromAddr(remoteAddr)
}

// IsTrustedProxy checks if an IP is a trusted proxy
func (d *ipDetector) IsTrustedProxy(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check individual IPs
	if d.trustedProxies[ip] {
		return true
	}

	// Check CIDR ranges
	for _, cidr := range d.trustedCIDRs {
		if cidr.Contains(parsedIP) {
			return true
		}
	}

	return false
}

// IsWhitelisted checks if an IP is in the whitelist
func (d *ipDetector) IsWhitelisted(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check individual IPs
	if d.whitelistIPs[ip] {
		return true
	}

	// Check CIDR ranges
	for _, cidr := range d.whitelistCIDRs {
		if cidr.Contains(parsedIP) {
			return true
		}
	}

	return false
}

// parseIPFromHeader extracts the first valid non-trusted IP from a header value
func (d *ipDetector) parseIPFromHeader(headerValue string) string {
	// Handle X-Forwarded-For format: "client, proxy1, proxy2"
	ips := strings.Split(headerValue, ",")

	for _, ipStr := range ips {
		ip := strings.TrimSpace(ipStr)

		// Skip empty values
		if ip == "" {
			continue
		}

		// Extract IP from "ip:port" format
		if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
			// Check if this looks like IPv6
			if strings.Count(ip, ":") > 1 {
				// IPv6 address, don't split on colon
			} else {
				// IPv4 with port, remove port
				ip = ip[:colonIndex]
			}
		}

		// Validate IP format
		if parsedIP := net.ParseIP(ip); parsedIP != nil {
			// Skip trusted proxies and private IPs (unless from trusted proxy)
			if !d.IsTrustedProxy(ip) && !d.isPrivateIP(parsedIP) {
				return ip
			}

			// If it's a private IP but we're behind a trusted proxy, use it
			if d.isPrivateIP(parsedIP) {
				return ip
			}
		}
	}

	return ""
}

// extractIPFromAddr extracts IP from "ip:port" format
func (d *ipDetector) extractIPFromAddr(addr string) string {
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}

// isPrivateIP checks if an IP is in private ranges
func (d *ipDetector) isPrivateIP(ip net.IP) bool {
	// Private IPv4 ranges
	private4 := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
	}

	// Private IPv6 ranges
	private6 := []string{
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}

	ranges := private4
	if ip.To4() == nil {
		ranges = append(ranges, private6...)
	}

	for _, cidr := range ranges {
		if _, network, err := net.ParseCIDR(cidr); err == nil {
			if network.Contains(ip) {
				return true
			}
		}
	}

	return false
}

// GetDefaultTrustedProxies returns common trusted proxy configurations
func GetDefaultTrustedProxies() []string {
	return []string{
		"127.0.0.1",
		"::1",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		// Cloudflare IP ranges (commonly used)
		"173.245.48.0/20",
		"103.21.244.0/22",
		"103.22.200.0/22",
		"103.31.4.0/22",
		"141.101.64.0/18",
		"108.162.192.0/18",
		"190.93.240.0/20",
		"188.114.96.0/20",
		"197.234.240.0/22",
		"198.41.128.0/17",
		"162.158.0.0/15",
		"104.16.0.0/13",
		"104.24.0.0/14",
		"172.64.0.0/13",
		"131.0.72.0/22",
	}
}

// GetDefaultIPHeaders returns the default priority order for IP headers
func GetDefaultIPHeaders() []string {
	return []string{
		"CF-Connecting-IP",
		"X-Real-IP",
		"X-Forwarded-For",
		"X-Forwarded",
		"Forwarded-For",
		"Forwarded",
	}
}
