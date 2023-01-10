package datastream

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// S3Connector provides details about the Amazon S3 destination in a stream
	// See: https://techdocs.akamai.com/datastream2/v2/reference/post-stream
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
	// See: https://techdocs.akamai.com/datastream2/v2/reference/post-stream
	AzureConnector struct {
		DestinationType DestinationType `json:"destinationType"`
		AccessKey       string          `json:"accessKey"`
		AccountName     string          `json:"accountName"`
		DisplayName     string          `json:"displayName"`
		ContainerName   string          `json:"containerName"`
		Path            string          `json:"path"`
	}

	// DatadogConnector provides detailed information about Datadog destination
	// See: https://techdocs.akamai.com/datastream2/v2/reference/post-stream
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
	// See: https://techdocs.akamai.com/datastream2/v2/reference/post-stream
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
	// See: https://techdocs.akamai.com/datastream2/v2/reference/post-stream
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
	// See: https://techdocs.akamai.com/datastream2/v2/reference/post-stream
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
	// See: https://techdocs.akamai.com/datastream2/v2/reference/post-stream
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
	// See: https://techdocs.akamai.com/datastream2/v2/reference/post-stream
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
	// See: https://techdocs.akamai.com/datastream2/v2/reference/post-stream
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
	// See: https://techdocs.akamai.com/datastream2/v2/reference/post-stream
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
	// See: https://techdocs.akamai.com/datastream2/v2/reference/post-stream
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

	// DestinationType is used to create an "enum" of possible DestinationTypes
	DestinationType string

	// AuthenticationType is used to create an "enum" of possible AuthenticationTypes of the CustomHTTPSConnector
	AuthenticationType string
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

	// AuthenticationTypeNone const
	AuthenticationTypeNone AuthenticationType = "NONE"
	// AuthenticationTypeBasic const
	AuthenticationTypeBasic AuthenticationType = "BASIC"
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
