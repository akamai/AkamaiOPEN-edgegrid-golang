package mtlskeystore

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	sess, err := session.New()
	require.NoError(t, err)
	tests := map[string]struct {
		options  []Option
		expected *mtlskeystore
	}{
		"no options provided, return default": {
			options: nil,
			expected: &mtlskeystore{
				Session: sess,
			},
		},
		"option provided, overwrite session": {
			options: []Option{func(c *mtlskeystore) {
				c.Session = nil
			}},
			expected: &mtlskeystore{
				Session: nil,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess, test.options...)
			assert.Equal(t, res, test.expected)
		})
	}
}
