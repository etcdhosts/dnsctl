package output

import (
	"fmt"
	"strings"

	client "github.com/etcdhosts/client-go/v2"
)

// FormatRecordAttrs formats record attributes for display.
func FormatRecordAttrs(r client.Record) string {
	var attrs []string

	if r.Weight != 1 {
		attrs = append(attrs, fmt.Sprintf("weight=%d", r.Weight))
	}
	if r.TTL > 0 {
		attrs = append(attrs, fmt.Sprintf("ttl=%d", r.TTL))
	}
	if r.Health != nil {
		attrs = append(attrs, fmt.Sprintf("hc=%s", FormatHealthCheck(r.Health)))
	}

	if len(attrs) == 0 {
		return ""
	}
	return "[" + strings.Join(attrs, ", ") + "]"
}

// FormatHealthCheck formats a health check for display.
func FormatHealthCheck(h *client.Health) string {
	if h.Type == client.CheckICMP {
		return "icmp"
	}
	if h.Path != "" {
		return fmt.Sprintf("%s:%d%s", h.Type, h.Port, h.Path)
	}
	return fmt.Sprintf("%s:%d", h.Type, h.Port)
}
