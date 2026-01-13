// Package output provides formatted output utilities.
package output

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Stringer is an interface for types that can be converted to string.
type Stringer interface {
	String() string
}

// Format represents an output format.
type Format string

const (
	FormatHosts Format = "hosts"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

// Print outputs data in the specified format.
func Print(data any, format Format) error {
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	case FormatYAML:
		return yaml.NewEncoder(os.Stdout).Encode(data)
	default:
		if s, ok := data.(Stringer); ok {
			fmt.Print(s.String())
		}
		return nil
	}
}
