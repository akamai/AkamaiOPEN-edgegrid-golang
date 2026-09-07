package v0

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func TestAPI_Validate(t *testing.T) {
	tests := map[string]struct {
		req       RegisterAPIRequest
		withError func(*testing.T, error)
	}{
		"min ok": {
			req: RegisterAPIRequest{
				ContractID: "Contract-1",
				GroupID:    1,
				APIAttributes: APIAttributes{
					Name:      "name",
					Hostnames: []string{"akamai.com"},
				},
			},
		},
		"empty": {
			req: RegisterAPIRequest{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "ContractID: cannot be blank\nGroupID: cannot be blank\nHostnames: cannot be blank\nName: cannot be blank", err.Error())
			},
		},
		"security schemes - in : invalid enum value": {
			req: RegisterAPIRequest{
				APIAttributes: APIAttributes{
					SecuritySchemes: &SecuritySchemes{
						APIKey: &SecurityScheme{
							In: ptr.To(SecuritySchemeLocation("HEADER")),
						},
					},
				}},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "In: value 'HEADER' is invalid. Must be one of: 'cookie', 'header', 'query'")
			},
		},
		"constraints - ConsumeType : invalid enum value": {
			req: RegisterAPIRequest{
				APIAttributes: APIAttributes{
					Constraints: &Constraints{
						RequestBody: &ConstraintsRequestBody{
							ConsumeType: []ConsumeType{"JSON", "XML", "urlencoded", "ANY"},
						},
					},
				}},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Constraints: {\n\tRequestBody: {\n\t\t0: value 'JSON' is invalid. Must be one of: 'json', 'xml', 'urlencoded', 'any'\n\t\t1: value 'XML' is invalid. Must be one of: 'json', 'xml', 'urlencoded', 'any'\n\t\t3: value 'ANY' is invalid. Must be one of: 'json', 'xml', 'urlencoded', 'any'\n\t}\n}\nContractID: cannot be blank\nGroupID: cannot be blank\nHostnames: cannot be blank\nName: cannot be blank")
			},
		},
		"resources - resource Method in : invalid enum value": {
			req: RegisterAPIRequest{
				ContractID: "3-XXXXXX",
				GroupID:    33333,
				APIAttributes: APIAttributes{
					Name:      "bookstore API",
					Hostnames: []string{"akamai.com"},
					Resources: orderedmap.New[string, Resource](orderedmap.WithInitialData[string, Resource](
						orderedmap.Pair[string, Resource]{
							Key: "/books",
							Value: Resource{
								Name:        "Books Resource",
								Description: ptr.To("Books Resource description"),
								Post: &Method{
									Parameters: []Parameter{
										{
											Name:        "limit",
											In:          "QUERY",
											Type:        "integer",
											Required:    true,
											Description: ptr.To("limit parameter"),
											Minimum:     ptr.To(float32(1)),
											Maximum:     ptr.To(float32(2)),
										},
										{
											Name:        "query",
											In:          "query",
											Type:        "string",
											Required:    true,
											Description: ptr.To("query parameter"),
											MinLength:   ptr.To(int64(1)),
											MaxLength:   ptr.To(int64(2)),
										},
									},
									RequestBody: orderedmap.New[string, Property](orderedmap.WithInitialData[string, Property](
										orderedmap.Pair[string, Property]{
											Key: "json",
											Value: Property{
												Name:        "Book Body",
												Type:        "object",
												Required:    true,
												Description: ptr.To("Json body desciption"),
												MaxBodySize: ptr.To(MaxBodySize("16kb")),
												Properties: []Property{
													{
														Name:      "name",
														Type:      "string",
														Required:  true,
														MinLength: ptr.To(int64(1)),
														MaxLength: ptr.To(int64(200)),
													},
													{
														Name:      "limit",
														Type:      "NUMBER",
														Required:  true,
														MinLength: ptr.To(int64(1)),
														MaxLength: ptr.To(int64(200)),
													},
												},
											},
										},
									)),
									Responses: &Responses{
										Headers: []Parameter{
											{
												Name:     "Set-Cookie",
												In:       "HEADER",
												Type:     "STRING",
												Required: true,
											},
											{
												Name:     "Max-Forwards",
												In:       "header",
												Type:     "integer",
												Required: true,
											},
										},
										Contents: []ResponseContent{
											{
												StatusCodes: []int64{20},
												JSON: &Property{
													Name:        "application/json",
													Type:        "ARRAY",
													Required:    false,
													MaxBodySize: ptr.To(MaxBodySize("16KB")),
													Items: &Items{
														Type: "STRING",
														Properties: []Property{
															{
																Name: "name",
																Type: "STRING",
															},
															{
																Name: "name",
																Type: "OBJECT",
															},
															{
																Name: "last_name",
																Type: "string",
															},
															{
																Name: "last_name",
																Type: "object",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						})),
				}},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "/books: {\n\tPost: {\n\t\tParameters[0]: {\n\t\t\tIn: value 'QUERY' is invalid. Must be one of: 'cookie', 'query', 'header', 'path'\n\t\t}\n\t\tResponses: {\n\t\t\tContents[0]: {\n\t\t\t\tJSON: {\n\t\t\t\t\tItems: {\n\t\t\t\t\t\tProperties[0]: {\n\t\t\t\t\t\t\tType: value 'STRING' is invalid. Must be one of: 'number', 'integer', 'string', 'boolean', 'object', 'array'\n\t\t\t\t\t\t}\n\t\t\t\t\t\tProperties[1]: {\n\t\t\t\t\t\t\tType: value 'OBJECT' is invalid. Must be one of: 'number', 'integer', 'string', 'boolean', 'object', 'array'\n\t\t\t\t\t\t}\n\t\t\t\t\t\tType: value 'STRING' is invalid. Must be one of: 'number', 'integer', 'string', 'boolean', 'object', 'array'\n\t\t\t\t\t}\n\t\t\t\t\tType: value 'ARRAY' is invalid. Must be one of: 'number', 'integer', 'string', 'boolean', 'object', 'array'\n\t\t\t\t}\n\t\t\t}\n\t\t\tHeaders[0]: {\n\t\t\t\tIn: value 'HEADER' is invalid. Must be one of: 'cookie', 'query', 'header', 'path'\n\t\t\t\tType: value 'STRING' is invalid. Must be one of: 'number', 'integer', 'string', 'boolean'\n\t\t\t}\n\t\t}\n\t\tjson: {\n\t\t\tMaxBodySize: value '16kb' is invalid. Must be one of: '6KB', '8KB', '12KB', '16KB' \n\t\t\tProperties[1]: {\n\t\t\t\tType: value 'NUMBER' is invalid. Must be one of: 'number', 'integer', 'string', 'boolean', 'object', 'array'\n\t\t\t}\n\t\t}\n\t}\n}")
			},
		},
		"resources - resource method RequestBody in : invalid enum value": {
			req: RegisterAPIRequest{
				ContractID: "3-XXXXXX",
				GroupID:    33333,
				APIAttributes: APIAttributes{
					Name:      "bookstore API",
					Hostnames: []string{"akamai.com"},
					Resources: orderedmap.New[string, Resource](orderedmap.WithInitialData[string, Resource](
						orderedmap.Pair[string, Resource]{
							Key: "/books",
							Value: Resource{
								Name:        "Books Resource",
								Description: ptr.To("Books Resource description"),
								Post: &Method{
									Parameters: []Parameter{
										{
											Name:        "limit",
											In:          "query",
											Type:        "integer",
											Required:    true,
											Description: ptr.To("limit parameter"),
											Minimum:     ptr.To(float32(1)),
											Maximum:     ptr.To(float32(2)),
										},
									},
									RequestBody: orderedmap.New[string, Property](orderedmap.WithInitialData[string, Property](
										orderedmap.Pair[string, Property]{
											Key: "json",
											Value: Property{
												Name:        "Book Body",
												Type:        "OBJECT",
												Required:    true,
												Description: ptr.To("Json body desciption"),
												Properties: []Property{
													{
														Name:      "name",
														Type:      "string",
														Required:  true,
														MinLength: ptr.To(int64(1)),
														MaxLength: ptr.To(int64(200)),
													},
												},
											},
										},
									)),
									Responses: &Responses{
										Contents: []ResponseContent{
											{
												StatusCodes: []int64{20},
												JSON: &Property{
													Name:     "application/json",
													Type:     "array",
													Required: false,
													Items: &Items{
														Type: "string",
													},
												},
											},
										},
									},
								},
							},
						})),
				}},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "/books: {\n\tPost: {\n\t\tjson: {\n\t\t\tType: value 'OBJECT' is invalid. Must be one of: 'number', 'integer', 'string', 'boolean', 'object', 'array'\n\t\t}\n\t}\n}")
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestVersioning_Validate(t *testing.T) {
	tests := map[string]struct {
		req       Versioning
		withError func(*testing.T, error)
	}{
		"min ok": {
			req: Versioning{},
		},
		"versioning - invalid 'in' ": {
			req: Versioning{
				In: ptr.To(VersioningLocation("invalid-location")),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "In: value 'invalid-location' is invalid. Must be one of: 'header', 'path', 'query'.", err.Error())
			},
		},
		"versioning - invalid uppercased enum value": {
			req: Versioning{
				In: ptr.To(VersioningLocation("HEADER")),
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "In: value 'HEADER' is invalid. Must be one of: 'header', 'path', 'query'.", err.Error())
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestConstraints_Validate(t *testing.T) {
	tests := map[string]struct {
		req       Constraints
		withError func(*testing.T, error)
	}{
		"min ok": {
			req: Constraints{},
		},
		"request body constraints - invalid consumeType ": {
			req: Constraints{
				RequestBody: &ConstraintsRequestBody{
					ConsumeType: []ConsumeType{"invalid"},
				},
			},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "RequestBody: (ConsumeType: (0: value 'invalid' is invalid. Must be one of: 'json', 'xml', 'urlencoded', 'any'.).).", err.Error())
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestParameter_Validate(t *testing.T) {
	tests := map[string]struct {
		req       Parameter
		withError func(*testing.T, error)
	}{
		"min ok": {
			req: Parameter{
				Name: "Name",
				Type: ParameterTypeInteger,
				In:   ParameterLocationPath,
			},
		},
		"empty": {
			req: Parameter{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "Name: cannot be blank; Type: cannot be blank.", err.Error())
			},
		},
		"invalid type": {
			req: Parameter{
				Type: "invalid-type",
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Type: value 'invalid-type' is invalid. Must be one of: 'number', 'integer', 'string', 'boolean'.")
			},
		},
		"invalid location": {
			req: Parameter{
				In: "invalid-location",
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "In: value 'invalid-location' is invalid. Must be one of: 'cookie', 'query', 'header', 'path'")
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestProperty_Validate(t *testing.T) {
	tests := map[string]struct {
		req       Property
		withError func(*testing.T, error)
	}{
		"min ok": {
			req: Property{
				Name: "Name",
				Type: PropertyTypeObject,
			},
		},
		"empty": {
			req: Property{},
			withError: func(t *testing.T, err error) {
				assert.Equal(t, "Name: cannot be blank; Type: cannot be blank.", err.Error())
			},
		},
		"invalid type": {
			req: Property{
				Type: "invalid-type",
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Type: value 'invalid-type' is invalid. Must be one of: 'number', 'integer', 'string', 'boolean', 'object', 'array'.")
			},
		},
		"invalid maxBodySize": {
			req: Property{
				MaxBodySize: ptr.To(MaxBodySize("invalid-body-size")),
			},
			withError: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "MaxBodySize: value 'invalid-body-size' is invalid. Must be one of: '6KB', '8KB', '12KB', '16KB'")
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.req.Validate()
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestResources_OrderShouldBePreservedDuringSerialization(t *testing.T) {
	var api = APIAttributes{}
	input := []byte(loadJson("testdata/api_resources_ordering.json"))
	err := json.Unmarshal(input, &api)
	if err != nil {
		panic(err)
	}
	output, err := json.Marshal(api)

	if err != nil {
		panic(err)
	}
	assert.Equal(t, string(input), string(output))
}

func TestResources_InsertOrderShouldBePreserved(t *testing.T) {
	var api = APIAttributes{
		Resources: orderedmap.New[string, Resource](),
	}
	var keys = []string{"c", "v", "b", "n", "m", "a", "s", "d", "e", "q", "g"}

	for _, value := range keys {
		api.Resources.Set(value, Resource{})
	}
	i := 0
	for pair := api.Resources.Oldest(); pair != nil; pair = pair.Next() {
		assert.Equal(t, keys[i], pair.Key)
		i++
	}
}

func TestRequestContentTypes_InsertOrderShouldBePreserved(t *testing.T) {
	var method = Method{
		RequestBody: orderedmap.New[string, Property](),
	}
	var keys = []string{"c", "v", "b", "n", "m", "a", "s", "d", "e", "q", "g"}

	for _, value := range keys {
		method.RequestBody.Set(value, Property{})
	}
	i := 0
	for pair := method.RequestBody.Oldest(); pair != nil; pair = pair.Next() {
		assert.Equal(t, keys[i], pair.Key)
		i++
	}
}

func TestRegisterAPI(t *testing.T) {
	tests := map[string]struct {
		body                RegisterAPIRequest
		expectedPath        string
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		withError           func(*testing.T, error)
		expectedResult      *API
	}{
		"201 created - required parameters only": {
			responseStatus: http.StatusCreated,
			body: RegisterAPIRequest{
				ContractID: "3-XXXXXX",
				GroupID:    33333,
				APIAttributes: APIAttributes{
					Name:      "bookstore API",
					Hostnames: []string{"akamai.com"},
				},
			},
			expectedRequestBody: `{
    "name": "bookstore API",
    "contractId": "3-XXXXXX",
    "groupId": 33333,
    "hostnames": [
        "akamai.com"
    ]
}`,
			expectedPath: "/api-definitions/v0/endpoints",
			responseBody: loadJson("testdata/book_store_api_required_only.json"),
			expectedResult: &API{
				RegisterAPIRequest: RegisterAPIRequest{
					APIAttributes: APIAttributes{
						Name:                      "bookstore API",
						Hostnames:                 []string{"akamai.com"},
						MatchCaseSensitive:        false,
						EnableAPIGateway:          false,
						MatchPathSegmentParameter: false,
						GraphQL:                   false,
					},
					ContractID: "3-XXXXXX",
					GroupID:    33333,
				},
				ID:            ptr.To(int64(52)),
				RecordVersion: ptr.To(int64(1)),
			},
		},
		"201 created": {
			responseStatus:      http.StatusCreated,
			body:                bookStoreAPI.RegisterAPIRequest,
			expectedResult:      &bookStoreAPI,
			expectedPath:        "/api-definitions/v0/endpoints",
			expectedRequestBody: loadJson("testdata/book_store_api.json"),
			responseBody:        loadJson("testdata/book_store_api.json"),
		},
		"400 bad request": {
			responseStatus:      http.StatusBadRequest,
			body:                bookStoreAPI.RegisterAPIRequest,
			expectedResult:      &bookStoreAPI,
			expectedPath:        "/api-definitions/v0/endpoints",
			expectedRequestBody: loadJson("testdata/book_store_api.json"),
			responseBody:        loadJson("testdata/400_bad_request.json"),
			withError: func(t *testing.T, err error) {
				want := &error400BadRequest
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, test.expectedRequestBody, string(requestBody))
				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.RegisterAPI(context.Background(), test.body)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, RegisterAPIResponse(*test.expectedResult), *result)
		})
	}
}

func TestUpdateAPIVersion(t *testing.T) {
	tests := map[string]struct {
		body                API
		expectedPath        string
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		withError           func(*testing.T, error)
		expectedResult      *API
	}{
		"200 OK": {
			responseStatus:      http.StatusOK,
			body:                bookStoreAPI,
			expectedResult:      &bookStoreAPI,
			expectedPath:        "/api-definitions/v0/endpoints/1/versions/2",
			expectedRequestBody: loadJson("testdata/book_store_api.json"),
			responseBody:        loadJson("testdata/book_store_api.json"),
		},
		"404 Not Found": {
			responseStatus: http.StatusNotFound,
			body:           bookStoreAPI,
			expectedResult: &bookStoreAPI,
			expectedPath:   "/api-definitions/v0/endpoints/1/versions/2",
			responseBody:   loadJson("testdata/404_not_found.json"),
			withError: func(t *testing.T, err error) {
				want := &error404NotFound
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"409 Conflict": {
			responseStatus: http.StatusConflict,
			body:           bookStoreAPI,
			expectedResult: &bookStoreAPI,
			expectedPath:   "/api-definitions/v0/endpoints/1/versions/2",
			responseBody:   loadJson("testdata/409_conflict.json"),
			withError: func(t *testing.T, err error) {
				want := &error409Conflict
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)
				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				if test.expectedRequestBody != "" {
					assert.JSONEq(t, test.expectedRequestBody, string(requestBody))
				}
				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateAPIVersion(context.Background(), UpdateAPIVersionRequest{
				ID:      1,
				Version: 2,
				Body:    UpdateAPIVersionRequestBody(test.body),
			})
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, UpdateAPIVersionResponse(*test.expectedResult), *result)
		})
	}
}

func TestGetAPIVersion(t *testing.T) {
	tests := map[string]struct {
		expectedPath        string
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		withError           func(*testing.T, error)
		expectedResult      *API
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			expectedResult: &bookStoreAPI,
			expectedPath:   "/api-definitions/v0/endpoints/1/versions/2",
			responseBody:   loadJson("testdata/book_store_api.json"),
		},
		"404 Not Found": {
			responseStatus: http.StatusNotFound,
			expectedResult: &bookStoreAPI,
			expectedPath:   "/api-definitions/v0/endpoints/1/versions/2",
			responseBody:   loadJson("testdata/404_not_found.json"),
			withError: func(t *testing.T, err error) {
				want := &error404NotFound
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"409 Conflict": {
			responseStatus: http.StatusConflict,
			expectedResult: &bookStoreAPI,
			expectedPath:   "/api-definitions/v0/endpoints/1/versions/2",
			responseBody:   loadJson("testdata/409_conflict.json"),
			withError: func(t *testing.T, err error) {
				want := &error409Conflict
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				if test.expectedRequestBody != "" {
					assert.JSONEq(t, test.expectedRequestBody, string(requestBody))
				}
				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.GetAPIVersion(context.Background(), GetAPIVersionRequest{
				ID:      1,
				Version: 2,
			})
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, GetAPIVersionResponse(*test.expectedResult), *result)
		})
	}
}

func TestFromOpenAPIFile(t *testing.T) {
	tests := map[string]struct {
		expectedPath   string
		responseStatus int
		responseBody   string
		withError      func(*testing.T, error)
		expectedResult *FromOpenAPIFileResponse
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			expectedResult: &FromOpenAPIFileResponse{
				Problems: []Error{},
				API:      bookStoreAPI.APIAttributes,
			},
			expectedPath: "/api-definitions/v0/endpoints/openapi",
			responseBody: loadJson("testdata/book_store_api_from_openapi.json"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPost, r.Method)
				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Contains(t, string(requestBody), "Content-Disposition: form-data; name=\"importFile\"; filename=\"api.json\"")
				assert.Contains(t, string(requestBody), "Content-Type: application/octet-stream")
				assert.Contains(t, string(requestBody), "Content-Disposition: form-data; name=\"root\"")
				assert.Contains(t, string(requestBody), "api.yaml")
				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.FromOpenAPIFile(context.Background(), FromOpenAPIFileRequest{
				Content:  []byte("zip archive with Open API Files"),
				RootFile: ptr.To("api.yaml"),
			})
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}

func TestToOpenAPI(t *testing.T) {
	tests := map[string]struct {
		expectedPath        string
		expectedRequestBody string
		responseStatus      int
		responseBody        string
		withError           func(*testing.T, error)
		expectedResult      ToOpenAPIFileResponse
	}{
		"200 OK": {
			responseStatus: http.StatusOK,
			expectedResult: ToOpenAPIFileResponse("OpenAPI File"),
			expectedPath:   "/api-definitions/v0/endpoints/1/versions/2/openapi",
			responseBody:   "OpenAPI File",
		},
		"404 Not Found": {
			responseStatus: http.StatusNotFound,
			expectedResult: ToOpenAPIFileResponse("OpenAPI File"),
			expectedPath:   "/api-definitions/v0/endpoints/1/versions/2/openapi",
			responseBody:   loadJson("testdata/404_not_found.json"),
			withError: func(t *testing.T, err error) {
				want := &error404NotFound
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				requestBody, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				if test.expectedRequestBody != "" {
					assert.JSONEq(t, test.expectedRequestBody, string(requestBody))
				}
				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			defer mockServer.Close()
			client := mockAPIClient(t, mockServer)
			result, err := client.ToOpenAPIFile(context.Background(), ToOpenAPIFileRequest{
				ID:      1,
				Version: 2,
			})
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResult, *result)
		})
	}
}

var bookStoreAPI = API{
	RegisterAPIRequest: RegisterAPIRequest{
		ContractID: "3-XXXXXX",
		GroupID:    33333,
		APIAttributes: APIAttributes{
			Name:                      "bookstore API",
			Hostnames:                 []string{"akamai.com"},
			BasePath:                  ptr.To("/api"),
			Tags:                      []string{"Tag1", "Tag2"},
			Description:               ptr.To("desc"),
			MatchCaseSensitive:        true,
			EnableAPIGateway:          true,
			MatchPathSegmentParameter: true,
			GraphQL:                   true,
			SecuritySchemes: &SecuritySchemes{
				APIKey: &SecurityScheme{
					In:   ptr.To(SecuritySchemeLocationHeader),
					Name: ptr.To("Authorization"),
				},
			},
			Constraints: &Constraints{
				EnforceOn: &EnforceOn{
					Request:  ptr.To(true),
					Response: ptr.To(true),
					UndefinedMethods: &UndefinedMethods{
						Get: true,
					},
					UndefinedParameters: &UndefinedParameters{
						RequestCookie: true,
					},
				},
				RequestBody: &ConstraintsRequestBody{
					ConsumeType:     []ConsumeType{ConsumeTypeJSON, ConsumeTypeXML},
					MaxBodySize:     ptr.To(int64(256)),
					MaxNestingDepth: ptr.To(int64(10)),
					Properties: &ConstraintsRequestBodyProperties{
						MaxStringLength: ptr.To(int64(1000)),
						MaxIntegerValue: ptr.To(int64(12345678)),
						MaxCount:        ptr.To(int64(200)),
						MaxNameLength:   ptr.To(int64(50)),
					},
				},
			},
			Versioning: &Versioning{
				In:    ptr.To(VersioningLocationHeader),
				Name:  ptr.To("Version"),
				Value: ptr.To("1"),
			},
			Resources: orderedmap.New[string, Resource](orderedmap.WithInitialData[string, Resource](
				orderedmap.Pair[string, Resource]{
					Key: "/books",
					Value: Resource{
						Name:        "Books Resource",
						Description: ptr.To("Books Resource description"),
						Post: &Method{
							Parameters: []Parameter{
								{
									Name:        "limit",
									In:          "query",
									Type:        "integer",
									Required:    true,
									Description: ptr.To("limit parameter"),
									Minimum:     ptr.To(float32(1)),
									Maximum:     ptr.To(float32(2)),
								},
								{
									Name:        "query",
									In:          "query",
									Type:        "string",
									Required:    true,
									Description: ptr.To("query parameter"),
									MinLength:   ptr.To(int64(1)),
									MaxLength:   ptr.To(int64(2)),
								},
							},
							RequestBody: orderedmap.New[string, Property](orderedmap.WithInitialData[string, Property](
								orderedmap.Pair[string, Property]{
									Key: "json",
									Value: Property{
										Name:        "Book Body",
										Type:        "object",
										Required:    true,
										Description: ptr.To("Json body desciption"),
										Properties: []Property{
											{
												Name:      "name",
												Type:      "string",
												Required:  true,
												MinLength: ptr.To(int64(1)),
												MaxLength: ptr.To(int64(200)),
											},
											{
												Name:     "tags",
												Type:     "array",
												Required: true,
												Items: &Items{
													Type: "object",
													Properties: []Property{
														{
															Name: "id",
															Type: "string",
														},
													},
												},
											},
										},
										XML: &XML{
											Attribute: ptr.To(true),
											Wrapped:   ptr.To(true),
											Namespace: ptr.To("akamai.com/schema"),
											Name:      ptr.To("BookRoot"),
											Prefix:    ptr.To("akam"),
										},
									},
								},
							)),
							Responses: &Responses{
								Headers: []Parameter{
									{
										Name:     "Set-Cookie",
										Type:     "string",
										Required: true,
									},
								},
								Contents: []ResponseContent{
									{
										StatusCodes: []int64{20},
										JSON: &Property{
											Name:     "application/json",
											Type:     "array",
											Required: false,
											Items: &Items{
												Type: "string",
											},
										},
									},
								},
							},
							Constraints: &MethodConstraints{
								EnforceOn: &MethodEnforceOn{
									UndefinedParameters: &UndefinedParameters{
										RequestBody: true,
									},
								},
							},
						},
					},
				})),
		},
	},
}

var error400BadRequest = Error{
	Type:     "https://problems.luna.akamaiapis.net/api-definitions/error-types/BAD-REQUEST",
	Title:    "Bad Request",
	Instance: "22f95e07-00f8-4dfc-a903-ba9b1e2999bf",
	Detail:   "Bad Request",
	Status:   http.StatusBadRequest,
	Severity: ptr.To("ERROR"),
}

var error404NotFound = Error{
	Type:     "https://problems.luna.akamaiapis.net/api-definitions/error-types/NOT-FOUND",
	Title:    "Not Found",
	Instance: "22f95e07-00f8-4dfc-a903-ba9b1e2999bf",
	Detail:   "Invalid version provided.",
	Status:   http.StatusNotFound,
	Severity: ptr.To("ERROR"),
}

var error409Conflict = Error{
	Type:     "/api-definitions/error-types/CONCURRENT-MODIFICATION-ERROR",
	Title:    "Concurrent Modification Error",
	Instance: "TestInstance123",
	Detail:   "API Endpoint does not allow concurrent modification. Please get the latest API Definition and try again.",
	Status:   http.StatusConflict,
	Severity: ptr.To("ERROR"),
}

func mockAPIClient(t *testing.T, mockServer *httptest.Server) APIDefinitions {
	serverURL, err := url.Parse(mockServer.URL)
	require.NoError(t, err)
	certPool := x509.NewCertPool()
	certPool.AddCert(mockServer.Certificate())
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}
	s, err := session.New(session.WithClient(httpClient), session.WithSigner(&edgegrid.Config{Host: serverURL.Host}))
	assert.NoError(t, err)
	return Client(s)
}

func TestClient(t *testing.T) {
	sess, err := session.New()
	require.NoError(t, err)
	tests := map[string]struct {
		options  []Option
		expected *apidefinitions
	}{
		"no options provided, return default": {
			options: nil,
			expected: &apidefinitions{
				Session: sess,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess, test.options...)
			assert.Equal(t, res, test.expected)
		})
	}
}

func loadJson(path string) string {
	contents, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	buf := bytes.Buffer{}
	if err := json.Compact(&buf, contents); err != nil {
		panic(fmt.Sprintf("%s: %s", err, string(contents)))
	}
	return buf.String()
}
