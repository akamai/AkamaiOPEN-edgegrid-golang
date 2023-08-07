package botman

type (
	// ValidationResponse is used to represent validation member in the botman api response
	ValidationResponse struct {
		Errors   []ValidationDetail `json:"errors"`
		Notices  []ValidationDetail `json:"notices"`
		Warnings []ValidationDetail `json:"warnings"`
	}

	// ValidationDetail is used to represent validation details
	ValidationDetail struct {
		Title  string `json:"title"`
		Type   string `json:"type"`
		Detail string `json:"detail"`
	}
)
