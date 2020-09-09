package papi

type (
	// PropertyCloneFrom optionally identifies another property instance to clone when making a POST request to create a new property
	PropertyCloneFrom struct {
		CloneFromVersionEtag string `json:"cloneFromVersionEtag"`
		CopyHostnames        bool   `json:"copyHostnames"`
		PropertyID           string `json:"propertyId"`
		Version              int    `json:"version"`
	}

	// Property contains configuration data to apply to edge content.
	Property struct {
		AccountID         string             `json:"accountId"`
		AssetID           string             `json:"assetId"`
		CloneFrom         *PropertyCloneFrom `json:"cloneFrom,omitempty`
		ContactID         string             `json:"contractId"`
		GroupID           string             `json:"groupId"`
		LatestVersion     int                `json:"latestVersion"`
		Note              string             `json:"note"`
		ProductID         string             `json:"productId"`
		ProductionVersion *int               `json:"productionVersion,omitempty`
		PropertyID        string             `json:"propertyId"`
		PropertyName      string             `json:"propertyName"`
		RuleFormat        string             `json:"ruleFormat"`
		StagingVersion    *int               `json:"stagingVersion,omitempty"`
	}
)
