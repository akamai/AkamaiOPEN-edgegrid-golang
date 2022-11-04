package papi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPapi_GetProperties(t *testing.T) {
	tests := map[string]struct {
		request          GetPropertiesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPropertiesResponse
		withError        error
	}{
		"200 OK": {
			request: GetPropertiesRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"properties": {
		"items": [
			{
				"accountId": "act_1-1TJZFB",
				"contractId": "ctr_1-1TJZH5",
				"groupId": "grp_15166",
				"propertyId": "prp_175780",
				"propertyName": "example.com",
				"latestVersion": 2,
				"stagingVersion": 1,
				"productId": "prp_175780",
				"productionVersion": null,
				"assetId": "aid_101",
				"note": "Notes about example.com"
			}
		]
	}
}`,
			expectedPath: "/papi/v1/properties?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &GetPropertiesResponse{
				Properties: PropertiesItems{Items: []*Property{
					{
						AccountID:         "act_1-1TJZFB",
						ContractID:        "ctr_1-1TJZH5",
						GroupID:           "grp_15166",
						PropertyID:        "prp_175780",
						ProductID:         "prp_175780",
						PropertyName:      "example.com",
						LatestVersion:     2,
						StagingVersion:    tools.IntPtr(1),
						ProductionVersion: nil,
						AssetID:           "aid_101",
						Note:              "Notes about example.com",
					},
				}},
			},
		},
		"500 internal server error": {
			request: GetPropertiesRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching properties",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching properties",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request: GetPropertiesRequest{
				GroupID: "grp_15166",
			},
			responseStatus: http.StatusInternalServerError,
			withError:      ErrStructValidation,
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
			result, err := client.GetProperties(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_GetProperty(t *testing.T) {
	tests := map[string]struct {
		request          GetPropertyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPropertyResponse
		withError        error
	}{
		"200 OK": {
			request: GetPropertyRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				PropertyID: "prp_175780",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"properties": {
		"items": [
			{
				"accountId": "act_1-1TJZFB",
				"contractId": "ctr_1-1TJZH5",
				"groupId": "grp_15166",
				"propertyId": "prp_175780",
				"propertyName": "example.com",
				"latestVersion": 2,
				"stagingVersion": 1,
				"productId": "prp_175780",
				"productionVersion": null,
				"assetId": "aid_101",
				"note": "Notes about example.com"
			}
		]
	}
}`,
			expectedPath: "/papi/v1/properties/prp_175780?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &GetPropertyResponse{
				Properties: PropertiesItems{Items: []*Property{
					{
						AccountID:         "act_1-1TJZFB",
						ContractID:        "ctr_1-1TJZH5",
						GroupID:           "grp_15166",
						PropertyID:        "prp_175780",
						ProductID:         "prp_175780",
						PropertyName:      "example.com",
						LatestVersion:     2,
						StagingVersion:    tools.IntPtr(1),
						ProductionVersion: nil,
						AssetID:           "aid_101",
						Note:              "Notes about example.com",
					},
				}},
				Property: &Property{

					AccountID:         "act_1-1TJZFB",
					ContractID:        "ctr_1-1TJZH5",
					GroupID:           "grp_15166",
					PropertyID:        "prp_175780",
					ProductID:         "prp_175780",
					PropertyName:      "example.com",
					LatestVersion:     2,
					StagingVersion:    tools.IntPtr(1),
					ProductionVersion: nil,
					AssetID:           "aid_101",
					Note:              "Notes about example.com",
				}},
		},
		"Property not found": {
			request: GetPropertyRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				PropertyID: "prp_175780",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"properties": {
		"items": [
		]
	}
}`,
			expectedPath: "/papi/v1/properties/prp_175780?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			withError:    ErrNotFound,
		},
		"500 internal server error": {
			request: GetPropertyRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				PropertyID: "prp_175780",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching properties",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/prp_175780?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching properties",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request: GetPropertyRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
			},
			withError: ErrStructValidation,
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
			result, err := client.GetProperty(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_CreateProperty(t *testing.T) {
	tests := map[string]struct {
		request          CreatePropertyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *CreatePropertyResponse
		withError        error
	}{
		"201 created": {
			request: CreatePropertyRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Property: PropertyCreate{
					ProductID:    "prd_Alta",
					PropertyName: "my.new.property.com",
					CloneFrom: &PropertyCloneFrom{
						PropertyID: "prp_1234",
						Version:    1,
					},
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
	"propertyLink": "/papi/v1/properties/prp_173137?contractId=ctr_1-1TJZH5&groupId=grp_15225"
}`,
			expectedPath: "/papi/v1/properties?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &CreatePropertyResponse{
				PropertyID:   "prp_173137",
				PropertyLink: "/papi/v1/properties/prp_173137?contractId=ctr_1-1TJZH5&groupId=grp_15225",
			},
		},
		"500 internal server error": {
			request: CreatePropertyRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Property: PropertyCreate{
					ProductID:    "prd_Alta",
					PropertyName: "my.new.property.com",
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating property",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating property",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request: CreatePropertyRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				Property: PropertyCreate{
					ProductID: "prd_Alta",
				},
			},
			withError: ErrStructValidation,
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
			result, err := client.CreateProperty(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestPapi_RemoveProperty(t *testing.T) {
	tests := map[string]struct {
		request          RemovePropertyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *RemovePropertyResponse
		withError        error
	}{
		"200 OK": {
			request: RemovePropertyRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				PropertyID: "prp_175780",
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
	"message": "Deletion Successful."
}`,
			expectedPath: "/papi/v1/properties/prp_175780?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			expectedResponse: &RemovePropertyResponse{
				Message: "Deletion Successful.",
			},
		},
		"500 internal server error": {
			request: RemovePropertyRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
				PropertyID: "prp_175780",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error removing property",
    "status": 500
}`,
			expectedPath: "/papi/v1/properties/prp_175780?contractId=ctr_1-1TJZFW&groupId=grp_15166",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error removing property",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request: RemovePropertyRequest{
				ContractID: "ctr_1-1TJZFW",
				GroupID:    "grp_15166",
			},
			withError: ErrStructValidation,
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
			result, err := client.RemoveProperty(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
