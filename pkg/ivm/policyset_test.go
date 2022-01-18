package ivm

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

func TestListPolicySets(t *testing.T) {
	tests := map[string]struct {
		params           ListPolicySetsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []PolicySet
		withError        error
	}{
		"200 OK - both networks": {
			params: ListPolicySetsRequest{
				Contract: "3-WNKXX1",
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "name": "terraform_beta_v2",
        "id": "570e9090-5dbe-11ec-8a0a-71665789c1d8",
        "type": "IMAGE",
        "region": "US",
        "properties": [
            "jsmith_sqa2"
        ],
        "user": "jsmith",
        "lastModified": "2021-12-15 15:47:42+0000"
    },
    {
        "name": "my_example_token",
        "id": "57c73ae0-5204-11ec-8a0a-71665789c1d8",
        "type": "IMAGE",
        "region": "US",
        "properties": [],
        "user": "jsmith",
        "lastModified": "2021-11-30 17:38:35+0000"
    },
    {
        "name": "terraform_demo-1104268",
        "id": "terraform_demo-1104268",
        "type": "IMAGE",
        "region": "US",
        "properties": [
            "jsmith_sqa2"
        ],
        "user": "System",
        "lastModified": "2021-12-15 15:51:54+0000"
    }
]`,
			expectedPath: "/imaging/v2/policysets/",
			expectedResponse: []PolicySet{
				{
					Name:         "terraform_beta_v2",
					ID:           "570e9090-5dbe-11ec-8a0a-71665789c1d8",
					Type:         "IMAGE",
					Region:       "US",
					Properties:   []string{"jsmith_sqa2"},
					User:         "jsmith",
					LastModified: "2021-12-15 15:47:42+0000",
				},
				{
					Name:         "my_example_token",
					ID:           "57c73ae0-5204-11ec-8a0a-71665789c1d8",
					Type:         "IMAGE",
					Region:       "US",
					Properties:   []string{},
					User:         "jsmith",
					LastModified: "2021-11-30 17:38:35+0000",
				},
				{
					Name:         "terraform_demo-1104268",
					ID:           "terraform_demo-1104268",
					Type:         "IMAGE",
					Region:       "US",
					Properties:   []string{"jsmith_sqa2"},
					User:         "System",
					LastModified: "2021-12-15 15:51:54+0000",
				},
			},
		},
		"200 OK - staging network": {
			params: ListPolicySetsRequest{
				Contract: "3-WNKXX1",
				Network:  NetworkStaging,
			},
			responseStatus: http.StatusOK,
			responseBody: `
[
    {
        "name": "terraform_beta_v2",
        "id": "570e9090-5dbe-11ec-8a0a-71665789c1d8",
        "type": "IMAGE",
        "region": "US",
        "properties": [
            "jsmith_sqa2"
        ],
        "user": "jsmith",
        "lastModified": "2021-12-15 15:47:42+0000"
    },
    {
        "name": "terraform_demo-1104268",
        "id": "terraform_demo-1104268",
        "type": "IMAGE",
        "region": "US",
        "properties": [
            "jsmith_sqa2"
        ],
        "user": "System",
        "lastModified": "2021-12-15 15:51:54+0000"
    }
]`,
			expectedPath: "/imaging/v2/network/staging/policysets/",
			expectedResponse: []PolicySet{
				{
					Name:         "terraform_beta_v2",
					ID:           "570e9090-5dbe-11ec-8a0a-71665789c1d8",
					Type:         "IMAGE",
					Region:       "US",
					Properties:   []string{"jsmith_sqa2"},
					User:         "jsmith",
					LastModified: "2021-12-15 15:47:42+0000",
				},
				{
					Name:         "terraform_demo-1104268",
					ID:           "terraform_demo-1104268",
					Type:         "IMAGE",
					Region:       "US",
					Properties:   []string{"jsmith_sqa2"},
					User:         "System",
					LastModified: "2021-12-15 15:51:54+0000",
				},
			},
		},
		"400 Bad request": {
			params: ListPolicySetsRequest{
				Contract: "3-WNKXX1",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
"type": "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1004",
"title": "Bad Request",
"instance": "52a21f40-9861-4d35-95d0-a603c85cb2ad",
"status": 400,
"detail": "A contract must be specified using the Contract header.",
"problemId": "52a21f40-9861-4d35-95d0-a603c85cb2ad"
}`,
			expectedPath: "/imaging/v2/policysets/",
			withError: &Error{
				Type:      "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1004",
				Title:     "Bad Request",
				Instance:  "52a21f40-9861-4d35-95d0-a603c85cb2ad",
				Status:    400,
				Detail:    "A contract must be specified using the Contract header.",
				ProblemID: "52a21f40-9861-4d35-95d0-a603c85cb2ad",
			},
		},
		"401 Not authorized": {
			params: ListPolicySetsRequest{
				Contract: "3-WNKXX1",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
"type": "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
"title": "Not authorized",
"status": 401,
"detail": "Inactive client token",
"instance": "https://akaa-mgfkwp3rw4k2whym-eyn4wdjeur5lz37c.luna-dev.akamaiapis.net/imaging/v2/network/staging/policysets/",
"method": "GET",
"serverIp": "104.81.220.242",
"clientIp": "22.22.22.22",
"requestId": "124cc33c",
"requestTime": "2022-01-12T16:53:44Z"
}`,
			expectedPath: "/imaging/v2/policysets/",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
				Title:       "Not authorized",
				Status:      401,
				Detail:      "Inactive client token",
				Instance:    "https://akaa-mgfkwp3rw4k2whym-eyn4wdjeur5lz37c.luna-dev.akamaiapis.net/imaging/v2/network/staging/policysets/",
				Method:      "GET",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "124cc33c",
				RequestTime: "2022-01-12T16:53:44Z",
			},
		},
		"403 Forbidden": {
			params: ListPolicySetsRequest{
				Contract: "3-WNKXX1",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
    "title": "Forbidden",
    "status": 403,
    "detail": "The client does not have the grant needed for the request",
    "instance": "https://akaa-75xqbs7cot5jts7q-yjgmpc4nakckpt44.luna-dev.akamaiapis.net/imaging/v2/network/staging/policysets/",
    "authzRealm": "xlnausb2pov4jvuo.kh6prvchvniqn5zp",
    "method": "GET",
    "serverIp": "104.81.220.242",
    "clientIp": "22.22.22.22",
    "requestId": "1254027a",
    "requestTime": "2022-01-12T16:56:56Z"
}`,
			expectedPath: "/imaging/v2/policysets/",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "https://akaa-75xqbs7cot5jts7q-yjgmpc4nakckpt44.luna-dev.akamaiapis.net/imaging/v2/network/staging/policysets/",
				AuthzRealm:  "xlnausb2pov4jvuo.kh6prvchvniqn5zp",
				Method:      "GET",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "1254027a",
				RequestTime: "2022-01-12T16:56:56Z",
			},
		},
		// 500
		"invalid network": {
			params: ListPolicySetsRequest{
				Contract: "3-WNKXX1",
				Network:  "foo",
			},
			withError: ErrStructValidation,
		},
		"missing contract": {
			params: ListPolicySetsRequest{
				Network: NetworkProduction,
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
			result, err := client.ListPolicySets(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetPolicySet(t *testing.T) {
	tests := map[string]struct {
		params           GetPolicySetRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *PolicySet
		withError        error
	}{
		"200 OK for both networks": {
			params: GetPolicySetRequest{
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
				Contract:    "3-WNKXX1",
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"id": "570f9090-5dbe-11ec-8a0a-71665789c1d8",
				"name": "terraform_beta_v2",
				"region": "US",
				"type": "IMAGE",
				"properties": [
					"jsmith_sqa2"
				],
				"user": "jsmith",
				"lastModified": "2021-12-15 15:47:42+0000"
			}`,
			expectedPath: "/imaging/v2/policysets/570f9090-5dbe-11ec-8a0a-71665789c1d8",
			expectedResponse: &PolicySet{
				ID:           "570f9090-5dbe-11ec-8a0a-71665789c1d8",
				Name:         "terraform_beta_v2",
				Region:       "US",
				Type:         "IMAGE",
				Properties:   []string{"jsmith_sqa2"},
				User:         "jsmith",
				LastModified: "2021-12-15 15:47:42+0000",
			},
		},
		"200 OK for both staging network": {
			params: GetPolicySetRequest{
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
				Contract:    "3-WNKXX1",
				Network:     NetworkStaging,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"id": "570f9090-5dbe-11ec-8a0a-71665789c1d8",
				"name": "terraform_beta_v2",
				"region": "US",
				"type": "IMAGE",
				"properties": [
					"jsmith_sqa2"
				],
				"user": "jsmith",
				"lastModified": "2021-12-15 15:47:42+0000"
			}`,
			expectedPath: "/imaging/v2/network/staging/policysets/570f9090-5dbe-11ec-8a0a-71665789c1d8",
			expectedResponse: &PolicySet{
				ID:           "570f9090-5dbe-11ec-8a0a-71665789c1d8",
				Name:         "terraform_beta_v2",
				Region:       "US",
				Type:         "IMAGE",
				Properties:   []string{"jsmith_sqa2"},
				User:         "jsmith",
				LastModified: "2021-12-15 15:47:42+0000",
			},
		},
		"404 Not found": {
			params: GetPolicySetRequest{
				PolicySetID: "570f9090-5dbe-11ec-8a0a",
				Contract:    "3-WNKXX1",
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
{
    "type": "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_9000",
    "title": "Not Found",
    "instance": "21bde25b-c9a5-4987-b1d5-0c3b92f77b2e",
    "status": 404,
    "detail": "Policy set does not exist.",
    "extensionFields": {
        "requestId": "5f94ea1284bf1800"
    },
    "problemId": "21bde25b-c9a5-4987-b1d5-0c3b92f77b2e",
    "requestId": "5f94ea1284bf1800"
}`,
			expectedPath: "/imaging/v2/policysets/570f9090-5dbe-11ec-8a0a",
			withError: &Error{
				Type:     "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_9000",
				Title:    "Not Found",
				Instance: "21bde25b-c9a5-4987-b1d5-0c3b92f77b2e",
				Status:   404,
				Detail:   "Policy set does not exist.",
				ExtensionFields: map[string]string{
					"requestId": "5f94ea1284bf1800",
				},
				ProblemID: "21bde25b-c9a5-4987-b1d5-0c3b92f77b2e",
				RequestID: "5f94ea1284bf1800",
			},
		},
		// 500
		"missing Policy set Id": {
			params: GetPolicySetRequest{
				Contract: "3-WNKXX1",
				Network:  NetworkProduction,
			},
			withError: ErrStructValidation,
		},
		"invalid network": {
			params: GetPolicySetRequest{
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
				Contract:    "3-WNKXX1",
				Network:     "foo",
			},
			withError: ErrStructValidation,
		},
		"missing contract": {
			params: GetPolicySetRequest{
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
				Network:     NetworkProduction,
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
			result, err := client.GetPolicySet(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreatePolicySet(t *testing.T) {
	tests := map[string]struct {
		params              CreatePolicySetRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *PolicySet
		withError           error
	}{
		"201 created": {
			params: CreatePolicySetRequest{
				Contract: "3-WNKXX1",
				CreatePolicySet: CreatePolicySet{
					Name:   "my_example_token",
					Type:   "IMAGE",
					Region: "US",
				},
			},
			expectedRequestBody: `{"name":"my_example_token","region":"US","type":"IMAGE"}`,
			responseStatus:      http.StatusCreated,
			responseBody: `{
				"name": "my_example_token",
				"id": "29467ef0-751d-11ec-7a0a-71665789c1d8",
				"type": "IMAGE",
				"region": "US",
				"properties": [],
				"user": "ftzgvvigljhoq5ia",
				"lastModified": "2022-01-14 09:34:25+0000"
			}`,
			expectedPath: "/imaging/v2/policysets/",
			expectedResponse: &PolicySet{
				Name:         "my_example_token",
				ID:           "29467ef0-751d-11ec-7a0a-71665789c1d8",
				Type:         "IMAGE",
				Region:       "US",
				Properties:   []string{},
				User:         "ftzgvvigljhoq5ia",
				LastModified: "2022-01-14 09:34:25+0000",
			},
		},
		"400 Bad request": {
			params: CreatePolicySetRequest{
				Contract: "3-WNKXX1",
				CreatePolicySet: CreatePolicySet{
					Name:   "my_example_token",
					Type:   "IMAGE",
					Region: "US",
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
				{
					"type": "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1004",
					"title": "Bad Request",
					"instance": "5ea0274b-2322-4a0a-92ee-fabaa5a84d41",
					"status": 400,
					"detail": "A contract must be specified using the Contract header.",
					"problemId": "5ea0274b-2322-4a0a-92ee-fabaa5a84d41"
				}`,
			expectedPath: "/imaging/v2/policysets/",
			withError: &Error{
				Type:      "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1004",
				Title:     "Bad Request",
				Instance:  "5ea0274b-2322-4a0a-92ee-fabaa5a84d41",
				Status:    400,
				Detail:    "A contract must be specified using the Contract header.",
				ProblemID: "5ea0274b-2322-4a0a-92ee-fabaa5a84d41",
			},
		},
		"missing Policy set Id": {
			params: CreatePolicySetRequest{
				Contract: "3-WNKXX1",
			},
			withError: ErrStructValidation,
		},
		"missing contract": {
			params: CreatePolicySetRequest{
				CreatePolicySet: CreatePolicySet{
					Name:   "my_example_token",
					Type:   "IMAGE",
					Region: "US",
				},
			},
			withError: ErrStructValidation,
		},
		"missing name": {
			params: CreatePolicySetRequest{
				Contract: "3-WNKXX1",
				CreatePolicySet: CreatePolicySet{
					Type:   "IMAGE",
					Region: "US",
				},
			},
			withError: ErrStructValidation,
		},
		"missing type": {
			params: CreatePolicySetRequest{
				Contract: "3-WNKXX1",
				CreatePolicySet: CreatePolicySet{
					Name:   "my_example_token",
					Region: "US",
				},
			},
			withError: ErrStructValidation,
		},
		"invalid type": {
			params: CreatePolicySetRequest{
				Contract: "3-WNKXX1",
				CreatePolicySet: CreatePolicySet{
					Name:   "my_example_token",
					Type:   "INVALID",
					Region: "US",
				},
			},
			withError: ErrStructValidation,
		},
		"missing region": {
			params: CreatePolicySetRequest{
				Contract: "3-WNKXX1",
				CreatePolicySet: CreatePolicySet{
					Name: "my_example_token",
					Type: "IMAGE",
				},
			},
			withError: ErrStructValidation,
		},
		"invalid region": {
			params: CreatePolicySetRequest{
				Contract: "3-WNKXX1",
				CreatePolicySet: CreatePolicySet{
					Name:   "my_example_token",
					Type:   "IMAGE",
					Region: "INVALID",
				},
			},
			withError: ErrStructValidation,
		},
		"validation error": {
			params:    CreatePolicySetRequest{},
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

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreatePolicySet(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdatePolicySet(t *testing.T) {
	tests := map[string]struct {
		params              UpdatePolicySetRequest
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		expectedPath        string
		expectedResponse    *PolicySet
		withError           error
	}{
		"200 updated": {
			params: UpdatePolicySetRequest{
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
				Contract:    "3-WNKXX1",
				UpdatePolicySet: UpdatePolicySet{
					Name:   "my_renamed_token_2",
					Region: "US",
				},
			},
			expectedRequestBody: `{"name":"my_renamed_token_2","region":"US"}`,
			responseStatus:      http.StatusOK,
			responseBody: `{
				"name": "my_renamed_token",
				"id": "1385f880-7477-11ec-8a0a-71665789c1d8",
				"type": "IMAGE",
				"region": "US",
				"properties": [],
				"user": "ftzgvvigljhoq5ib",
				"lastModified": "2022-01-14 10:35:11+0000"
			}`,
			expectedPath: "/imaging/v2/policysets/570f9090-5dbe-11ec-8a0a-71665789c1d8",
			expectedResponse: &PolicySet{
				Name:         "my_renamed_token",
				ID:           "1385f880-7477-11ec-8a0a-71665789c1d8",
				Type:         "IMAGE",
				Region:       "US",
				Properties:   []string{},
				User:         "ftzgvvigljhoq5ib",
				LastModified: "2022-01-14 10:35:11+0000",
			},
		},
		"400 Bad request": {
			params: UpdatePolicySetRequest{
				PolicySetID: "second",
				Contract:    "3-WNKXX1",
				UpdatePolicySet: UpdatePolicySet{
					Name:   "my_renamed_token",
					Region: "US",
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
				{
 				    "type": "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_9000",
					"title": "Bad Request",
					"instance": "b10b7635-b1ab-4742-bd88-82a6fcdd791a",
					"status": 400,
					"detail": "Policy Set does not exist.",
					"extensionFields": {
						"requestId": "5f9605d49c62d759"
					},
					"problemId": "b10b7635-b1ab-4742-bd88-82a6fcdd791a",
					"requestId": "5f9605d49c62d759"
				}`,
			expectedPath: "/imaging/v2/policysets/second",
			withError: &Error{
				Type:     "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_9000",
				Title:    "Bad Request",
				Instance: "b10b7635-b1ab-4742-bd88-82a6fcdd791a",
				Status:   400,
				Detail:   "Policy Set does not exist.",
				ExtensionFields: map[string]string{
					"requestId": "5f9605d49c62d759",
				},
				ProblemID: "b10b7635-b1ab-4742-bd88-82a6fcdd791a",
				RequestID: "5f9605d49c62d759",
			},
		},
		"500 Internal server error": {
			params: UpdatePolicySetRequest{
				PolicySetID: "second",
				Contract:    "3-WNKXX1",
				UpdatePolicySet: UpdatePolicySet{
					Name:   "my_renamed_token",
					Region: "EMEA",
				},
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
				{
					"type": "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_9000",
					"title": "Internal Server Error",
					"instance": "3692e138-f0b7-4479-b12a-48cdee93cf4d",
					"status": 500,
					"detail": "An unexpected error has occurred",
					"extensionFields": {
						"requestId": "5f960786419d5e8d"
					},
					"problemId": "3692e138-f0b7-4479-b12a-48cdee93cf4d",
					"requestId": "5f960786419d5e8d"
				}`,
			expectedPath: "/imaging/v2/policysets/second",
			withError: &Error{
				Type:     "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_9000",
				Title:    "Internal Server Error",
				Instance: "3692e138-f0b7-4479-b12a-48cdee93cf4d",
				Status:   500,
				Detail:   "An unexpected error has occurred",
				ExtensionFields: map[string]string{
					"requestId": "5f960786419d5e8d",
				},
				ProblemID: "3692e138-f0b7-4479-b12a-48cdee93cf4d",
				RequestID: "5f960786419d5e8d",
			},
		},
		"missing Policy set Id": {
			params: UpdatePolicySetRequest{
				Contract: "3-WNKXX1",
			},
			withError: ErrStructValidation,
		},
		"missing contract": {
			params: UpdatePolicySetRequest{
				UpdatePolicySet: UpdatePolicySet{
					Name:   "my_example_token",
					Region: "US",
				},
			},
			withError: ErrStructValidation,
		},
		"missing name": {
			params: UpdatePolicySetRequest{
				Contract: "3-WNKXX1",
				UpdatePolicySet: UpdatePolicySet{
					Region: "US",
				},
			},
			withError: ErrStructValidation,
		},
		"missing region": {
			params: UpdatePolicySetRequest{
				Contract: "3-WNKXX1",
				UpdatePolicySet: UpdatePolicySet{
					Name: "my_example_token",
				},
			},
			withError: ErrStructValidation,
		},
		"invalid region": {
			params: UpdatePolicySetRequest{
				Contract: "3-WNKXX1",
				UpdatePolicySet: UpdatePolicySet{
					Name:   "my_example_token",
					Region: "INVALID",
				},
			},
			withError: ErrStructValidation,
		},
		"validation error": {
			params:    UpdatePolicySetRequest{},
			withError: ErrStructValidation,
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
			result, err := client.UpdatePolicySet(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeletePolicySet(t *testing.T) {
	tests := map[string]struct {
		params              DeletePolicySetRequest
		expectedRequestBody string
		responseStatus      int
		expectedPath        string
		responseBody        string
		withError           error
	}{
		"204 no content (deleted)": {
			params: DeletePolicySetRequest{
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
				Contract:    "3-WNKXX1",
			},
			responseStatus: http.StatusNoContent,
			expectedPath:   "/imaging/v2/policysets/570f9090-5dbe-11ec-8a0a-71665789c1d8",
		},
		"400 Bad request": {
			params: DeletePolicySetRequest{
				PolicySetID: "second",
				Contract:    "3-WNKXX1",
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
				{
					"type": "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1004",
					"title": "Bad Request",
					"instance": "1a7387db-7056-4ae9-8cdf-5f23e0645487",
					"status": 400,
					"detail": "A contract must be specified using the Contract header.",
					"problemId": "1a7387db-7056-4ae9-8cdf-5f23e0645487"
				}`,
			expectedPath: "/imaging/v2/policysets/second",
			withError: &Error{
				Type:      "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1004",
				Title:     "Bad Request",
				Instance:  "1a7387db-7056-4ae9-8cdf-5f23e0645487",
				Status:    400,
				Detail:    "A contract must be specified using the Contract header.",
				ProblemID: "1a7387db-7056-4ae9-8cdf-5f23e0645487",
			},
		},
		"401 Not authorized": {
			params: DeletePolicySetRequest{
				PolicySetID: "second",
				Contract:    "3-WNKAA1",
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `
				{
					"type": "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
					"title": "Not authorized",
					"status": 401,
					"detail": "Inactive client token",
					"instance": "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/imaging/v2/policysets/058d9da0-7477-11ec-8a0a-71665789c1d8",
					"method": "DELETE",
					"serverIp": "104.81.220.242",
					"clientIp": "22.22.22.22",
					"requestId": "981a7be",
					"requestTime": "2022-01-14T09:22:54Z"
				}`,
			expectedPath: "/imaging/v2/policysets/second",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authn/deny",
				Title:       "Not authorized",
				Status:      401,
				Detail:      "Inactive client token",
				Instance:    "https://akaa-p3wvjp6bqtotgpjh-fbk2vczjtq7b5l6a.luna-dev.akamaiapis.net/imaging/v2/policysets/058d9da0-7477-11ec-8a0a-71665789c1d8",
				Method:      "DELETE",
				ServerIP:    "104.81.220.242",
				ClientIP:    "22.22.22.22",
				RequestID:   "981a7be",
				RequestTime: "2022-01-14T09:22:54Z",
			},
		},
		"404 Not found": {
			params: DeletePolicySetRequest{
				PolicySetID: "second",
				Contract:    "3-WNKXX1",
			},
			responseStatus: http.StatusNotFound,
			responseBody: `
				{
					"type": "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_3001",
					"title": "Not Found",
					"instance": "3da0fcf0-4fb9-4b80-986c-4f9993436189",
					"status": 404,
					"detail": "That policy set does not exist.",
					"problemId": "3da0fcf0-4fb9-4b80-986c-4f9993436189"
				}`,
			expectedPath: "/imaging/v2/policysets/second",
			withError: &Error{
				Type:      "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_3001",
				Title:     "Not Found",
				Instance:  "3da0fcf0-4fb9-4b80-986c-4f9993436189",
				Status:    404,
				Detail:    "That policy set does not exist.",
				ProblemID: "3da0fcf0-4fb9-4b80-986c-4f9993436189",
			},
		},
		// 500
		"missing Policy set Id": {
			params: DeletePolicySetRequest{
				Contract: "3-WNKXX1",
			},
			withError: ErrStructValidation,
		},
		"missing contract": {
			params: DeletePolicySetRequest{
				PolicySetID: "570f9090-5dbe-11ec-8a0a-71665789c1d8",
			},
			withError: ErrStructValidation,
		},
		"validation error": {
			params:    DeletePolicySetRequest{},
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

				if len(test.expectedRequestBody) > 0 {
					body, err := ioutil.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, test.expectedRequestBody, string(body))
				}
			}))
			client := mockAPIClient(t, mockServer)
			err := client.DeletePolicySet(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
