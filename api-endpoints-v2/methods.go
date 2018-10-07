package apiendpoints

type Methods []Method

type Method struct {
	APIResourceMethodID      int          `json:"apiResourceMethodId"`
	APIResourceMethod        string       `json:"apiResourceMethod"`
	APIResourceMethodLogicID int          `json:"apiResourceMethodLogicId"`
	APIParameters            []Parameters `json:"apiParameters"`
}

type MethodValue string

const (
	MethodGet     string = "get"
	MethodPost    string = "post"
	MethodPut     string = "put"
	MethodDelete  string = "delete"
	MethodHead    string = "head"
	MethodPatch   string = "patch"
	MethodOptions string = "options"
)
