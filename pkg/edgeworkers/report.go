package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Reports is an edgeworkers reports API interface
	Reports interface {
		// GetReport gets details for an EdgeWorker
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/reportsreportid#get-report
		GetReport(context.Context, GetReportRequest) (*GetReportResponse, error)

		// ListReports lists EdgeWorker reports
		//
		// See: https://techdocs.akamai.com/edgeworkers/reference/reports
		ListReports(context.Context) (*ListReportsResponse, error)
	}

	// GetReportRequest contains parameters used to get an EdgeWorker report
	GetReportRequest struct {
		ReportID int
		Start    string
		// If end date is not provided, then API will assign current date by default
		End          string
		EdgeWorker   string
		Status       *string
		EventHandler *string
	}

	// GetReportResponse represents a response object returned by GetReport
	GetReportResponse struct {
		ReportID    int          `json:"reportId"`
		Name        string       `json:"name"`
		Description string       `json:"description"`
		Start       string       `json:"start"`
		End         string       `json:"end"`
		Data        []ReportData `json:"data"`
	}

	// ReportData represents report data
	ReportData struct {
		EdgeWorkerID int  `json:"edgeWorkerId"`
		Data         Data `json:"data"`
	}

	// Data represents data object
	Data struct {
		OnClientRequest  OnClientRequest  `json:"onClientRequest"`
		OnOriginRequest  OnOriginRequest  `json:"onOriginRequest"`
		OnOriginResponse OnOriginResponse `json:"onOriginResponse"`
		OnClientResponse OnClientResponse `json:"onClientResponse"`
		ResponseProvider ResponseProvider `json:"responseProvider"`
		Init             Init             `json:"init"`
	}

	// OnClientRequest represents OnClientRequest list
	OnClientRequest []OnRequestAndResponse
	// OnOriginRequest represents OnOriginRequest list
	OnOriginRequest []OnRequestAndResponse
	// OnOriginResponse represents OnOriginResponse list
	OnOriginResponse []OnRequestAndResponse
	// OnClientResponse represents OnClientResponse list
	OnClientResponse []OnRequestAndResponse
	// ResponseProvider represents ResponseProvider list
	ResponseProvider []OnRequestAndResponse
	// Init represents Init list
	Init []InitObject

	// OnRequestAndResponse represents object structure for OnClientRequest, OnOriginRequest, OnOriginResponse,
	// OnClientResponse, ResponseProvider fields
	OnRequestAndResponse struct {
		StartDateTime     string   `json:"startDateTime"`
		EdgeWorkerVersion string   `json:"edgeWorkerVersion"`
		ExecDuration      Duration `json:"execDuration"`
		Invocations       int      `json:"invocations"`
	}

	// InitObject represents object structure for Init field
	InitObject struct {
		StartDateTime     string   `json:"startDateTime"`
		EdgeWorkerVersion string   `json:"edgeWorkerVersion"`
		InitDuration      Duration `json:"initDuration"`
		Invocations       int      `json:"invocations"`
	}

	// Duration represents Duration object
	Duration struct {
		Avg int `json:"avg"`
		Min int `json:"min"`
		Max int `json:"max"`
	}

	// ListReportsResponse represents list of report types
	ListReportsResponse struct {
		Reports []ReportResponse `json:"reports"`
	}

	// ReportResponse represents report type object
	ReportResponse struct {
		ReportID    int    `json:"reportId"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Unavailable bool   `json:"unavailable"`
	}
)

// These constants represent all the valid report statuses
const (
	StatusSuccess                   = "success"
	StatusGenericError              = "genericError"
	StatusUnknownEdgeWorkerID       = "unknownEdgeWorkerId"
	StatusUnimplementedEventHandler = "unimplementedEventHandler"
	StatusRuntimeError              = "runtimeError"
	StatusExecutionError            = "executionError"
	StatusTimeoutError              = "timeoutError"
	StatusResourceLimitHit          = "resourceLimitHit"
	StatusCPUTimeoutError           = "cpuTimeoutError"
	StatusWallTimeoutError          = "wallTimeoutError"
	StatusInitCPUTimeoutError       = "initCpuTimeoutError"
	StatusInitWallTimeoutError      = "initWallTimeoutError"
)

// These constants represent all the valid report event handlers
const (
	EventHandlerOnClientRequest  = "onClientRequest"
	EventHandlerOnOriginRequest  = "onOriginRequest"
	EventHandlerOnOriginResponse = "onOriginResponse"
	EventHandlerOnClientResponse = "onClientResponse"
	EventHandlerResponseProvider = "responseProvider"
)

var (
	// ErrGetReport is returned in case an error occurs on GetReport operation
	ErrGetReport = errors.New("get an EdgeWorker report")
	// ErrListReports is returned in case an error occurs on ListReports operation
	ErrListReports = errors.New("get EdgeWorker reports")
)

// Validate validates GetReportRequest
func (r GetReportRequest) Validate() error {
	return validation.Errors{
		"ReportID": validation.Validate(r.ReportID, validation.Required),
		"Start": validation.Validate(r.Start, validation.Required, validation.Date("2006-01-02T15:04:05.999Z").Error(
			fmt.Sprintf("value '%s' is invalid. It must have format '2006-01-02T15:04:05.999Z'", r.Start))),
		"End": validation.Validate(r.End, validation.Date("2006-01-02T15:04:05.999Z").Error(
			fmt.Sprintf("value '%s' is invalid. It must have format '2006-01-02T15:04:05.999Z'", r.End))),
		"EdgeWorker": validation.Validate(r.EdgeWorker, validation.Required),
		"Status": validation.Validate(r.Status, validation.NilOrNotEmpty, validation.In(StatusSuccess, StatusGenericError, StatusUnknownEdgeWorkerID, StatusUnimplementedEventHandler,
			StatusRuntimeError, StatusExecutionError, StatusTimeoutError, StatusResourceLimitHit, StatusCPUTimeoutError, StatusWallTimeoutError, StatusInitCPUTimeoutError,
			StatusInitWallTimeoutError).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'success', 'genericError', "+
			"'unknownEdgeWorkerId', 'unimplementedEventHandler', 'runtimeError', 'executionError', 'timeoutError', "+
			"'resourceLimitHit', 'cpuTimeoutError', 'wallTimeoutError', 'initCpuTimeoutError', 'initWallTimeoutError'", stringFromPtr(r.Status)))),
		"EventHandler": validation.Validate(r.EventHandler, validation.NilOrNotEmpty, validation.In(EventHandlerOnClientRequest, EventHandlerOnOriginRequest, EventHandlerOnOriginResponse,
			EventHandlerOnClientResponse, EventHandlerResponseProvider).Error(fmt.Sprintf("value '%s' is invalid. Must be one of: 'onClientRequest', "+
			"'onOriginRequest', 'onOriginResponse', 'onClientResponse', 'responseProvider'", stringFromPtr(r.EventHandler)))),
	}.Filter()
}

func stringFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (e *edgeworkers) GetReport(ctx context.Context, params GetReportRequest) (*GetReportResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("GetReport")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetReport, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/edgeworkers/v1/reports/%d", params.ReportID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetReport, err)
	}

	q := uri.Query()
	q.Add("edgeWorker", params.EdgeWorker)
	q.Add("start", params.Start)
	if params.End != "" {
		q.Add("end", params.End)
	}
	if params.Status != nil {
		status := *params.Status
		q.Add("status", status)
	}
	if params.EventHandler != nil {
		status := *params.EventHandler
		q.Add("eventHandler", status)
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetReport, err)
	}

	var result GetReportResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetReport, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetReport, e.Error(resp))
	}

	return &result, nil
}

func (e *edgeworkers) ListReports(ctx context.Context) (*ListReportsResponse, error) {
	logger := e.Log(ctx)
	logger.Debug("ListReports")

	uri := fmt.Sprintf("/edgeworkers/v1/reports")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListReports, err)
	}

	var result ListReportsResponse
	resp, err := e.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListReports, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListReports, e.Error(resp))
	}

	return &result, nil
}
