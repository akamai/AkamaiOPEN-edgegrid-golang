package purgecache

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TO BE UNCOMMENTED WHEN FIRST ENDPOINTS WILL BE ADDED to PURGE CACHE
//func mockAPIClient(t *testing.T, mockServer *httptest.Server) PurgeCache {
//	serverURL, err := url.Parse(mockServer.URL)
//	require.NoError(t, err)
//	certPool := x509.NewCertPool()
//	certPool.AddCert(mockServer.Certificate())
//	httpClient := &http.Client{
//		Transport: &http.Transport{
//			TLSClientConfig: &tls.Config{
//				RootCAs: certPool,
//			},
//		},
//	}
//	s, err := session.New(session.WithClient(httpClient), session.WithSigner(&edgegrid.Config{Host: serverURL.Host}))
//	assert.NoError(t, err)
//	return Client(s)
//}

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
			res := Client(sess, test.options...)
			assert.Equal(t, test.expected, res)
		})
	}
}
