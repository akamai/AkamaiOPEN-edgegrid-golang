package datastream

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Stream is a ds stream operations API interface
	Stream interface {
		// CreateStream creates a stream
		//
		// See: https://techdocs.akamai.com/datastream2/reference/post-stream
		CreateStream(context.Context, CreateStreamRequest) (*StreamUpdate, error)

		// GetStream gets stream details
		//
		// See: https://techdocs.akamai.com/datastream2/reference/get-stream
		GetStream(context.Context, GetStreamRequest) (*DetailedStreamVersion, error)

		// UpdateStream updates a stream
		//
		// See: https://techdocs.akamai.com/datastream2/reference/put-stream
		UpdateStream(context.Context, UpdateStreamRequest) (*StreamUpdate, error)

		// DeleteStream deletes a stream
		//
		// See: https://techdocs.akamai.com/datastream2/reference/delete-stream
		DeleteStream(context.Context, DeleteStreamRequest) (*DeleteStreamResponse, error)

		// ListStreams retrieves list of streams
		//
		// See: https://techdocs.akamai.com/datastream2/reference/get-streams
		ListStreams(context.Context, ListStreamsRequest) ([]StreamDetails, error)
	}

	// DetailedStreamVersion is returned from GetStream
	DetailedStreamVersion struct {
		ActivationStatus ActivationStatus   `json:"activationStatus"`
		Config           Config             `json:"config"`
		Connectors       []ConnectorDetails `json:"connectors"`
		ContractID       string             `json:"contractId"`
		CreatedBy        string             `json:"createdBy"`
		CreatedDate      string             `json:"createdDate"`
		Datasets         []DataSets         `json:"datasets"`
		EmailIDs         string             `json:"emailIds"`
		Errors           []Errors           `json:"errors"`
		GroupID          int                `json:"groupId"`
		GroupName        string             `json:"groupName"`
		ModifiedBy       string             `json:"modifiedBy"`
		ModifiedDate     string             `json:"modifiedDate"`
		ProductID        string             `json:"productId"`
		ProductName      string             `json:"productName"`
		Properties       []Property         `json:"properties"`
		StreamID         int64              `json:"streamId"`
		StreamName       string             `json:"streamName"`
		StreamType       StreamType         `json:"streamType"`
		StreamVersionID  int64              `json:"streamVersionId"`
		TemplateName     TemplateName       `json:"templateName"`
	}

	// ConnectorDetails provides detailed information about the connector’s configuration in the stream
	ConnectorDetails struct {
		AuthenticationType AuthenticationType `json:"authenticationType"`
		ConnectorID        int                `json:"connectorId"`
		CompressLogs       bool               `json:"compressLogs"`
		ConnectorName      string             `json:"connectorName"`
		ConnectorType      ConnectorType      `json:"connectorType"`
		Path               string             `json:"path"`
		URL                string             `json:"url"`
		Endpoint           string             `json:"endpoint"`
		IndexName          string             `json:"indexName"`
		ServiceAccountName string             `json:"serviceAccountName"`
		ProjectID          string             `json:"projectId"`
		Service            string             `json:"service"`
		Bucket             string             `json:"bucket"`
		Tags               string             `json:"tags"`
		Region             string             `json:"region"`
		AccountName        string             `json:"accountName"`
		Namespace          string             `json:"namespace"`
		ContainerName      string             `json:"containerName"`
		Source             string             `json:"source"`
		ContentType        string             `json:"contentType"`
		CustomHeaderName   string             `json:"customHeaderName"`
		CustomHeaderValue  string             `json:"customHeaderValue"`
		TLSHostname        string             `json:"tlsHostname"`
		MTLS               string             `json:"mTLS"`
	}

	// StreamConfiguration is used in CreateStream as a request body
	StreamConfiguration struct {
		ActivateNow     bool                `json:"activateNow"`
		Config          Config              `json:"config"`
		Connectors      []AbstractConnector `json:"connectors,omitempty"`
		ContractID      string              `json:"contractId"`
		DatasetFieldIDs []int               `json:"datasetFieldIds"`
		EmailIDs        string              `json:"emailIds,omitempty"`
		GroupID         *int                `json:"groupId"`
		PropertyIDs     []int               `json:"propertyIds"`
		StreamName      string              `json:"streamName"`
		StreamType      StreamType          `json:"streamType"`
		TemplateName    TemplateName        `json:"templateName"`
	}

	// Config of the configuration of log lines, names of the files sent to a destination, and delivery frequency for these files
	Config struct {
		Delimiter        *DelimiterType `json:"delimiter,omitempty"`
		Format           FormatType     `json:"format,omitempty"`
		Frequency        Frequency      `json:"frequency"`
		UploadFilePrefix string         `json:"uploadFilePrefix,omitempty"`
		UploadFileSuffix string         `json:"uploadFileSuffix,omitempty"`
	}

	// The Frequency of collecting logs from each uploader and sending these logs to a destination.
	Frequency struct {
		TimeInSec TimeInSec `json:"timeInSec"`
	}

	// DataSets is a list of fields selected from the associated template that the stream monitors in logs
	DataSets struct {
		DatasetFields           []DatasetFields `json:"datasetFields"`
		DatasetGroupDescription string          `json:"datasetGroupDescription"`
		DatasetGroupName        string          `json:"datasetGroupName"`
	}

	// DatasetFields is list of data set fields selected from the associated template that the stream monitors in logs
	DatasetFields struct {
		DatasetFieldID          int    `json:"datasetFieldId"`
		DatasetFieldDescription string `json:"datasetFieldDescription"`
		DatasetFieldJsonKey     string `json:"datasetFieldJsonKey"`
		DatasetFieldName        string `json:"datasetFieldName"`
		Order                   int    `json:"order"`
	}

	// Errors associated to the stream
	Errors struct {
		Detail string `json:"detail"`
		Title  string `json:"title"`
		Type   string `json:"type"`
	}

	// Property identifies the properties monitored in the stream.
	Property struct {
		Hostnames    []string `json:"hostnames"`
		ProductID    string   `json:"productId"`
		ProductName  string   `json:"productName"`
		PropertyID   int      `json:"propertyId"`
		PropertyName string   `json:"propertyName"`
	}

	// ActivationStatus is used to create an enum of possible ActivationStatus values
	ActivationStatus string

	// AbstractConnector is an interface for all Connector types
	AbstractConnector interface {
		SetConnectorType()
		Validate() error
	}

	// DelimiterType enum
	DelimiterType string

	// FormatType enum
	FormatType string

	// TemplateName enum
	TemplateName string

	// StreamType enum
	StreamType string

	// TimeInSec enum
	TimeInSec int

	// CreateStreamRequest is passed to CreateStream
	CreateStreamRequest struct {
		StreamConfiguration StreamConfiguration
	}

	// GetStreamRequest is passed to GetStream
	GetStreamRequest struct {
		StreamID int64
		Version  *int64
	}

	// UpdateStreamRequest is passed to UpdateStream
	UpdateStreamRequest struct {
		StreamID            int64
		StreamConfiguration StreamConfiguration
	}

	// StreamUpdate contains information about stream ID and version
	StreamUpdate struct {
		StreamVersionKey StreamVersionKey `json:"streamVersionKey"`
	}

	// StreamVersionKey contains information about stream ID and version
	StreamVersionKey struct {
		StreamID        int64 `json:"streamId"`
		StreamVersionID int64 `json:"streamVersionId"`
	}

	// DeleteStreamRequest is passed to DeleteStream
	DeleteStreamRequest struct {
		StreamID int64
	}

	// DeleteStreamResponse is returned from DeleteStream
	DeleteStreamResponse struct {
		Message string `json:"message"`
	}

	// ListStreamsRequest is passed to ListStreams
	ListStreamsRequest struct {
		GroupID *int
	}

	// StreamDetails list is returned from ListStreams method
	StreamDetails struct {
		ActivationStatus ActivationStatus `json:"activationStatus"`
		Archived         bool             `json:"archived"`
		Connectors       string           `json:"connectors"`
		ContractID       string           `json:"contractId"`
		CreatedBy        string           `json:"createdBy"`
		CreatedDate      string           `json:"createdDate"`
		CurrentVersionID int64            `json:"currentVersionId"`
		Errors           []Errors         `json:"errors"`
		GroupID          int              `json:"groupId"`
		GroupName        string           `json:"groupName"`
		Properties       []Property       `json:"properties"`
		StreamID         int64            `json:"streamId"`
		StreamName       string           `json:"streamName"`
		StreamTypeName   string           `json:"streamTypeName"`
		StreamVersionID  int64            `json:"streamVersionId"`
	}
)

const (
	// ActivationStatusActivated const
	ActivationStatusActivated ActivationStatus = "ACTIVATED"
	// ActivationStatusDeactivated const
	ActivationStatusDeactivated ActivationStatus = "DEACTIVATED"
	// ActivationStatusActivating const
	ActivationStatusActivating ActivationStatus = "ACTIVATING"
	// ActivationStatusDeactivating const state
	ActivationStatusDeactivating ActivationStatus = "DEACTIVATING"
	// ActivationStatusInactive const
	ActivationStatusInactive ActivationStatus = "INACTIVE"

	// StreamTypeRawLogs const
	StreamTypeRawLogs StreamType = "RAW_LOGS"

	// TemplateNameEdgeLogs const
	TemplateNameEdgeLogs TemplateName = "EDGE_LOGS"

	// DelimiterTypeSpace const
	DelimiterTypeSpace DelimiterType = "SPACE"

	// FormatTypeStructured const
	FormatTypeStructured FormatType = "STRUCTURED"
	// FormatTypeJson const
	FormatTypeJson FormatType = "JSON"

	// TimeInSec30 const
	TimeInSec30 TimeInSec = 30
	// TimeInSec60 const
	TimeInSec60 TimeInSec = 60
)

// Validate validates CreateStreamRequest
func (r CreateStreamRequest) Validate() error {
	return validation.Errors{
		"StreamConfiguration.Config":                     validation.Validate(r.StreamConfiguration.Config, validation.Required),
		"StreamConfiguration.Config.Delimiter":           validation.Validate(r.StreamConfiguration.Config.Delimiter, validation.When(r.StreamConfiguration.Config.Format == FormatTypeStructured, validation.Required, validation.In(DelimiterTypeSpace)), validation.When(r.StreamConfiguration.Config.Format == FormatTypeJson, validation.Nil)),
		"StreamConfiguration.Config.Format":              validation.Validate(r.StreamConfiguration.Config.Format, validation.Required, validation.In(FormatTypeStructured, FormatTypeJson), validation.When(r.StreamConfiguration.Config.Delimiter != nil, validation.Required, validation.In(FormatTypeStructured))),
		"StreamConfiguration.Config.Frequency":           validation.Validate(r.StreamConfiguration.Config.Frequency, validation.Required),
		"StreamConfiguration.Config.Frequency.TimeInSec": validation.Validate(r.StreamConfiguration.Config.Frequency.TimeInSec, validation.Required, validation.In(TimeInSec30, TimeInSec60)),
		"StreamConfiguration.Connectors":                 validation.Validate(r.StreamConfiguration.Connectors, validation.Required, validation.Length(1, 1)),
		"StreamConfiguration.ContractId":                 validation.Validate(r.StreamConfiguration.ContractID, validation.Required),
		"StreamConfiguration.DatasetFields":              validation.Validate(r.StreamConfiguration.DatasetFieldIDs, validation.Required),
		"StreamConfiguration.GroupID":                    validation.Validate(r.StreamConfiguration.GroupID, validation.Required),
		"StreamConfiguration.PropertyIDs":                validation.Validate(r.StreamConfiguration.PropertyIDs, validation.Required),
		"StreamConfiguration.StreamName":                 validation.Validate(r.StreamConfiguration.StreamName, validation.Required),
		"StreamConfiguration.StreamType":                 validation.Validate(r.StreamConfiguration.StreamType, validation.Required, validation.In(StreamTypeRawLogs)),
		"StreamConfiguration.TemplateName":               validation.Validate(r.StreamConfiguration.TemplateName, validation.Required, validation.In(TemplateNameEdgeLogs)),
	}.Filter()
}

// Validate validates GetStreamRequest
func (r GetStreamRequest) Validate() error {
	return validation.Errors{
		"streamId": validation.Validate(r.StreamID, validation.Required),
	}.Filter()
}

// Validate validates UpdateStreamRequest
func (r UpdateStreamRequest) Validate() error {
	return validation.Errors{
		"StreamConfiguration.Config":                     validation.Validate(r.StreamConfiguration.Config, validation.Required),
		"StreamConfiguration.Config.Delimiter":           validation.Validate(r.StreamConfiguration.Config.Delimiter, validation.When(r.StreamConfiguration.Config.Format == FormatTypeStructured, validation.Required, validation.In(DelimiterTypeSpace)), validation.When(r.StreamConfiguration.Config.Format == FormatTypeJson, validation.Nil)),
		"StreamConfiguration.Config.Format":              validation.Validate(r.StreamConfiguration.Config.Format, validation.In(FormatTypeStructured, FormatTypeJson)),
		"StreamConfiguration.Config.Frequency":           validation.Validate(r.StreamConfiguration.Config.Frequency, validation.Required),
		"StreamConfiguration.Config.Frequency.TimeInSec": validation.Validate(r.StreamConfiguration.Config.Frequency.TimeInSec, validation.Required, validation.In(TimeInSec30, TimeInSec60)),
		"StreamConfiguration.Connectors":                 validation.Validate(r.StreamConfiguration.Connectors, validation.When(r.StreamConfiguration.Connectors != nil, validation.Length(1, 1))),
		"StreamConfiguration.ContractId":                 validation.Validate(r.StreamConfiguration.ContractID, validation.Required),
		"StreamConfiguration.DatasetFields":              validation.Validate(r.StreamConfiguration.DatasetFieldIDs, validation.Required),
		"StreamConfiguration.GroupID":                    validation.Validate(r.StreamConfiguration.GroupID, validation.Nil),
		"StreamConfiguration.PropertyIDs":                validation.Validate(r.StreamConfiguration.PropertyIDs, validation.Required),
		"StreamConfiguration.StreamName":                 validation.Validate(r.StreamConfiguration.StreamName, validation.Required),
		"StreamConfiguration.StreamType":                 validation.Validate(r.StreamConfiguration.StreamType, validation.Required, validation.In(StreamTypeRawLogs)),
		"StreamConfiguration.TemplateName":               validation.Validate(r.StreamConfiguration.TemplateName, validation.Required, validation.In(TemplateNameEdgeLogs)),
	}.Filter()
}

// Validate validates DeleteStreamRequest
func (r DeleteStreamRequest) Validate() error {
	return validation.Errors{
		"streamId": validation.Validate(r.StreamID, validation.Required),
	}.Filter()
}

var (
	// ErrCreateStream represents error when creating stream fails
	ErrCreateStream = errors.New("creating stream")
	// ErrGetStream represents error when fetching stream fails
	ErrGetStream = errors.New("fetching stream information")
	// ErrUpdateStream represents error when updating stream fails
	ErrUpdateStream = errors.New("updating stream")
	// ErrDeleteStream represents error when deleting stream fails
	ErrDeleteStream = errors.New("deleting stream")
	// ErrListStreams represents error when listing streams fails
	ErrListStreams = errors.New("listing streams")
)

func (d *ds) CreateStream(ctx context.Context, params CreateStreamRequest) (*StreamUpdate, error) {
	logger := d.Log(ctx)
	logger.Debug("CreateStream")

	setConnectorTypes(&params.StreamConfiguration)
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateStream, ErrStructValidation, err)
	}

	uri := "/datastream-config-api/v1/log/streams"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateStream, err)
	}

	var rval StreamUpdate
	resp, err := d.Exec(req, &rval, params.StreamConfiguration)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateStream, err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrCreateStream, d.Error(resp))
	}

	return &rval, nil
}

func (d *ds) GetStream(ctx context.Context, params GetStreamRequest) (*DetailedStreamVersion, error) {
	logger := d.Log(ctx)
	logger.Debug("GetStream")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetStream, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v1/log/streams/%d",
		params.StreamID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetStream, err)
	}

	if params.Version != nil {
		query := uri.Query()
		query.Add("version", strconv.FormatInt(*params.Version, 10))
		uri.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetStream, err)
	}

	var rval DetailedStreamVersion
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetStream, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetStream, d.Error(resp))
	}

	return &rval, nil
}

func (d *ds) UpdateStream(ctx context.Context, params UpdateStreamRequest) (*StreamUpdate, error) {
	logger := d.Log(ctx)
	logger.Debug("UpdateStream")

	setConnectorTypes(&params.StreamConfiguration)
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateStream, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v1/log/streams/%d",
		params.StreamID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateStream, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateStream, err)
	}

	var rval StreamUpdate
	resp, err := d.Exec(req, &rval, params.StreamConfiguration)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateStream, err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("%s: %w", ErrUpdateStream, d.Error(resp))
	}

	return &rval, nil
}

func (d *ds) DeleteStream(ctx context.Context, params DeleteStreamRequest) (*DeleteStreamResponse, error) {
	logger := d.Log(ctx)
	logger.Debug("DeleteStream")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrDeleteStream, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v1/log/streams/%d",
		params.StreamID),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrDeleteStream, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeleteStream, err)
	}

	var rval DeleteStreamResponse
	resp, err := d.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeleteStream, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrDeleteStream, d.Error(resp))
	}

	return &rval, nil
}

func (d *ds) ListStreams(ctx context.Context, params ListStreamsRequest) ([]StreamDetails, error) {
	logger := d.Log(ctx)
	logger.Debug("ListStreams")

	uri, err := url.Parse("/datastream-config-api/v1/log/streams")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListStreams, err)
	}

	q := uri.Query()
	if params.GroupID != nil {
		q.Add("groupId", fmt.Sprintf("%d", *params.GroupID))
	}

	uri.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListStreams, err)
	}

	var result []StreamDetails
	resp, err := d.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListStreams, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListStreams, d.Error(resp))
	}

	return result, nil
}

func setConnectorTypes(configuration *StreamConfiguration) {
	for _, connector := range configuration.Connectors {
		connector.SetConnectorType()
	}
}
