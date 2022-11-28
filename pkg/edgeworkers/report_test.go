package edgeworkers

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

func TestGetSummaryReport(t *testing.T) {
	tests := map[string]struct {
		params           GetSummaryReportRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetSummaryReportResponse
		withError        error
	}{
		"200 OK - get summary overview report": {
			params: GetSummaryReportRequest{
				Start:      "2022-01-10T03:00:00Z",
				End:        "2022-01-14T13:22:28Z",
				EdgeWorker: "1234",
				Status:     tools.StringPtr(StatusSuccess),
			},
			expectedPath: "/edgeworkers/v1/reports/1?edgeWorker=1234&end=2022-01-14T13%3A22%3A28Z&start=2022-01-10T03%3A00%3A00Z&status=success",
			responseBody: `
{
    "reportId": 1,
    "name": "Overall summary",
    "description": "This report contains an overview of other reports.",
    "start": "2022-01-10T03:00:00Z",
    "end": "2022-01-14T13:22:28Z",
    "data": {
        "memory": {
            "avg": 3607.069168,
            "min": 0.0,
            "max": 458044.0
        },
        "successes": {
            "total": 88119
        },
        "initDuration": {
            "avg": 1.2559162420382166,
            "min": 0.33,
            "max": 28.975
        },
        "execDuration": {
            "avg": 0.11508884576538543,
            "min": 0.005,
            "max": 9.415
        },
        "errors": {
            "total": 0
        },
        "invocations": {
            "total": 88119
        }
    }
}`,
			responseStatus: http.StatusOK,
			expectedResponse: &GetSummaryReportResponse{
				ReportID:    1,
				Name:        "Overall summary",
				Description: "This report contains an overview of other reports.",
				Start:       "2022-01-10T03:00:00Z",
				End:         "2022-01-14T13:22:28Z",
				Data: DataSummary{
					Memory: &Summary{
						Avg: 3607.069168,
						Min: 0.0,
						Max: 458044.0,
					},
					Successes: &Total{Total: 88119},
					InitDuration: &Summary{
						Avg: 1.2559162420382166,
						Min: 0.33,
						Max: 28.975,
					},
					ExecDuration: &Summary{
						Avg: 0.11508884576538543,
						Min: 0.005,
						Max: 9.415,
					},
					Errors:      &Total{Total: 0},
					Invocations: &Total{Total: 88119},
				},
			},
		},
		"missing mandatory params": {
			params:    GetSummaryReportRequest{},
			withError: ErrStructValidation,
		},
		"invalid status": {
			params: GetSummaryReportRequest{
				Start:      "2021-12-04T00:00:00Z",
				End:        "2022-01-01T00:00:00Z",
				EdgeWorker: "37017",
				Status:     tools.StringPtr("not_valid"),
			},
			withError: ErrStructValidation,
		},
		"empty status": {
			params: GetSummaryReportRequest{
				Start:      "2021-12-04T00:00:00Z",
				End:        "2022-01-01T00:00:00Z",
				EdgeWorker: "37017",
				Status:     tools.StringPtr(""),
			},
			expectedPath: "/edgeworkers/v1/reports/1?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
			withError:    ErrStructValidation,
		},
		"invalid event handler": {
			params: GetSummaryReportRequest{
				Start:        "2021-12-04T00:00:00Z",
				End:          "2022-01-01T00:00:00Z",
				EdgeWorker:   "37017",
				EventHandler: tools.StringPtr("not_valid"),
			},
			withError: ErrStructValidation,
		},
		"empty event handler": {
			params: GetSummaryReportRequest{
				Start:        "2021-12-04T00:00:00Z",
				End:          "2022-01-01T00:00:00Z",
				EdgeWorker:   "37017",
				EventHandler: tools.StringPtr(""),
			},
			expectedPath: "/edgeworkers/v1/reports/1?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
			withError:    ErrStructValidation,
		},
		"500 internal server error - get group which does not exist": {
			params: GetSummaryReportRequest{
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
    "instance": "host_name/edgeworkers/v1/reports/1",
    "method": "GET",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2022-01-17T00:00:00Z"
}`,
			expectedPath: "/edgeworkers/v1/reports/1?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/resource-impl/forward-origin-error",
				Title:       "Server Error",
				Status:      500,
				Instance:    "host_name/edgeworkers/v1/reports/1",
				Method:      "GET",
				ServerIP:    "104.81.220.111",
				ClientIP:    "89.64.55.111",
				RequestID:   "a73affa111",
				RequestTime: "2022-01-17T00:00:00Z",
			},
		},
		"403 Forbidden - incorrect credentials": {
			params: GetSummaryReportRequest{
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
    "instance": "host_name/edgeworkers/v1/reports/1?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
    "authzRealm": "scuomder224df6ct.dkekfr3qqg4dghpj",
    "method": "GET",
    "serverIp": "104.81.220.111",
    "clientIp": "89.64.55.111",
    "requestId": "a73affa111",
    "requestTime": "2022-01-10T10:56:32Z"
}`,
			expectedPath: "/edgeworkers/v1/reports/1?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
			withError: &Error{
				Type:        "https://problems.luna-dev.akamaiapis.net/-/pep-authz/deny",
				Title:       "Forbidden",
				Status:      403,
				Detail:      "The client does not have the grant needed for the request",
				Instance:    "host_name/edgeworkers/v1/reports/1?edgeWorker=37017&end=2022-01-01T00%3A00%3A00Z&start=2021-12-04T00%3A00%3A00Z",
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
			result, err := client.GetSummaryReport(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

func TestGetReport(t *testing.T) {
	tests := map[string]struct {
		params           GetReportRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetReportResponse
		withError        error
	}{
		"200 OK - get report with id 2": {
			params: GetReportRequest{
				ReportID:   2,
				Start:      "2021-12-04T00:00:00Z",
				End:        "2022-01-17T11:28:11Z",
				EdgeWorker: "37017",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/edgeworkers/v1/reports/2?edgeWorker=37017&end=2022-01-17T11%3A28%3A11Z&start=2021-12-04T00%3A00%3A00Z",
			responseBody: `{
  "reportId": 2,
  "name": "Initialization and execution times by EdgeWorker ID and event handler",
  "description": "This report lists execution and initialization times, grouped by EdgeWorker ID and event handler.",
  "start": "2022-01-10T03:00:00Z",
  "end": "2022-01-17T10:43:15Z",
  "data": [
    {
      "edgeWorkerId": 37017,
      "data": {
        "onClientRequest": [
          {
			"startDateTime": "2022-01-10T03:00:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 8,
			"execDuration": {
				"avg": 0.113125,
				"min": 0.095,
				"max": 0.148
			}
		  },
		  {
			"startDateTime": "2022-01-10T03:05:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 15,
			"execDuration": {
				"avg": 0.14859999999999998,
				"min": 0.087,
				"max": 0.348
			}
		  }
        ],
        "onOriginRequest": [
          {
			"startDateTime": "2022-01-10T03:45:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"execDuration": {
				"avg": 0.0,
				"min": 0.0,
				"max": 0.0
			}
		  },
		  {
			"startDateTime": "2022-01-10T06:00:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"execDuration": {
				"avg": 0.0,
				"min": 0.0,
				"max": 0.0
			}
		  }
        ],
        "onOriginResponse": [
          {
			"startDateTime": "2022-01-10T03:45:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"execDuration": {
				"avg": 0.0,
				"min": 0.0,
				"max": 0.0
			}
		  },
		  {
			"startDateTime": "2022-01-10T06:00:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"execDuration": {
				"avg": 0.0,
				"min": 0.0,
				"max": 0.0
			}
		  }
        ],
        "onClientResponse": [
          {
			"startDateTime": "2022-01-10T03:00:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 8,
			"execDuration": {
				"avg": 0.029625,
				"min": 0.018,
				"max": 0.059
			}
		  },
		  {
			"startDateTime": "2022-01-10T03:05:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 15,
			"execDuration": {
				"avg": 0.031933333,
				"min": 0.018,
				"max": 0.062
			}
		  }
        ],
        "responseProvider": [
          {
			"startDateTime": "2022-01-10T03:45:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"execDuration": {
				"avg": 0.0,
				"min": 0.0,
				"max": 0.0
			}
		  },
		  {
			"startDateTime": "2022-01-10T06:00:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"execDuration": {
				"avg": 0.0,
				"min": 0.0,
				"max": 0.0
			}
		  }
        ],
        "init": [
          {
            "startDateTime": "2022-01-10T03:05:00Z",
            "edgeWorkerVersion": "10.18",
			"invocations": 7,
            "initDuration": {
              "avg": 1.03,
			  "min": 1.03,
			  "max": 1.03
            }
          },
          {
			"startDateTime": "2022-01-10T05:10:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"initDuration": {
				"avg": 0.729,
				"min": 0.729,
				"max": 0.729
			}
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
				Start:       "2022-01-10T03:00:00Z",
				End:         "2022-01-17T10:43:15Z",
				Data: []ReportData{
					{
						EdgeWorkerID: 37017,
						Data: Data{
							OnClientRequest: &OnClientRequest{
								{
									StartDateTime:     "2022-01-10T03:00:00Z",
									EdgeWorkerVersion: "10.18",
									ExecDuration: &Summary{
										Avg: 0.113125,
										Min: 0.095,
										Max: 0.148,
									},
									Invocations: 8,
								},
								{
									StartDateTime:     "2022-01-10T03:05:00Z",
									EdgeWorkerVersion: "10.18",
									ExecDuration: &Summary{
										Avg: 0.14859999999999998,
										Min: 0.087,
										Max: 0.348,
									},
									Invocations: 15,
								},
							},
							OnOriginRequest: &OnOriginRequest{
								{
									StartDateTime:     "2022-01-10T03:45:00Z",
									EdgeWorkerVersion: "10.18",
									ExecDuration: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
									Invocations: 1,
								},
								{
									StartDateTime:     "2022-01-10T06:00:00Z",
									EdgeWorkerVersion: "10.18",
									ExecDuration: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
									Invocations: 1,
								},
							},
							OnOriginResponse: &OnOriginResponse{
								{
									StartDateTime:     "2022-01-10T03:45:00Z",
									EdgeWorkerVersion: "10.18",
									ExecDuration: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
									Invocations: 1,
								},
								{
									StartDateTime:     "2022-01-10T06:00:00Z",
									EdgeWorkerVersion: "10.18",
									ExecDuration: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
									Invocations: 1,
								},
							},
							OnClientResponse: &OnClientResponse{
								{
									StartDateTime:     "2022-01-10T03:00:00Z",
									EdgeWorkerVersion: "10.18",
									ExecDuration: &Summary{
										Avg: 0.029625,
										Min: 0.018,
										Max: 0.059,
									},
									Invocations: 8,
								},
								{
									StartDateTime:     "2022-01-10T03:05:00Z",
									EdgeWorkerVersion: "10.18",
									ExecDuration: &Summary{
										Avg: 0.031933333,
										Min: 0.018,
										Max: 0.062,
									},
									Invocations: 15,
								},
							},
							ResponseProvider: &ResponseProvider{
								{
									StartDateTime:     "2022-01-10T03:45:00Z",
									EdgeWorkerVersion: "10.18",
									ExecDuration: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
									Invocations: 1,
								},
								{
									StartDateTime:     "2022-01-10T06:00:00Z",
									EdgeWorkerVersion: "10.18",
									ExecDuration: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
									Invocations: 1,
								},
							},
							Init: &Init{
								{
									StartDateTime:     "2022-01-10T03:05:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       7,
									InitDuration: Summary{
										Avg: 1.03,
										Min: 1.03,
										Max: 1.03,
									},
								},
								{
									StartDateTime:     "2022-01-10T05:10:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									InitDuration: Summary{
										Avg: 0.729,
										Min: 0.729,
										Max: 0.729,
									},
								},
							},
						},
					},
				},
			},
		},
		"200 OK - get report with id 3": {
			params: GetReportRequest{
				ReportID:   3,
				Start:      "2021-12-04T00:00:00Z",
				End:        "2022-01-17T11:28:11Z",
				EdgeWorker: "37017",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/edgeworkers/v1/reports/3?edgeWorker=37017&end=2022-01-17T11%3A28%3A11Z&start=2021-12-04T00%3A00%3A00Z",
			responseBody: `
{
    "reportId": 3,
    "name": "Execution statuses by EdgeWorker ID and event handler",
    "description": "This report lists execution statuses, grouped by EdgeWorker ID and event handler.",
    "start": "2022-01-10T03:00:00Z",
    "end": "2022-01-17T11:28:11Z",
    "data": [
    {
      "edgeWorkerId": 37017,
      "data": {
		"responseProvider": [
		  {
			"startDateTime": "2022-01-10T03:45:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"status": "unimplementedEventHandler"
		  },
		  {
			"startDateTime": "2022-01-10T06:00:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"status": "unimplementedEventHandler"
		  }
        ],
	    "onOriginRequest": [
		  {
			"startDateTime": "2022-01-10T03:45:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"status": "unimplementedEventHandler"
		  },
		  {
			"startDateTime": "2022-01-10T06:00:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"status": "unimplementedEventHandler"
		  }
        ],
		"onClientResponse": [
		  {
			"startDateTime": "2022-01-10T03:00:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 8,
			"status": "success"
		  },
		  {
			"startDateTime": "2022-01-10T03:05:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 15,
			"status": "success"
		  }
        ],
		"onOriginResponse": [
		  {
			"startDateTime": "2022-01-10T03:45:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"status": "unimplementedEventHandler"
		  },
		  {
			"startDateTime": "2022-01-10T06:00:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 1,
			"status": "unimplementedEventHandler"
		  }
        ],
        "onClientRequest": [
		  {
			"startDateTime": "2022-01-10T03:00:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 8,
			"status": "success"
		  },
		  {
			"startDateTime": "2022-01-10T03:05:00Z",
			"edgeWorkerVersion": "10.18",
			"invocations": 15,
			"status": "success"
		  }
        ]
      }
    }
  ]
}`,
			expectedResponse: &GetReportResponse{
				ReportID:    3,
				Name:        "Execution statuses by EdgeWorker ID and event handler",
				Description: "This report lists execution statuses, grouped by EdgeWorker ID and event handler.",
				Start:       "2022-01-10T03:00:00Z",
				End:         "2022-01-17T11:28:11Z",
				Data: []ReportData{
					{
						EdgeWorkerID: 37017,
						Data: Data{
							ResponseProvider: &ResponseProvider{
								{
									StartDateTime:     "2022-01-10T03:45:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Status:            tools.StringPtr("unimplementedEventHandler"),
								},
								{
									StartDateTime:     "2022-01-10T06:00:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Status:            tools.StringPtr("unimplementedEventHandler"),
								},
							},
							OnOriginRequest: &OnOriginRequest{
								{
									StartDateTime:     "2022-01-10T03:45:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Status:            tools.StringPtr("unimplementedEventHandler"),
								},
								{
									StartDateTime:     "2022-01-10T06:00:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Status:            tools.StringPtr("unimplementedEventHandler"),
								},
							},
							OnClientResponse: &OnClientResponse{
								{
									StartDateTime:     "2022-01-10T03:00:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       8,
									Status:            tools.StringPtr("success"),
								},
								{
									StartDateTime:     "2022-01-10T03:05:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       15,
									Status:            tools.StringPtr("success"),
								},
							},
							OnOriginResponse: &OnOriginResponse{
								{
									StartDateTime:     "2022-01-10T03:45:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Status:            tools.StringPtr("unimplementedEventHandler"),
								},
								{
									StartDateTime:     "2022-01-10T06:00:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Status:            tools.StringPtr("unimplementedEventHandler"),
								},
							},
							OnClientRequest: &OnClientRequest{
								{
									StartDateTime:     "2022-01-10T03:00:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       8,
									Status:            tools.StringPtr("success"),
								},
								{
									StartDateTime:     "2022-01-10T03:05:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       15,
									Status:            tools.StringPtr("success"),
								},
							},
						},
					},
				},
			},
		},
		"200 OK - get report with id 4": {
			params: GetReportRequest{
				ReportID:   4,
				Start:      "2022-01-10T03:00:00Z",
				End:        "2022-01-17T12:07:48Z",
				EdgeWorker: "37017",
			},
			responseStatus: http.StatusOK,
			expectedPath:   "/edgeworkers/v1/reports/4?edgeWorker=37017&end=2022-01-17T12%3A07%3A48Z&start=2022-01-10T03%3A00%3A00Z",
			responseBody: `
{
    "reportId": 4,
    "name": "Memory usage by EdgeWorker ID and event handler",
    "description": "This report lists memory usage, grouped by EdgeWorker ID and event handler.",
    "start": "2022-01-10T03:00:00Z",
    "end": "2022-01-17T12:07:48Z",
	"data": [
	    {
            "edgeWorkerId": 37017,
            "data": {
                "responseProvider": [
                    {
                        "startDateTime": "2022-01-10T03:45:00Z",
                        "edgeWorkerVersion": "10.18",
                        "invocations": 1,
                        "memory": {
                            "avg": 0.0,
                            "min": 0.0,
                            "max": 0.0
                        }
                    },
                    {
                        "startDateTime": "2022-01-10T06:00:00Z",
                        "edgeWorkerVersion": "10.18",
                        "invocations": 1,
                        "memory": {
                            "avg": 0.0,
                            "min": 0.0,
                            "max": 0.0
                        }
                    }
				],
				"onOriginRequest": [
				  {
					"startDateTime": "2022-01-10T03:45:00Z",
					"edgeWorkerVersion": "10.18",
					"invocations": 1,
					"memory": {
						"avg": 0.0,
						"min": 0.0,
						"max": 0.0
					}
				  },
				  {
					"startDateTime": "2022-01-10T06:00:00Z",
					"edgeWorkerVersion": "10.18",
					"invocations": 1,
					"memory": {
						"avg": 0.0,
						"min": 0.0,
						"max": 0.0
					}
				  }
			    ],
				"onClientResponse": [
				  {
					"startDateTime": "2022-01-10T03:00:00Z",
					"edgeWorkerVersion": "10.18",
					"invocations": 8,
					"memory": {
						"avg": 0.0,
						"min": 0.0,
						"max": 0.0
					}
				  },
				  {
					"startDateTime": "2022-01-10T03:05:00Z",
					"edgeWorkerVersion": "10.18",
					"invocations": 15,
					"memory": {
						"avg": 0.0,
						"min": 0.0,
						"max": 0.0
					}
				  }
				],
			    "onOriginResponse": [
				  {
					"startDateTime": "2022-01-10T03:45:00Z",
					"edgeWorkerVersion": "10.18",
					"invocations": 1,
					"memory": {
						"avg": 0.0,
						"min": 0.0,
						"max": 0.0
					}
				  },
				  {
					"startDateTime": "2022-01-10T06:00:00Z",
					"edgeWorkerVersion": "10.18",
					"invocations": 1,
					"memory": {
						"avg": 0.0,
						"min": 0.0,
						"max": 0.0
					}
				  }
				],
				"onClientRequest": [
				  {
					"startDateTime": "2022-01-10T03:00:00Z",
					"edgeWorkerVersion": "10.18",
					"invocations": 8,
					"memory": {
						"avg": 1480.0,
						"min": 1256.0,
						"max": 1704.0
					}
				  },
				  {
					"startDateTime": "2022-01-10T03:05:00Z",
					"edgeWorkerVersion": "10.18",
					"invocations": 15,
					"memory": {
						"avg": 1608.266667,
						"min": 1256.0,
						"max": 2544.0
					}
				  }
				]
			}
	    }
    ]
}`,
			expectedResponse: &GetReportResponse{
				ReportID:    4,
				Name:        "Memory usage by EdgeWorker ID and event handler",
				Description: "This report lists memory usage, grouped by EdgeWorker ID and event handler.",
				Start:       "2022-01-10T03:00:00Z",
				End:         "2022-01-17T12:07:48Z",
				Data: []ReportData{
					{
						EdgeWorkerID: 37017,
						Data: Data{
							ResponseProvider: &ResponseProvider{
								{
									StartDateTime:     "2022-01-10T03:45:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Memory: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
								},
								{
									StartDateTime:     "2022-01-10T06:00:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Memory: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
								},
							},
							OnOriginRequest: &OnOriginRequest{
								{
									StartDateTime:     "2022-01-10T03:45:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Memory: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
								},
								{
									StartDateTime:     "2022-01-10T06:00:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Memory: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
								},
							},
							OnClientResponse: &OnClientResponse{
								{
									StartDateTime:     "2022-01-10T03:00:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       8,
									Memory: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
								},
								{
									StartDateTime:     "2022-01-10T03:05:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       15,
									Memory: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
								},
							},
							OnOriginResponse: &OnOriginResponse{
								{
									StartDateTime:     "2022-01-10T03:45:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Memory: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
								},
								{
									StartDateTime:     "2022-01-10T06:00:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       1,
									Memory: &Summary{
										Avg: 0.0,
										Min: 0.0,
										Max: 0.0,
									},
								},
							},
							OnClientRequest: &OnClientRequest{
								{
									StartDateTime:     "2022-01-10T03:00:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       8,
									Memory: &Summary{
										Avg: 1480.0,
										Min: 1256.0,
										Max: 1704.0,
									},
								},
								{
									StartDateTime:     "2022-01-10T03:05:00Z",
									EdgeWorkerVersion: "10.18",
									Invocations:       15,
									Memory: &Summary{
										Avg: 1608.266667,
										Min: 1256.0,
										Max: 2544.0,
									},
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
