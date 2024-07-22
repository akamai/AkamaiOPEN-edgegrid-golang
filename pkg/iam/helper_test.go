package iam

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestIAMListAllowedCPCodes(t *testing.T) {
	tests := map[string]struct {
		params           ListAllowedCPCodesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse ListAllowedCPCodesResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: ListAllowedCPCodesRequest{
				UserName: "jsmith",
				ListAllowedCPCodesRequestBody: ListAllowedCPCodesRequestBody{
					ClientType: "CLIENT",
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `[
  {
    "name": "Stream Analyzer (36915)",
    "value": 36915
  },
  {
    "name": "plopessa-uvod-ns (373118)",
    "value": 373118
  },
  {
    "name": "ArunNS (866797)",
    "value": 866797
  },
  {
    "name": "1234 (933076)",
    "value": 933076
  }
]`,
			expectedPath: "/identity-management/v3/users/jsmith/allowed-cpcodes",
			expectedResponse: ListAllowedCPCodesResponse{
				{
					Name:  "Stream Analyzer (36915)",
					Value: 36915,
				},
				{
					Name:  "plopessa-uvod-ns (373118)",
					Value: 373118,
				},
				{
					Name:  "ArunNS (866797)",
					Value: 866797,
				},
				{
					Name:  "1234 (933076)",
					Value: 933076,
				},
			},
		},
		"200 OK with groups": {
			params: ListAllowedCPCodesRequest{
				UserName: "jsmith",
				ListAllowedCPCodesRequestBody: ListAllowedCPCodesRequestBody{
					ClientType: "SERVICE_ACCOUNT",
					Groups: []AllowedCPCodesGroup{
						{
							GroupID: 1,
						},
					},
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `[
  {
    "name": "Stream Analyzer (36915)",
    "value": 36915
  },
  {
    "name": "plopessa-uvod-ns (373118)",
    "value": 373118
  },
  {
    "name": "ArunNS (866797)",
    "value": 866797
  },
  {
    "name": "1234 (933076)",
    "value": 933076
  }
]`,
			expectedPath: "/identity-management/v3/users/jsmith/allowed-cpcodes",
			expectedResponse: ListAllowedCPCodesResponse{
				{
					Name:  "Stream Analyzer (36915)",
					Value: 36915,
				},
				{
					Name:  "plopessa-uvod-ns (373118)",
					Value: 373118,
				},
				{
					Name:  "ArunNS (866797)",
					Value: 866797,
				},
				{
					Name:  "1234 (933076)",
					Value: 933076,
				},
			},
		},
		"500 internal server error": {
			params: ListAllowedCPCodesRequest{
				UserName: "jsmith",
				ListAllowedCPCodesRequestBody: ListAllowedCPCodesRequestBody{
					ClientType: "CLIENT",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
				"type": "internal_error",
				"title": "Internal Server Error",
				"detail": "Error making request",
				"status": 500
				}`,
			expectedPath: "/identity-management/v3/users/jsmith/allowed-cpcodes",
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
		"missing user name and client type": {
			params: ListAllowedCPCodesRequest{},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "list allowed CP codes: struct validation:\nClientType: cannot be blank\nUserName: cannot be blank")
			},
		},
		"group is required for client type SERVICE_ACCOUNT": {
			params: ListAllowedCPCodesRequest{
				UserName: "jsmith",
				ListAllowedCPCodesRequestBody: ListAllowedCPCodesRequestBody{
					ClientType: "SERVICE_ACCOUNT",
				},
			},
			withError: func(t *testing.T, err error) {
				want := ErrStructValidation
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
				assert.Contains(t, err.Error(), "list allowed CP codes: struct validation:\nGroups: cannot be blank")
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
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.ListAllowedCPCodes(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
