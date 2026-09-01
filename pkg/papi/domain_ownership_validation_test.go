package papi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapiValidateDomainsOwnership(t *testing.T) {
	tests := map[string]struct {
		params              ValidateDomainsOwnershipRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *ValidateDomainsOwnershipResponse
		withError           func(*testing.T, error)
	}{
		"200 OK with refreshToken set to true": {
			params: ValidateDomainsOwnershipRequest{
				RefreshToken: true,
				Body:         ValidateDomainsOwnershipRequestBody{Hostnames: []string{"example.com", "www.example.com"}},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "accid1",
    "hostnames":  [
		{
			"hostname": "example.com",
			"domainValidationStatus": "VALIDATION_IN_PROGRESS",
			"validationScope": "HOST",
			"challengeTokenExpiryDate": "2024-12-31T23:59:59Z",
			"validationCname": {
				"hostname": "akamai-challenge.example.com",
				"target": "example.com.akamai-domain-validation.com"
			},
			"validationTxt": {
				"hostname": "akamai-challenge-txt.example.com",
				"challengeToken": "token123"
			},
			"validationHttp": {
				"fileContentMethod": {
					"url": "http://example.com/.well-known/akamai-challenge",
					"body": "validation-body"
				},
				"redirectMethod": {
					"httpRedirectFrom": "http://example.com/.well-known/redirect",
					"httpRedirectTo": "http://validation.example.com"
				}
			}
		},
		{
			"hostname": "www.example.com",
			"domainValidationStatus": "VALIDATED"
		}
	]
}`,
			expectedPath:        "/papi/v1/domain-challenges?refreshToken=true",
			expectedRequestBody: `{"hostnames":["example.com","www.example.com"]}`,
			expectedResponse: &ValidateDomainsOwnershipResponse{
				AccountID: "accid1",
				Hostnames: []HostnameValidationDetails{
					{
						Hostname:                 "example.com",
						DomainValidationStatus:   "VALIDATION_IN_PROGRESS",
						ChallengeTokenExpiryDate: ptr.To(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
						ValidationScope:          ptr.To("HOST"),
						ValidationCname: &ValidationCname{
							Hostname: "akamai-challenge.example.com",
							Target:   "example.com.akamai-domain-validation.com",
						},
						ValidationTXT: &ValidationTXT{
							Hostname:       "akamai-challenge-txt.example.com",
							ChallengeToken: "token123",
						},
						ValidationHTTP: &ValidationHTTP{
							FileContentMethod: FileContentMethod{
								URL:  "http://example.com/.well-known/akamai-challenge",
								Body: "validation-body",
							},
							RedirectMethod: RedirectMethod{
								HTTPRedirectFrom: "http://example.com/.well-known/redirect",
								HTTPRedirectTo:   "http://validation.example.com",
							},
						},
					},
					{
						Hostname:               "www.example.com",
						DomainValidationStatus: "VALIDATED",
					},
				},
			},
		},
		"200 OK with refreshToken set to false": {
			params: ValidateDomainsOwnershipRequest{
				RefreshToken: false,
				Body:         ValidateDomainsOwnershipRequestBody{Hostnames: []string{"example.com", "www.example.com"}},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "accountId": "accid1",
    "hostnames":  [
		{
			"hostname": "example.com",
			"domainValidationStatus": "VALIDATION_IN_PROGRESS",
			"validationScope": "HOST",
			"challengeTokenExpiryDate": "2024-12-31T23:59:59Z",
			"validationCname": {
				"hostname": "akamai-challenge.example.com",
				"target": "example.com.akamai-domain-validation.com"
			},
			"validationTxt": {
				"hostname": "akamai-challenge-txt.example.com",
				"challengeToken": "token123"
			},
			"validationHttp": {
				"fileContentMethod": {
					"url": "http://example.com/.well-known/akamai-challenge",
					"body": "validation-body"
				},
				"redirectMethod": {
					"httpRedirectFrom": "http://example.com/.well-known/redirect",
					"httpRedirectTo": "http://validation.example.com"
				}
			}
		},
		{
			"hostname": "www.example.com",
			"domainValidationStatus": "VALIDATED"
		}
	]
}`,
			expectedPath:        "/papi/v1/domain-challenges?refreshToken=false",
			expectedRequestBody: `{"hostnames":["example.com","www.example.com"]}`,
			expectedResponse: &ValidateDomainsOwnershipResponse{
				AccountID: "accid1",
				Hostnames: []HostnameValidationDetails{
					{
						Hostname:                 "example.com",
						DomainValidationStatus:   "VALIDATION_IN_PROGRESS",
						ChallengeTokenExpiryDate: ptr.To(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
						ValidationScope:          ptr.To("HOST"),
						ValidationCname: &ValidationCname{
							Hostname: "akamai-challenge.example.com",
							Target:   "example.com.akamai-domain-validation.com",
						},
						ValidationTXT: &ValidationTXT{
							Hostname:       "akamai-challenge-txt.example.com",
							ChallengeToken: "token123",
						},
						ValidationHTTP: &ValidationHTTP{
							FileContentMethod: FileContentMethod{
								URL:  "http://example.com/.well-known/akamai-challenge",
								Body: "validation-body",
							},
							RedirectMethod: RedirectMethod{
								HTTPRedirectFrom: "http://example.com/.well-known/redirect",
								HTTPRedirectTo:   "http://validation.example.com",
							},
						},
					},
					{
						Hostname:               "www.example.com",
						DomainValidationStatus: "VALIDATED",
					},
				},
			},
		},
		"validation error - empty hostnames": {
			params: ValidateDomainsOwnershipRequest{
				Body: ValidateDomainsOwnershipRequestBody{Hostnames: []string{}},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "validating domains ownership: struct validation: Hostnames: cannot be blank")
			},
		},
		"validation error - nil hostnames": {
			params: ValidateDomainsOwnershipRequest{
				Body: ValidateDomainsOwnershipRequestBody{Hostnames: nil},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "validating domains ownership: struct validation: Hostnames: cannot be blank")
			},
		},
		"validation error - empty string in hostnames": {
			params: ValidateDomainsOwnershipRequest{
				Body: ValidateDomainsOwnershipRequestBody{Hostnames: []string{"example.com", ""}},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "validating domains ownership: struct validation: 1: cannot be blank")
			},
		},
		"500 internal server error": {
			params: ValidateDomainsOwnershipRequest{
				Body: ValidateDomainsOwnershipRequestBody{Hostnames: []string{"example.com"}},
			},
			responseStatus:      http.StatusInternalServerError,
			expectedPath:        "/papi/v1/domain-challenges?refreshToken=false",
			expectedRequestBody: `{"hostnames":["example.com"]}`,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error validating domains",
    "status": 500
}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error validating domains",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.ErrorIs(t, err, ErrValidateDomainsOwnership)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)

				if len(test.expectedRequestBody) > 0 {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}

				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ValidateDomainsOwnership(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
