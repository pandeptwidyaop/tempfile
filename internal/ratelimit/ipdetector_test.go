package ratelimit

import (
	"testing"
)

func TestIPDetector_GetRealIP(t *testing.T) {
	tests := []struct {
		name           string
		trustedProxies []string
		headers        map[string]string
		remoteAddr     string
		expected       string
	}{
		{
			name:           "Direct connection",
			trustedProxies: []string{"127.0.0.1"},
			headers:        map[string]string{},
			remoteAddr:     "203.0.113.1:12345",
			expected:       "203.0.113.1",
		},
		{
			name:           "Cloudflare CF-Connecting-IP",
			trustedProxies: []string{"127.0.0.1"},
			headers: map[string]string{
				"CF-Connecting-IP": "203.0.113.1",
				"X-Forwarded-For":  "203.0.113.1, 172.16.0.1",
			},
			remoteAddr: "172.16.0.1:12345",
			expected:   "203.0.113.1",
		},
		{
			name:           "X-Real-IP header",
			trustedProxies: []string{"127.0.0.1"},
			headers: map[string]string{
				"X-Real-IP": "203.0.113.1",
			},
			remoteAddr: "127.0.0.1:12345",
			expected:   "203.0.113.1",
		},
		{
			name:           "X-Forwarded-For with multiple IPs",
			trustedProxies: []string{"127.0.0.1", "172.16.0.1"},
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.1, 172.16.0.1, 127.0.0.1",
			},
			remoteAddr: "127.0.0.1:12345",
			expected:   "203.0.113.1",
		},
		{
			name:           "X-Forwarded-For with trusted proxy CIDR",
			trustedProxies: []string{"172.16.0.0/12"},
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.1, 172.16.0.1",
			},
			remoteAddr: "172.16.0.1:12345",
			expected:   "203.0.113.1",
		},
		{
			name:           "Private IP from trusted proxy",
			trustedProxies: []string{"127.0.0.1"},
			headers: map[string]string{
				"X-Forwarded-For": "192.168.1.100",
			},
			remoteAddr: "127.0.0.1:12345",
			expected:   "192.168.1.100",
		},
		{
			name:           "IPv6 address",
			trustedProxies: []string{"::1"},
			headers: map[string]string{
				"X-Forwarded-For": "2001:db8::1",
			},
			remoteAddr: "[::1]:12345",
			expected:   "2001:db8::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := NewIPDetector(tt.trustedProxies, GetDefaultIPHeaders())
			result := detector.GetRealIP(tt.headers, tt.remoteAddr)

			if result != tt.expected {
				t.Errorf("GetRealIP() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIPDetector_IsTrustedProxy(t *testing.T) {
	tests := []struct {
		name           string
		trustedProxies []string
		ip             string
		expected       bool
	}{
		{
			name:           "Exact IP match",
			trustedProxies: []string{"127.0.0.1", "172.16.0.1"},
			ip:             "127.0.0.1",
			expected:       true,
		},
		{
			name:           "CIDR match",
			trustedProxies: []string{"172.16.0.0/12"},
			ip:             "172.16.0.1",
			expected:       true,
		},
		{
			name:           "CIDR no match",
			trustedProxies: []string{"172.16.0.0/12"},
			ip:             "192.168.1.1",
			expected:       false,
		},
		{
			name:           "IPv6 exact match",
			trustedProxies: []string{"::1"},
			ip:             "::1",
			expected:       true,
		},
		{
			name:           "IPv6 CIDR match",
			trustedProxies: []string{"2001:db8::/32"},
			ip:             "2001:db8::1",
			expected:       true,
		},
		{
			name:           "Not trusted",
			trustedProxies: []string{"127.0.0.1"},
			ip:             "203.0.113.1",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := NewIPDetector(tt.trustedProxies, GetDefaultIPHeaders())
			result := detector.IsTrustedProxy(tt.ip)

			if result != tt.expected {
				t.Errorf("IsTrustedProxy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIPDetector_IsWhitelisted(t *testing.T) {
	tests := []struct {
		name         string
		whitelistIPs []string
		ip           string
		expected     bool
	}{
		{
			name:         "Exact IP match",
			whitelistIPs: []string{"127.0.0.1", "203.0.113.1"},
			ip:           "203.0.113.1",
			expected:     true,
		},
		{
			name:         "CIDR match",
			whitelistIPs: []string{"192.168.0.0/16"},
			ip:           "192.168.1.100",
			expected:     true,
		},
		{
			name:         "CIDR no match",
			whitelistIPs: []string{"192.168.0.0/16"},
			ip:           "10.0.0.1",
			expected:     false,
		},
		{
			name:         "IPv6 exact match",
			whitelistIPs: []string{"2001:db8::1"},
			ip:           "2001:db8::1",
			expected:     true,
		},
		{
			name:         "IPv6 CIDR match",
			whitelistIPs: []string{"2001:db8::/32"},
			ip:           "2001:db8::100",
			expected:     true,
		},
		{
			name:         "Not whitelisted",
			whitelistIPs: []string{"127.0.0.1"},
			ip:           "203.0.113.1",
			expected:     false,
		},
		{
			name:         "Empty whitelist",
			whitelistIPs: []string{},
			ip:           "203.0.113.1",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := NewIPDetectorWithWhitelist(
				[]string{"127.0.0.1"},
				GetDefaultIPHeaders(),
				tt.whitelistIPs,
			)

			if ipDet, ok := detector.(*ipDetector); ok {
				result := ipDet.IsWhitelisted(tt.ip)
				if result != tt.expected {
					t.Errorf("IsWhitelisted() = %v, want %v", result, tt.expected)
				}
			} else {
				t.Fatal("Expected *ipDetector type")
			}
		})
	}
}

func TestNewIPDetectorWithWhitelist(t *testing.T) {
	trustedProxies := []string{"127.0.0.1", "172.16.0.0/12"}
	headers := []string{"X-Real-IP", "X-Forwarded-For"}
	whitelistIPs := []string{"203.0.113.1", "192.168.0.0/16"}

	detector := NewIPDetectorWithWhitelist(trustedProxies, headers, whitelistIPs)

	if detector == nil {
		t.Fatal("NewIPDetectorWithWhitelist() returned nil")
	}

	// Test that it implements IPDetector interface
	var _ IPDetector = detector

	// Test whitelist functionality
	if ipDet, ok := detector.(*ipDetector); ok {
		if !ipDet.IsWhitelisted("203.0.113.1") {
			t.Error("Expected IP to be whitelisted")
		}

		if !ipDet.IsWhitelisted("192.168.1.100") {
			t.Error("Expected IP in CIDR range to be whitelisted")
		}

		if ipDet.IsWhitelisted("10.0.0.1") {
			t.Error("Expected IP to not be whitelisted")
		}
	} else {
		t.Fatal("Expected *ipDetector type")
	}
}

func TestGetDefaultTrustedProxies(t *testing.T) {
	proxies := GetDefaultTrustedProxies()

	if len(proxies) == 0 {
		t.Error("GetDefaultTrustedProxies() should return non-empty slice")
	}

	// Check that localhost is included
	found := false
	for _, proxy := range proxies {
		if proxy == "127.0.0.1" {
			found = true
			break
		}
	}

	if !found {
		t.Error("GetDefaultTrustedProxies() should include 127.0.0.1")
	}
}

func TestGetDefaultIPHeaders(t *testing.T) {
	headers := GetDefaultIPHeaders()

	if len(headers) == 0 {
		t.Error("GetDefaultIPHeaders() should return non-empty slice")
	}

	// Check that CF-Connecting-IP is first (highest priority)
	if headers[0] != "CF-Connecting-IP" {
		t.Errorf("GetDefaultIPHeaders() first header should be CF-Connecting-IP, got %s", headers[0])
	}
}

func BenchmarkIPDetector_GetRealIP(b *testing.B) {
	detector := NewIPDetector(GetDefaultTrustedProxies(), GetDefaultIPHeaders())
	headers := map[string]string{
		"X-Forwarded-For": "203.0.113.1, 172.16.0.1, 127.0.0.1",
		"X-Real-IP":       "203.0.113.1",
	}
	remoteAddr := "127.0.0.1:12345"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector.GetRealIP(headers, remoteAddr)
	}
}
