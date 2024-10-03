// Package test contains utility code used in tests
package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// NewTimeFromString returns a time value parsed from a string
// in the RFC3339Nano format. Note that it cuts off trailing zeros in the milliseconds part, which
// might cause issues in IAM endpoints which do not accept the time format without the milliseconds part.
//
// Example: if "2025-10-11T23:06:59.000Z" is used, the actual value that will be sent is "2025-10-11T23:06:59Z".
func NewTimeFromString(t *testing.T, s string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339Nano, s)
	require.NoError(t, err)
	return parsedTime
}
