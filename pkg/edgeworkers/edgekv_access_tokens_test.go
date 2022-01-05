package edgeworkers

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEdgeKVAccessToken(t *testing.T) {
	tests := map[string]struct {
		params              CreateEdgeKVAccessTokenRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *CreateEdgeKVAccessTokenResponse
		withError           error
	}{
		"200 OK - create token": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: false,
				AllowOnStaging:    true,
				Expiry:            "2022-03-30",
				Name:              "devexp-token-1",
				NamespacePermissions: NamespacePermissions{
					"default":            []Permission{"r", "w", "d"},
					"devexp-jsmith-test": []Permission{"r", "w"},
				},
			},
			expectedRequestBody: `{"allowOnProduction":false,"allowOnStaging":true,"expiry":"2022-03-30","name":"devexp-token-1","namespacePermissions":{"default":["r","w","d"],"devexp-jsmith-test":["r","w"]}}`,
			responseStatus:      http.StatusOK,
			responseBody: `
{
    "name": "devexp-token-1",
    "uuid": "1ab0e94b-c47e-568e-ab3e-1921ffcefe0c",
    "expiry": "2022-03-30",
    "value": "eyJ0eXAiOxJKV1QxLCJhbGciOiJSUzI1NiJ9.eyJld2lkcyI6ImFsbCIsInN1YiI6IjUwMCIsIm5hbWVzcGFjZS1kZWZhdWx0IjpbInIiLCJkIiwidyJdLCJjcGMiOiI5NzEwNTIiLCJpc3MiOiJha2FtYWkuY29tL0VkZ2VEQi9QdWxzYXIvdjAuMTEuMCIsIm5hbWVzcGFjZS1kZXZleHAtcm9iZXJ0by10ZXN0IjpbInIiLCJ3Il0sImV4cCI6MTY0ODY4NDc5OSwiZW52IjpbInAiLCJzIl0sImlhdCI6MTY0MDg1ODIzNywianRpIjoiMTBiMGU5NGItYzQ3ZS01NjhlLWFiM2UtMTkyMWZmY2VmZTBjIiwicmVxaWQiOiJha2FtYWkiLCJub2VjbCI6dHJ1ZX0.AZfP-VFqDKNWcu1Or73EFfjG_GBDdJUP81Zs0BnNs_bScc8oyBAEiBjxwEsUxrvRRr7rSu-BxFjiDpxx5DlfbgEwd8H2DFV08cfQFqs7aab4WYLrx4ZweD9Hbg2gGLA-dRAbtSrq_FQKQysOvO2ymPn13E78PvK96t8r4cnN1irXbfyBUOXOE3OVOAKsk-w0Ig7qFDa_4o6YyDMPTpwEQ34T1cVqRYStIVzjSaCwgSfdaQG5qzTzTlFoDzG24tz8YlLgoM5OQf9xgsTsisCOF2jf44VWMu2S0e6MIC5gg7zXx7X2t59Y8TsAd0VqqB37y0AzEXkJblbZUlO9HcGebg"
}`,
			expectedPath: "/edgekv/v1/tokens",
			expectedResponse: &CreateEdgeKVAccessTokenResponse{
				Name:   "devexp-token-1",
				UUID:   "1ab0e94b-c47e-568e-ab3e-1921ffcefe0c",
				Expiry: "2022-03-30",
				Value:  "eyJ0eXAiOxJKV1QxLCJhbGciOiJSUzI1NiJ9.eyJld2lkcyI6ImFsbCIsInN1YiI6IjUwMCIsIm5hbWVzcGFjZS1kZWZhdWx0IjpbInIiLCJkIiwidyJdLCJjcGMiOiI5NzEwNTIiLCJpc3MiOiJha2FtYWkuY29tL0VkZ2VEQi9QdWxzYXIvdjAuMTEuMCIsIm5hbWVzcGFjZS1kZXZleHAtcm9iZXJ0by10ZXN0IjpbInIiLCJ3Il0sImV4cCI6MTY0ODY4NDc5OSwiZW52IjpbInAiLCJzIl0sImlhdCI6MTY0MDg1ODIzNywianRpIjoiMTBiMGU5NGItYzQ3ZS01NjhlLWFiM2UtMTkyMWZmY2VmZTBjIiwicmVxaWQiOiJha2FtYWkiLCJub2VjbCI6dHJ1ZX0.AZfP-VFqDKNWcu1Or73EFfjG_GBDdJUP81Zs0BnNs_bScc8oyBAEiBjxwEsUxrvRRr7rSu-BxFjiDpxx5DlfbgEwd8H2DFV08cfQFqs7aab4WYLrx4ZweD9Hbg2gGLA-dRAbtSrq_FQKQysOvO2ymPn13E78PvK96t8r4cnN1irXbfyBUOXOE3OVOAKsk-w0Ig7qFDa_4o6YyDMPTpwEQ34T1cVqRYStIVzjSaCwgSfdaQG5qzTzTlFoDzG24tz8YlLgoM5OQf9xgsTsisCOF2jf44VWMu2S0e6MIC5gg7zXx7X2t59Y8TsAd0VqqB37y0AzEXkJblbZUlO9HcGebg",
			},
		},
		"at least one allow is required": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: false,
				AllowOnStaging:    false,
				Expiry:            "2022-03-30",
				Name:              "name",
				NamespacePermissions: NamespacePermissions{
					"default":            []Permission{"r", "w", "d"},
					"devexp-jsmith-test": []Permission{"r", "w"},
				},
			},
			withError: ErrStructValidation,
		},
		"missing Name": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: true,
				AllowOnStaging:    true,
				Expiry:            "2022-03-30",
				Name:              "",
				NamespacePermissions: NamespacePermissions{
					"default":            []Permission{"r", "w", "d"},
					"devexp-jsmith-test": []Permission{"r", "w"},
				},
			}, withError: ErrStructValidation,
		},
		"invalid date": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: true,
				AllowOnStaging:    true,
				Expiry:            "30/09/2021",
				Name:              "name",
				NamespacePermissions: NamespacePermissions{
					"default":            []Permission{"r", "w", "d"},
					"devexp-jsmith-test": []Permission{"r", "w"},
				},
			}, withError: ErrStructValidation,
		},
		"invalid permission": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: true,
				AllowOnStaging:    true,
				Expiry:            "2022-03-30",
				Name:              "devexp-token-1",
				NamespacePermissions: NamespacePermissions{
					"default": []Permission{"a", "w", "d"},
				},
			}, withError: ErrStructValidation,
		},
		"empty namespace": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: true,
				AllowOnStaging:    true,
				Expiry:            "2022-03-30",
				Name:              "devexp-token-1",
				NamespacePermissions: NamespacePermissions{
					"": []Permission{"r", "w", "d"},
				},
			}, withError: ErrStructValidation,
		},
		"missing permission": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: true,
				AllowOnStaging:    true,
				Expiry:            "2022-03-30",
				Name:              "devexp-token-1",
				NamespacePermissions: NamespacePermissions{
					"default": []Permission{},
				},
			}, withError: ErrStructValidation,
		},
		"missing NamespacePermissions": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: true,
				AllowOnStaging:    true,
				Expiry:            "2022-03-30",
				Name:              "devexp-token-1",
			}, withError: ErrStructValidation,
		},
		"400 bad request": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: true,
				AllowOnStaging:    true,
				Expiry:            "2022-03-30",
				Name:              "devexp-token-1",
				NamespacePermissions: NamespacePermissions{
					"default": []Permission{"r", "w", "d"},
				},
			},
			responseStatus: http.StatusConflict,
			responseBody: `
{
    "detail": "Invalid permission",
    "errorCode": "EKV_2000",
    "instance": "/edgeKV/error-instances/1f2a46ed-b6e8-4f50-b4db-102e260c1753",
    "status": 400,
    "title": "Bad Request",
    "type": "https://learn.akamai.com",
    "additionalDetail": {
        "requestId": "f60f61cda34a0657"
    }
}`,
			expectedPath: "/edgekv/v1/tokens",
			withError: &Error{
				Detail:    "Invalid permission",
				ErrorCode: "EKV_2000",
				Instance:  "/edgeKV/error-instances/1f2a46ed-b6e8-4f50-b4db-102e260c1753",
				Status:    400,
				Title:     "Bad Request",
				Type:      "https://learn.akamai.com",
				AdditionalDetail: Additional{
					RequestID: "f60f61cda34a0657",
				},
			},
		},
		"401 Not authorized - incorrect credentials": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: true,
				AllowOnStaging:    true,
				Expiry:            "2022-03-30",
				Name:              "devexp-token-1",
				NamespacePermissions: NamespacePermissions{
					"default":            []Permission{"r", "w", "d"},
					"devexp-jsmith-test": []Permission{"r", "w"},
				},
			},
			responseStatus: http.StatusUnauthorized,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
    "title": "Not authorized",
    "status": 401,
    "detail": "Inactive client token",
    "instance": "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/edgekv/v1/tokens",
    "method": "POST",
    "serverIp": "104.81.220.242",
    "clientIp": "22.22.22.22",
    "requestId": "1e7f0f0f",
    "requestTime": "2021-12-30T14:12:50Z"
}`,
			expectedPath: "/edgekv/v1/tokens",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
				Title:       "Not authorized",
				Status:      401,
				Detail:      "Inactive client token",
				Instance:    "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/edgekv/v1/tokens",
				Method:      "POST",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "1e7f0f0f",
				RequestTime: "2021-12-30T14:12:50Z",
			},
		},
		"409 duplicated token name": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: true,
				AllowOnStaging:    true,
				Expiry:            "2022-03-30",
				Name:              "devexp-token-1",
				NamespacePermissions: NamespacePermissions{
					"default":            []Permission{"r", "w", "d"},
					"devexp-jsmith-test": []Permission{"r", "w"},
				},
			},
			responseStatus: http.StatusConflict,
			responseBody: `
{
    "detail": "Token with name devexp-token-1 is already stored.",
    "errorCode": "EKV_3000",
    "instance": "/edgeKV/error-instances/e82edcd9-498e-42f9-a078-6d9c4f9dbcb9",
    "status": 409,
    "title": "Conflict",
    "type": "https://learn.akamai.com",
    "additionalDetail": {
        "requestId": "bc7561cda1f3021b"
    }
}`,
			expectedPath: "/edgekv/v1/tokens",
			withError: &Error{
				Detail:    "Token with name devexp-token-1 is already stored.",
				ErrorCode: "EKV_3000",
				Instance:  "/edgeKV/error-instances/e82edcd9-498e-42f9-a078-6d9c4f9dbcb9",
				Status:    409,
				Title:     "Conflict",
				Type:      "https://learn.akamai.com",
				AdditionalDetail: Additional{
					RequestID: "bc7561cda1f3021b",
				},
			},
		},
		"500 internal server error": {
			params: CreateEdgeKVAccessTokenRequest{
				AllowOnProduction: true,
				AllowOnStaging:    true,
				Expiry:            "2022-03-30",
				Name:              "devexp-token-1",
				NamespacePermissions: NamespacePermissions{
					"default":            []Permission{"r", "w", "d"},
					"devexp-jsmith-test": []Permission{"r", "w"},
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "detail": "An internal error occurred.",
    "errorCode": "EKV_0000",
    "instance": "/edgeKV/error-instances/e9bc19b5-ec1e-485d-80d0-20237a928684",
    "status": 500,
    "title": "Internal Server Error",
    "type": "https://learn.akamai.com",
    "additionalDetail": {
        "requestId": "b2f461d47426558c"
    }
}`,
			expectedPath: "/edgekv/v1/tokens",
			withError: &Error{
				Detail:    "An internal error occurred.",
				ErrorCode: "EKV_0000",
				Instance:  "/edgeKV/error-instances/e9bc19b5-ec1e-485d-80d0-20237a928684",
				Status:    500,
				Title:     "Internal Server Error",
				Type:      "https://learn.akamai.com",
				AdditionalDetail: Additional{
					RequestID: "b2f461d47426558c",
				},
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
			result, err := client.CreateEdgeKVAccessToken(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetEdgeKVAccessToken(t *testing.T) {
	tests := map[string]struct {
		params           GetEdgeKVAccessTokenRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetEdgeKVAccessTokenResponse
		withError        error
	}{
		"200 OK - get token": {
			params: GetEdgeKVAccessTokenRequest{
				TokenName: "devexp-token-1",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "name": "devexp-token-1",
    "uuid": "10b0e94b-c47e-568e-ab3e-1921ffcefe0c",
    "expiry": "2022-03-30",
    "value": "eyJ0eXAxOxJKV1QxLCJhbGciOiJSUzI1NiJ9.eyJld2lkcyI6ImFsbCIsInN1YiI6IjUwMCIsIm5hbWVzcGFjZS1kZWZhdWx0IjpbInIiLCJkIiwidyJdLCJjcGMiOiI5NzEwNTIiLCJpc3MiOiJha2FtYWkuY29tL0VkZ2VEQi9QdWxzYXIvdjAuMTEuMCIsIm5hbWVzcGFjZS1kZXZleHAtcm9iZXJ0by10ZXN0IjpbInIiLCJ3Il0sImV4cCI6MTY0ODY4NDc5OSwiZW52IjpbInAiLCJzIl0sImlhdCI6MTY0MDg1ODIzNywianRpIjoiMTBiMGU5NGItYzQ3ZS01NjhlLWFiM2UtMTkyMWZmY2VmZTBjIiwicmVxaWQiOiJha2FtYWkiLCJub2VjbCI6dHJ1ZX0.AZfP-VFqDKNWcu1Or73EFfjG_GBDdJUP81Zs0BnNs_bScc8oyBAEiBjxwEsUxrvRRr7rSu-BxFjiDpxx5DlfbgEwd8H2DFV08cfQFqs7aab4WYLrx4ZweD9Hbg2gGLA-dRAbtSrq_FQKQysOvO2ymPn13E78PvK96t8r4cnN1irXbfyBUOXOE3OVOAKsk-w0Ig7qFDa_4o6YyDMPTpwEQ34T1cVqRYStIVzjSaCwgSfdaQG5qzTzTlFoDzG24tz8YlLgoM5OQf9xgsTsisCOF2jf44VWMu2S0e6MIC5gg7zXx7X2t59Y8TsAd0VqqB37y0AzEXkJblbZUlO9HcGebg"
}`,
			expectedPath: "/edgekv/v1/tokens/devexp-token-1",
			expectedResponse: &GetEdgeKVAccessTokenResponse{
				Name:   "devexp-token-1",
				UUID:   "10b0e94b-c47e-568e-ab3e-1921ffcefe0c",
				Expiry: "2022-03-30",
				Value:  "eyJ0eXAxOxJKV1QxLCJhbGciOiJSUzI1NiJ9.eyJld2lkcyI6ImFsbCIsInN1YiI6IjUwMCIsIm5hbWVzcGFjZS1kZWZhdWx0IjpbInIiLCJkIiwidyJdLCJjcGMiOiI5NzEwNTIiLCJpc3MiOiJha2FtYWkuY29tL0VkZ2VEQi9QdWxzYXIvdjAuMTEuMCIsIm5hbWVzcGFjZS1kZXZleHAtcm9iZXJ0by10ZXN0IjpbInIiLCJ3Il0sImV4cCI6MTY0ODY4NDc5OSwiZW52IjpbInAiLCJzIl0sImlhdCI6MTY0MDg1ODIzNywianRpIjoiMTBiMGU5NGItYzQ3ZS01NjhlLWFiM2UtMTkyMWZmY2VmZTBjIiwicmVxaWQiOiJha2FtYWkiLCJub2VjbCI6dHJ1ZX0.AZfP-VFqDKNWcu1Or73EFfjG_GBDdJUP81Zs0BnNs_bScc8oyBAEiBjxwEsUxrvRRr7rSu-BxFjiDpxx5DlfbgEwd8H2DFV08cfQFqs7aab4WYLrx4ZweD9Hbg2gGLA-dRAbtSrq_FQKQysOvO2ymPn13E78PvK96t8r4cnN1irXbfyBUOXOE3OVOAKsk-w0Ig7qFDa_4o6YyDMPTpwEQ34T1cVqRYStIVzjSaCwgSfdaQG5qzTzTlFoDzG24tz8YlLgoM5OQf9xgsTsisCOF2jf44VWMu2S0e6MIC5gg7zXx7X2t59Y8TsAd0VqqB37y0AzEXkJblbZUlO9HcGebg",
			},
		},
		"missing token name": {
			params:    GetEdgeKVAccessTokenRequest{},
			withError: ErrStructValidation,
		},
		"403 Forbidden - incorrect credentials": {
			params: GetEdgeKVAccessTokenRequest{
				TokenName: "devexp-token-1",
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
    "title": "Not authorized",
    "status": 401,
    "detail": "Inactive client token",
    "instance": "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/edgekv/v1/tokens/devexp-token-99",
    "method": "GET",
    "serverIp": "104.81.220.242",
    "clientIp": "22.22.22.22",
    "requestId": "cb5cd20",
    "requestTime": "2022-01-03T07:46:28Z"
}`,
			expectedPath: "/edgekv/v1/tokens/devexp-token-1",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
				Title:       "Not authorized",
				Status:      401,
				Detail:      "Inactive client token",
				Instance:    "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/edgekv/v1/tokens/devexp-token-99",
				Method:      "GET",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "cb5cd20",
				RequestTime: "2022-01-03T07:46:28Z",
			},
		},
		"404 Not Found - Token doesn't exist": {
			params: GetEdgeKVAccessTokenRequest{
				TokenName: "devexp-token-99",
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "detail": "Token with name devexp-token-99 does not exist.",
    "errorCode": "EKV_3000",
    "instance": "/edgeKV/error-instances/add4ab5a-48b0-4350-aa8b-7f64e9b6a5ea",
    "status": 404,
    "title": "Not Found",
    "type": "https://learn.akamai.com",
    "additionalDetail": {
        "requestId": "ae9061cddea87d94"
    }
}`,
			expectedPath: "/edgekv/v1/tokens/devexp-token-99",
			withError: &Error{
				Detail:    "Token with name devexp-token-99 does not exist.",
				ErrorCode: "EKV_3000",
				Instance:  "/edgeKV/error-instances/add4ab5a-48b0-4350-aa8b-7f64e9b6a5ea",
				Status:    404,
				Title:     "Not Found",
				Type:      "https://learn.akamai.com",
				AdditionalDetail: Additional{
					RequestID: "ae9061cddea87d94",
				},
			},
		},
		"500 Internal server error": {
			params: GetEdgeKVAccessTokenRequest{
				TokenName: ";",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
    "title": "Server Error",
    "status": 500,
    "instance": "https://akaa-7udtftgmvpnmsbwx-noxd5uwfehzxv4rj.luna-dev.akamaiapis.net/edgekv/v1/tokens/;",
    "method": "GET",
    "serverIp": "104.81.220.242",
    "clientIp": "22.22.22.22",
    "requestId": "e98b01a",
    "requestTime": "2022-01-03T11:13:00Z"
}`,
			expectedPath: "/edgekv/v1/tokens/;",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "https://akaa-7udtftgmvpnmsbwx-noxd5uwfehzxv4rj.luna-dev.akamaiapis.net/edgekv/v1/tokens/;",
				Method:      "GET",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "e98b01a",
				RequestTime: "2022-01-03T11:13:00Z",
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
			result, err := client.GetEdgeKVAccessToken(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListEdgeKVAccessTokens(t *testing.T) {
	tests := map[string]struct {
		params           ListEdgeKVAccessTokensRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListEdgeKVAccessTokensResponse
		withError        error
	}{
		"200 OK - list EdgeKV tokens": {
			responseStatus: http.StatusOK,
			responseBody: `
{
    "tokens": [
        {
            "name": "my_token",
            "uuid": "8301fef4-80e5-5efb-9bfb-8f5869a5df7b",
            "expiry": "2022-03-30"
        },
        {
            "name": "token1",
            "uuid": "5b5d3bfb-8d2e-5fbb-858d-33807edc9554",
            "expiry": "2022-01-22"
        },
        {
            "name": "token2",
            "uuid": "62181cfe-268a-5302-8834-67c67ec86efd",
            "expiry": "2022-01-22"
        },
        {
            "name": "token3",
            "uuid": "edb02678-ae1c-564c-8f73-c977ffdfe016",
            "expiry": "2022-01-22"
        }
    ]
}`,
			expectedPath: "/edgekv/v1/tokens",
			expectedResponse: &ListEdgeKVAccessTokensResponse{
				[]EdgeKVAccessToken{
					{
						Name:   "my_token",
						UUID:   "8301fef4-80e5-5efb-9bfb-8f5869a5df7b",
						Expiry: "2022-03-30",
					},
					{
						Name:   "token1",
						UUID:   "5b5d3bfb-8d2e-5fbb-858d-33807edc9554",
						Expiry: "2022-01-22",
					},
					{
						Name:   "token2",
						UUID:   "62181cfe-268a-5302-8834-67c67ec86efd",
						Expiry: "2022-01-22",
					},
					{
						Name:   "token3",
						UUID:   "edb02678-ae1c-564c-8f73-c977ffdfe016",
						Expiry: "2022-01-22",
					},
				},
			},
		},
		"200 OK - list EdgeKV tokens including expired": {
			params: ListEdgeKVAccessTokensRequest{
				IncludeExpired: true,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "tokens": [
        {
            "name": "my_token",
            "uuid": "8301fef4-80e5-5efb-9bfb-8f5869a5df7b",
            "expiry": "2022-03-30"
        },
        {
            "name": "token1",
            "uuid": "5b5d3bfb-8d2e-5fbb-858d-33807edc9554",
            "expiry": "2022-01-22"
        },
        {
            "name": "token2",
            "uuid": "62181cfe-268a-5302-8834-67c67ec86efd",
            "expiry": "2022-01-22"
        },
        {
            "name": "token3",
            "uuid": "edb02678-ae1c-564c-8f73-c977ffdfe016",
            "expiry": "2022-01-22"
        },
        {
            "name": "preexistingTokenTest",
            "uuid": "7a14da8c-1709-570b-9535-2cc6e2ee5a8a",
            "expiry": "2021-12-21"
        }
    ]
}`,
			expectedPath: "/edgekv/v1/tokens?includeExpired=true",
			expectedResponse: &ListEdgeKVAccessTokensResponse{
				[]EdgeKVAccessToken{
					{
						Name:   "my_token",
						UUID:   "8301fef4-80e5-5efb-9bfb-8f5869a5df7b",
						Expiry: "2022-03-30",
					},
					{
						Name:   "token1",
						UUID:   "5b5d3bfb-8d2e-5fbb-858d-33807edc9554",
						Expiry: "2022-01-22",
					},
					{
						Name:   "token2",
						UUID:   "62181cfe-268a-5302-8834-67c67ec86efd",
						Expiry: "2022-01-22",
					},
					{
						Name:   "token3",
						UUID:   "edb02678-ae1c-564c-8f73-c977ffdfe016",
						Expiry: "2022-01-22",
					},
					{
						Name:   "preexistingTokenTest",
						UUID:   "7a14da8c-1709-570b-9535-2cc6e2ee5a8a",
						Expiry: "2021-12-21",
					},
				},
			},
		},
		"401 Forbidden - incorrect credentials": {
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
    "title": "Not authorized",
    "status": 401,
    "detail": "Inactive client token",
    "instance": "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/edgekv/v1/tokens",
    "method": "GET",
    "serverIp": "104.81.220.242",
    "clientIp": "22.22.22.22",
    "requestId": "d64edd6",
    "requestTime": "2022-01-03T09:01:30Z"
}`,
			expectedPath: "/edgekv/v1/tokens",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
				Title:       "Not authorized",
				Status:      401,
				Detail:      "Inactive client token",
				Instance:    "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/edgekv/v1/tokens",
				Method:      "GET",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "d64edd6",
				RequestTime: "2022-01-03T09:01:30Z",
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
    "title": "Server Error",
    "status": 500,
    "instance": "https://akaa-7udtftgmvpnmsbwx-noxd5uwfehzxv4rj.luna-dev.akamaiapis.net/edgekv/v1/tokens/;",
    "method": "GET",
    "serverIp": "104.81.220.242",
    "clientIp": "22.22.22.22",
    "requestId": "e98b01a",
    "requestTime": "2022-01-03T11:13:00Z"
}`,
			expectedPath: "/edgekv/v1/tokens",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "https://akaa-7udtftgmvpnmsbwx-noxd5uwfehzxv4rj.luna-dev.akamaiapis.net/edgekv/v1/tokens/;",
				Method:      "GET",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "e98b01a",
				RequestTime: "2022-01-03T11:13:00Z",
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
			result, err := client.ListEdgeKVAccessTokens(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeleteEdgeKVAccessToken(t *testing.T) {
	tests := map[string]struct {
		params           DeleteEdgeKVAccessTokenRequest
		withError        error
		expectedPath     string
		responseStatus   int
		responseBody     string
		expectedResponse *DeleteEdgeKVAccessTokenResponse
	}{
		"200 Deleted": {
			params: DeleteEdgeKVAccessTokenRequest{
				TokenName: "devexp-token-3",
			},
			expectedPath: "/edgekv/v1/tokens/devexp-token-3",
			responseBody: `
{
    "name": "devexp-token-3",
    "uuid": "cc0a9045-e654-5f17-9b37-6ab6e565803f"
}`,
			expectedResponse: &DeleteEdgeKVAccessTokenResponse{
				Name: "devexp-token-3",
				UUID: "cc0a9045-e654-5f17-9b37-6ab6e565803f",
			},
			responseStatus: http.StatusOK,
		},
		"missing token name": {
			params: DeleteEdgeKVAccessTokenRequest{
				TokenName: "",
			},
			withError: ErrStructValidation,
		},
		"401 not authorized": {
			responseStatus: http.StatusUnauthorized,
			params: DeleteEdgeKVAccessTokenRequest{
				TokenName: "devexp-token-99",
			},
			expectedPath: "/edgekv/v1/tokens/devexp-token-99",
			responseBody: `{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
    "title": "Not authorized",
    "status": 401,
    "detail": "Inactive client token",
    "instance": "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/edgekv/v1/tokens/devexp-token-3",
    "method": "DELETE",
    "serverIp": "104.81.220.242",
    "clientIp": "22.22.22.22",
    "requestId": "ddc683c",
    "requestTime": "2022-01-03T09:51:55Z"
}`,
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
				Title:       "Not authorized",
				Status:      401,
				Detail:      "Inactive client token",
				Instance:    "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/edgekv/v1/tokens/devexp-token-3",
				Method:      "DELETE",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "ddc683c",
				RequestTime: "2022-01-03T09:51:55Z",
			},
		},
		"404 Not Found": {
			params: DeleteEdgeKVAccessTokenRequest{
				TokenName: "devexp-token-99",
			},
			expectedPath:   "/edgekv/v1/tokens/devexp-token-99",
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "detail": "Token with name devexp-token-99 does not exist.",
    "errorCode": "EKV_3000",
    "instance": "/edgeKV/error-instances/d4d7171f-2ef9-4e60-96ba-1ad74e35bb39",
    "status": 404,
    "title": "Not Found",
    "type": "https://learn.akamai.com",
    "additionalDetail": {
        "requestId": "a46f61d2c9539c77"
    }
}`,
			withError: &Error{
				Detail:    "Token with name devexp-token-99 does not exist.",
				ErrorCode: "EKV_3000",
				Instance:  "/edgeKV/error-instances/d4d7171f-2ef9-4e60-96ba-1ad74e35bb39",
				Status:    404,
				Title:     "Not Found",
				Type:      "https://learn.akamai.com",
				AdditionalDetail: Additional{
					RequestID: "a46f61d2c9539c77",
				},
			},
		},
		"500 internal server error": {
			params: DeleteEdgeKVAccessTokenRequest{
				TokenName: ";",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
    "title": "Server Error",
    "status": 500,
    "instance": "https://akaa-7udtftgmvpnmsbwx-noxd5uwfehzxv4rj.luna-dev.akamaiapis.net/edgekv/v1/tokens/;",
    "method": "DELETE",
    "serverIp": "104.81.220.242",
    "clientIp": "22.22.22.22",
    "requestId": "e6f4e86",
    "requestTime": "2022-01-03T10:55:00Z"
}`,
			expectedPath: "/edgekv/v1/tokens/;",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "https://akaa-7udtftgmvpnmsbwx-noxd5uwfehzxv4rj.luna-dev.akamaiapis.net/edgekv/v1/tokens/;",
				Method:      "DELETE",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "e6f4e86",
				RequestTime: "2022-01-03T10:55:00Z",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			_, err := client.DeleteEdgeKVAccessToken(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				if test.responseStatus != 0 {
					assert.Contains(t, err.Error(), strconv.FormatInt(int64(test.responseStatus), 10))
				}

				return
			}
			require.NoError(t, err)
		})
	}
}
