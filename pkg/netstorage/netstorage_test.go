package netstorage

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockAPIClient(t *testing.T, mockServer *httptest.Server) NetStorage {
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
