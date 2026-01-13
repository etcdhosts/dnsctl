package output

import (
	"net"
	"testing"

	client "github.com/etcdhosts/client-go/v2"
)

func TestFormatRecordAttrs(t *testing.T) {
	tests := []struct {
		name     string
		record   client.Record
		expected string
	}{
		{
			name: "default values",
			record: client.Record{
				Hostname: "test.local",
				IP:       net.ParseIP("192.168.1.1"),
				Weight:   1,
			},
			expected: "",
		},
		{
			name: "custom weight",
			record: client.Record{
				Hostname: "test.local",
				IP:       net.ParseIP("192.168.1.1"),
				Weight:   5,
			},
			expected: "[weight=5]",
		},
		{
			name: "custom TTL",
			record: client.Record{
				Hostname: "test.local",
				IP:       net.ParseIP("192.168.1.1"),
				Weight:   1,
				TTL:      300,
			},
			expected: "[ttl=300]",
		},
		{
			name: "weight and TTL",
			record: client.Record{
				Hostname: "test.local",
				IP:       net.ParseIP("192.168.1.1"),
				Weight:   3,
				TTL:      600,
			},
			expected: "[weight=3, ttl=600]",
		},
		{
			name: "with health check TCP",
			record: client.Record{
				Hostname: "test.local",
				IP:       net.ParseIP("192.168.1.1"),
				Weight:   1,
				Health:   &client.Health{Type: client.CheckTCP, Port: 8080},
			},
			expected: "[hc=tcp:8080]",
		},
		{
			name: "with health check HTTP",
			record: client.Record{
				Hostname: "test.local",
				IP:       net.ParseIP("192.168.1.1"),
				Weight:   1,
				Health:   &client.Health{Type: client.CheckHTTP, Port: 80, Path: "/health"},
			},
			expected: "[hc=http:80/health]",
		},
		{
			name: "all attributes",
			record: client.Record{
				Hostname: "test.local",
				IP:       net.ParseIP("192.168.1.1"),
				Weight:   2,
				TTL:      120,
				Health:   &client.Health{Type: client.CheckICMP},
			},
			expected: "[weight=2, ttl=120, hc=icmp]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatRecordAttrs(tt.record)
			if result != tt.expected {
				t.Errorf("FormatRecordAttrs() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatHealthCheck(t *testing.T) {
	tests := []struct {
		name     string
		health   *client.Health
		expected string
	}{
		{
			name:     "ICMP",
			health:   &client.Health{Type: client.CheckICMP},
			expected: "icmp",
		},
		{
			name:     "TCP",
			health:   &client.Health{Type: client.CheckTCP, Port: 3306},
			expected: "tcp:3306",
		},
		{
			name:     "HTTP without path",
			health:   &client.Health{Type: client.CheckHTTP, Port: 8080},
			expected: "http:8080",
		},
		{
			name:     "HTTP with path",
			health:   &client.Health{Type: client.CheckHTTP, Port: 80, Path: "/ready"},
			expected: "http:80/ready",
		},
		{
			name:     "HTTPS with path",
			health:   &client.Health{Type: client.CheckHTTPS, Port: 443, Path: "/healthz"},
			expected: "https:443/healthz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatHealthCheck(tt.health)
			if result != tt.expected {
				t.Errorf("FormatHealthCheck() = %q, want %q", result, tt.expected)
			}
		})
	}
}
