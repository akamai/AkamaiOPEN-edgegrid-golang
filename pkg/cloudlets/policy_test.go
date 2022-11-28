package cloudlets

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/tools"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestListPolicies(t *testing.T) {
	tests := map[string]struct {
		params           ListPoliciesRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse []Policy
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params:         ListPoliciesRequest{},
			responseStatus: http.StatusOK,
			responseBody: `
[
	{
        "location": "/cloudlets/api/v2/policies/1001",
        "serviceVersion": null,
        "apiVersion": "2.0",
        "policyId": 1001,
        "cloudletId": 99,
        "cloudletCode": "CC",
        "groupId": 1234,
        "name": "CookieCutter",
        "description": "Custom cookie cutter",
        "propertyName": "www.example.org",
        "createdBy": "sjones",
        "createDate": 1400535431324,
        "lastModifiedBy": "sjones",
        "lastModifiedDate": 1441829042000,
        "activations": [
            {
                "serviceVersion": null,
                "apiVersion": "2.0",
                "network": "prod",
                "policyInfo": {
                    "policyId": 1001,
                    "name": "CookieCutter",
                    "version": 2,
                    "status": "INACTIVE",
                    "statusDetail": "waiting to complete tests in test environment",
                    "activatedBy": "jsmith",
                    "activationDate": 1441829042000
                },
                "propertyInfo": {
                    "name": "www.example.org",
                    "version": 2,
                    "groupId": 1234,
                    "status": "INACTIVE",
                    "activatedBy": "sjones",
                    "activationDate": 1441137842000
                }
            },
            {
                "serviceVersion": null,
                "apiVersion": "2.0",
                "network": "STAGING",
                "policyInfo": {
                    "policyId": 1001,
                    "name": "CookieCutter",
                    "version": 22,
                    "status": "ACTIVE",
                    "statusDetail": "testing",
                    "activatedBy": "jsmith",
                    "activationDate": 1400535431000
                },
                "propertyInfo": {
                    "name": "www.example.org",
                    "version": 22,
                    "groupId": 1234,
                    "status": "ACTIVE",
                    "activatedBy": "jsmith",
                    "activationDate": 1441137842000
                }
            }
        ]
    }
]`,
			expectedPath: "/cloudlets/api/v2/policies?includeDeleted=false&offset=0",
			expectedResponse: []Policy{
				{
					Location:         "/cloudlets/api/v2/policies/1001",
					APIVersion:       "2.0",
					PolicyID:         1001,
					CloudletID:       99,
					CloudletCode:     "CC",
					GroupID:          1234,
					Name:             "CookieCutter",
					Description:      "Custom cookie cutter",
					CreatedBy:        "sjones",
					CreateDate:       1400535431324,
					LastModifiedBy:   "sjones",
					LastModifiedDate: 1441829042000,
					Activations: []PolicyActivation{
						{
							APIVersion: "2.0",
							Network:    "prod",
							PolicyInfo: PolicyInfo{
								PolicyID:       1001,
								Name:           "CookieCutter",
								Version:        2,
								Status:         "INACTIVE",
								StatusDetail:   "waiting to complete tests in test environment",
								ActivatedBy:    "jsmith",
								ActivationDate: 1441829042000,
							},
							PropertyInfo: PropertyInfo{
								Name:           "www.example.org",
								Version:        2,
								GroupID:        1234,
								Status:         "INACTIVE",
								ActivatedBy:    "sjones",
								ActivationDate: 1441137842000,
							},
						},
						{
							APIVersion: "2.0",
							Network:    "staging",
							PolicyInfo: PolicyInfo{
								PolicyID:       1001,
								Name:           "CookieCutter",
								Version:        22,
								Status:         "ACTIVE",
								StatusDetail:   "testing",
								ActivatedBy:    "jsmith",
								ActivationDate: 1400535431000,
							},
							PropertyInfo: PropertyInfo{
								Name:           "www.example.org",
								Version:        22,
								GroupID:        1234,
								Status:         "ACTIVE",
								ActivatedBy:    "jsmith",
								ActivationDate: 1441137842000,
							},
						},
					},
				},
			},
		},
		"200 OK with params": {
			params: ListPoliciesRequest{
				CloudletID:     tools.Int64Ptr(2),
				IncludeDeleted: true,
				Offset:         4,
				PageSize:       tools.IntPtr(5),
			},
			responseStatus:   http.StatusOK,
			responseBody:     `[]`,
			expectedPath:     "/cloudlets/api/v2/policies?cloudletId=2&includeDeleted=true&offset=4&pageSize=5",
			expectedResponse: []Policy{},
		},
		"500 internal server error": {
			params:         ListPoliciesRequest{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error making request",
  "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies?includeDeleted=false&offset=0",
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
			result, err := client.ListPolicies(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetPolicy(t *testing.T) {
	tests := map[string]struct {
		policyID         int64
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Policy
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			policyID:       1001,
			responseStatus: http.StatusOK,
			responseBody: `
    {
        "location": "/cloudlets/api/v2/policies/1001",
        "serviceVersion": null,
        "apiVersion": "2.0",
        "policyId": 1001,
        "cloudletId": 99,
        "cloudletCode": "CC",
        "groupId": 1234,
        "name": "CookieCutter",
        "description": "Custom cookie cutter",
        "propertyName": "www.example.org",
        "createdBy": "sjones",
        "createDate": 1400535431324,
        "lastModifiedBy": "sjones",
        "lastModifiedDate": 1441829042000,
        "activations": [
            {
                "serviceVersion": null,
                "apiVersion": "2.0",
                "network": "prod",
                "policyInfo": {
                    "policyId": 1001,
                    "name": "CookieCutter",
                    "version": 2,
                    "status": "INACTIVE",
                    "statusDetail": "waiting to complete tests in test environment",
                    "activatedBy": "jsmith",
                    "activationDate": 1441829042000
                },
                "propertyInfo": {
                    "name": "www.example.org",
                    "version": 2,
                    "groupId": 1234,
                    "status": "INACTIVE",
                    "activatedBy": "sjones",
                    "activationDate": 1441137842000
                }
            },
            {
                "serviceVersion": null,
                "apiVersion": "2.0",
                "network": "staging",
                "policyInfo": {
                    "policyId": 1001,
                    "name": "CookieCutter",
                    "version": 22,
                    "status": "ACTIVE",
                    "statusDetail": "testing",
                    "activatedBy": "jsmith",
                    "activationDate": 1400535431000
                },
                "propertyInfo": {
                    "name": "www.example.org",
                    "version": 22,
                    "groupId": 1234,
                    "status": "ACTIVE",
                    "activatedBy": "jsmith",
                    "activationDate": 1441137842000
                }
            }
        ]
    }`,
			expectedPath: "/cloudlets/api/v2/policies/1001",
			expectedResponse: &Policy{
				Location:         "/cloudlets/api/v2/policies/1001",
				APIVersion:       "2.0",
				PolicyID:         1001,
				CloudletID:       99,
				CloudletCode:     "CC",
				GroupID:          1234,
				Name:             "CookieCutter",
				Description:      "Custom cookie cutter",
				CreatedBy:        "sjones",
				CreateDate:       1400535431324,
				LastModifiedBy:   "sjones",
				LastModifiedDate: 1441829042000,
				Activations: []PolicyActivation{
					{
						APIVersion: "2.0",
						Network:    "prod",
						PolicyInfo: PolicyInfo{
							PolicyID:       1001,
							Name:           "CookieCutter",
							Version:        2,
							Status:         "INACTIVE",
							StatusDetail:   "waiting to complete tests in test environment",
							ActivatedBy:    "jsmith",
							ActivationDate: 1441829042000,
						},
						PropertyInfo: PropertyInfo{
							Name:           "www.example.org",
							Version:        2,
							GroupID:        1234,
							Status:         "INACTIVE",
							ActivatedBy:    "sjones",
							ActivationDate: 1441137842000,
						},
					},
					{
						APIVersion: "2.0",
						Network:    "staging",
						PolicyInfo: PolicyInfo{
							PolicyID:       1001,
							Name:           "CookieCutter",
							Version:        22,
							Status:         "ACTIVE",
							StatusDetail:   "testing",
							ActivatedBy:    "jsmith",
							ActivationDate: 1400535431000,
						},
						PropertyInfo: PropertyInfo{
							Name:           "www.example.org",
							Version:        22,
							GroupID:        1234,
							Status:         "ACTIVE",
							ActivatedBy:    "jsmith",
							ActivationDate: 1441137842000,
						},
					},
				},
			},
		},
		"500 internal server error": {
			policyID:       1,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies/1",
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
			result, err := client.GetPolicy(context.Background(), GetPolicyRequest{PolicyID: test.policyID})
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestCreatePolicy(t *testing.T) {
	tests := map[string]struct {
		request          CreatePolicyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Policy
		withError        error
	}{
		"201 created": {
			request: CreatePolicyRequest{
				CloudletID: 0,
				GroupID:    35730,
				Name:       "TestName1",
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
    "activations": [],
    "apiVersion": "2.0",
    "cloudletCode": "ER",
    "cloudletId": 0,
    "createDate": 1629299944251,
    "createdBy": "jsmith",
    "deleted": false,
    "description": null,
    "groupId": 35730,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629299944251,
    "location": "/cloudlets/api/v2/policies/276858",
    "name": "TestName1",
    "policyId": 276858,
    "propertyName": null,
    "serviceVersion": null
}`,
			expectedPath: "/cloudlets/api/v2/policies",
			expectedResponse: &Policy{
				APIVersion:       "2.0",
				CloudletCode:     "ER",
				CloudletID:       0,
				CreateDate:       1629299944251,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "",
				GroupID:          35730,
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629299944251,
				Location:         "/cloudlets/api/v2/policies/276858",
				Name:             "TestName1",
				PolicyID:         276858,
				Activations:      []PolicyActivation{},
			},
		},
		"500 internal server error": {
			request: CreatePolicyRequest{
				CloudletID: 0,
				GroupID:    35730,
				Name:       "TestName1",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error creating enrollment",
  "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating enrollment",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			request:   CreatePolicyRequest{},
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
			result, err := client.CreatePolicy(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDeletePolicy(t *testing.T) {
	tests := map[string]struct {
		policyID       int64
		responseStatus int
		responseBody   string
		expectedPath   string
		withError      error
	}{
		"204 no content": {
			policyID:       276858,
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			expectedPath:   "/cloudlets/api/v2/policies/276858",
		},
		"500 internal server error": {
			policyID:       0,
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error creating enrollment",
  "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies/0",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating enrollment",
				StatusCode: http.StatusInternalServerError,
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
			err := client.RemovePolicy(context.Background(), RemovePolicyRequest{PolicyID: test.policyID})
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdatePolicy(t *testing.T) {
	tests := map[string]struct {
		request          UpdatePolicyRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *Policy
		withError        error
	}{
		"200 updated": {
			request: UpdatePolicyRequest{
				UpdatePolicy: UpdatePolicy{
					GroupID:     35730,
					Name:        "TestName1Updated",
					Description: "Description",
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusOK,
			responseBody: `{
    "activations": [],
    "apiVersion": "2.0",
    "cloudletCode": "ER",
    "cloudletId": 0,
    "createDate": 1629299944251,
    "createdBy": "jsmith",
    "deleted": false,
    "description": "Description",
    "groupId": 35730,
    "lastModifiedBy": "jsmith",
    "lastModifiedDate": 1629370566748,
    "location": "/cloudlets/api/v2/policies/276858",
    "name": "TestName1Updated",
    "policyId": 276858,
    "propertyName": null,
    "serviceVersion": null
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858",
			expectedResponse: &Policy{
				APIVersion:       "2.0",
				CloudletCode:     "ER",
				CloudletID:       0,
				CreateDate:       1629299944251,
				CreatedBy:        "jsmith",
				Deleted:          false,
				Description:      "Description",
				GroupID:          35730,
				LastModifiedBy:   "jsmith",
				LastModifiedDate: 1629370566748,
				Location:         "/cloudlets/api/v2/policies/276858",
				Name:             "TestName1Updated",
				PolicyID:         276858,
				Activations:      []PolicyActivation{},
			},
		},
		"500 internal server error": {
			request: UpdatePolicyRequest{
				UpdatePolicy: UpdatePolicy{
					GroupID: 35730,
					Name:    "TestName1Updated",
				},
				PolicyID: 276858,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error creating enrollment",
  "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/policies/276858",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating enrollment",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			expectedPath: "/cloudlets/api/v2/policies/0",
			request: UpdatePolicyRequest{
				UpdatePolicy: UpdatePolicy{
					Name: "A B",
				},
			},
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
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdatePolicy(context.Background(), test.request)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
