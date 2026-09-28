package cloudcertificates

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Trimmed PEM material shared across the cloudcertificates test suite.
const (
	rsaCSRPEM        = "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n"
	ecdsaCSRPEM      = "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n"
	rsaCertPEM       = "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"
	ecdsaCertPEM     = "-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n"
	rsaTrustChainPEM = "-----BEGIN CERTIFICATE-----\nRSA-TRUST-CHAIN\n-----END CERTIFICATE-----\n"
)

func mockAPIClient(t *testing.T, mockServer *httptest.Server) CloudCertificates {
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
	t.Parallel()
	sess, err := session.New()
	require.NoError(t, err)
	tests := map[string]struct {
		options  []Option
		expected *cloudcertificates
	}{
		"no options provided, return default": {
			options: nil,
			expected: &cloudcertificates{
				Session: sess,
			},
		},
		"option provided, overwrite session": {
			options: []Option{func(c *cloudcertificates) {
				c.Session = nil
			}},
			expected: &cloudcertificates{
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
