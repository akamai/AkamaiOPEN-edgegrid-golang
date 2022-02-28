package edgeworkers

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSecureToken(t *testing.T) {
	tests := map[string]struct {
		params              CreateSecureTokenRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *CreateSecureTokenResponse
		withError           error
	}{
		"201 Created - create secure token": {
			params: CreateSecureTokenRequest{
				ACL:      "/*",
				Expiry:   15,
				Hostname: "test.devexp.akamai.com",
			},
			expectedRequestBody: `{"acl":"/*","expiry":15,"hostname":"test.devexp.akamai.com"}`,
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "akamaiEwTrace": "st=1641295764~exp=1641296664~acl=/*~hmac=f6d18857a6c738664b65a59036ac6f8348abe6b34a9503ec1262f18f114ed43f"
}`,
			expectedPath: "/edgeworkers/v1/secure-token",
			expectedResponse: &CreateSecureTokenResponse{
				AkamaiEWTrace: "st=1641295764~exp=1641296664~acl=/*~hmac=f6d18857a6c738664b65a59036ac6f8348abe6b34a9503ec1262f18f114ed43f",
			},
		},
		"201 Created - create secure token with hostname only": {
			params: CreateSecureTokenRequest{
				Hostname: "test.devexp.akamai.com",
			},
			expectedRequestBody: `{"hostname":"test.devexp.akamai.com"}`,
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "akamaiEwTrace": "st=1641295764~exp=1641296664~acl=/*~hmac=f6d18857a6c738664b65a59036ac6f8348abe6b34a9503ec1262f18f114ed43f"
}`,
			expectedPath: "/edgeworkers/v1/secure-token",
			expectedResponse: &CreateSecureTokenResponse{
				AkamaiEWTrace: "st=1641295764~exp=1641296664~acl=/*~hmac=f6d18857a6c738664b65a59036ac6f8348abe6b34a9503ec1262f18f114ed43f",
			},
		},
		"201 Created - create secure token with hostname and propertyId": {
			params: CreateSecureTokenRequest{
				Hostname:   "test.devexp.akamai.com",
				PropertyID: "200153206",
			},
			expectedRequestBody: `{"hostname":"test.devexp.akamai.com","propertyId":"200153206"}`,
			responseStatus:      http.StatusCreated,
			responseBody: `
{
    "akamaiEwTrace": "st=1641295764~exp=1641296664~acl=/*~hmac=f6d18857a6c738664b65a59036ac6f8348abe6b34a9503ec1262f18f114ed43f"
}`,
			expectedPath: "/edgeworkers/v1/secure-token",
			expectedResponse: &CreateSecureTokenResponse{
				AkamaiEWTrace: "st=1641295764~exp=1641296664~acl=/*~hmac=f6d18857a6c738664b65a59036ac6f8348abe6b34a9503ec1262f18f114ed43f",
			},
		},
		"validation error - empty request": {
			params:    CreateSecureTokenRequest{},
			withError: ErrStructValidation,
		},
		"validation error - both ALC and URL": {
			params: CreateSecureTokenRequest{
				ACL:      "/*",
				Expiry:   15,
				Hostname: "test.devexp.akamai.com",
				URL:      "/",
			},
			withError: ErrStructValidation,
		},
		"validation error - invalid expiry": {
			params: CreateSecureTokenRequest{
				ACL:      "/*",
				Expiry:   1440,
				Hostname: "test.devexp.akamai.com",
			},
			withError: ErrStructValidation,
		},
		"401 unauthorized": {
			params: CreateSecureTokenRequest{
				ACL:      "/*",
				Expiry:   15,
				Hostname: "test.devexp.akamai.com",
				Network:  "STAGING",
			},
			responseStatus: http.StatusUnauthorized,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
    "title": "Not authorized",
    "status": 401,
    "detail": "Inactive client token",
    "instance": "https://akaa-mgfkwp3rw4k2whym-eyn4wdjeur5lz37d.luna-dev.akamaiapis.net/edgeworkers/v1/secure-token",
    "method": "POST",
    "serverIp": "104.81.220.242",
    "clientIp": "22.22.22.22",
    "requestId": "17f6b2bc",
    "requestTime": "2022-01-04T10:31:23Z"
}`,
			expectedPath: "/edgeworkers/v1/secure-token",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
				Title:       "Not authorized",
				Status:      401,
				Detail:      "Inactive client token",
				Instance:    "https://akaa-mgfkwp3rw4k2whym-eyn4wdjeur5lz37d.luna-dev.akamaiapis.net/edgeworkers/v1/secure-token",
				Method:      "POST",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "17f6b2bc",
				RequestTime: "2022-01-04T10:31:23Z",
			},
		},
		"403 Forbidden - incorrect credentials": {
			params: CreateSecureTokenRequest{
				ACL:      "/*",
				Expiry:   15,
				Hostname: "test.devexp.akamai.com",
				Network:  "STAGING",
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
    "title": "Forbidden",
    "status": 403,
    "detail": "The client does not have the grant needed for the request",
    "instance": "https://akaa-xfaqsq2csihdccx5-4osos3xx73uxd2if.luna-dev.akamaiapis.net/edgeworkers/v1/secure-token",
    "authzRealm": "b7iuwfuwdvstkoil.dhxzzfwdsq2jgp7w",
    "method": "POST",
    "serverIp": "104.81.220.242",
    "clientIp": "22.22.22.22",
    "requestId": "1801a12b",
    "requestTime": "2022-01-04T10:36:06Z"
}`,
			expectedPath: "/edgeworkers/v1/secure-token",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "https://akaa-xfaqsq2csihdccx5-4osos3xx73uxd2if.luna-dev.akamaiapis.net/edgeworkers/v1/secure-token",
				AuthzRealm:  "b7iuwfuwdvstkoil.dhxzzfwdsq2jgp7w",
				Method:      "POST",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "1801a12b",
				RequestTime: "2022-01-04T10:36:06Z",
			},
		},
		"404 Not found": {
			params: CreateSecureTokenRequest{
				ACL:      "/*",
				Expiry:   15,
				Hostname: "some1.test",
				Network:  "STAGING",
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/edgeworkers/error-types/secret-key-not-found",
    "title": "Rest API Error",
    "instance": "eb764a5e-f375-4959-9e4d-b3a70d28721d",
    "status": 404,
    "detail": "Secret key could not be found.",
    "errorCode": "EW2301"
}`,
			expectedPath: "/edgeworkers/v1/secure-token",
			withError: &Error{
				Type:      "/edgeworkers/error-types/secret-key-not-found",
				Title:     "Rest API Error",
				Instance:  "eb764a5e-f375-4959-9e4d-b3a70d28721d",
				Status:    404,
				Detail:    "Secret key could not be found.",
				ErrorCode: "EW2301",
			},
		},
		"500 internal server error": {
			params: CreateSecureTokenRequest{
				ACL:      "/*",
				Expiry:   15,
				Hostname: "test.devexp.akamai.com",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
  "title": "Server Error",
  "status": 500,
  "instance": "host_name/edgeworkers/v1/secure-token",
  "method": "POST",
  "serverIp": "104.81.220.111",
  "clientIp": "89.64.55.111",
  "requestId": "a73affa111",
  "requestTime": "2021-12-13T13:32:37Z"
}`,
			expectedPath: "/edgeworkers/v1/secure-token",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/secure-token",
				Method:      "POST",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2021-12-13T13:32:37Z",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateSecureToken(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
