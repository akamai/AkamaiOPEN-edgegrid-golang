package cloudaccess

type (
	// KeyLink contains hypermedia link for the key
	KeyLink struct {
		AccessKeyUID int64  `json:"accessKeyUid"`
		Link         string `json:"link"`
	}

	// KeyVersion holds details for a version of an access key
	KeyVersion struct {
		AccessKeyUID int64  `json:"accessKeyUid"`
		Link         string `json:"link"`
		Version      int64  `json:"version"`
	}

	// RequestInformation contains information about a request to create an access key
	RequestInformation struct {
		AccessKeyName        string         `json:"accessKeyName"`
		AuthenticationMethod AuthType       `json:"authenticationMethod"`
		ContractID           string         `json:"contractId"`
		GroupID              int32          `json:"groupId"`
		NetworkConfiguration *SecureNetwork `json:"networkConfiguration"`
	}

	// SecureNetwork contains additional information about network
	SecureNetwork struct {
		AdditionalCDN   CDNType     `json:"additionalCdn"`
		SecurityNetwork NetworkType `json:"securityNetwork"`
	}

	// CDNType is a type of additionalCdn
	CDNType string

	// NetworkType is a type of securityNetwork
	NetworkType string

	// AuthType is a type of authentication
	AuthType string

	// ProcessingType is a type of ProcessingStatus
	ProcessingType string
)

const (
	// ChinaCDN represents CDN value of "CHINA_CDN"
	ChinaCDN CDNType = "CHINA_CDN"
	// RussiaCDN represents CDN value of "RUSSIA_CDN"
	RussiaCDN CDNType = "RUSSIA_CDN"

	// NetworkEnhanced represents Network value of "ENHANCED_TLS"
	NetworkEnhanced NetworkType = "ENHANCED_TLS"
	// NetworkStandard represents Network value of "STANDARD_TLS"
	NetworkStandard NetworkType = "STANDARD_TLS"

	// AuthAWS represents Authentication value of "AWS4_HMAC_SHA256"
	AuthAWS AuthType = "AWS4_HMAC_SHA256"
	// AuthGOOG represents Authentication value of "GOOG4_HMAC_SHA256"
	AuthGOOG AuthType = "GOOG4_HMAC_SHA256"

	// ProcessingInProgress represents ProcessingStatus value of 'IN_PROGRESS'
	ProcessingInProgress ProcessingType = "IN_PROGRESS"
	// ProcessingFailed represents ProcessingStatus value of 'FAILED'
	ProcessingFailed ProcessingType = "FAILED"
	// ProcessingDone represents ProcessingStatus value of 'DONE'
	ProcessingDone ProcessingType = "DONE"
)
