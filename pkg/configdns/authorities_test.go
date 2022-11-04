package dns

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDns_GetAuthorities(t *testing.T) {
	tests := map[string]struct {
		contractID       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *AuthorityResponse
		withError        error
	}{
		"200 OK": {
			contractID:     "9-9XXXXX",
			responseStatus: http.StatusOK,
			responseBody: `
{
	"contracts": [
        {
            "contractId": "9-9XXXXX",
            "authorities": [
                "a1-118.akam.net.",
                "a2-64.akam.net.",
                "a6-66.akam.net.",
                "a18-67.akam.net.",
                "a7-64.akam.net.",
                "a11-64.akam.net."
            ]
        }
    ]
}`,
			expectedPath: "/config-dns/v2/data/authorities?contractIds=9-9XXXXX",
			expectedResponse: &AuthorityResponse{
				Contracts: []Contract{
					{
						ContractID: "9-9XXXXX",
						Authorities: []string{
							"a1-118.akam.net.",
							"a2-64.akam.net.",
							"a6-66.akam.net.",
							"a18-67.akam.net.",
							"a7-64.akam.net.",
							"a11-64.akam.net.",
						},
					},
				},
			},
		},
		"Missing arguments": {
			contractID:     "",
			responseStatus: http.StatusOK,
			responseBody:   "",
			expectedPath:   "/config-dns/v2/data/authorities?contractIds=9-9XXXXX",
			withError:      ErrBadRequest,
		},
		"500 internal server error": {
			contractID:     "9-9XXXXX",
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching authorities",
    "status": 500
}`,
			expectedPath: "/config-dns/v2/data/authorities?contractIds=9-9XXXXX",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching authorities",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetAuthorities(context.Background(), test.contractID)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_GetNameServerRecordList(t *testing.T) {
	tests := map[string]struct {
		contractID       string
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []string
		withError        error
	}{
		"test with valid arguments": {
			contractID:     "9-9XXXXX",
			responseStatus: http.StatusOK,
			responseBody: `
{
	"contracts": [
        {
            "contractId": "9-9XXXXX",
            "authorities": [
                "a1-118.akam.net.",
                "a2-64.akam.net.",
                "a6-66.akam.net.",
                "a18-67.akam.net.",
                "a7-64.akam.net.",
                "a11-64.akam.net."
            ]
        }
    ]
}`,
			expectedResponse: []string{"a1-118.akam.net.", "a2-64.akam.net.", "a6-66.akam.net.", "a18-67.akam.net.", "a7-64.akam.net.", "a11-64.akam.net."},
			expectedPath:     "/config-dns/v2/data/authorities?contractIds=9-9XXXXX",
		},
		"test with missing arguments": {
			contractID:   "",
			expectedPath: "/config-dns/v2/data/authorities?contractIds=9-9XXXXX",
			withError:    ErrBadRequest,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetNameServerRecordList(context.Background(), test.contractID)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDns_NewAuthorityResponse(t *testing.T) {
	client := Client(session.Must(session.New()))

	resp := client.NewAuthorityResponse(context.Background(), "empty")

	assert.NotNil(t, resp)
}
