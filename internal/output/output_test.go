package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

type mockStringer struct {
	value string
}

func (m mockStringer) String() string {
	return m.value
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestPrint_Hosts(t *testing.T) {
	data := mockStringer{value: "192.168.1.1 test.local\n"}

	output := captureStdout(func() {
		_ = Print(data, FormatHosts)
	})

	if output != "192.168.1.1 test.local\n" {
		t.Errorf("Print(hosts) = %q, want %q", output, "192.168.1.1 test.local\n")
	}
}

func TestPrint_JSON(t *testing.T) {
	data := map[string]string{"key": "value"}

	output := captureStdout(func() {
		_ = Print(data, FormatJSON)
	})

	if !strings.Contains(output, `"key": "value"`) {
		t.Errorf("Print(json) = %q, expected JSON with key:value", output)
	}
}

func TestPrint_YAML(t *testing.T) {
	data := map[string]string{"key": "value"}

	output := captureStdout(func() {
		_ = Print(data, FormatYAML)
	})

	if !strings.Contains(output, "key: value") {
		t.Errorf("Print(yaml) = %q, expected YAML with key:value", output)
	}
}

func TestPrint_NonStringer(t *testing.T) {
	// When data doesn't implement Stringer and format is hosts, nothing should print
	data := map[string]int{"count": 42}

	output := captureStdout(func() {
		_ = Print(data, FormatHosts)
	})

	if output != "" {
		t.Errorf("Print(hosts, non-stringer) = %q, want empty string", output)
	}
}

func TestFormat_String(t *testing.T) {
	tests := []struct {
		format   Format
		expected string
	}{
		{FormatHosts, "hosts"},
		{FormatJSON, "json"},
		{FormatYAML, "yaml"},
	}

	for _, tt := range tests {
		if string(tt.format) != tt.expected {
			t.Errorf("Format %v = %q, want %q", tt.format, string(tt.format), tt.expected)
		}
	}
}
