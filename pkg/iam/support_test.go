package iam

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIAM_GetPasswordPolicy(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPasswordPolicyResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v3/user-admin/common/password-policy",
			responseBody: `
{
	  "caseDif": 0,
	  "maxRepeating": 1,
	  "minDigits": 1,
	  "minLength": 1,
	  "minLetters": 1,
	  "minNonAlpha": 0,
	  "minReuse": 1,
	  "pwclass": "test_class",
	  "rotateFrequency": 10
}
`,
			expectedResponse: &GetPasswordPolicyResponse{
				CaseDiff:        0,
				MaxRepeating:    1,
				MinDigits:       1,
				MinLength:       1,
				MinLetters:      1,
				MinNonAlpha:     0,
				MinReuse:        1,
				PwClass:         "test_class",
				RotateFrequency: 10,
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			expectedPath:   "/identity-management/v3/user-admin/common/password-policy",
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetPasswordPolicy(context.Background())
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_SupportedCountries(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
			[
				"Greece",
				"Greenland",
				"Grenada"
			]`,
			expectedPath:     "/identity-management/v3/user-admin/common/countries",
			expectedResponse: []string{"Greece", "Greenland", "Grenada"},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/common/countries",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.SupportedCountries(context.Background())
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_SupportedTimezones(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []Timezone
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
			[
				{
					"timezone": "Asia/Rangoon",
					"description": "Asia/Rangoon GMT+6",
					"offset": "+6",
					"posix": "Asia/Rangoon"
				}
			]`,
			expectedPath: "/identity-management/v3/user-admin/common/timezones",
			expectedResponse: []Timezone{
				{
					Timezone:    "Asia/Rangoon",
					Description: "Asia/Rangoon GMT+6",
					Offset:      "+6",
					Posix:       "Asia/Rangoon",
				},
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/common/timezones",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.SupportedTimezones(context.Background())
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_SupportedContactTypes(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
[
    "Billing",
    "Security"
]`,
			expectedPath:     "/identity-management/v3/user-admin/common/contact-types",
			expectedResponse: []string{"Billing", "Security"},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/common/contact-types",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.SupportedContactTypes(context.Background())
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_SupportedLanguages(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
[
    "Deutsch",
    "English"
]`,
			expectedPath:     "/identity-management/v3/user-admin/common/supported-languages",
			expectedResponse: []string{"Deutsch", "English"},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/common/supported-languages",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.SupportedLanguages(context.Background())
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_ListProducts(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
[
    "EdgeComputing for Java",
    "Streaming"
]`,
			expectedPath:     "/identity-management/v3/user-admin/common/notification-products",
			expectedResponse: []string{"EdgeComputing for Java", "Streaming"},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/common/notification-products",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListProducts(context.Background())
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_ListTimeoutPolicies(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []TimeoutPolicy
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "name": "after15Minutes",
        "value": 900
    },
    {
        "name": "after30Minutes",
        "value": 1800
    }
]`,
			expectedPath: "/identity-management/v3/user-admin/common/timeout-policies",
			expectedResponse: []TimeoutPolicy{
				{
					Name:  "after15Minutes",
					Value: 900,
				},
				{
					Name:  "after30Minutes",
					Value: 1800,
				},
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/common/timeout-policies",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListTimeoutPolicies(context.Background())
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_ListStates(t *testing.T) {
	tests := map[string]struct {
		params           ListStatesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []string
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListStatesRequest{
				Country: "canada",
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
	"AB",
	"BC"
]`,
			expectedPath:     "/identity-management/v3/user-admin/common/countries/canada/states",
			expectedResponse: []string{"AB", "BC"},
		},
		"500 internal server error": {
			params: ListStatesRequest{
				Country: "canada",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/user-admin/common/countries/canada/states",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"missing country": {
			params: ListStatesRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListStates(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, result)
		})
	}
}

func TestIAM_ListAccountSwitchKeys(t *testing.T) {
	tests := map[string]struct {
		params           ListAccountSwitchKeysRequest
		responseStatus   int
		expectedPath     string
		responseBody     string
		expectedResponse ListAccountSwitchKeysResponse
		withError        func(*testing.T, error)
	}{
		"200 OK with specified client": {
			params: ListAccountSwitchKeysRequest{
				ClientID: "test1234",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v3/api-clients/test1234/account-switch-keys",
			responseBody: `
[
  {
    "accountName": "Test Name A",
    "accountSwitchKey": "ABC-123"
  },
  {
    "accountName": "Test Name A",
    "accountSwitchKey": "ABCD-1234"
  },
  {
    "accountName": "Test Name B",
    "accountSwitchKey": "ABCDE-12345"
  }
]
`,
			expectedResponse: ListAccountSwitchKeysResponse{
				AccountSwitchKey{
					AccountName:      "Test Name A",
					AccountSwitchKey: "ABC-123",
				},
				AccountSwitchKey{
					AccountName:      "Test Name A",
					AccountSwitchKey: "ABCD-1234",
				},
				AccountSwitchKey{
					AccountName:      "Test Name B",
					AccountSwitchKey: "ABCDE-12345",
				},
			},
		},
		"200 OK without specified client": {
			params:         ListAccountSwitchKeysRequest{},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v3/api-clients/self/account-switch-keys",
			responseBody: `
[
  {
    "accountName": "Test Name A",
    "accountSwitchKey": "ABC-123"
  },
  {
    "accountName": "Test Name A",
    "accountSwitchKey": "ABCD-1234"
  },
  {
    "accountName": "Test Name B",
    "accountSwitchKey": "ABCDE-12345"
  }
]
`,
			expectedResponse: ListAccountSwitchKeysResponse{
				AccountSwitchKey{
					AccountName:      "Test Name A",
					AccountSwitchKey: "ABC-123",
				},
				AccountSwitchKey{
					AccountName:      "Test Name A",
					AccountSwitchKey: "ABCD-1234",
				},
				AccountSwitchKey{
					AccountName:      "Test Name B",
					AccountSwitchKey: "ABCDE-12345",
				},
			},
		},
		"200 OK - no account switch keys": {
			params: ListAccountSwitchKeysRequest{
				ClientID: "test1234",
			},
			responseStatus:   http.StatusOK,
			expectedPath:     "/identity-management/v3/api-clients/test1234/account-switch-keys",
			responseBody:     `[]`,
			expectedResponse: ListAccountSwitchKeysResponse{},
		},
		"200 OK with query param": {
			params: ListAccountSwitchKeysRequest{
				ClientID: "test1234",
				Search:   "Name A",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/identity-management/v3/api-clients/test1234/account-switch-keys?search=Name+A",
			responseBody: `
[
  {
    "accountName": "Test Name A",
    "accountSwitchKey": "ABC-123"
  },
  {
    "accountName": "Test Name A",
    "accountSwitchKey": "ABCD-1234"
  }
]
`,
			expectedResponse: ListAccountSwitchKeysResponse{
				AccountSwitchKey{
					AccountName:      "Test Name A",
					AccountSwitchKey: "ABC-123",
				},
				AccountSwitchKey{
					AccountName:      "Test Name A",
					AccountSwitchKey: "ABCD-1234",
				},
			},
		},
		"404 not found": {
			params: ListAccountSwitchKeysRequest{
				ClientID: "test12344",
			},
			responseStatus: http.StatusNotFound,
			expectedPath:   "/identity-management/v3/api-clients/test12344/account-switch-keys",
			responseBody: `
{
	"instances": "",
    "type": "/identity-management/error-types/2",
    "status": 404,
    "title": "invalid open identity",
	"detail": ""
}				
`,
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "/identity-management/error-types/2",
					Title:      "invalid open identity",
					StatusCode: http.StatusNotFound,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"500 internal server error": {
			params: ListAccountSwitchKeysRequest{
				ClientID: "test12344",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}`,
			expectedPath: "/identity-management/v3/api-clients/test12344/account-switch-keys",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				_, err := w.Write([]byte(tc.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			users, err := client.ListAccountSwitchKeys(context.Background(), tc.params)
			if tc.withError != nil {
				tc.withError(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, users)
		})
	}
}
