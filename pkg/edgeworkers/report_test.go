package edgeworkers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetReport(t *testing.T) {
	tests := map[string]struct {
		params           GetReportRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetReportResponse
		withError        error
	}{
		"200 OK - get report": {
			params: GetReportRequest{
				ReportID:   2,
				Start:      "2021-12-04T00:00:00Z",
				End:        "2022-01-01T15:00:00Z",
				EdgeWorker: "37017",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/edgeworkers/v1/reports/2?edgeWorker=37017&end=2022-01-01T15%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
			responseBody: `{
  "reportId": 2,
  "name": "Initialization and execution times by EdgeWorker ID and event handler",
  "description": "This report lists execution and initialization times, grouped by EdgeWorker ID and event handler.",
  "start": "2020-08-17T00:00:00Z",
  "end": "2020-08-17T20:00:00Z",
  "data": [
    {
      "edgeWorkerId": 37017,
      "data": {
        "onClientRequest": [
          {
            "startDateTime": "2016-08-17T00:00:00Z",
            "edgeWorkerVersion": "0.8",
            "execDuration": {
              "avg": 654,
              "min": 554,
              "max": 754
            },
            "invocations": 7
          },
          {
            "startDateTime": "2016-08-17T00:05:00Z",
            "edgeWorkerVersion": "0.9",
            "execDuration": {
              "avg": 453,
              "min": 454,
              "max": 459
            },
            "invocations": 5
          }
        ],
        "onOriginRequest": [
          {
            "startDateTime": "2016-08-17T04:00:00Z",
            "edgeWorkerVersion": "1.0",
            "execDuration": {
              "avg": 234,
              "min": 134,
              "max": 334
            },
            "invocations": 6
          }
        ],
        "onOriginResponse": [
          {
            "startDateTime": "2016-08-17T04:00:00Z",
            "edgeWorkerVersion": "1.0",
            "execDuration": {
              "avg": 324,
              "min": 224,
              "max": 424
            },
            "invocations": 2
          }
        ],
        "onClientResponse": [
          {
            "startDateTime": "2016-08-17T04:00:10Z",
            "edgeWorkerVersion": "1.0",
            "execDuration": {
              "avg": 275,
              "min": 175,
              "max": 375
            },
            "invocations": 6
          }
        ],
        "responseProvider": [
          {
            "startDateTime": "2016-08-17T04:00:00Z",
            "edgeWorkerVersion": "1.0",
            "execDuration": {
              "avg": 324,
              "min": 224,
              "max": 424
            },
            "invocations": 2
          }
        ],
        "init": [
          {
            "startDateTime": "2016-08-17T00:00:00Z",
            "edgeWorkerVersion": "0.8",
            "initDuration": {
              "avg": 654,
              "min": 554,
              "max": 754
            },
            "invocations": 7
          }
        ]
      }
    }
  ]
}`,
			expectedResponse: &GetReportResponse{
				ReportID:    2,
				Name:        "Initialization and execution times by EdgeWorker ID and event handler",
				Description: "This report lists execution and initialization times, grouped by EdgeWorker ID and event handler.",
				Start:       "2020-08-17T00:00:00Z",
				End:         "2020-08-17T20:00:00Z",
				Data: []ReportData{
					{
						EdgeWorkerID: 37017,
						Data: Data{
							OnClientRequest: []OnRequestAndResponse{
								{
									StartDateTime:     "2016-08-17T00:00:00Z",
									EdgeWorkerVersion: "0.8",
									ExecDuration: Duration{
										Avg: 654,
										Min: 554,
										Max: 754,
									},
									Invocations: 7,
								},
								{
									StartDateTime:     "2016-08-17T00:05:00Z",
									EdgeWorkerVersion: "0.9",
									ExecDuration: Duration{
										Avg: 453,
										Min: 454,
										Max: 459,
									},
									Invocations: 5,
								},
							},
							OnOriginRequest: []OnRequestAndResponse{
								{
									StartDateTime:     "2016-08-17T04:00:00Z",
									EdgeWorkerVersion: "1.0",
									ExecDuration: Duration{
										Avg: 234,
										Min: 134,
										Max: 334,
									},
									Invocations: 6,
								},
							},
							OnOriginResponse: []OnRequestAndResponse{
								{
									StartDateTime:     "2016-08-17T04:00:00Z",
									EdgeWorkerVersion: "1.0",
									ExecDuration: Duration{
										Avg: 324,
										Min: 224,
										Max: 424,
									},
									Invocations: 2,
								},
							},
							OnClientResponse: []OnRequestAndResponse{
								{
									StartDateTime:     "2016-08-17T04:00:10Z",
									EdgeWorkerVersion: "1.0",
									ExecDuration: Duration{
										Avg: 275,
										Min: 175,
										Max: 375,
									},
									Invocations: 6,
								},
							},
							ResponseProvider: []OnRequestAndResponse{
								{
									StartDateTime:     "2016-08-17T04:00:00Z",
									EdgeWorkerVersion: "1.0",
									ExecDuration: Duration{
										Avg: 324,
										Min: 224,
										Max: 424,
									},
									Invocations: 2,
								},
							},
							Init: []InitObject{
								{
									StartDateTime:     "2016-08-17T00:00:00Z",
									EdgeWorkerVersion: "0.8",
									InitDuration: Duration{
										Avg: 654,
										Min: 554,
										Max: 754,
									},
									Invocations: 7,
								},
							},
						},
					},
				},
			},
		},
		"200 OK - get report without data": {
			params: GetReportRequest{
				ReportID:     4,
				Start:        "2021-12-04T00:00:00Z",
				End:          "2022-01-01T00:00:00Z",
				EdgeWorker:   "37017",
				Status:       tools.StringPtr(StatusSuccess),
				EventHandler: tools.StringPtr(EventHandlerOnClientRequest),
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/edgeworkers/v1/reports/4?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&eventHandler=onClientRequest&start=2021-12-04T00%3A00%3A00Z&status=success",
			responseBody: `{
    "reportId": 4,
    "name": "Memory usage by EdgeWorker ID and event handler",
    "description": "This report lists memory usage, grouped by EdgeWorker ID and event handler.",
    "start": "2021-12-04T00:00:00Z",
    "end": "2022-01-01T00:00:00Z",
    "data": []
}`,
			expectedResponse: &GetReportResponse{
				ReportID:    4,
				Name:        "Memory usage by EdgeWorker ID and event handler",
				Description: "This report lists memory usage, grouped by EdgeWorker ID and event handler.",
				Start:       "2021-12-04T00:00:00Z",
				End:         "2022-01-01T00:00:00Z",
				Data:        []ReportData{},
			},
		},
		"missing mandatory params": {
			params:    GetReportRequest{},
			withError: ErrStructValidation,
		},
		"invalid status": {
			params: GetReportRequest{
				ReportID:   4,
				Start:      "2021-12-04T00:00:00Z",
				End:        "2022-01-01T00:00:00Z",
				EdgeWorker: "37017",
				Status:     tools.StringPtr("not_valid"),
			},
			withError: ErrStructValidation,
		},
		"empty status": {
			params: GetReportRequest{
				ReportID:   4,
				Start:      "2021-12-04T00:00:00Z",
				End:        "2022-01-01T00:00:00Z",
				EdgeWorker: "37017",
				Status:     tools.StringPtr(""),
			},
			expectedPath: "/edgeworkers/v1/reports/4?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
			withError:    ErrStructValidation,
		},
		"invalid event handler": {
			params: GetReportRequest{
				ReportID:     4,
				Start:        "2021-12-04T00:00:00Z",
				End:          "2022-01-01T00:00:00Z",
				EdgeWorker:   "37017",
				EventHandler: tools.StringPtr("not_valid"),
			},
			withError: ErrStructValidation,
		},
		"empty event handler": {
			params: GetReportRequest{
				ReportID:     4,
				Start:        "2021-12-04T00:00:00Z",
				End:          "2022-01-01T00:00:00Z",
				EdgeWorker:   "37017",
				EventHandler: tools.StringPtr(""),
			},
			expectedPath: "/edgeworkers/v1/reports/4?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
			withError:    ErrStructValidation,
		},
		"500 internal server error - get group which does not exist": {
			params: GetReportRequest{
				ReportID:   4,
				Start:      "2021-12-04T00:00:00Z",
				End:        "2022-01-01T00:00:00Z",
				EdgeWorker: "37017",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
    "title": "Server Error",
    "status": 500,
    "instance": "host_name/edgeworkers/v1/reports/4",
    "method": "GET",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2022-01-01T00:00:00Z"
}`,
			expectedPath: "/edgeworkers/v1/reports/4?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/reports/4",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2022-01-01T00:00:00Z",
			},
		},
		"403 Forbidden - incorrect credentials": {
			params: GetReportRequest{
				ReportID:   4,
				Start:      "2021-12-04T00:00:00Z",
				End:        "2022-01-01T00:00:00Z",
				EdgeWorker: "37017",
			},
			responseStatus: http.StatusForbidden,
			responseBody: `
{
    "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
    "title": "Forbidden",
    "status": 403,
    "detail": "The client does not have the grant needed for the request",
    "instance": "host_name/edgeworkers/v1/reports/4?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
    "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
    "method": "GET",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2022-01-10T10:56:32Z"
}`,
			expectedPath: "/edgeworkers/v1/reports/4?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/reports/4?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2022-01-10T10:56:32Z",
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
			result, err := client.GetReport(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestListReports(t *testing.T) {
	tests := map[string]struct {
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ListReportsResponse
		withError        error
	}{
		"200 OK - list EdgeWorker reports": {
			responseStatus: http.StatusOK,
			responseBody: `
{
    "reports": [
        {
            "reportId": 1,
            "name": "Overall summary",
            "description": "This report contains an overview of other reports."
        },
        {
            "reportId": 2,
            "name": "Initialization and execution times by EdgeWorker ID and event handler",
            "description": "This report lists execution and initialization times, grouped by EdgeWorker ID and event handler."
        },
        {
            "reportId": 3,
            "name": "Execution statuses by EdgeWorker ID and event handler",
            "description": "This report lists execution statuses, grouped by EdgeWorker ID and event handler."
        },
        {
            "reportId": 4,
            "name": "Memory usage by EdgeWorker ID and event handler",
            "description": "This report lists memory usage, grouped by EdgeWorker ID and event handler."
        }
    ]
}`,
			expectedPath: "/edgeworkers/v1/reports",
			expectedResponse: &ListReportsResponse{[]ReportResponse{
				{
					ReportID:    1,
					Name:        "Overall summary",
					Description: "This report contains an overview of other reports.",
				},
				{
					ReportID:    2,
					Name:        "Initialization and execution times by EdgeWorker ID and event handler",
					Description: "This report lists execution and initialization times, grouped by EdgeWorker ID and event handler.",
				},
				{
					ReportID:    3,
					Name:        "Execution statuses by EdgeWorker ID and event handler",
					Description: "This report lists execution statuses, grouped by EdgeWorker ID and event handler.",
				},
				{
					ReportID:    4,
					Name:        "Memory usage by EdgeWorker ID and event handler",
					Description: "This report lists memory usage, grouped by EdgeWorker ID and event handler.",
				},
			}},
		},
		"403 Forbidden - incorrect credentials": {
			responseStatus: http.StatusForbidden,
			responseBody: `
{
   "type": "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
   "title": "Forbidden",
   "status": 403,
   "detail": "The client does not have the grant needed for the request",
   "instance": "host_name/edgeworkers/v1/reports",
   "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
   "method": "GET",
   "serverIp": "104.81.220.111",
   "clientIp": "89.64.55.111",
   "requestId": "a73affa111",
   "requestTime": "2022-01-10T12:31:29Z"
}`,
			expectedPath: "/edgeworkers/v1/reports",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/reports",
				AuthzRealm:  "scuomder224df6ct.dkekfr3qqg4dghpj",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2022-01-10T12:31:29Z",
			},
		},
		"500 internal server error": {
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
   "type": "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
   "title": "Server Error",
   "status": 500,
   "instance": "host_name/edgeworkers/v1/reports",
   "method": "GET",
   "serverIp": "104.81.220.111",
   "clientIp": "89.64.55.111",
   "requestId": "a73affa111",
   "requestTime": "2022-01-10T12:31:29Z"
}`,
			expectedPath: "/edgeworkers/v1/reports",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/reports",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2022-01-10T12:31:29Z",
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
			result, err := client.ListReports(context.Background())
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
