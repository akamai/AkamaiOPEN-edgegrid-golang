package datastream

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// DetailedStreamVersion is returned from GetStream
	DetailedStreamVersion struct {
		ContractID            string                 `json:"contractId"`
		CreatedBy             string                 `json:"createdBy"`
		CreatedDate           string                 `json:"createdDate"`
		CollectMidgress       bool                   `json:"collectMidgress,omitempty"`
		DatasetFields         []DataSetField         `json:"datasetFields"`
		DeliveryConfiguration DeliveryConfiguration  `json:"deliveryConfiguration"`
		Destination           Destination            `json:"destination"`
		GroupID               int                    `json:"groupId,omitempty"`
		LatestVersion         int                    `json:"latestVersion"`
		ModifiedBy            string                 `json:"modifiedBy"`
		ModifiedDate          string                 `json:"modifiedDate"`
		NotificationEmails    []string               `json:"notificationEmails"`
		ProductID             string                 `json:"productId"`
		Properties            []Property             `json:"properties"`
		AppSecConfigs         []AppSecConfig         `json:"appSecConfigs"`
		AnswerXServiceIDs     []AnswerXServiceDetail `json:"serviceSubletterIds"`
		StreamID              int64                  `json:"streamId"`
		StreamName            string                 `json:"streamName"`
		StreamVersion         int64                  `json:"streamVersion"`
		StreamStatus          StreamStatus           `json:"streamStatus"`
		IntegrationType       string                 `json:"integrationType"`
		SamplingPercentage    int                    `json:"samplingPercentage"`
		LogType               LogType                `json:"-"`
	}

	// Destination provides detailed information about the destination’s configuration in the stream
	Destination struct {
		AuthenticationType AuthenticationType `json:"authenticationType"`
		CompressLogs       bool               `json:"compressLogs"`
		DestinationType    DestinationType    `json:"destinationType"`
		DisplayName        string             `json:"displayName"`
		Path               string             `json:"path"`
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
		ContractID            string                `json:"contractId,omitempty"`
		CollectMidgress       bool                  `json:"collectMidgress,omitempty"`
		DatasetFields         []DatasetFieldID      `json:"datasetFields,omitempty"`
		Destination           AbstractConnector     `json:"destination"`
		DeliveryConfiguration DeliveryConfiguration `json:"deliveryConfiguration"`
		GroupID               int                   `json:"groupId,omitempty"`
		NotificationEmails    []string              `json:"notificationEmails,omitempty"`
		Properties            []PropertyID          `json:"properties,omitempty"`
		AppSecConfigs         []AppSecConfigID      `json:"appSecConfigs,omitempty"`
		AnswerXServiceIDs     []AnswerXServiceID    `json:"serviceSubletterIds,omitempty"`
		StreamName            string                `json:"streamName"`
		SamplingPercentage    int                   `json:"samplingPercentage,omitempty"`
	}

	// DeliveryConfiguration of the configuration of log lines, names of the files sent to a destination, and delivery frequency for these files
	DeliveryConfiguration struct {
		Delimiter        *DelimiterType `json:"fieldDelimiter,omitempty"`
		Format           FormatType     `json:"format"`
		Frequency        Frequency      `json:"frequency"`
		UploadFilePrefix string         `json:"uploadFilePrefix,omitempty"`
		UploadFileSuffix string         `json:"uploadFileSuffix,omitempty"`
	}

	// The Frequency of collecting logs from each uploader and sending these logs to a destination.
	Frequency struct {
		IntervalInSeconds IntervalInSeconds `json:"intervalInSeconds"`
	}

	// DataSets is a list of fields selected from the associated template that the stream monitors in logs
	DataSets struct {
		DataSetFields []DataSetField `json:"datasetFields"`
	}

	// DataSetField is a data set field selected from the associated template that the stream monitors in logs
	DataSetField struct {
		DatasetFieldID          int    `json:"datasetFieldId"`
		DatasetFieldDescription string `json:"datasetFieldDescription"`
		DatasetFieldJsonKey     string `json:"datasetFieldJsonKey"`
		DatasetFieldName        string `json:"datasetFieldName"`
		DatasetFieldGroup       string `json:"datasetFieldGroup"`
	}

	// DatasetFieldID is a dataset field value used in create stream request
	DatasetFieldID struct {
		DatasetFieldID int `json:"datasetFieldId"`
	}

	// Property identifies brief info about the properties monitored in the stream.
	Property struct {
		PropertyID      int    `json:"propertyId"`
		PropertyName    string `json:"propertyName"`
		IntegrationType string `json:"integrationType"`
	}

	// AppSecConfig holds the AppSec configuration associated with a stream.
	AppSecConfig struct {
		AppSecID   int    `json:"appSecId"`
		AppSecName string `json:"appSecName"`
	}

	// PropertyID identifies property details required in the create stream request.
	PropertyID struct {
		PropertyID int `json:"propertyId"`
	}

	// AppSecConfigID identifies AppSec config details required in the create stream request.
	AppSecConfigID struct {
		AppSecID int `json:"appSecId"`
	}

	// AnswerXServiceDetail contains information about an AnswerX service ID available for streaming.
	AnswerXServiceDetail struct {
		// SSID is the AnswerX service ID.
		SSID int64 `json:"ssid"`

		// Name is the display name of the AnswerX service.
		Name string `json:"name"`

		// Product is the product associated with the AnswerX service.
		Product string `json:"product"`
	}

	// AnswerXServiceID identifies an AnswerX service ID required in the create/update stream request.
	AnswerXServiceID struct {
		// SSID is the AnswerX service identifier to associate with the stream.
		SSID int64 `json:"ssid"`
	}

	// StreamStatus is used to create an enum of possible StreamStatus values
	StreamStatus string

	// AbstractConnector is an interface for all Connector types
	AbstractConnector interface {
		SetDestinationType()
		Validate() error
	}

	// DelimiterType enum
	DelimiterType string

	// FormatType enum
	FormatType string

	// IntervalInSeconds enum
	IntervalInSeconds int

	// CreateStreamRequest is passed to CreateStream
	CreateStreamRequest struct {
		StreamConfiguration StreamConfiguration
		Activate            bool
		LogType             LogType
	}

	// GetStreamRequest is passed to GetStream
	GetStreamRequest struct {
		StreamID int64
		Version  *int64
		LogType  LogType
	}

	// UpdateStreamRequest is passed to UpdateStream
	UpdateStreamRequest struct {
		StreamID            int64
		StreamConfiguration StreamConfiguration
		Activate            bool
		LogType             LogType
	}

	// StreamUpdate contains information about stream ID and version
	StreamUpdate struct {
		StreamID      int64 `json:"streamId"`
		StreamVersion int64 `json:"streamVersion"`
	}

	// DeleteStreamRequest is passed to DeleteStream
	DeleteStreamRequest struct {
		StreamID int64
		LogType  LogType
	}

	// ListStreamsRequest is passed to ListStreams
	ListStreamsRequest struct {
		GroupID *int
		LogType LogType
	}

	// StreamDetails contains information about stream
	StreamDetails struct {
		ContractID         string                 `json:"contractId"`
		CreatedBy          string                 `json:"createdBy"`
		CreatedDate        string                 `json:"createdDate"`
		GroupID            int                    `json:"groupId"`
		LatestVersion      int64                  `json:"latestVersion"`
		ModifiedBy         string                 `json:"modifiedBy"`
		ModifiedDate       string                 `json:"modifiedDate"`
		Properties         []Property             `json:"properties"`
		ProductID          string                 `json:"productId"`
		StreamID           int64                  `json:"streamId"`
		StreamName         string                 `json:"streamName"`
		StreamStatus       StreamStatus           `json:"streamStatus"`
		StreamVersion      int64                  `json:"streamVersion"`
		IntegrationType    string                 `json:"integrationType"`
		SamplingPercentage int                    `json:"samplingPercentage"`
		LogType            LogType                `json:"-"`
		AppSecConfigs      []AppSecConfig         `json:"appSecConfigs"`
		AnswerXServiceIDs  []AnswerXServiceDetail `json:"serviceSubletterIds"`
	}

	// LogType enumeration
	LogType string
)

const (
	// LogTypeCDN is the log type for CDN streams.
	LogTypeCDN LogType = "CDN"
	// LogTypeAppSec is the log type for AppSec streams.
	LogTypeAppSec LogType = "APPSEC"
	// LogTypeAnswerX is the log type for ANSWERX streams.
	LogTypeAnswerX LogType = "ANSWERX"
)

const (
	// StreamStatusActivated const
	StreamStatusActivated StreamStatus = "ACTIVATED"
	// StreamStatusDeactivated const
	StreamStatusDeactivated StreamStatus = "DEACTIVATED"
	// StreamStatusActivating const
	StreamStatusActivating StreamStatus = "ACTIVATING"
	// StreamStatusDeactivating const state
	StreamStatusDeactivating StreamStatus = "DEACTIVATING"
	// StreamStatusInactive const
	StreamStatusInactive StreamStatus = "INACTIVE"

	// DelimiterTypeSpace const
	DelimiterTypeSpace DelimiterType = "SPACE"

	// FormatTypeStructured const
	FormatTypeStructured FormatType = "STRUCTURED"
	// FormatTypeJson const
	FormatTypeJson FormatType = "JSON"

	// IntervalInSeconds30 const
	IntervalInSeconds30 IntervalInSeconds = 30
	// IntervalInSeconds60 const
	IntervalInSeconds60 IntervalInSeconds = 60
)

// ToPathValue converts the LogType enum to the lowercase string used in URL paths.
func (logType LogType) ToPathValue() string {
	return strings.ToLower(string(logType))
}

// Validate validates ListStreamsRequest.
func (r ListStreamsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		// check that the LogType field is not empty and is a valid value.
		"LogType": validation.Validate(r.LogType, validation.Required, validation.In(LogTypeCDN, LogTypeAppSec, LogTypeAnswerX)),
	})
}

// Validate validates CreateStreamRequest
func (r CreateStreamRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"StreamConfiguration.DeliveryConfiguration":                             validation.Validate(r.StreamConfiguration.DeliveryConfiguration, validation.Required),
		"StreamConfiguration.DeliveryConfiguration.Delimiter":                   validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Delimiter, validation.When(r.StreamConfiguration.DeliveryConfiguration.Format == FormatTypeStructured, validation.Required, validation.In(DelimiterTypeSpace)), validation.When(r.StreamConfiguration.DeliveryConfiguration.Format == FormatTypeJson, validation.Nil)),
		"StreamConfiguration.DeliveryConfiguration.Format":                      validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Format, validation.Required, validation.In(FormatTypeStructured, FormatTypeJson), validation.When(r.StreamConfiguration.DeliveryConfiguration.Delimiter != nil, validation.Required, validation.In(FormatTypeStructured))),
		"StreamConfiguration.DeliveryConfiguration.Frequency":                   validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Frequency, validation.Required),
		"StreamConfiguration.DeliveryConfiguration.Frequency.IntervalInSeconds": validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Frequency.IntervalInSeconds, validation.Required, validation.In(IntervalInSeconds30, IntervalInSeconds60)),
		"StreamConfiguration.Destination":                                       validation.Validate(r.StreamConfiguration.Destination, validation.Required),
		"StreamConfiguration.ContractId":                                        validation.Validate(r.StreamConfiguration.ContractID, validation.When(r.LogType != LogTypeCDN, validation.Required)),
		"StreamConfiguration.DatasetFields":                                     validation.Validate(r.StreamConfiguration.DatasetFields, validation.When(r.LogType == LogTypeCDN, validation.Required), validation.When(r.LogType == LogTypeAppSec, validation.Empty), validation.When(r.LogType == LogTypeAnswerX, validation.Required)),
		"StreamConfiguration.GroupID":                                           validation.Validate(r.StreamConfiguration.GroupID, validation.When(r.LogType == LogTypeCDN, validation.When(r.StreamConfiguration.GroupID != 0, validation.Min(1))), validation.When(r.LogType != LogTypeCDN, validation.Required, validation.Min(1))),
		"StreamConfiguration.Properties":                                        validation.Validate(r.StreamConfiguration.Properties, validation.When(r.LogType == LogTypeCDN, validation.Required)),
		"StreamConfiguration.AppSecConfigs":                                     validation.Validate(r.StreamConfiguration.AppSecConfigs, validation.When(r.LogType == LogTypeAppSec, validation.Required)),
		"StreamConfiguration.AnswerXServiceIDs":                                 validation.Validate(r.StreamConfiguration.AnswerXServiceIDs, validation.When(r.LogType == LogTypeAnswerX, validation.Required)),
		"StreamConfiguration.StreamName":                                        validation.Validate(r.StreamConfiguration.StreamName, validation.Required),
		"StreamConfiguration.SamplingPercentage":                                validation.Validate(r.StreamConfiguration.SamplingPercentage, validation.When(r.StreamConfiguration.SamplingPercentage != 0, validation.Min(1), validation.Max(100))),
		"LogType":                                                               validation.Validate(r.LogType, validation.Required, validation.In(LogTypeCDN, LogTypeAppSec, LogTypeAnswerX)),
	})
}

// Validate validates GetStreamRequest
func (r GetStreamRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"streamId": validation.Validate(r.StreamID, validation.Required),
		"LogType":  validation.Validate(r.LogType, validation.Required, validation.In(LogTypeCDN, LogTypeAppSec, LogTypeAnswerX)),
	})
}

// Validate validates UpdateStreamRequest
func (r UpdateStreamRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"StreamConfiguration.DeliveryConfiguration":                             validation.Validate(r.StreamConfiguration.DeliveryConfiguration, validation.Required),
		"StreamConfiguration.DeliveryConfiguration.Delimiter":                   validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Delimiter, validation.When(r.StreamConfiguration.DeliveryConfiguration.Format == FormatTypeStructured, validation.Required, validation.In(DelimiterTypeSpace)), validation.When(r.StreamConfiguration.DeliveryConfiguration.Format == FormatTypeJson, validation.Nil)),
		"StreamConfiguration.DeliveryConfiguration.Format":                      validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Format, validation.In(FormatTypeStructured, FormatTypeJson)),
		"StreamConfiguration.DeliveryConfiguration.Frequency":                   validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Frequency, validation.Required),
		"StreamConfiguration.DeliveryConfiguration.Frequency.IntervalInSeconds": validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Frequency.IntervalInSeconds, validation.Required, validation.In(IntervalInSeconds30, IntervalInSeconds60)),
		"StreamConfiguration.Destination":                                       validation.Validate(r.StreamConfiguration.Destination, validation.Required),
		"StreamConfiguration.ContractId":                                        validation.Validate(r.StreamConfiguration.ContractID, validation.When(r.LogType != LogTypeCDN, validation.Required)),
		"StreamConfiguration.DatasetFields":                                     validation.Validate(r.StreamConfiguration.DatasetFields, validation.When(r.LogType == LogTypeCDN, validation.Required), validation.When(r.LogType == LogTypeAppSec, validation.Empty), validation.When(r.LogType == LogTypeAnswerX, validation.Required)),
		"StreamConfiguration.GroupID":                                           validation.Validate(r.StreamConfiguration.GroupID, validation.When(r.LogType == LogTypeCDN, validation.When(r.StreamConfiguration.GroupID != 0, validation.Min(1))), validation.When(r.LogType != LogTypeCDN, validation.In(0))),
		"StreamConfiguration.Properties":                                        validation.Validate(r.StreamConfiguration.Properties, validation.When(r.LogType == LogTypeCDN, validation.Required)),
		"StreamConfiguration.AppSecConfigs":                                     validation.Validate(r.StreamConfiguration.AppSecConfigs, validation.When(r.LogType == LogTypeAppSec, validation.Required)),
		"StreamConfiguration.AnswerXServiceIDs":                                 validation.Validate(r.StreamConfiguration.AnswerXServiceIDs, validation.When(r.LogType == LogTypeAnswerX, validation.Required)),
		"StreamConfiguration.StreamName":                                        validation.Validate(r.StreamConfiguration.StreamName, validation.Required),
		"StreamConfiguration.SamplingPercentage":                                validation.Validate(r.StreamConfiguration.SamplingPercentage, validation.When(r.StreamConfiguration.SamplingPercentage != 0, validation.Min(1), validation.Max(100))),
		"LogType":                                                               validation.Validate(r.LogType, validation.Required, validation.In(LogTypeCDN, LogTypeAppSec, LogTypeAnswerX)),
	})
}

// Validate validates DeleteStreamRequest
func (r DeleteStreamRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"streamId": validation.Validate(r.StreamID, validation.Required),
		"LogType":  validation.Validate(r.LogType, validation.Required, validation.In(LogTypeCDN, LogTypeAppSec, LogTypeAnswerX)),
	})
}

func (d *ds) CreateStream(ctx context.Context, params CreateStreamRequest) (*DetailedStreamVersion, error) {
	logger := d.Log(ctx)
	logger.Debug("CreateStream")

	setDestinationType(&params.StreamConfiguration)
	params.StreamConfiguration.ContractID = strings.TrimSpace(params.StreamConfiguration.ContractID)
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateStream, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/datastream-config-api/v3/log/%s/streams", params.LogType.ToPathValue()))
	if err != nil {
		return nil, fmt.Errorf("%w: parsing URL: %s", ErrCreateStream, err)
	}

	q := uri.Query()
	q.Add("activate", fmt.Sprintf("%t", params.Activate))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateStream, err)
	}

	var rval DetailedStreamVersion
	resp, err := d.Exec(req, &rval, params.StreamConfiguration)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateStream, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrCreateStream, d.Error(resp))
	}

	// echo the log type back in the response since the API does not return it.
	rval.LogType = params.LogType

	return &rval, nil
}

func (d *ds) GetStream(ctx context.Context, params GetStreamRequest) (*DetailedStreamVersion, error) {
	logger := d.Log(ctx)
	logger.Debug("GetStream")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetStream, ErrStructValidation, err)
	}

	path := fmt.Sprintf("/datastream-config-api/v3/log/%s/streams/%d", params.LogType.ToPathValue(), params.StreamID)
	uri, err := url.Parse(path)
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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetStream, d.Error(resp))
	}

	// propagate log type to stream details since GetStreamRequest requires log type and the API does not return it
	rval.LogType = params.LogType

	return &rval, nil
}

func (d *ds) UpdateStream(ctx context.Context, params UpdateStreamRequest) (*DetailedStreamVersion, error) {
	logger := d.Log(ctx)
	logger.Debug("UpdateStream")

	setDestinationType(&params.StreamConfiguration)
	params.StreamConfiguration.ContractID = strings.TrimSpace(params.StreamConfiguration.ContractID)
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateStream, ErrStructValidation, err)
	}

	path := fmt.Sprintf("/datastream-config-api/v3/log/%s/streams/%d", params.LogType.ToPathValue(), params.StreamID)
	uri, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrUpdateStream, err)
	}

	q := uri.Query()
	q.Add("activate", fmt.Sprintf("%t", params.Activate))
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateStream, err)
	}

	var rval DetailedStreamVersion
	resp, err := d.Exec(req, &rval, params.StreamConfiguration)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateStream, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateStream, d.Error(resp))
	}

	// echo the log type back in the response since the API does not return it.
	rval.LogType = params.LogType

	return &rval, nil
}

func (d *ds) DeleteStream(ctx context.Context, params DeleteStreamRequest) error {
	logger := d.Log(ctx)
	logger.Debug("DeleteStream")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteStream, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v3/log/%s/streams/%d",
		params.LogType.ToPathValue(), params.StreamID),
	)
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrDeleteStream, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteStream, err)
	}

	resp, err := d.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteStream, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s: %w", ErrDeleteStream, d.Error(resp))
	}

	return nil
}

func (d *ds) ListStreams(ctx context.Context, params ListStreamsRequest) ([]StreamDetails, error) {
	logger := d.Log(ctx)
	logger.Debug("ListStreams")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListStreams, ErrStructValidation, err)
	}

	path := fmt.Sprintf("/datastream-config-api/v3/log/%s/streams", params.LogType.ToPathValue())
	uri, err := url.Parse(path)
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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListStreams, d.Error(resp))
	}

	// propagate log type to stream details since ListStreamsRequest requires log type and the API does not return it
	for i := range result {
		result[i].LogType = params.LogType
	}

	return result, nil
}

func setDestinationType(configuration *StreamConfiguration) {
	configuration.Destination.SetDestinationType()
}
