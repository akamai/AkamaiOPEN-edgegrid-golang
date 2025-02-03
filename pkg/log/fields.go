package log

import (
	"slices"
)

// Fields is a helper type used for logging
type Fields map[string]interface{}

// Fields implements Fielder.
func (f Fields) Fields() Fields {
	return f
}

// Get returns flatten Fields map
func (f Fields) Get() []any {

	keys := make([]string, 0, len(f))
	for key := range f {
		keys = append(keys, key)
	}
	// Sort keys for consistent output
	slices.Sort(keys)

	out := make([]any, 0, len(f)*2)
	for _, key := range keys {
		out = append(out, key, f[key])
	}
	return out
}

// Fielder interface allows for creating custom Fields
type Fielder interface {
	Fields() Fields
}
