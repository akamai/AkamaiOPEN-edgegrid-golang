package apidefinitions

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyVersion(t *testing.T) {
	tests := map[string]struct {
		request             VerifyVersionRequest
		responseBody        string
		responseStatus      int
		withError           error
		expectedPath        string
		expectedRequestBody string
		expectedResponse    VerifyVersionResponse
	}{
		"200 ok": {
			request: VerifyVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: VerifyVersionRequestBody{
					Networks: []NetworkType{ActivationNetworkStaging},
				},
			},
			responseBody: `[
	{
		"severity":"ERROR",
 		"detail": "You shall not pass"
	},
	{
		"severity":"WARNING",
 		"detail": "You shall not pass"
	}
]`,
			responseStatus: http.StatusOK,
			expectedRequestBody: `
{
    "networks": [
        "STAGING"
    ]
}`,
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/activate/verify",
			expectedResponse: VerifyVersionResponse{
				{
					Severity: SeverityError,
					Detail:   "You shall not pass",
				},
				{
					Severity: SeverityWarning,
					Detail:   "You shall not pass",
				},
			},
		},
		"404 not found": {
			request: VerifyVersionRequest{
				VersionNumber: 4,
				APIEndpointID: 987,
				Body: VerifyVersionRequestBody{
					Networks: []NetworkType{ActivationNetworkStaging},
				},
			},
			responseBody: `
{
    "type": "https://problems.luna.akamaiapis.net/api-definitions/error-types/NOT-FOUND",
    "status": 404,
    "title": "Not Found",
    "instance": "22f95e07-00f8-4dfc-a903-ba9b1e2999bf",
    "detail": "Invalid version provided.",
    "severity": "ERROR"
}`,
			responseStatus: http.StatusNotFound,
			expectedPath:   "/api-definitions/v2/endpoints/987/versions/4/activate/verify",
			withError: &Error{
				Type:     "https://problems.luna.akamaiapis.net/api-definitions/error-types/NOT-FOUND",
				Status:   http.StatusNotFound,
				Title:    "Not Found",
				Instance: "22f95e07-00f8-4dfc-a903-ba9b1e2999bf",
				Detail:   "Invalid version provided.",
				Severity: ptr.To("ERROR"),
			},
		},
		"incorrect network": {
			request: VerifyVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: VerifyVersionRequestBody{
					Networks: []NetworkType{"TYPO"},
				},
			},
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/activate",
			withError:    ErrStructValidation,
		},
		"missing network": {
			request: VerifyVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body:          VerifyVersionRequestBody{},
			},
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/activate",
			withError:    ErrStructValidation,
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
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.VerifyVersion(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestActivateVersion(t *testing.T) {
	tests := map[string]struct {
		request             ActivateVersionRequest
		responseBody        string
		responseStatus      int
		withError           error
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *ActivateVersionResponse
	}{
		"200 ok": {
			request: ActivateVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks:               []NetworkType{ActivationNetworkStaging},
					Notes:                  "some notes",
					NotificationRecipients: []string{"devteam@domain.com"},
				},
			},
			responseBody: `
{
    "notes": "some notes",
    "networks": [
        "STAGING"
    ],
    "notificationRecipients": [
        "devteam@domain.com"
    ]
}`,
			responseStatus: http.StatusOK,
			expectedRequestBody: `
{
    "notes": "some notes",
    "networks": [
        "STAGING"
    ],
    "notificationRecipients": [
        "devteam@domain.com"
    ]
}`,
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/activate",
			expectedResponse: &ActivateVersionResponse{
				Networks:               []NetworkType{ActivationNetworkStaging},
				Notes:                  "some notes",
				NotificationRecipients: []string{"devteam@domain.com"},
			},
		},
		"200 only mandatory fields": {
			request: ActivateVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks: []NetworkType{ActivationNetworkProduction},
				},
			},
			responseBody: `
{
    "networks": [
        "PRODUCTION"
    ],
	"notes": null,
	"notificationRecipients": null
}`,
			responseStatus: http.StatusOK,
			expectedRequestBody: `
{
    "networks": [
        "PRODUCTION"
    ]
}`,
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/activate",
			expectedResponse: &ActivateVersionResponse{
				Networks: []NetworkType{ActivationNetworkProduction},
			},
		},
		"400 already active": {
			request: ActivateVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks:               []NetworkType{ActivationNetworkStaging},
					Notes:                  "some notes",
					NotificationRecipients: []string{"devteam@domain.com"},
				},
			},
			responseBody: `
{
    "type": "/api-definitions/error-types/endpoint-version-already-active",
    "status": 400,
    "title": "Version active",
    "instance": "dc9fc264-e64d-4d94-bce2-4ac4093f4435",
    "detail": "You cannot activate the version because it’s already active on STAGING.",
    "severity": "ERROR",
    "endpointId": 987,
    "endpointName": "dxetest",
    "versionNumber": 1,
    "network": "STAGING"
}`,
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/api-definitions/v2/endpoints/987/versions/1/activate",
			withError: &Error{
				Type:          "/api-definitions/error-types/endpoint-version-already-active",
				Status:        http.StatusBadRequest,
				Title:         "Version active",
				Instance:      "dc9fc264-e64d-4d94-bce2-4ac4093f4435",
				Detail:        "You cannot activate the version because it’s already active on STAGING.",
				Severity:      ptr.To("ERROR"),
				EndpointID:    ptr.To(int64(987)),
				EndpointName:  ptr.To("dxetest"),
				VersionNumber: ptr.To(int64(1)),
				Network:       ptr.To("STAGING"),
			},
		},
		"404 not found": {
			request: ActivateVersionRequest{
				VersionNumber: 4,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks:               []NetworkType{ActivationNetworkStaging},
					Notes:                  "some notes",
					NotificationRecipients: []string{"devteam@domain.com"},
				},
			},
			responseBody: `
{
    "type": "https://problems.luna.akamaiapis.net/api-definitions/error-types/NOT-FOUND",
    "status": 404,
    "title": "Not Found",
    "instance": "22f95e07-00f8-4dfc-a903-ba9b1e2999bf",
    "detail": "Invalid version provided.",
    "severity": "ERROR"
}`,
			responseStatus: http.StatusNotFound,
			expectedPath:   "/api-definitions/v2/endpoints/987/versions/4/activate",
			withError: &Error{
				Type:     "https://problems.luna.akamaiapis.net/api-definitions/error-types/NOT-FOUND",
				Status:   http.StatusNotFound,
				Title:    "Not Found",
				Instance: "22f95e07-00f8-4dfc-a903-ba9b1e2999bf",
				Detail:   "Invalid version provided.",
				Severity: ptr.To("ERROR"),
			},
		},
		"incorrect network": {
			request: ActivateVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks:               []NetworkType{"TYPO"},
					Notes:                  "some notes",
					NotificationRecipients: []string{"devteam@domain.com"},
				},
			},
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/activate",
			withError:    ErrStructValidation,
		},
		"missing network": {
			request: ActivateVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Notes:                  "some notes",
					NotificationRecipients: []string{"devteam@domain.com"},
				},
			},
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/activate",
			withError:    ErrStructValidation,
		},
		"incorrect recipient": {
			request: ActivateVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks:               []NetworkType{ActivationNetworkStaging},
					Notes:                  "some notes",
					NotificationRecipients: []string{"devteam"},
				},
			},
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/activate",
			withError:    ErrStructValidation,
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
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ActivateVersion(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeactivateVersion(t *testing.T) {
	tests := map[string]struct {
		request             DeactivateVersionRequest
		responseBody        string
		responseStatus      int
		withError           error
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *DeactivateVersionResponse
	}{
		"200 ok": {
			request: DeactivateVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks:               []NetworkType{ActivationNetworkStaging},
					Notes:                  "some notes",
					NotificationRecipients: []string{"devteam@domain.com"},
				},
			},
			responseBody: `
{
    "notes": "some notes",
    "networks": [
        "STAGING"
    ],
    "notificationRecipients": [
        "devteam@domain.com"
    ]
}`,
			responseStatus: http.StatusOK,
			expectedRequestBody: `
{
    "notes": "some notes",
    "networks": [
        "STAGING"
    ],
    "notificationRecipients": [
        "devteam@domain.com"
    ]
}`,
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/deactivate",
			expectedResponse: &DeactivateVersionResponse{
				Networks:               []NetworkType{ActivationNetworkStaging},
				Notes:                  "some notes",
				NotificationRecipients: []string{"devteam@domain.com"},
			},
		},
		"200 only mandatory fields": {
			request: DeactivateVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks: []NetworkType{ActivationNetworkProduction},
				},
			},
			responseBody: `
{
	"notes": null,
    "networks": [
        "PRODUCTION"
    ],
	"notificationRecipients": null
}`,
			responseStatus: http.StatusOK,
			expectedRequestBody: `
{
    "networks": [
        "PRODUCTION"
    ]
}`,
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/deactivate",
			expectedResponse: &DeactivateVersionResponse{
				Networks: []NetworkType{ActivationNetworkProduction},
			},
		},
		"400 not active": {
			request: DeactivateVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks:               []NetworkType{ActivationNetworkStaging},
					Notes:                  "some notes",
					NotificationRecipients: []string{"devteam@domain.com"},
				},
			},
			responseBody: `
{
	"type": "/api-definitions/error-types/endpoint-version-not-active",
    "status": 400,
    "title": "Version not active",
    "instance": "3692ad83-2ac3-4912-aca0-a5b93eac84c2",
    "detail": "You cannot deactivate the version because it's not active on STAGING.",
    "severity": "ERROR",
    "endpointId": 987,
    "endpointName": "dxetest",
    "versionNumber": 1,
    "network": "STAGING"
}`,
			responseStatus: http.StatusBadRequest,
			expectedPath:   "/api-definitions/v2/endpoints/987/versions/1/deactivate",
			withError: &Error{
				Type:          "/api-definitions/error-types/endpoint-version-not-active",
				Status:        http.StatusBadRequest,
				Title:         "Version not active",
				Instance:      "3692ad83-2ac3-4912-aca0-a5b93eac84c2",
				Detail:        "You cannot deactivate the version because it's not active on STAGING.",
				Severity:      ptr.To("ERROR"),
				EndpointID:    ptr.To(int64(987)),
				EndpointName:  ptr.To("dxetest"),
				VersionNumber: ptr.To(int64(1)),
				Network:       ptr.To("STAGING"),
			},
		},
		"404 not found": {
			request: DeactivateVersionRequest{
				VersionNumber: 4,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks:               []NetworkType{ActivationNetworkStaging},
					Notes:                  "some notes",
					NotificationRecipients: []string{"devteam@domain.com"},
				},
			},
			responseBody: `
{
    "type": "https://problems.luna.akamaiapis.net/api-definitions/error-types/NOT-FOUND",
    "status": 404,
    "title": "Not Found",
    "instance": "22f95e07-00f8-4dfc-a903-ba9b1e2999bf",
    "detail": "Invalid version provided.",
    "severity": "ERROR"
}`,
			responseStatus: http.StatusNotFound,
			expectedPath:   "/api-definitions/v2/endpoints/987/versions/4/deactivate",
			withError: &Error{
				Type:     "https://problems.luna.akamaiapis.net/api-definitions/error-types/NOT-FOUND",
				Status:   http.StatusNotFound,
				Title:    "Not Found",
				Instance: "22f95e07-00f8-4dfc-a903-ba9b1e2999bf",
				Detail:   "Invalid version provided.",
				Severity: ptr.To("ERROR"),
			},
		},
		"incorrect network": {
			request: DeactivateVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks:               []NetworkType{"TYP0"},
					Notes:                  "some notes",
					NotificationRecipients: []string{"devteam@domain.com"},
				},
			},
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/deactivate",
			withError:    ErrStructValidation,
		},
		"incorrect recipient": {
			request: DeactivateVersionRequest{
				VersionNumber: 1,
				APIEndpointID: 987,
				Body: ActivationRequestBody{
					Networks:               []NetworkType{ActivationNetworkStaging},
					Notes:                  "some notes",
					NotificationRecipients: []string{"devteam"},
				},
			},
			expectedPath: "/api-definitions/v2/endpoints/987/versions/1/deactivate",
			withError:    ErrStructValidation,
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
				if test.expectedRequestBody != "" {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.JSONEq(t, test.expectedRequestBody, string(body))
				}
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.DeactivateVersion(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
