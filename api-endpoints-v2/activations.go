package apiendpoints

type Activations struct {
	Networks               []NetworkValue `json:"networks"`
	NotificationRecipients []string       `json:"notificationRecipients"`
	Notes                  string         `json:"notes"`
}

// NetworkValue is used to create an "enum" of possible Activations.Networks[] values
type NetworkValue string

const (
	NetworkStaging    NetworkValue = "STAGING"
	NetworkProduction NetworkValue = "PRODUCTION"
)
