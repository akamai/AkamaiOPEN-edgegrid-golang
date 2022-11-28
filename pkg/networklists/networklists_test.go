package networklists

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockAPIClient(t *testing.T, mockServer *httptest.Server) NTWRKLISTS {
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

func dummyOpt() Option {
	return func(*networklists) {

	}
}

// loadFixtureBytes returns the entire contents of the given file as a byte slice
func loadFixtureBytes(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// compactJSON converts a JSON-encoded byte slice to a compact form (so our JSON fixtures can be readable)
func compactJSON(encoded []byte) string {
	buf := bytes.Buffer{}
	if err := json.Compact(&buf, encoded); err != nil {
		panic(fmt.Sprintf("%s: %s", err, string(encoded)))
	}

	return buf.String()
}

func TestClient(t *testing.T) {
	sess, err := session.New()
	require.NoError(t, err)
	tests := map[string]struct {
		options  []Option
		expected *networklists
	}{
		"no options provided, return default": {
			options: nil,
			expected: &networklists{
				Session: sess,
			},
		},
		"dummy option": {
			options: []Option{dummyOpt()},
			expected: &networklists{
				Session: sess,
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
