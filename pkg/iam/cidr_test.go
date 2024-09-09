package iam

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCIDRBlocks(t *testing.T) {
	tests := map[string]struct {
		params           ListCIDRBlocksRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListCIDRBlocksResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:         ListCIDRBlocksRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "cidrBlockId": 1,
        "enabled": true,
        "comments": "abc",
        "cidrBlock": "1.2.3.4/8",
        "createdDate": "2024-06-17T08:46:41.000Z",
        "createdBy": "johndoe",
        "modifiedDate": "2024-06-17T08:46:41.000Z",
        "modifiedBy": "johndoe"
    },
    {
        "cidrBlockId": 2,
        "enabled": false,
        "comments": null,
        "cidrBlock": "2.4.8.16/32",
        "createdDate": "2024-06-25T06:14:36.000Z",
        "createdBy": "johndoe",
        "modifiedDate": "2024-06-25T06:14:36.000Z",
        "modifiedBy": "johndoe"
    }
]`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist?actions=false",
			expectedResponse: &ListCIDRBlocksResponse{
				{
					CIDRBlockID:  1,
					Enabled:      true,
					Comments:     "abc",
					CIDRBlock:    "1.2.3.4/8",
					CreatedDate:  test.NewTimeFromString(t, "2024-06-17T08:46:41.000Z"),
					CreatedBy:    "johndoe",
					ModifiedDate: test.NewTimeFromString(t, "2024-06-17T08:46:41.000Z"),
					ModifiedBy:   "johndoe",
					Actions:      nil,
				},
				{
					CIDRBlockID:  2,
					Enabled:      false,
					Comments:     "",
					CIDRBlock:    "2.4.8.16/32",
					CreatedDate:  test.NewTimeFromString(t, "2024-06-25T06:14:36.000Z"),
					CreatedBy:    "johndoe",
					ModifiedDate: test.NewTimeFromString(t, "2024-06-25T06:14:36.000Z"),
					ModifiedBy:   "johndoe",
					Actions:      nil,
				},
			},
		},
		"200 with actions": {
			params:         ListCIDRBlocksRequest{Actions: true},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "cidrBlockId": 1,
        "enabled": true,
        "comments": "abc",
        "cidrBlock": "1.2.3.4/8",
        "createdDate": "2024-06-17T08:46:41.000Z",
        "createdBy": "johndoe",
        "modifiedDate": "2024-06-17T08:46:41.000Z",
        "modifiedBy": "johndoe",
        "actions": {
            "edit": true,
            "delete": true
        }
    },
    {
        "cidrBlockId": 2,
        "enabled": false,
        "comments": null,
        "cidrBlock": "2.4.8.16/32",
        "createdDate": "2024-06-25T06:14:36.000Z",
        "createdBy": "johndoe",
        "modifiedDate": "2024-06-25T06:14:36.000Z",
        "modifiedBy": "johndoe",
        "actions": {
            "edit": true,
            "delete": true
        }
    }
]`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist?actions=true",
			expectedResponse: &ListCIDRBlocksResponse{
				{
					CIDRBlockID:  1,
					Enabled:      true,
					Comments:     "abc",
					CIDRBlock:    "1.2.3.4/8",
					CreatedDate:  test.NewTimeFromString(t, "2024-06-17T08:46:41.000Z"),
					CreatedBy:    "johndoe",
					ModifiedDate: test.NewTimeFromString(t, "2024-06-17T08:46:41.000Z"),
					ModifiedBy:   "johndoe",
					Actions: &CIDRActions{
						Edit:   true,
						Delete: true,
					},
				},
				{
					CIDRBlockID:  2,
					Enabled:      false,
					Comments:     "",
					CIDRBlock:    "2.4.8.16/32",
					CreatedDate:  test.NewTimeFromString(t, "2024-06-25T06:14:36.000Z"),
					CreatedBy:    "johndoe",
					ModifiedDate: test.NewTimeFromString(t, "2024-06-25T06:14:36.000Z"),
					ModifiedBy:   "johndoe",
					Actions: &CIDRActions{
						Edit:   true,
						Delete: true,
					},
				},
			},
		},
		"500 internal server error": {
			params:         ListCIDRBlocksRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}
`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist?actions=false",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
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
			result, err := client.ListCIDRBlocks(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreateCIDRBlock(t *testing.T) {
	tests := map[string]struct {
		params              CreateCIDRBlockRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *CreateCIDRBlockResponse
		withError           func(*testing.T, error)
	}{
		"201 created": {
			params: CreateCIDRBlockRequest{
				CIDRBlock: "1.2.3.4/32",
				Comments:  "abc",
				Enabled:   true,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "cidrBlockId": 1234,
    "enabled": true,
    "comments": "abc",
    "cidrBlock": "1.2.3.4/32",
    "createdDate": "2024-07-15T13:53:49.000Z",
    "createdBy": "johndoe",
    "modifiedDate": "2024-07-15T13:53:49.000Z",
    "modifiedBy": "johndoe"
}`,
			expectedPath:        "/identity-management/v3/user-admin/ip-acl/allowlist",
			expectedRequestBody: `{"cidrBlock":"1.2.3.4/32","comments":"abc","enabled":true}`,
			expectedResponse: &CreateCIDRBlockResponse{
				CIDRBlockID:  1234,
				Enabled:      true,
				Comments:     "abc",
				CIDRBlock:    "1.2.3.4/32",
				CreatedDate:  test.NewTimeFromString(t, "2024-07-15T13:53:49.000Z"),
				CreatedBy:    "johndoe",
				ModifiedDate: test.NewTimeFromString(t, "2024-07-15T13:53:49.000Z"),
				ModifiedBy:   "johndoe",
				Actions:      nil,
			},
		},
		"201 without comment": {
			params: CreateCIDRBlockRequest{
				CIDRBlock: "1.2.3.4/32",
				Enabled:   true,
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "cidrBlockId": 1234,
    "enabled": true,
    "comments": null,
    "cidrBlock": "1.2.3.4/32",
    "createdDate": "2024-07-15T13:53:49.000Z",
    "createdBy": "johndoe",
    "modifiedDate": "2024-07-15T13:53:49.000Z",
    "modifiedBy": "johndoe"
}`,
			expectedPath:        "/identity-management/v3/user-admin/ip-acl/allowlist",
			expectedRequestBody: `{"cidrBlock":"1.2.3.4/32","enabled":true}`,
			expectedResponse: &CreateCIDRBlockResponse{
				CIDRBlockID:  1234,
				Enabled:      true,
				Comments:     "",
				CIDRBlock:    "1.2.3.4/32",
				CreatedDate:  test.NewTimeFromString(t, "2024-07-15T13:53:49.000Z"),
				CreatedBy:    "johndoe",
				ModifiedDate: test.NewTimeFromString(t, "2024-07-15T13:53:49.000Z"),
				ModifiedBy:   "johndoe",
				Actions:      nil,
			},
		},
		"missing required parameters": {
			params: CreateCIDRBlockRequest{},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
				assert.Contains(t, err.Error(), "create CIDR block: struct validation: CIDRBlock: cannot be blank")
			},
		},
		"403 - incorrect cidrblock": {
			params: CreateCIDRBlockRequest{
				CIDRBlock: "1.2.3.4:32",
				Comments:  "abc",
				Enabled:   true,
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "/ip-acl/error-types/1013",
    "httpStatus": 403,
    "title": "CIDR format not correct",
    "detail": "invalid cidrblock/ip",
    "instance": "",
    "errors": []
}`,
			expectedPath:        "/identity-management/v3/user-admin/ip-acl/allowlist",
			expectedRequestBody: `{"cidrBlock":"1.2.3.4:32","comments":"abc","enabled":true}`,
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "/ip-acl/error-types/1013",
					HTTPStatus: http.StatusForbidden,
					Title:      "CIDR format not correct",
					Detail:     "invalid cidrblock/ip",
					StatusCode: http.StatusForbidden,
					Instance:   "",
					Errors:     json.RawMessage("[]"),
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
		"500 internal server error": {
			params: CreateCIDRBlockRequest{
				CIDRBlock: "1.2.3.4/32",
				Comments:  "abc",
				Enabled:   true,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error making request",
    "status": 500
}
`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
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
			result, err := client.CreateCIDRBlock(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetCIDRBlock(t *testing.T) {
	tests := map[string]struct {
		params           GetCIDRBlockRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetCIDRBlockResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:         GetCIDRBlockRequest{CIDRBlockID: 1},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "cidrBlockId": 1,
    "enabled": true,
    "comments": "abc",
    "cidrBlock": "1.2.3.4/8",
    "createdDate": "2024-06-17T08:46:41.000Z",
    "createdBy": "johndoe",
    "modifiedDate": "2024-06-17T08:46:41.000Z",
    "modifiedBy": "johndoe"
}`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist/1?actions=false",
			expectedResponse: &GetCIDRBlockResponse{
				CIDRBlockID:  1,
				Enabled:      true,
				Comments:     "abc",
				CIDRBlock:    "1.2.3.4/8",
				CreatedDate:  test.NewTimeFromString(t, "2024-06-17T08:46:41.000Z"),
				CreatedBy:    "johndoe",
				ModifiedDate: test.NewTimeFromString(t, "2024-06-17T08:46:41.000Z"),
				ModifiedBy:   "johndoe",
				Actions:      nil,
			},
		},
		"200 with actions": {
			params:         GetCIDRBlockRequest{CIDRBlockID: 1, Actions: true},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "cidrBlockId": 1,
    "enabled": true,
    "comments": "abc",
    "cidrBlock": "1.2.3.4/8",
    "createdDate": "2024-06-17T08:46:41.000Z",
    "createdBy": "johndoe",
    "modifiedDate": "2024-06-17T08:46:41.000Z",
    "modifiedBy": "johndoe",
    "actions": {
        "edit": true,
        "delete": true
    }
}`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist/1?actions=true",
			expectedResponse: &GetCIDRBlockResponse{
				CIDRBlockID:  1,
				Enabled:      true,
				Comments:     "abc",
				CIDRBlock:    "1.2.3.4/8",
				CreatedDate:  test.NewTimeFromString(t, "2024-06-17T08:46:41.000Z"),
				CreatedBy:    "johndoe",
				ModifiedDate: test.NewTimeFromString(t, "2024-06-17T08:46:41.000Z"),
				ModifiedBy:   "johndoe",
				Actions: &CIDRActions{
					Edit:   true,
					Delete: true,
				},
			},
		},
		"missing required parameters": {
			params: GetCIDRBlockRequest{},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
				assert.Contains(t, err.Error(), "get CIDR block: struct validation: CIDRBlockID: cannot be blank")
			},
		},
		"incorrect parameters": {
			params: GetCIDRBlockRequest{CIDRBlockID: -1},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
				assert.Contains(t, err.Error(), "get CIDR block: struct validation: CIDRBlockID: must be no less than 1")
			},
		},
		"404 no such block": {
			params:         GetCIDRBlockRequest{CIDRBlockID: 9000},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/ip-acl/error-types/1010",
    "httpStatus": 404,
    "title": "no data found",
    "detail": "",
    "instance": "",
    "errors": []
}`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist/9000?actions=false",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "/ip-acl/error-types/1010",
					HTTPStatus: http.StatusNotFound,
					Title:      "no data found",
					Detail:     "",
					StatusCode: http.StatusNotFound,
					Instance:   "",
					Errors:     json.RawMessage("[]"),
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
		"500 internal server error": {
			params:         GetCIDRBlockRequest{CIDRBlockID: 1},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
			"type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist/1?actions=false",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
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
			result, err := client.GetCIDRBlock(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdateCIDRBlock(t *testing.T) {
	tests := map[string]struct {
		params              UpdateCIDRBlockRequest
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedRequestBody string
		expectedResponse    *UpdateCIDRBlockResponse
		withError           func(*testing.T, error)
	}{
		"200 OK": {
			params: UpdateCIDRBlockRequest{
				CIDRBlockID: 1,
				Body: UpdateCIDRBlockRequestBody{
					CIDRBlock: "1.2.3.4/32",
					Comments:  "abc - updated",
					Enabled:   false,
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
   "cidrBlockId": 1234,
   "enabled": false,
   "comments": "abc - updated",
   "cidrBlock": "1.2.3.4/32",
   "createdDate": "2024-07-15T13:53:49.000Z",
   "createdBy": "johndoe",
   "modifiedDate": "2024-07-16T13:53:49.000Z",
   "modifiedBy": "johndoe"
}`,
			expectedPath:        "/identity-management/v3/user-admin/ip-acl/allowlist/1",
			expectedRequestBody: `{"cidrBlock":"1.2.3.4/32","comments":"abc - updated","enabled":false}`,
			expectedResponse: &UpdateCIDRBlockResponse{
				CIDRBlockID:  1234,
				Enabled:      false,
				Comments:     "abc - updated",
				CIDRBlock:    "1.2.3.4/32",
				CreatedDate:  test.NewTimeFromString(t, "2024-07-15T13:53:49.000Z"),
				CreatedBy:    "johndoe",
				ModifiedDate: test.NewTimeFromString(t, "2024-07-16T13:53:49.000Z"),
				ModifiedBy:   "johndoe",
				Actions:      nil,
			},
		},
		"missing required parameters": {
			params: UpdateCIDRBlockRequest{},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
				assert.Contains(t, err.Error(), "update CIDR block: struct validation: CIDRBlock: cannot be blank\nCIDRBlockID: cannot be blank")
			},
		},
		"invalid required parameters": {
			params: UpdateCIDRBlockRequest{
				CIDRBlockID: -1,
				Body: UpdateCIDRBlockRequestBody{
					CIDRBlock: "1.2.3.4/32",
					Comments:  "abc - updated",
					Enabled:   false,
				},
			},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
				assert.Contains(t, err.Error(), "update CIDR block: struct validation: CIDRBlockID: must be no less than 1")
			},
		},
		"500 internal server error": {
			params: UpdateCIDRBlockRequest{
				CIDRBlockID: 1,
				Body: UpdateCIDRBlockRequestBody{
					CIDRBlock: "1.2.3.4/32",
					Comments:  "abc - updated",
					Enabled:   false,
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
			"type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist/1",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
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
			result, err := client.UpdateCIDRBlock(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeleteCIDRBlocks(t *testing.T) {
	tests := map[string]struct {
		params         DeleteCIDRBlockRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 No content": {
			params:         DeleteCIDRBlockRequest{CIDRBlockID: 1},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/identity-management/v3/user-admin/ip-acl/allowlist/1",
		},
		"missing required parameters": {
			params: DeleteCIDRBlockRequest{},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
				assert.Contains(t, err.Error(), "delete CIDR block: struct validation: CIDRBlockID: cannot be blank")
			},
		},
		"incorrect parameters": {
			params: DeleteCIDRBlockRequest{CIDRBlockID: -1},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
				assert.Contains(t, err.Error(), "delete CIDR block: struct validation: CIDRBlockID: must be no less than 1")
			},
		},
		"404 no such block": {
			params:         DeleteCIDRBlockRequest{CIDRBlockID: 9000},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "/ip-acl/error-types/1010",
    "httpStatus": 404,
    "title": "no data found",
    "detail": "",
    "instance": "",
    "errors": []
}`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist/9000",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "/ip-acl/error-types/1010",
					HTTPStatus: http.StatusNotFound,
					Title:      "no data found",
					Detail:     "",
					StatusCode: http.StatusNotFound,
					Instance:   "",
					Errors:     json.RawMessage("[]"),
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
		"500 internal server error": {
			params:         DeleteCIDRBlockRequest{CIDRBlockID: 1},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
		{
		   "type": "internal_error",
		   "title": "Internal Server Error",
		   "detail": "Error making request",
		   "status": 500
		}
		`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist/1",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
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
			err := client.DeleteCIDRBlock(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateCIDRBlocks(t *testing.T) {
	tests := map[string]struct {
		params         ValidateCIDRBlockRequest
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      func(*testing.T, error)
	}{
		"204 No content": {
			params:         ValidateCIDRBlockRequest{CIDRBlock: "1.2.3.4/32"},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/identity-management/v3/user-admin/ip-acl/allowlist/validate?cidrblock=1.2.3.4%2F32",
		},
		"missing required parameters": {
			params: ValidateCIDRBlockRequest{},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
				assert.Contains(t, err.Error(), "validate CIDR block: struct validation: CIDRBlock: cannot be blank")
			},
		},
		"400 invalid": {
			params:         ValidateCIDRBlockRequest{CIDRBlock: "abc"},
			responseStatus: http.StatusBadRequest,
			responseBody: `
{
    "type": "/ip-acl/error-types/1013",
    "httpStatus": 400,
    "title": "CIDR format not correct",
    "detail": "invalid cidr format",
    "instance": "",
    "errors": []
}`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist/validate?cidrblock=abc",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "/ip-acl/error-types/1013",
					HTTPStatus: http.StatusBadRequest,
					Title:      "CIDR format not correct",
					Detail:     "invalid cidr format",
					StatusCode: http.StatusBadRequest,
					Instance:   "",
					Errors:     json.RawMessage("[]"),
				}
				assert.Equal(t, true, err.Is(e))
			},
		},
		"500 internal server error": {
			params:         ValidateCIDRBlockRequest{CIDRBlock: "1.2.3.4/32"},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
				{
				   "type": "internal_error",
				   "title": "Internal Server Error",
				   "detail": "Error making request",
				   "status": 500
				}
				`,
			expectedPath: "/identity-management/v3/user-admin/ip-acl/allowlist/validate?cidrblock=1.2.3.4%2F32",
			withError: func(t *testing.T, e error) {
				err := Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.Equal(t, true, err.Is(e))
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
			err := client.ValidateCIDRBlock(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
