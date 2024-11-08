package datastream

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// DetailedStreamVersion is returned from GetStream
	DetailedStreamVersion struct {
		ContractID            string                `json:"contractId"`
		CreatedBy             string                `json:"createdBy"`
		CreatedDate           string                `json:"createdDate"`
		CollectMidgress       bool                  `json:"collectMidgress,omitempty"`
		DatasetFields         []DataSetField        `json:"datasetFields"`
		DeliveryConfiguration DeliveryConfiguration `json:"deliveryConfiguration"`
		Destination           Destination           `json:"destination"`
		GroupID               int                   `json:"groupId,omitempty"`
		LatestVersion         int                   `json:"latestVersion"`
		ModifiedBy            string                `json:"modifiedBy"`
		ModifiedDate          string                `json:"modifiedDate"`
		NotificationEmails    []string              `json:"notificationEmails"`
		ProductID             string                `json:"productId"`
		Properties            []Property            `json:"properties"`
		StreamID              int64                 `json:"streamId"`
		StreamName            string                `json:"streamName"`
		StreamVersion         int64                 `json:"streamVersion"`
		StreamStatus          StreamStatus          `json:"streamStatus"`
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
		ContractID            string                `json:"contractId"`
		CollectMidgress       bool                  `json:"collectMidgress,omitempty"`
		DatasetFields         []DatasetFieldID      `json:"datasetFields"`
		Destination           AbstractConnector     `json:"destination"`
		DeliveryConfiguration DeliveryConfiguration `json:"deliveryConfiguration"`
		GroupID               int                   `json:"groupId,omitempty"`
		NotificationEmails    []string              `json:"notificationEmails,omitempty"`
		Properties            []PropertyID          `json:"properties"`
		StreamName            string                `json:"streamName"`
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
		PropertyID   int    `json:"propertyId"`
		PropertyName string `json:"propertyName"`
	}

	// PropertyID identifies property details required in the create stream request.
	PropertyID struct {
		PropertyID int `json:"propertyId"`
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
		Activate            bool
	}

	// StreamUpdate contains information about stream ID and version
	StreamUpdate struct {
		StreamID      int64 `json:"streamId"`
		StreamVersion int64 `json:"streamVersion"`
	}

	// DeleteStreamRequest is passed to DeleteStream
	DeleteStreamRequest struct {
		StreamID int64
	}

	// ListStreamsRequest is passed to ListStreams
	ListStreamsRequest struct {
		GroupID *int
	}

	// StreamDetails contains information about stream
	StreamDetails struct {
		ContractID    string       `json:"contractId"`
		CreatedBy     string       `json:"createdBy"`
		CreatedDate   string       `json:"createdDate"`
		GroupID       int          `json:"groupId"`
		LatestVersion int64        `json:"latestVersion"`
		ModifiedBy    string       `json:"modifiedBy"`
		ModifiedDate  string       `json:"modifiedDate"`
		Properties    []Property   `json:"properties"`
		ProductID     string       `json:"productId"`
		StreamID      int64        `json:"streamId"`
		StreamName    string       `json:"streamName"`
		StreamStatus  StreamStatus `json:"streamStatus"`
		StreamVersion int64        `json:"streamVersion"`
	}
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

// Validate validates CreateStreamRequest
func (r CreateStreamRequest) Validate() error {
	return validation.Errors{
		"StreamConfiguration.DeliveryConfiguration":                             validation.Validate(r.StreamConfiguration.DeliveryConfiguration, validation.Required),
		"StreamConfiguration.DeliveryConfiguration.Delimiter":                   validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Delimiter, validation.When(r.StreamConfiguration.DeliveryConfiguration.Format == FormatTypeStructured, validation.Required, validation.In(DelimiterTypeSpace)), validation.When(r.StreamConfiguration.DeliveryConfiguration.Format == FormatTypeJson, validation.Nil)),
		"StreamConfiguration.DeliveryConfiguration.Format":                      validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Format, validation.Required, validation.In(FormatTypeStructured, FormatTypeJson), validation.When(r.StreamConfiguration.DeliveryConfiguration.Delimiter != nil, validation.Required, validation.In(FormatTypeStructured))),
		"StreamConfiguration.DeliveryConfiguration.Frequency":                   validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Frequency, validation.Required),
		"StreamConfiguration.DeliveryConfiguration.Frequency.IntervalInSeconds": validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Frequency.IntervalInSeconds, validation.Required, validation.In(IntervalInSeconds30, IntervalInSeconds60)),
		"StreamConfiguration.Destination":                                       validation.Validate(r.StreamConfiguration.Destination, validation.Required),
		"StreamConfiguration.ContractId":                                        validation.Validate(r.StreamConfiguration.ContractID, validation.Required),
		"StreamConfiguration.DatasetFields":                                     validation.Validate(r.StreamConfiguration.DatasetFields, validation.Required),
		"StreamConfiguration.GroupID":                                           validation.Validate(r.StreamConfiguration.GroupID, validation.Required, validation.Min(1)),
		"StreamConfiguration.Properties":                                        validation.Validate(r.StreamConfiguration.Properties, validation.Required),
		"StreamConfiguration.StreamName":                                        validation.Validate(r.StreamConfiguration.StreamName, validation.Required),
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
		"StreamConfiguration.DeliveryConfiguration":                             validation.Validate(r.StreamConfiguration.DeliveryConfiguration, validation.Required),
		"StreamConfiguration.DeliveryConfiguration.Delimiter":                   validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Delimiter, validation.When(r.StreamConfiguration.DeliveryConfiguration.Format == FormatTypeStructured, validation.Required, validation.In(DelimiterTypeSpace)), validation.When(r.StreamConfiguration.DeliveryConfiguration.Format == FormatTypeJson, validation.Nil)),
		"StreamConfiguration.DeliveryConfiguration.Format":                      validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Format, validation.In(FormatTypeStructured, FormatTypeJson)),
		"StreamConfiguration.DeliveryConfiguration.Frequency":                   validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Frequency, validation.Required),
		"StreamConfiguration.DeliveryConfiguration.Frequency.IntervalInSeconds": validation.Validate(r.StreamConfiguration.DeliveryConfiguration.Frequency.IntervalInSeconds, validation.Required, validation.In(IntervalInSeconds30, IntervalInSeconds60)),
		"StreamConfiguration.Destination":                                       validation.Validate(r.StreamConfiguration.Destination, validation.Required),
		"StreamConfiguration.ContractId":                                        validation.Validate(r.StreamConfiguration.ContractID, validation.Required),
		"StreamConfiguration.DatasetFields":                                     validation.Validate(r.StreamConfiguration.DatasetFields, validation.Required),
		"StreamConfiguration.GroupID":                                           validation.Validate(r.StreamConfiguration.GroupID, validation.In(0)),
		"StreamConfiguration.Properties":                                        validation.Validate(r.StreamConfiguration.Properties, validation.Required),
		"StreamConfiguration.StreamName":                                        validation.Validate(r.StreamConfiguration.StreamName, validation.Required),
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

func (d *ds) CreateStream(ctx context.Context, params CreateStreamRequest) (*DetailedStreamVersion, error) {
	logger := d.Log(ctx)
	logger.Debug("CreateStream")

	setDestinationType(&params.StreamConfiguration)
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateStream, ErrStructValidation, err)
	}

	uri, err := url.Parse("/datastream-config-api/v2/log/streams")
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

	return &rval, nil
}

func (d *ds) GetStream(ctx context.Context, params GetStreamRequest) (*DetailedStreamVersion, error) {
	logger := d.Log(ctx)
	logger.Debug("GetStream")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetStream, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v2/log/streams/%d",
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
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetStream, d.Error(resp))
	}

	return &rval, nil
}

func (d *ds) UpdateStream(ctx context.Context, params UpdateStreamRequest) (*DetailedStreamVersion, error) {
	logger := d.Log(ctx)
	logger.Debug("UpdateStream")

	setDestinationType(&params.StreamConfiguration)
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateStream, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/datastream-config-api/v2/log/streams/%d", params.StreamID))

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

	return &rval, nil
}

func (d *ds) DeleteStream(ctx context.Context, params DeleteStreamRequest) error {
	logger := d.Log(ctx)
	logger.Debug("DeleteStream")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteStream, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf(
		"/datastream-config-api/v2/log/streams/%d",
		params.StreamID),
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

	uri, err := url.Parse("/datastream-config-api/v2/log/streams")
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

	return result, nil
}

func setDestinationType(configuration *StreamConfiguration) {
	configuration.Destination.SetDestinationType()
}
