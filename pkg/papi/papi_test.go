package papi

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockAPIClient(t *testing.T, mockServer *httptest.Server) PAPI {
	serverURL, err := url.Parse(mockServer.URL)
	require.NoError(t, err)
	certPool := x509.NewCertPool()
	certPool.AddCert(mockServer.Certificate())
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}
	s, err := session.New(session.WithClient(httpClient), session.WithSigner(&edgegrid.Config{Host: serverURL.Host}))
	assert.NoError(t, err)
	return Client(s)
}

func TestClient(t *testing.T) {
	sess, err := session.New()
	require.NoError(t, err)
	tests := map[string]struct {
		options  []Option
		expected *papi
	}{
		"no options provided, return default": {
			options: nil,
			expected: &papi{
				Session:     sess,
				usePrefixes: true,
			},
		},
		"papi prefixes set to false": {
			options: []Option{WithUsePrefixes(false)},
			expected: &papi{
				Session:     sess,
				usePrefixes: false,
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
