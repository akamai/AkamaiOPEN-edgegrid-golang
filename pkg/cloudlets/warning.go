package cloudlets

type (
	// Warning represents warning information regarding the policy version and loadbalancer version request
	Warning struct {
		Detail      string `json:"detail,omitempty"`
		JSONPointer string `json:"jsonPointer,omitempty"`
		Status      int    `json:"status,omitempty"`
		Title       string `json:"title,omitempty"`
		Type        string `json:"type,omitempty"`
	}
)
