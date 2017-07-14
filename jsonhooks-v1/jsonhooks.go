// Package jsonhooks adds hooks that are automatically called before JSON marshaling (PreMarshalJSON) and
// after JSON unmarshaling (PostUnmarshalJSON). It does not do so recursively.
package jsonhooks

import (
	"encoding/json"
)

// Marshal wraps encoding/json.Marshal, calls v.PreMarshalJSON() if it exists
func Marshal(v interface{}) ([]byte, error) {
	if _, ok := v.(PreJSONMarshaler); ok {
		err := v.(PreJSONMarshaler).PreMarshalJSON()
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(v)
}

// Unmarshal wraps encoding/json.Unmarshal, calls v.PostUnmarshalJSON() if it exists
func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	if _, ok := v.(PostJSONUnmarshaler); ok {
		err := v.(PostJSONUnmarshaler).PostUnmarshalJSON()
		if err != nil {
			return err
		}
	}

	return nil
}

// PreJSONMarshaler infers support for the PreMarshalJSON pre-hook
type PreJSONMarshaler interface {
	PreMarshalJSON() error
}

// ImplementsPreJSONMarshaler checks for support for the PreMarshalJSON pre-hook
func ImplementsPreJSONMarshaler(v interface{}) bool {
	_, ok := v.(PreJSONMarshaler)
	return ok
}

// PostJSONUnmarshaler infers support for the PostUnmarshalJSON post-hook
type PostJSONUnmarshaler interface {
	PostUnmarshalJSON() error
}

// ImplementsPostJSONUnmarshaler checks for support for the PostUnmarshalJSON post-hook
func ImplementsPostJSONUnmarshaler(v interface{}) (interface{}, bool) {
	v, ok := v.(PostJSONUnmarshaler)
	return v, ok
}
