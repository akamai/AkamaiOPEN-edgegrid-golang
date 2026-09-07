package purgecache

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Parallel()
	sess, err := session.New()
	require.NoError(t, err)
	tests := map[string]struct {
		options  []Option
		expected *purgecache
	}{
		"no options provided, return default": {
			options: nil,
			expected: &purgecache{
				Session: sess,
			},
		},
		"option provided, overwrite session": {
			options: []Option{func(c *purgecache) {
				c.Session = nil
			}},
			expected: &purgecache{
				Session: nil,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := Client(sess, test.options...)
			assert.Equal(t, test.expected, res)
		})
	}
}
