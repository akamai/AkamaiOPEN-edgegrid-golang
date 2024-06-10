// Package ptr helps with creating pointers to values of any type.
package ptr

// To returns a pointer to a copy of a value of any type.
func To[T any](t T) *T {
	return &t
}
