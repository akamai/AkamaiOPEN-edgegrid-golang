package papi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseLinkParse(t *testing.T) {
	tests := map[string]struct {
		given     string
		expected  string
		withError bool
	}{
		"valid URL passed": {
			given:    "/papi/v1/cpcodes/123?contractId=contract-1TJZFW&groupId=group",
			expected: "123",
		},
		"invalid URL passed": {
			given:     ":",
			withError: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := ResponseLinkParse(test.given)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}
