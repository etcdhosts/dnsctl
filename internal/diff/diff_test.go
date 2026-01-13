package diff

import (
	"testing"
)

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single line no newline",
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			name:     "single line with newline",
			input:    "hello\n",
			expected: []string{"hello"},
		},
		{
			name:     "multiple lines",
			input:    "line1\nline2\nline3\n",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "multiple lines no trailing newline",
			input:    "line1\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitLines(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("SplitLines(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("SplitLines(%q)[%d] = %q, want %q", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestCompute(t *testing.T) {
	tests := []struct {
		name     string
		old      []string
		newData  []string
		expected []Line
	}{
		{
			name:     "identical",
			old:      []string{"a", "b", "c"},
			newData:  []string{"a", "b", "c"},
			expected: []Line{{OpEqual, "a"}, {OpEqual, "b"}, {OpEqual, "c"}},
		},
		{
			name:     "all different",
			old:      []string{"a", "b"},
			newData:  []string{"x", "y"},
			expected: []Line{{OpDelete, "a"}, {OpDelete, "b"}, {OpInsert, "x"}, {OpInsert, "y"}},
		},
		{
			name:     "insert at end",
			old:      []string{"a", "b"},
			newData:  []string{"a", "b", "c"},
			expected: []Line{{OpEqual, "a"}, {OpEqual, "b"}, {OpInsert, "c"}},
		},
		{
			name:     "delete from end",
			old:      []string{"a", "b", "c"},
			newData:  []string{"a", "b"},
			expected: []Line{{OpEqual, "a"}, {OpEqual, "b"}, {OpDelete, "c"}},
		},
		{
			name:     "insert in middle",
			old:      []string{"a", "c"},
			newData:  []string{"a", "b", "c"},
			expected: []Line{{OpEqual, "a"}, {OpInsert, "b"}, {OpEqual, "c"}},
		},
		{
			name:     "delete from middle",
			old:      []string{"a", "b", "c"},
			newData:  []string{"a", "c"},
			expected: []Line{{OpEqual, "a"}, {OpDelete, "b"}, {OpEqual, "c"}},
		},
		{
			name:     "empty old",
			old:      []string{},
			newData:  []string{"a", "b"},
			expected: []Line{{OpInsert, "a"}, {OpInsert, "b"}},
		},
		{
			name:     "empty new",
			old:      []string{"a", "b"},
			newData:  []string{},
			expected: []Line{{OpDelete, "a"}, {OpDelete, "b"}},
		},
		{
			name:     "both empty",
			old:      []string{},
			newData:  []string{},
			expected: []Line{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Compute(tt.old, tt.newData)
			if len(result) != len(tt.expected) {
				t.Errorf("Compute() got %d lines, want %d", len(result), len(tt.expected))
				t.Errorf("got: %v", result)
				t.Errorf("want: %v", tt.expected)
				return
			}
			for i := range result {
				if result[i].Op != tt.expected[i].Op || result[i].Text != tt.expected[i].Text {
					t.Errorf("Compute()[%d] = {%v, %q}, want {%v, %q}",
						i, result[i].Op, result[i].Text, tt.expected[i].Op, tt.expected[i].Text)
				}
			}
		})
	}
}

func TestCompute_ComplexCase(t *testing.T) {
	old := []string{
		"192.168.1.1 web.local",
		"192.168.1.2 db.local",
		"192.168.1.3 cache.local",
	}
	newData := []string{
		"192.168.1.1 web.local",
		"192.168.1.4 api.local",
		"192.168.1.3 cache.local",
	}

	result := Compute(old, newData)

	// Verify structure: web.local unchanged, db.local deleted, api.local inserted, cache.local unchanged
	var deletes, inserts, equals int
	for _, line := range result {
		switch line.Op {
		case OpDelete:
			deletes++
		case OpInsert:
			inserts++
		case OpEqual:
			equals++
		}
	}

	if equals != 2 {
		t.Errorf("expected 2 equal lines, got %d", equals)
	}
	if deletes != 1 {
		t.Errorf("expected 1 delete, got %d", deletes)
	}
	if inserts != 1 {
		t.Errorf("expected 1 insert, got %d", inserts)
	}
}
