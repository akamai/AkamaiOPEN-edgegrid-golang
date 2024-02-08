package cloudlets

type (
	// Warning represents warning information about the policy version request.
	Warning struct {
		Detail      string `json:"detail,omitempty"`
		JSONPointer string `json:"jsonPointer,omitempty"`
		Status      int    `json:"status,omitempty"`
		Title       string `json:"title,omitempty"`
		Type        string `json:"type,omitempty"`
	}
)
