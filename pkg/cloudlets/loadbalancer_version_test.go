package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/tools"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestCreateLoadBalancerVersion(t *testing.T) {
	tests := map[string]struct {
		params           CreateLoadBalancerVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *LoadBalancerVersion
		withError        func(*testing.T, error)
	}{
		"201 Created": {
			params: CreateLoadBalancerVersionRequest{
				LoadBalancerVersion: LoadBalancerVersion{
					BalancingType: BalancingTypeWeighted,
					DataCenters: []DataCenter{
						{
							CloudService: false,
							Hostname:     "clorigin.example.com",
							LivenessHosts: []string{
								"clorigin3.www.example.com",
							},
							Latitude:  tools.Float64Ptr(102.78108),
							Longitude: tools.Float64Ptr(-116.07064),
							Continent: "NA",
							Country:   "US",
							OriginID:  "clorigin3",
							Percent:   tools.Float64Ptr(100.0),
						},
					},
					Deleted:     false,
					Description: "Test load balancing configuration.",
					Immutable:   false,
					LivenessSettings: &LivenessSettings{
						HostHeader: "clorigin3.www.example.com",
						AdditionalHeaders: map[string]string{
							"Host": "example.com",
							"test": "test",
						},
						Interval:         25,
						Path:             "/status",
						Port:             443,
						Protocol:         "HTTPS",
						Status3xxFailure: false,
						Status4xxFailure: true,
						Status5xxFailure: false,
						Timeout:          30,
					},
					OriginID: "clorigin3",
					Version:  4,
				},
				OriginID: "clorigin3",
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "balancingType": "WEIGHTED",
    "createdBy": "jjones",
    "createdDate": "2015-10-08T11:42:18.690Z",
    "dataCenters": [
        {
            "cloudService": false,
			"hostname": "clorigin.example.com",
            "livenessHosts": [
                "clorigin3.www.example.com"
            ],
            "latitude": 102.78108,
            "longitude": -116.07064,
            "continent": "NA",
            "country": "US",
            "originId": "clorigin3",
            "percent": 100.0
        }
    ],
    "deleted": false,
    "description": "Test load balancing configuration.",
    "immutable": false,
    "lastModifiedBy": "ejnovak",
    "lastModifiedDate": "2016-05-02T00:40:02.237Z",
    "livenessSettings": {
        "hostHeader": "clorigin3.www.example.com",
		"additionalHeaders": {
			"Host": "example.com",
			"test": "test"
		},
        "interval": 25,
        "path": "/status",
        "port": 443,
        "protocol": "HTTPS",
        "status3xxFailure": false,
        "status4xxFailure": true,
        "status5xxFailure": false,
        "timeout": 30
    },
    "originId": "clorigin3",
    "version": 4
}`,
			expectedPath: "/cloudlets/api/v2/origins/clorigin3/versions",
			expectedResponse: &LoadBalancerVersion{
				BalancingType: BalancingTypeWeighted,
				CreatedBy:     "jjones",
				CreatedDate:   "2015-10-08T11:42:18.690Z",
				DataCenters: []DataCenter{
					{
						CloudService: false,
						Hostname:     "clorigin.example.com",
						LivenessHosts: []string{
							"clorigin3.www.example.com",
						},
						Latitude:  tools.Float64Ptr(102.78108),
						Longitude: tools.Float64Ptr(-116.07064),
						Continent: "NA",
						Country:   "US",
						OriginID:  "clorigin3",
						Percent:   tools.Float64Ptr(100.0),
					},
				},
				Deleted:          false,
				Description:      "Test load balancing configuration.",
				Immutable:        false,
				LastModifiedBy:   "ejnovak",
				LastModifiedDate: "2016-05-02T00:40:02.237Z",
				LivenessSettings: &LivenessSettings{
					HostHeader: "clorigin3.www.example.com",
					AdditionalHeaders: map[string]string{
						"Host": "example.com",
						"test": "test",
					},
					Interval:         25,
					Path:             "/status",
					Port:             443,
					Protocol:         "HTTPS",
					Status3xxFailure: false,
					Status4xxFailure: true,
					Status5xxFailure: false,
					Timeout:          30,
				},
				OriginID: "clorigin3",
				Version:  4,
			},
		},
		"201 OK empty": {
			params: CreateLoadBalancerVersionRequest{
				LoadBalancerVersion: LoadBalancerVersion{},
				OriginID:            "clorigin3",
			},
			responseStatus: http.StatusCreated,
			responseBody: `
{
    "createdBy": "jjones",
    "createdDate": "2021-09-02T12:51:19.798Z",
    "deleted": false,
    "immutable": false,
    "lastModifiedBy": "jjones",
    "lastModifiedDate": "2021-09-02T12:51:19.798Z",
    "originId": "clorigin3",
    "version": 1,
    "warnings": [
        {
            "detail": "The Load Balancer type is not set",
            "title": "Validation Warning",
            "type": "/cloudlets/error-types/validation-warning",
            "jsonPointer": "/"
        },
        {
            "detail": "The Load Balancer must have at least one data center defined",
            "title": "Validation Warning",
            "type": "/cloudlets/error-types/validation-warning",
            "jsonPointer": "/"
        }
    ]
}`,
			expectedPath: "/cloudlets/api/v2/origins/clorigin3/versions",
			expectedResponse: &LoadBalancerVersion{
				CreatedBy:        "jjones",
				CreatedDate:      "2021-09-02T12:51:19.798Z",
				Deleted:          false,
				Immutable:        false,
				LastModifiedBy:   "jjones",
				LastModifiedDate: "2021-09-02T12:51:19.798Z",
				OriginID:         "clorigin3",
				Version:          1,
				Warnings: []Warning{
					{
						Detail:      "The Load Balancer type is not set",
						Title:       "Validation Warning",
						Type:        "/cloudlets/error-types/validation-warning",
						JSONPointer: "/",
					},
					{
						Detail:      "The Load Balancer must have at least one data center defined",
						Title:       "Validation Warning",
						Type:        "/cloudlets/error-types/validation-warning",
						JSONPointer: "/",
					},
				},
			},
		},
		"500 internal server error": {
			params: CreateLoadBalancerVersionRequest{
				OriginID:            "23",
				LoadBalancerVersion: LoadBalancerVersion{},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
	"type": "internal_error",
   "title": "Internal Server Error",
   "detail": "Error making request",
   "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/origins/23/versions",
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
				assert.Equal(t, http.MethodPost, r.Method)

				// Check if received body is valid json
				b, err := ioutil.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.True(t, json.Valid(b))

				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.CreateLoadBalancerVersion(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestDataCenterValidate(t *testing.T) {
	tests := map[string]struct {
		DataCenter
		withError error
	}{
		"valid data center": {
			DataCenter: DataCenter{
				CloudService: false,
				Hostname:     "clorigin.example.com",
				LivenessHosts: []string{
					"clorigin3.www.example.com",
				},
				Latitude:  tools.Float64Ptr(102.78108),
				Longitude: tools.Float64Ptr(-116.07064),
				Continent: "NA",
				Country:   "US",
				OriginID:  "clorigin3",
				Percent:   tools.Float64Ptr(100.0),
			},
		},
		"valid data center, minimal set of params": {
			DataCenter: DataCenter{
				Latitude:  tools.Float64Ptr(102.78108),
				Longitude: tools.Float64Ptr(-116.07064),
				Continent: "NA",
				Country:   "US",
				OriginID:  "clorigin3",
				Percent:   tools.Float64Ptr(100.0),
			},
		},
		"longitude, latitude and percent can be 0": {
			DataCenter: DataCenter{
				CloudService: false,
				Hostname:     "clorigin.example.com",
				LivenessHosts: []string{
					"clorigin3.www.example.com",
				},
				Latitude:  tools.Float64Ptr(0),
				Longitude: tools.Float64Ptr(0),
				Continent: "NA",
				Country:   "US",
				OriginID:  "clorigin3",
				Percent:   tools.Float64Ptr(0),
			},
		},
		"missing all required parameters error": {
			DataCenter: DataCenter{},
			withError: validation.Errors{
				"Continent": validation.ErrRequired,
				"Country":   validation.ErrRequired,
				"Latitude":  validation.ErrNotNilRequired,
				"Longitude": validation.ErrNotNilRequired,
				"OriginID":  validation.ErrRequired,
				"Percent":   validation.ErrNotNilRequired,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.DataCenter.Validate()
			if test.withError != nil {
				require.Error(t, err)
				assert.Equal(t, test.withError.Error(), err.Error(), "want: %s; got %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestLivenessSettingsValidate(t *testing.T) {
	tests := map[string]struct {
		LivenessSettings *LivenessSettings
		withError        error
	}{
		"path is required error": {
			LivenessSettings: &LivenessSettings{
				Port:     80,
				Protocol: "HTTP",
			},
			withError: validation.Errors{
				"Path": validation.ErrRequired,
			},
		},
		"RequestString is required error": {
			LivenessSettings: &LivenessSettings{
				Port:           1234,
				Protocol:       "TCP",
				ResponseString: "response",
			},
			withError: validation.Errors{
				"RequestString": validation.ErrRequired,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.LivenessSettings.Validate()
			if test.withError != nil && err != nil {
				assert.True(t, err.Error() == test.withError.Error(), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetLoadBalancerVersion(t *testing.T) {
	tests := map[string]struct {
		params           GetLoadBalancerVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *LoadBalancerVersion
		withError        error
	}{
		"200 OK": {
			params: GetLoadBalancerVersionRequest{
				OriginID: "clorigin3",
				Version:  4,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "balancingType": "WEIGHTED",
    "createdBy": "jjones",
    "createdDate": "2015-10-08T11:42:18.690Z",
    "dataCenters": [
        {
            "cloudService": false,
            "livenessHosts": [
                "clorigin3.www.example.com"
            ],
            "latitude": 102.78108,
            "longitude": -116.07064,
            "continent": "NA",
            "country": "US",
            "originId": "clorigin3",
            "percent": 100.0
        }
    ],
    "deleted": false,
    "description": "Test load balancing configuration.",
    "immutable": false,
    "lastModifiedBy": "ejnovak",
    "lastModifiedDate": "2016-05-02T00:40:02.237Z",
    "livenessSettings": {
        "hostHeader": "clorigin3.www.example.com",
        "interval": 25,
        "path": "/status",
        "port": 443,
        "protocol": "HTTPS",
        "status3xxFailure": false,
        "status4xxFailure": true,
        "status5xxFailure": false,
        "timeout": 30
    },
    "originId": "clorigin3",
    "version": 4
}`,
			expectedPath: "/cloudlets/api/v2/origins/clorigin3/versions/4",
			expectedResponse: &LoadBalancerVersion{
				BalancingType: BalancingTypeWeighted,
				CreatedBy:     "jjones",
				CreatedDate:   "2015-10-08T11:42:18.690Z",
				DataCenters: []DataCenter{
					{
						CloudService: false,
						LivenessHosts: []string{
							"clorigin3.www.example.com",
						},
						Latitude:  tools.Float64Ptr(102.78108),
						Longitude: tools.Float64Ptr(-116.07064),
						Continent: "NA",
						Country:   "US",
						OriginID:  "clorigin3",
						Percent:   tools.Float64Ptr(100.0),
					},
				},
				Deleted:          false,
				Description:      "Test load balancing configuration.",
				Immutable:        false,
				LastModifiedBy:   "ejnovak",
				LastModifiedDate: "2016-05-02T00:40:02.237Z",
				LivenessSettings: &LivenessSettings{
					HostHeader:       "clorigin3.www.example.com",
					Interval:         25,
					Path:             "/status",
					Port:             443,
					Protocol:         "HTTPS",
					Status3xxFailure: false,
					Status4xxFailure: true,
					Status5xxFailure: false,
					Timeout:          30,
				},
				OriginID: "clorigin3",
				Version:  4,
			},
		},
		"500 internal server error": {
			params: GetLoadBalancerVersionRequest{
				OriginID: "oid1",
				Version:  3,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
 "type": "internal_error",
 "title": "Internal Server Error",
 "detail": "Error creating enrollment",
 "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/origins/oid1/versions/3",
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
				assert.Equal(t, http.MethodGet, r.Method)

				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetLoadBalancerVersion(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestUpdateLoadBalancerVersion(t *testing.T) {
	tests := map[string]struct {
		params           UpdateLoadBalancerVersionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *LoadBalancerVersion
		withError        error
	}{
		"200 OK": {
			params: UpdateLoadBalancerVersionRequest{
				OriginID: "clorigin3",
				Version:  4,
				LoadBalancerVersion: LoadBalancerVersion{
					BalancingType: BalancingTypeWeighted,
					DataCenters: []DataCenter{
						{
							CloudService: false,
							Hostname:     "clorigin.example.com",
							LivenessHosts: []string{
								"clorigin3.www.example.com",
							},
							Latitude:  tools.Float64Ptr(102.78108),
							Longitude: tools.Float64Ptr(-116.07064),
							Continent: "NA",
							Country:   "US",
							OriginID:  "clorigin3",
							Percent:   tools.Float64Ptr(100.0),
						},
					},
					Deleted:     false,
					Description: "Test load balancing configuration.",
					Immutable:   false,
					LivenessSettings: &LivenessSettings{
						HostHeader:       "clorigin3.www.example.com",
						Interval:         25,
						Path:             "/status",
						Port:             443,
						Protocol:         "HTTPS",
						Status3xxFailure: false,
						Status4xxFailure: true,
						Status5xxFailure: false,
						Timeout:          30,
					},
					OriginID: "clorigin3",
					Version:  4,
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "balancingType": "WEIGHTED",
    "createdBy": "jjones",
    "createdDate": "2015-10-08T11:42:18.690Z",
    "dataCenters": [
        {
            "cloudService": false,
			"hostname": "clorigin.example.com",
            "livenessHosts": [
                "clorigin3.www.example.com"
            ],
            "latitude": 102.78108,
            "longitude": -116.07064,
            "continent": "NA",
            "country": "US",
            "originId": "clorigin3",
            "percent": 100.0
        }
    ],
    "deleted": false,
    "description": "Test load balancing configuration.",
    "immutable": false,
    "lastModifiedBy": "ejnovak",
    "lastModifiedDate": "2016-05-02T00:40:02.237Z",
    "livenessSettings": {
        "hostHeader": "clorigin3.www.example.com",
        "interval": 25,
        "path": "/status",
        "port": 443,
        "protocol": "HTTPS",
        "status3xxFailure": false,
        "status4xxFailure": true,
        "status5xxFailure": false,
        "timeout": 30
    },
    "originId": "clorigin3",
    "version": 4
}`,
			expectedPath: "/cloudlets/api/v2/origins/clorigin3/versions/4",
			expectedResponse: &LoadBalancerVersion{
				BalancingType: BalancingTypeWeighted,
				CreatedBy:     "jjones",
				CreatedDate:   "2015-10-08T11:42:18.690Z",
				DataCenters: []DataCenter{
					{
						CloudService: false,
						Hostname:     "clorigin.example.com",
						LivenessHosts: []string{
							"clorigin3.www.example.com",
						},
						Latitude:  tools.Float64Ptr(102.78108),
						Longitude: tools.Float64Ptr(-116.07064),
						Continent: "NA",
						Country:   "US",
						OriginID:  "clorigin3",
						Percent:   tools.Float64Ptr(100.0),
					},
				},
				Deleted:          false,
				Description:      "Test load balancing configuration.",
				Immutable:        false,
				LastModifiedBy:   "ejnovak",
				LastModifiedDate: "2016-05-02T00:40:02.237Z",
				LivenessSettings: &LivenessSettings{
					HostHeader:       "clorigin3.www.example.com",
					Interval:         25,
					Path:             "/status",
					Port:             443,
					Protocol:         "HTTPS",
					Status3xxFailure: false,
					Status4xxFailure: true,
					Status5xxFailure: false,
					Timeout:          30,
				},
				OriginID: "clorigin3",
				Version:  4,
			},
		},
		"500 internal server error": {
			params: UpdateLoadBalancerVersionRequest{
				OriginID: "clorigin3",
				Version:  4,
				LoadBalancerVersion: LoadBalancerVersion{
					BalancingType: BalancingTypeWeighted,
					DataCenters: []DataCenter{
						{
							CloudService: false,
							LivenessHosts: []string{
								"clorigin3.www.example.com",
							},
							Latitude:  tools.Float64Ptr(102.78108),
							Longitude: tools.Float64Ptr(-116.07064),
							Continent: "NA",
							Country:   "US",
							OriginID:  "clorigin3",
							Percent:   tools.Float64Ptr(100.0),
						},
					},
					Description: "Test load balancing configuration.",
					LivenessSettings: &LivenessSettings{
						HostHeader: "clorigin3.www.example.com",
						Path:       "/status",
						Port:       443,
						Protocol:   "HTTPS",
					},
					OriginID: "clorigin3",
					Version:  4,
				},
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
 "type": "internal_error",
 "title": "Internal Server Error",
 "detail": "Error creating enrollment",
 "status": 500
}`,
			expectedPath: "/cloudlets/api/v2/origins/clorigin3/versions/4",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating enrollment",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"validation error": {
			params: UpdateLoadBalancerVersionRequest{
				OriginID: "oid1",
				Version:  -1,
			},
			withError: ErrStructValidation,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodPut, r.Method)

				// Check if received body is valid json
				b, err := ioutil.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.True(t, json.Valid(b))

				w.WriteHeader(test.responseStatus)
				_, err = w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateLoadBalancerVersion(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListLoadBalancerVersions(t *testing.T) {
	tests := map[string]struct {
		listLoadBalancerVersionsRequest ListLoadBalancerVersionsRequest
		responseStatus                  int
		responseBody                    string
		expectedPath                    string
		expectedResponse                []LoadBalancerVersion
		withError                       func(*testing.T, error)
	}{
		"200 OK": {
			listLoadBalancerVersionsRequest: ListLoadBalancerVersionsRequest{
				OriginID: "clorigin3",
			},
			responseStatus: http.StatusOK,
			responseBody: `[
				{
					"createdBy": "jjones",
					"createdDate": "2015-10-08T11:42:18.690Z",
					"dataCenters": [
						{
							"cloudService": false,
							"livenessHosts": [
								"clorigin3.www.example.com"
							],
							"latitude": 102.78108,
							"longitude": -116.07064,
							"continent": "NA",
							"country": "US",
							"originId": "clorigin3",
							"percent": 100.0
						}
					],
					"deleted": false,
					"description": "Test load balancing configuration.",
					"immutable": false,
					"lastModifiedBy": "jsmith",
					"lastModifiedDate": "2016-05-02T00:40:02.237Z",
					"livenessSettings": {
						"hostHeader": "clorigin3.www.example.com",
						"path": "/status",
						"port": 443,
						"protocol": "HTTPS"
					},
					"originId": "clorigin3",
					"version": 2
				},
				{
					"createdBy": "jjones",
					"createdDate": "2015-10-08T11:42:18.690Z",
					"dataCenters": [
						{
							"cloudService": false,
							"livenessHosts": [
								"clorigin3.www.example.com"
							],
							"latitude": 102.78108,
							"longitude": -116.07064,
							"continent": "NA",
							"country": "US",
							"originId": "clorigin3",
							"percent": 100.0
						}
					],
					"deleted": false,
					"description": "Test load balancing configuration.",
					"immutable": false,
					"lastModifiedBy": "jsmith",
					"lastModifiedDate": "2016-05-02T00:40:02.237Z",
					"livenessSettings": {
						"hostHeader": "clorigin3.www.example.com",
						"path": "/status",
						"port": 443,
						"protocol": "HTTPS"
					},
					"originId": "clorigin3",
					"version": 1
				}
			]`,
			expectedPath: "/cloudlets/api/v2/origins/clorigin3/versions?includeModel=true",
			expectedResponse: []LoadBalancerVersion{
				{
					CreatedBy:   "jjones",
					CreatedDate: "2015-10-08T11:42:18.690Z",
					DataCenters: []DataCenter{
						{
							CloudService: false,
							LivenessHosts: []string{
								"clorigin3.www.example.com",
							},
							Latitude:  tools.Float64Ptr(102.78108),
							Longitude: tools.Float64Ptr(-116.07064),
							Continent: "NA",
							Country:   "US",
							OriginID:  "clorigin3",
							Percent:   tools.Float64Ptr(100.0),
						},
					},
					Deleted:          false,
					Description:      "Test load balancing configuration.",
					Immutable:        false,
					LastModifiedBy:   "jsmith",
					LastModifiedDate: "2016-05-02T00:40:02.237Z",
					OriginID:         "clorigin3",
					LivenessSettings: &LivenessSettings{
						HostHeader: "clorigin3.www.example.com",
						Path:       "/status",
						Port:       443,
						Protocol:   "HTTPS",
					},
					Version: 2,
				},
				{
					CreatedBy:   "jjones",
					CreatedDate: "2015-10-08T11:42:18.690Z",
					DataCenters: []DataCenter{
						{
							CloudService: false,
							LivenessHosts: []string{
								"clorigin3.www.example.com",
							},
							Latitude:  tools.Float64Ptr(102.78108),
							Longitude: tools.Float64Ptr(-116.07064),
							Continent: "NA",
							Country:   "US",
							OriginID:  "clorigin3",
							Percent:   tools.Float64Ptr(100.0),
						},
					},
					Deleted:          false,
					Description:      "Test load balancing configuration.",
					Immutable:        false,
					LastModifiedBy:   "jsmith",
					LastModifiedDate: "2016-05-02T00:40:02.237Z",
					OriginID:         "clorigin3",
					LivenessSettings: &LivenessSettings{
						HostHeader: "clorigin3.www.example.com",
						Path:       "/status",
						Port:       443,
						Protocol:   "HTTPS",
					},
					Version: 1,
				},
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
			result, err := client.ListLoadBalancerVersions(context.Background(), test.listLoadBalancerVersionsRequest)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
