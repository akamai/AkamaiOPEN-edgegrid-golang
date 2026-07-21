package datastream

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// S3Connector provides details about the Amazon S3 destination in a stream
	S3Connector struct {
		DestinationType DestinationType `json:"destinationType"`
		AccessKey       string          `json:"accessKey"`
		Bucket          string          `json:"bucket"`
		DisplayName     string          `json:"displayName"`
		Path            string          `json:"path"`
		Region          string          `json:"region"`
		SecretAccessKey string          `json:"secretAccessKey"`
	}

	// AzureConnector provides details about the Azure Storage destination configuration in a data stream
	AzureConnector struct {
		DestinationType DestinationType `json:"destinationType"`
		AccessKey       string          `json:"accessKey"`
		AccountName     string          `json:"accountName"`
		DisplayName     string          `json:"displayName"`
		ContainerName   string          `json:"containerName"`
		Path            string          `json:"path"`
	}

	// DatadogConnector provides detailed information about Datadog destination
	DatadogConnector struct {
		DestinationType DestinationType `json:"destinationType"`
		AuthToken       string          `json:"authToken"`
		CompressLogs    bool            `json:"compressLogs"`
		DisplayName     string          `json:"displayName"`
		Service         string          `json:"service,omitempty"`
		Source          string          `json:"source,omitempty"`
		Tags            string          `json:"tags,omitempty"`
		Endpoint        string          `json:"endpoint"`
	}

	// SplunkConnector provides detailed information about the Splunk destination
	SplunkConnector struct {
		DestinationType     DestinationType `json:"destinationType"`
		CompressLogs        bool            `json:"compressLogs"`
		DisplayName         string          `json:"displayName"`
		EventCollectorToken string          `json:"eventCollectorToken"`
		Endpoint            string          `json:"endpoint"`
		CustomHeaderName    string          `json:"customHeaderName,omitempty"`
		CustomHeaderValue   string          `json:"customHeaderValue,omitempty"`
		TLSHostname         string          `json:"tlsHostname,omitempty"`
		CACert              string          `json:"caCert,omitempty"`
		ClientCert          string          `json:"clientCert,omitempty"`
		ClientKey           string          `json:"clientKey,omitempty"`
	}

	// GCSConnector provides detailed information about the Google Cloud Storage destination
	GCSConnector struct {
		DestinationType    DestinationType `json:"destinationType"`
		Bucket             string          `json:"bucket"`
		DisplayName        string          `json:"displayName"`
		Path               string          `json:"path,omitempty"`
		PrivateKey         string          `json:"privateKey"`
		ProjectID          string          `json:"projectId"`
		ServiceAccountName string          `json:"serviceAccountName"`
	}

	// CustomHTTPSConnector provides detailed information about the custom HTTPS endpoint
	CustomHTTPSConnector struct {
		DestinationType    DestinationType    `json:"destinationType"`
		AuthenticationType AuthenticationType `json:"authenticationType"`
		CompressLogs       bool               `json:"compressLogs"`
		DisplayName        string             `json:"displayName"`
		Password           string             `json:"password,omitempty"`
		Endpoint           string             `json:"endpoint"`
		UserName           string             `json:"userName,omitempty"`
		ContentType        string             `json:"contentType,omitempty"`
		CustomHeaderName   string             `json:"customHeaderName,omitempty"`
		CustomHeaderValue  string             `json:"customHeaderValue,omitempty"`
		TLSHostname        string             `json:"tlsHostname,omitempty"`
		CACert             string             `json:"caCert,omitempty"`
		ClientCert         string             `json:"clientCert,omitempty"`
		ClientKey          string             `json:"clientKey,omitempty"`
	}

	// SumoLogicConnector provides detailed information about the Sumo Logic destination
	SumoLogicConnector struct {
		DestinationType   DestinationType `json:"destinationType"`
		CollectorCode     string          `json:"collectorCode"`
		CompressLogs      bool            `json:"compressLogs"`
		DisplayName       string          `json:"displayName"`
		Endpoint          string          `json:"endpoint"`
		ContentType       string          `json:"contentType,omitempty"`
		CustomHeaderName  string          `json:"customHeaderName,omitempty"`
		CustomHeaderValue string          `json:"customHeaderValue,omitempty"`
	}

	// OracleCloudStorageConnector provides details about the Oracle Cloud Storage destination
	OracleCloudStorageConnector struct {
		DestinationType DestinationType `json:"destinationType"`
		AccessKey       string          `json:"accessKey"`
		Bucket          string          `json:"bucket"`
		DisplayName     string          `json:"displayName"`
		Namespace       string          `json:"namespace"`
		Path            string          `json:"path"`
		Region          string          `json:"region"`
		SecretAccessKey string          `json:"secretAccessKey"`
	}

	// LogglyConnector contains details about Loggly destination.
	LogglyConnector struct {
		DestinationType   DestinationType `json:"destinationType"`
		DisplayName       string          `json:"displayName"`
		Endpoint          string          `json:"endpoint"`
		AuthToken         string          `json:"authToken"`
		Tags              string          `json:"tags,omitempty"`
		ContentType       string          `json:"contentType,omitempty"`
		CustomHeaderName  string          `json:"customHeaderName,omitempty"`
		CustomHeaderValue string          `json:"customHeaderValue,omitempty"`
	}

	// NewRelicConnector contains details about New Relic destination.
	NewRelicConnector struct {
		DestinationType   DestinationType `json:"destinationType"`
		DisplayName       string          `json:"displayName"`
		Endpoint          string          `json:"endpoint"`
		AuthToken         string          `json:"authToken"`
		ContentType       string          `json:"contentType,omitempty"`
		CustomHeaderName  string          `json:"customHeaderName,omitempty"`
		CustomHeaderValue string          `json:"customHeaderValue,omitempty"`
	}

	// ElasticsearchConnector contains details about Elasticsearch destination.
	ElasticsearchConnector struct {
		DestinationType   DestinationType `json:"destinationType"`
		DisplayName       string          `json:"displayName"`
		Endpoint          string          `json:"endpoint"`
		IndexName         string          `json:"indexName"`
		UserName          string          `json:"userName"`
		Password          string          `json:"password"`
		ContentType       string          `json:"contentType,omitempty"`
		CustomHeaderName  string          `json:"customHeaderName,omitempty"`
		CustomHeaderValue string          `json:"customHeaderValue,omitempty"`
		TLSHostname       string          `json:"tlsHostname,omitempty"`
		CACert            string          `json:"caCert,omitempty"`
		ClientCert        string          `json:"clientCert,omitempty"`
		ClientKey         string          `json:"clientKey,omitempty"`
	}

	// S3CompatibleConnector provides details about the S3 compatible destination in a stream
	S3CompatibleConnector struct {
		// DestinationType is the destination type's name. Set it to S3_COMPATIBLE for this destination type.
		DestinationType DestinationType `json:"destinationType"`
		// AccessKey is the access key for the destination.
		AccessKey string `json:"accessKey"`
		// Bucket is the bucket name for the destination.
		Bucket string `json:"bucket"`
		// DisplayName is the display name of the destination.
		DisplayName string `json:"displayName"`
		// Path is the path within the bucket where logs will be stored. It is optional.
		Path string `json:"path,omitempty"`
		// Region is the region where the bucket is located.
		Region string `json:"region"`
		// SecretAccessKey is the secret access key for the destination.
		SecretAccessKey string `json:"secretAccessKey"`
		// Endpoint is the endpoint URL of the destination.
		Endpoint string `json:"endpoint"`
	}

	// TrafficPeakConnector provides detailed information about the TrafficPeak endpoint
	TrafficPeakConnector struct {
		// DestinationType is the destination type's name. Set it to TRAFFICPEAK for this destination type.
		DestinationType DestinationType `json:"destinationType"`
		// AuthenticationType is the authentication type for the destination. Set it to BASIC.
		AuthenticationType AuthenticationType `json:"authenticationType"`
		// CompressLogs indicates whether to compress logs before sending them to the destination.
		CompressLogs bool `json:"compressLogs"`
		// DisplayName is the display name of the destination.
		DisplayName string `json:"displayName"`
		// Password is the password for the destination.
		Password string `json:"password"`
		// Endpoint is the endpoint URL of the destination.
		Endpoint string `json:"endpoint"`
		// UserName is the user name for the destination.
		UserName string `json:"userName"`
		// ContentType is the content type for the destination. Set it to application/json or application/json; charset=utf-8.
		ContentType TrafficPeakContentType `json:"contentType"`
		// CustomHeaderName is the custom header name for the destination. It is optional.
		CustomHeaderName string `json:"customHeaderName,omitempty"`
		// CustomHeaderValue is the custom header value for the destination. It is optional.
		CustomHeaderValue string `json:"customHeaderValue,omitempty"`
	}

	// DynatraceConnector contains details about Dynatrace destination.
	DynatraceConnector struct {
		// DestinationType is the destination type's name. Set it to DYNATRACE for this destination type.
		DestinationType DestinationType `json:"destinationType"`
		// DisplayName is the display name of the destination.
		DisplayName string `json:"displayName"`
		// Endpoint is the endpoint URL of the destination.
		Endpoint string `json:"endpoint"`
		// AuthToken is the authentication token for the destination.
		AuthToken string `json:"authToken"`
		// CustomHeaderName is the custom header name for the destination. It is optional.
		CustomHeaderName string `json:"customHeaderName,omitempty"`
		// CustomHeaderValue is the custom header value for the destination. It is optional.
		CustomHeaderValue string `json:"customHeaderValue,omitempty"`
	}

	// DestinationType is used to create an "enum" of possible DestinationTypes
	DestinationType string

	// AuthenticationType is used to create an "enum" of possible AuthenticationTypes of the CustomHTTPSConnector
	AuthenticationType string

	// TrafficPeakContentType is used to create an "enum" of possible Content types of the TrafficPeakConnector
	TrafficPeakContentType string
)

const (
	// DestinationTypeAzure const
	DestinationTypeAzure DestinationType = "AZURE"
	// DestinationTypeS3 const
	DestinationTypeS3 DestinationType = "S3"
	// DestinationTypeDataDog const
	DestinationTypeDataDog DestinationType = "DATADOG"
	// DestinationTypeSplunk const
	DestinationTypeSplunk DestinationType = "SPLUNK"
	// DestinationTypeGcs const
	DestinationTypeGcs DestinationType = "GCS"
	// DestinationTypeHTTPS const
	DestinationTypeHTTPS DestinationType = "HTTPS"
	// DestinationTypeSumoLogic const
	DestinationTypeSumoLogic DestinationType = "SUMO_LOGIC"
	// DestinationTypeOracle const
	DestinationTypeOracle DestinationType = "Oracle_Cloud_Storage"
	// DestinationTypeLoggly const
	DestinationTypeLoggly DestinationType = "LOGGLY"
	// DestinationTypeNewRelic const
	DestinationTypeNewRelic DestinationType = "NEWRELIC"
	// DestinationTypeElasticsearch const
	DestinationTypeElasticsearch DestinationType = "ELASTICSEARCH"
	// DestinationTypeS3Compatible const
	DestinationTypeS3Compatible DestinationType = "S3_COMPATIBLE"
	// DestinationTypeTrafficPeak const
	DestinationTypeTrafficPeak DestinationType = "TRAFFICPEAK"
	// DestinationTypeDynatrace const
	DestinationTypeDynatrace DestinationType = "DYNATRACE"

	// AuthenticationTypeNone const
	AuthenticationTypeNone AuthenticationType = "NONE"
	// AuthenticationTypeBasic const
	AuthenticationTypeBasic AuthenticationType = "BASIC"

	// TrafficPeakContentTypeJSON const
	TrafficPeakContentTypeJSON TrafficPeakContentType = "application/json"
	// TrafficPeakContentTypeJSONUTF8 const
	TrafficPeakContentTypeJSONUTF8 TrafficPeakContentType = "application/json; charset=utf-8"
)

var customHeaderNameRegexp = regexp.MustCompile("^[A-Za-z0-9_-]+$")

// SetDestinationType for S3Connector
func (c *S3Connector) SetDestinationType() {
	c.DestinationType = DestinationTypeS3
}

// Validate validates S3Connector
func (c *S3Connector) Validate() error {
	return validation.Errors{
		"DestinationType": validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeS3)),
		"AccessKey":       validation.Validate(c.AccessKey, validation.Required),
		"Bucket":          validation.Validate(c.Bucket, validation.Required),
		"DisplayName":     validation.Validate(c.DisplayName, validation.Required),
		"Path":            validation.Validate(c.Path, validation.Required),
		"Region":          validation.Validate(c.Region, validation.Required),
		"SecretAccessKey": validation.Validate(c.SecretAccessKey, validation.Required),
	}.Filter()
}

// SetDestinationType for AzureConnector
func (c *AzureConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeAzure
}

// Validate validates AzureConnector
func (c *AzureConnector) Validate() error {
	return validation.Errors{
		"DestinationType": validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeAzure)),
		"AccessKey":       validation.Validate(c.AccessKey, validation.Required),
		"AccountName":     validation.Validate(c.AccountName, validation.Required),
		"DisplayName":     validation.Validate(c.DisplayName, validation.Required),
		"ContainerName":   validation.Validate(c.ContainerName, validation.Required),
		"Path":            validation.Validate(c.Path, validation.Required),
	}.Filter()
}

// SetDestinationType for DatadogConnector
func (c *DatadogConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeDataDog
}

// Validate validates DatadogConnector
func (c *DatadogConnector) Validate() error {
	return validation.Errors{
		"DestinationType": validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeDataDog)),
		"AuthToken":       validation.Validate(c.AuthToken, validation.Required),
		"DisplayName":     validation.Validate(c.DisplayName, validation.Required),
		"Endpoint":        validation.Validate(c.Endpoint, validation.Required),
	}.Filter()
}

// SetDestinationType for SplunkConnector
func (c *SplunkConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeSplunk
}

// Validate validates SplunkConnector
func (c *SplunkConnector) Validate() error {
	return validation.Errors{
		"DestinationType":     validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeSplunk)),
		"DisplayName":         validation.Validate(c.DisplayName, validation.Required),
		"EventCollectorToken": validation.Validate(c.EventCollectorToken, validation.Required),
		"Endpoint":            validation.Validate(c.Endpoint, validation.Required),
		"CustomHeaderName":    validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue":   validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}

// SetDestinationType for GCSConnector
func (c *GCSConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeGcs
}

// Validate validates GCSConnector
func (c *GCSConnector) Validate() error {
	return validation.Errors{
		"DestinationType":    validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeGcs)),
		"Bucket":             validation.Validate(c.Bucket, validation.Required),
		"DisplayName":        validation.Validate(c.DisplayName, validation.Required),
		"PrivateKey":         validation.Validate(c.PrivateKey, validation.Required),
		"ProjectID":          validation.Validate(c.ProjectID, validation.Required),
		"ServiceAccountName": validation.Validate(c.ServiceAccountName, validation.Required),
	}.Filter()
}

// SetDestinationType for CustomHTTPSConnector
func (c *CustomHTTPSConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeHTTPS
}

// Validate validates CustomHTTPSConnector
func (c *CustomHTTPSConnector) Validate() error {
	return validation.Errors{
		"DestinationType":    validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeHTTPS)),
		"AuthenticationType": validation.Validate(c.AuthenticationType, validation.Required, validation.In(AuthenticationTypeBasic, AuthenticationTypeNone)),
		"DisplayName":        validation.Validate(c.DisplayName, validation.Required),
		"Endpoint":           validation.Validate(c.Endpoint, validation.Required),
		"UserName":           validation.Validate(c.UserName, validation.Required.When(c.AuthenticationType == AuthenticationTypeBasic)),
		"Password":           validation.Validate(c.Password, validation.Required.When(c.AuthenticationType == AuthenticationTypeBasic)),
		"CustomHeaderName":   validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue":  validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}

// SetDestinationType for SumoLogicConnector
func (c *SumoLogicConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeSumoLogic
}

// Validate validates SumoLogicConnector
func (c *SumoLogicConnector) Validate() error {
	return validation.Errors{
		"DestinationType":   validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeSumoLogic)),
		"CollectorCode":     validation.Validate(c.CollectorCode, validation.Required),
		"DisplayName":       validation.Validate(c.DisplayName, validation.Required),
		"Endpoint":          validation.Validate(c.Endpoint, validation.Required),
		"CustomHeaderName":  validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue": validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}

// SetDestinationType for OracleCloudStorageConnector
func (c *OracleCloudStorageConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeOracle
}

// Validate validates OracleCloudStorageConnector
func (c *OracleCloudStorageConnector) Validate() error {
	return validation.Errors{
		"DestinationType": validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeOracle)),
		"AccessKey":       validation.Validate(c.AccessKey, validation.Required),
		"Bucket":          validation.Validate(c.Bucket, validation.Required),
		"DisplayName":     validation.Validate(c.DisplayName, validation.Required),
		"Namespace":       validation.Validate(c.Namespace, validation.Required),
		"Path":            validation.Validate(c.Path, validation.Required),
		"Region":          validation.Validate(c.Region, validation.Required),
		"SecretAccessKey": validation.Validate(c.SecretAccessKey, validation.Required),
	}.Filter()
}

// SetDestinationType for LogglyConnector
func (c *LogglyConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeLoggly
}

// Validate validates LogglyConnector
func (c *LogglyConnector) Validate() error {
	return validation.Errors{
		"DestinationType":   validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeLoggly)),
		"DisplayName":       validation.Validate(c.DisplayName, validation.Required),
		"Endpoint":          validation.Validate(c.Endpoint, validation.Required),
		"AuthToken":         validation.Validate(c.AuthToken, validation.Required),
		"CustomHeaderName":  validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue": validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}

// SetDestinationType for NewRelicConnector
func (c *NewRelicConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeNewRelic
}

// Validate validates NewRelicConnector
func (c *NewRelicConnector) Validate() error {
	return validation.Errors{
		"DestinationType":   validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeNewRelic)),
		"DisplayName":       validation.Validate(c.DisplayName, validation.Required),
		"Endpoint":          validation.Validate(c.Endpoint, validation.Required),
		"AuthToken":         validation.Validate(c.AuthToken, validation.Required),
		"CustomHeaderName":  validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue": validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}

// SetDestinationType for ElasticsearchConnector
func (c *ElasticsearchConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeElasticsearch
}

// Validate validates ElasticsearchConnector
func (c *ElasticsearchConnector) Validate() error {
	return validation.Errors{
		"DestinationType":   validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeElasticsearch)),
		"DisplayName":       validation.Validate(c.DisplayName, validation.Required),
		"Endpoint":          validation.Validate(c.Endpoint, validation.Required),
		"UserName":          validation.Validate(c.UserName, validation.Required),
		"Password":          validation.Validate(c.Password, validation.Required),
		"IndexName":         validation.Validate(c.IndexName, validation.Required),
		"CustomHeaderName":  validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue": validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}

// SetDestinationType for S3CompatibleConnector
func (c *S3CompatibleConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeS3Compatible
}

// Validate validates S3CompatibleConnector
func (c *S3CompatibleConnector) Validate() error {
	return validation.Errors{
		"DestinationType": validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeS3Compatible)),
		"AccessKey":       validation.Validate(c.AccessKey, validation.Required),
		"Bucket":          validation.Validate(c.Bucket, validation.Required),
		"DisplayName":     validation.Validate(c.DisplayName, validation.Required),
		"Endpoint":        validation.Validate(c.Endpoint, validation.Required),
		"Region":          validation.Validate(c.Region, validation.Required),
		"SecretAccessKey": validation.Validate(c.SecretAccessKey, validation.Required),
	}.Filter()
}

// SetDestinationType for TrafficPeakConnector
func (c *TrafficPeakConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeTrafficPeak
}

// Validate validates TrafficPeakConnector
func (c *TrafficPeakConnector) Validate() error {
	return validation.Errors{
		"DestinationType":    validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeTrafficPeak)),
		"AuthenticationType": validation.Validate(c.AuthenticationType, validation.Required, validation.In(AuthenticationTypeBasic)),
		"DisplayName":        validation.Validate(c.DisplayName, validation.Required),
		"Endpoint":           validation.Validate(c.Endpoint, validation.Required),
		"UserName":           validation.Validate(c.UserName, validation.Required),
		"Password":           validation.Validate(c.Password, validation.Required),
		"CustomHeaderName":   validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue":  validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
		"ContentType":        validation.Validate(c.ContentType, validation.Required, validation.In(TrafficPeakContentTypeJSON, TrafficPeakContentTypeJSONUTF8)),
	}.Filter()
}

// SetDestinationType for DynatraceConnector
func (c *DynatraceConnector) SetDestinationType() {
	c.DestinationType = DestinationTypeDynatrace
}

// Validate validates DynatraceConnector
func (c *DynatraceConnector) Validate() error {
	return validation.Errors{
		"DestinationType":   validation.Validate(c.DestinationType, validation.Required, validation.In(DestinationTypeDynatrace)),
		"DisplayName":       validation.Validate(c.DisplayName, validation.Required),
		"Endpoint":          validation.Validate(c.Endpoint, validation.Required),
		"AuthToken":         validation.Validate(c.AuthToken, validation.Required),
		"CustomHeaderName":  validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue": validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}
