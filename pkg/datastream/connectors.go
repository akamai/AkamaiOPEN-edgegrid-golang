package datastream

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// S3Connector provides details about the Amazon S3 connector in a stream.
	//
	// See: https://techdocs.akamai.com/datastream2/docs/stream-amazon-s3
	S3Connector struct {
		ConnectorType   ConnectorType `json:"connectorType"`
		AccessKey       string        `json:"accessKey"`
		Bucket          string        `json:"bucket"`
		ConnectorName   string        `json:"connectorName"`
		Path            string        `json:"path"`
		Region          string        `json:"region"`
		SecretAccessKey string        `json:"secretAccessKey"`
	}

	// AzureConnector provides details about the Azure Storage connector configuration in a data stream.
	//
	// See: https://techdocs.akamai.com/datastream2/docs/stream-azure-storage
	AzureConnector struct {
		ConnectorType ConnectorType `json:"connectorType"`
		AccessKey     string        `json:"accessKey"`
		AccountName   string        `json:"accountName"`
		ConnectorName string        `json:"connectorName"`
		ContainerName string        `json:"containerName"`
		Path          string        `json:"path"`
	}

	// DatadogConnector provides detailed information about Datadog connector.
	//
	// See: https://techdocs.akamai.com/datastream2/docs/stream-datadog
	DatadogConnector struct {
		ConnectorType ConnectorType `json:"connectorType"`
		AuthToken     string        `json:"authToken"`
		CompressLogs  bool          `json:"compressLogs"`
		ConnectorName string        `json:"connectorName"`
		Service       string        `json:"service,omitempty"`
		Source        string        `json:"source,omitempty"`
		Tags          string        `json:"tags,omitempty"`
		URL           string        `json:"url"`
	}

	// SplunkConnector provides detailed information about the Splunk connector.
	//
	// See: https://techdocs.akamai.com/datastream2/docs/stream-splunk
	SplunkConnector struct {
		ConnectorType       ConnectorType `json:"connectorType"`
		CompressLogs        bool          `json:"compressLogs"`
		ConnectorName       string        `json:"connectorName"`
		EventCollectorToken string        `json:"eventCollectorToken"`
		URL                 string        `json:"url"`
		CustomHeaderName    string        `json:"customHeaderName,omitempty"`
		CustomHeaderValue   string        `json:"customHeaderValue,omitempty"`
		TLSHostname         string        `json:"tlsHostname,omitempty"`
		CACert              string        `json:"caCert,omitempty"`
		ClientCert          string        `json:"clientCert,omitempty"`
		ClientKey           string        `json:"clientKey,omitempty"`
	}

	// GCSConnector provides detailed information about the Google Cloud Storage connector.
	//
	// See: https://techdocs.akamai.com/datastream2/docs/stream-google-cloud
	GCSConnector struct {
		ConnectorType      ConnectorType `json:"connectorType"`
		Bucket             string        `json:"bucket"`
		ConnectorName      string        `json:"connectorName"`
		Path               string        `json:"path,omitempty"`
		PrivateKey         string        `json:"privateKey"`
		ProjectID          string        `json:"projectId"`
		ServiceAccountName string        `json:"serviceAccountName"`
	}

	// CustomHTTPSConnector provides detailed information about the custom HTTPS endpoint.
	//
	// See: https://techdocs.akamai.com/datastream2/docs/stream-custom-https
	CustomHTTPSConnector struct {
		ConnectorType      ConnectorType      `json:"connectorType"`
		AuthenticationType AuthenticationType `json:"authenticationType"`
		CompressLogs       bool               `json:"compressLogs"`
		ConnectorName      string             `json:"connectorName"`
		Password           string             `json:"password,omitempty"`
		URL                string             `json:"url"`
		UserName           string             `json:"userName,omitempty"`
		ContentType        string             `json:"contentType,omitempty"`
		CustomHeaderName   string             `json:"customHeaderName,omitempty"`
		CustomHeaderValue  string             `json:"customHeaderValue,omitempty"`
		TLSHostname        string             `json:"tlsHostname,omitempty"`
		CACert             string             `json:"caCert,omitempty"`
		ClientCert         string             `json:"clientCert,omitempty"`
		ClientKey          string             `json:"clientKey,omitempty"`
	}

	// SumoLogicConnector provides detailed information about the Sumo Logic connector.
	//
	// See: https://techdocs.akamai.com/datastream2/docs/stream-sumo-logic
	SumoLogicConnector struct {
		ConnectorType     ConnectorType `json:"connectorType"`
		CollectorCode     string        `json:"collectorCode"`
		CompressLogs      bool          `json:"compressLogs"`
		ConnectorName     string        `json:"connectorName"`
		Endpoint          string        `json:"endpoint"`
		ContentType       string        `json:"contentType,omitempty"`
		CustomHeaderName  string        `json:"customHeaderName,omitempty"`
		CustomHeaderValue string        `json:"customHeaderValue,omitempty"`
	}

	// OracleCloudStorageConnector provides details about the Oracle Cloud Storage connector.
	//
	// See: https://techdocs.akamai.com/datastream2/docs/stream-oracle-cloud
	OracleCloudStorageConnector struct {
		ConnectorType   ConnectorType `json:"connectorType"`
		AccessKey       string        `json:"accessKey"`
		Bucket          string        `json:"bucket"`
		ConnectorName   string        `json:"connectorName"`
		Namespace       string        `json:"namespace"`
		Path            string        `json:"path"`
		Region          string        `json:"region"`
		SecretAccessKey string        `json:"secretAccessKey"`
	}

	// LogglyConnector contains details about Loggly connector.
	//
	// See: https://techdocs.akamai.com/datastream2/docs/stream-loggly
	LogglyConnector struct {
		ConnectorType     ConnectorType `json:"connectorType"`
		ConnectorName     string        `json:"connectorName"`
		Endpoint          string        `json:"endpoint"`
		AuthToken         string        `json:"authToken"`
		Tags              string        `json:"tags,omitempty"`
		ContentType       string        `json:"contentType,omitempty"`
		CustomHeaderName  string        `json:"customHeaderName,omitempty"`
		CustomHeaderValue string        `json:"customHeaderValue,omitempty"`
	}

	// NewRelicConnector connector contains details about New Relic connector.
	//
	// See: https://techdocs.akamai.com/datastream2/docs/stream-new-relic
	NewRelicConnector struct {
		ConnectorType     ConnectorType `json:"connectorType"`
		ConnectorName     string        `json:"connectorName"`
		Endpoint          string        `json:"endpoint"`
		AuthToken         string        `json:"authToken"`
		ContentType       string        `json:"contentType,omitempty"`
		CustomHeaderName  string        `json:"customHeaderName,omitempty"`
		CustomHeaderValue string        `json:"customHeaderValue,omitempty"`
	}

	// ElasticsearchConnector contains details about Elasticsearch connector.
	//
	// See: https://techdocs.akamai.com/datastream2/docs/stream-elasticsearch
	ElasticsearchConnector struct {
		ConnectorType     ConnectorType `json:"connectorType"`
		ConnectorName     string        `json:"connectorName"`
		Endpoint          string        `json:"endpoint"`
		IndexName         string        `json:"indexName"`
		UserName          string        `json:"userName"`
		Password          string        `json:"password"`
		ContentType       string        `json:"contentType,omitempty"`
		CustomHeaderName  string        `json:"customHeaderName,omitempty"`
		CustomHeaderValue string        `json:"customHeaderValue,omitempty"`
		TLSHostname       string        `json:"tlsHostname,omitempty"`
		CACert            string        `json:"caCert,omitempty"`
		ClientCert        string        `json:"clientCert,omitempty"`
		ClientKey         string        `json:"clientKey,omitempty"`
	}

	// ConnectorType is used to create an "enum" of possible ConnectorTypes
	ConnectorType string

	// AuthenticationType is used to create an "enum" of possible AuthenticationTypes of the CustomHTTPSConnector
	AuthenticationType string
)

const (
	// ConnectorTypeAzure const
	ConnectorTypeAzure ConnectorType = "AZURE"
	// ConnectorTypeS3 const
	ConnectorTypeS3 ConnectorType = "S3"
	// ConnectorTypeDataDog const
	ConnectorTypeDataDog ConnectorType = "DATADOG"
	// ConnectorTypeSplunk const
	ConnectorTypeSplunk ConnectorType = "SPLUNK"
	// ConnectorTypeGcs const
	ConnectorTypeGcs ConnectorType = "GCS"
	// ConnectorTypeHTTPS const
	ConnectorTypeHTTPS ConnectorType = "HTTPS"
	// ConnectorTypeSumoLogic const
	ConnectorTypeSumoLogic ConnectorType = "SUMO_LOGIC"
	// ConnectorTypeOracle const
	ConnectorTypeOracle ConnectorType = "Oracle_Cloud_Storage"
	// ConnectorTypeLoggly const
	ConnectorTypeLoggly ConnectorType = "LOGGLY"
	// ConnectorTypeNewRelic const
	ConnectorTypeNewRelic ConnectorType = "NEWRELIC"
	// ConnectorTypeElasticsearch const
	ConnectorTypeElasticsearch ConnectorType = "ELASTICSEARCH"

	// AuthenticationTypeNone const
	AuthenticationTypeNone AuthenticationType = "NONE"
	// AuthenticationTypeBasic const
	AuthenticationTypeBasic AuthenticationType = "BASIC"
)

var customHeaderNameRegexp = regexp.MustCompile("^[A-Za-z0-9_-]+$")

// SetConnectorType for S3Connector
func (c *S3Connector) SetConnectorType() {
	c.ConnectorType = ConnectorTypeS3
}

// Validate validates S3Connector
func (c *S3Connector) Validate() error {
	return validation.Errors{
		"ConnectorType":   validation.Validate(c.ConnectorType, validation.Required, validation.In(ConnectorTypeS3)),
		"AccessKey":       validation.Validate(c.AccessKey, validation.Required),
		"Bucket":          validation.Validate(c.Bucket, validation.Required),
		"ConnectorName":   validation.Validate(c.ConnectorName, validation.Required),
		"Path":            validation.Validate(c.Path, validation.Required),
		"Region":          validation.Validate(c.Region, validation.Required),
		"SecretAccessKey": validation.Validate(c.SecretAccessKey, validation.Required),
	}.Filter()
}

// SetConnectorType for AzureConnector
func (c *AzureConnector) SetConnectorType() {
	c.ConnectorType = ConnectorTypeAzure
}

// Validate validates AzureConnector
func (c *AzureConnector) Validate() error {
	return validation.Errors{
		"ConnectorType": validation.Validate(c.ConnectorType, validation.Required, validation.In(ConnectorTypeAzure)),
		"AccessKey":     validation.Validate(c.AccessKey, validation.Required),
		"AccountName":   validation.Validate(c.AccountName, validation.Required),
		"ConnectorName": validation.Validate(c.ConnectorName, validation.Required),
		"ContainerName": validation.Validate(c.ContainerName, validation.Required),
		"Path":          validation.Validate(c.Path, validation.Required),
	}.Filter()
}

// SetConnectorType for DatadogConnector
func (c *DatadogConnector) SetConnectorType() {
	c.ConnectorType = ConnectorTypeDataDog
}

// Validate validates DatadogConnector
func (c *DatadogConnector) Validate() error {
	return validation.Errors{
		"ConnectorType": validation.Validate(c.ConnectorType, validation.Required, validation.In(ConnectorTypeDataDog)),
		"AuthToken":     validation.Validate(c.AuthToken, validation.Required),
		"ConnectorName": validation.Validate(c.ConnectorName, validation.Required),
		"URL":           validation.Validate(c.URL, validation.Required),
	}.Filter()
}

// SetConnectorType for SplunkConnector
func (c *SplunkConnector) SetConnectorType() {
	c.ConnectorType = ConnectorTypeSplunk
}

// Validate validates SplunkConnector
func (c *SplunkConnector) Validate() error {
	return validation.Errors{
		"ConnectorType":       validation.Validate(c.ConnectorType, validation.Required, validation.In(ConnectorTypeSplunk)),
		"ConnectorName":       validation.Validate(c.ConnectorName, validation.Required),
		"EventCollectorToken": validation.Validate(c.EventCollectorToken, validation.Required),
		"URL":                 validation.Validate(c.URL, validation.Required),
		"CustomHeaderName":    validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue":   validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}

// SetConnectorType for GCSConnector
func (c *GCSConnector) SetConnectorType() {
	c.ConnectorType = ConnectorTypeGcs
}

// Validate validates GCSConnector
func (c *GCSConnector) Validate() error {
	return validation.Errors{
		"ConnectorType":      validation.Validate(c.ConnectorType, validation.Required, validation.In(ConnectorTypeGcs)),
		"Bucket":             validation.Validate(c.Bucket, validation.Required),
		"ConnectorName":      validation.Validate(c.ConnectorName, validation.Required),
		"PrivateKey":         validation.Validate(c.PrivateKey, validation.Required),
		"ProjectID":          validation.Validate(c.ProjectID, validation.Required),
		"ServiceAccountName": validation.Validate(c.ServiceAccountName, validation.Required),
	}.Filter()
}

// SetConnectorType for CustomHTTPSConnector
func (c *CustomHTTPSConnector) SetConnectorType() {
	c.ConnectorType = ConnectorTypeHTTPS
}

// Validate validates CustomHTTPSConnector
func (c *CustomHTTPSConnector) Validate() error {
	return validation.Errors{
		"ConnectorType":      validation.Validate(c.ConnectorType, validation.Required, validation.In(ConnectorTypeHTTPS)),
		"AuthenticationType": validation.Validate(c.AuthenticationType, validation.Required, validation.In(AuthenticationTypeBasic, AuthenticationTypeNone)),
		"ConnectorName":      validation.Validate(c.ConnectorName, validation.Required),
		"URL":                validation.Validate(c.URL, validation.Required),
		"UserName":           validation.Validate(c.UserName, validation.Required.When(c.AuthenticationType == AuthenticationTypeBasic)),
		"Password":           validation.Validate(c.Password, validation.Required.When(c.AuthenticationType == AuthenticationTypeBasic)),
		"CustomHeaderName":   validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue":  validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}

// SetConnectorType for SumoLogicConnector
func (c *SumoLogicConnector) SetConnectorType() {
	c.ConnectorType = ConnectorTypeSumoLogic
}

// Validate validates SumoLogicConnector
func (c *SumoLogicConnector) Validate() error {
	return validation.Errors{
		"ConnectorType":     validation.Validate(c.ConnectorType, validation.Required, validation.In(ConnectorTypeSumoLogic)),
		"CollectorCode":     validation.Validate(c.CollectorCode, validation.Required),
		"ConnectorName":     validation.Validate(c.ConnectorName, validation.Required),
		"Endpoint":          validation.Validate(c.Endpoint, validation.Required),
		"CustomHeaderName":  validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue": validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}

// SetConnectorType for OracleCloudStorageConnector
func (c *OracleCloudStorageConnector) SetConnectorType() {
	c.ConnectorType = ConnectorTypeOracle
}

// Validate validates OracleCloudStorageConnector
func (c *OracleCloudStorageConnector) Validate() error {
	return validation.Errors{
		"ConnectorType":   validation.Validate(c.ConnectorType, validation.Required, validation.In(ConnectorTypeOracle)),
		"AccessKey":       validation.Validate(c.AccessKey, validation.Required),
		"Bucket":          validation.Validate(c.Bucket, validation.Required),
		"ConnectorName":   validation.Validate(c.ConnectorName, validation.Required),
		"Namespace":       validation.Validate(c.Namespace, validation.Required),
		"Path":            validation.Validate(c.Path, validation.Required),
		"Region":          validation.Validate(c.Region, validation.Required),
		"SecretAccessKey": validation.Validate(c.SecretAccessKey, validation.Required),
	}.Filter()
}

// SetConnectorType for LogglyConnector
func (c *LogglyConnector) SetConnectorType() {
	c.ConnectorType = ConnectorTypeLoggly
}

// Validate validates LogglyConnector
func (c *LogglyConnector) Validate() error {
	return validation.Errors{
		"ConnectorType":     validation.Validate(c.ConnectorType, validation.Required, validation.In(ConnectorTypeLoggly)),
		"ConnectorName":     validation.Validate(c.ConnectorName, validation.Required),
		"Endpoint":          validation.Validate(c.Endpoint, validation.Required),
		"AuthToken":         validation.Validate(c.AuthToken, validation.Required),
		"CustomHeaderName":  validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue": validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}

// SetConnectorType for NewRelicConnector
func (c *NewRelicConnector) SetConnectorType() {
	c.ConnectorType = ConnectorTypeNewRelic
}

// Validate validates NewRelicConnector
func (c *NewRelicConnector) Validate() error {
	return validation.Errors{
		"ConnectorType":     validation.Validate(c.ConnectorType, validation.Required, validation.In(ConnectorTypeNewRelic)),
		"ConnectorName":     validation.Validate(c.ConnectorName, validation.Required),
		"Endpoint":          validation.Validate(c.Endpoint, validation.Required),
		"AuthToken":         validation.Validate(c.AuthToken, validation.Required),
		"CustomHeaderName":  validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue": validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}

// SetConnectorType for ElasticsearchConnector
func (c *ElasticsearchConnector) SetConnectorType() {
	c.ConnectorType = ConnectorTypeElasticsearch
}

// Validate validates ElasticsearchConnector
func (c *ElasticsearchConnector) Validate() error {
	return validation.Errors{
		"ConnectorType":     validation.Validate(c.ConnectorType, validation.Required, validation.In(ConnectorTypeElasticsearch)),
		"ConnectorName":     validation.Validate(c.ConnectorName, validation.Required),
		"Endpoint":          validation.Validate(c.Endpoint, validation.Required),
		"UserName":          validation.Validate(c.UserName, validation.Required),
		"Password":          validation.Validate(c.Password, validation.Required),
		"IndexName":         validation.Validate(c.IndexName, validation.Required),
		"CustomHeaderName":  validation.Validate(c.CustomHeaderName, validation.Required.When(c.CustomHeaderValue != ""), validation.When(c.CustomHeaderName != "", validation.Match(customHeaderNameRegexp))),
		"CustomHeaderValue": validation.Validate(c.CustomHeaderValue, validation.Required.When(c.CustomHeaderName != "")),
	}.Filter()
}
