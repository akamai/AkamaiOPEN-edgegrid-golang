package tools

// IntPtr returns the address of the int
func IntPtr(i int) *int {
	return &i
}
