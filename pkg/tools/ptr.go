// Package tools contains utilities used in EdgeGrid
package tools

// BoolPtr returns the address of the bool
//
// Deprecated: this function will be removed in a future release. Use [ptr.To] instead.
func BoolPtr(b bool) *bool {
	return &b
}

// IntPtr returns the address of the int
//
// Deprecated: this function will be removed in a future release. Use [ptr.To] instead.
func IntPtr(i int) *int {
	return &i
}

// Int64Ptr returns the address of the int64
//
// Deprecated: this function will be removed in a future release. Use [ptr.To] instead.
func Int64Ptr(i int64) *int64 {
	return &i
}

// Float32Ptr returns the address of the float32
//
// Deprecated: this function will be removed in a future release. Use [ptr.To] instead.
func Float32Ptr(f float32) *float32 {
	return &f
}

// Float64Ptr returns the address of the float64
//
// Deprecated: this function will be removed in a future release. Use [ptr.To] instead.
func Float64Ptr(f float64) *float64 {
	return &f
}

// StringPtr returns the address of the string
//
// Deprecated: this function will be removed in a future release. Use [ptr.To] instead.
func StringPtr(s string) *string {
	return &s
}
