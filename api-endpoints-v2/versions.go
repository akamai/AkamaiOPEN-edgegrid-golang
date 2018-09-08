package apiendpoints

type Versions struct {
	APIEndPointID   int       `json:"apiEndPointId"`
	APIEndPointName string    `json:"apiEndPointName"`
	APIVersions     []Version `json:"apiVersions"`
}

type Version struct {
	CreatedBy            string       `json:"createdBy"`
	CreateDate           string       `json:"createDate"`
	UpdateDate           string       `json:"updateDate"`
	UpdatedBy            string       `json:"updatedBy"`
	APIEndPointVersionID int          `json:"apiEndPointVersionId"`
	BasePath             string       `json:"basePath"`
	Description          *string      `json:"description`
	BasedOn              *int         `json:"basedOn"`
	StagingStatus        *StatusValue `json:"stagingStatus"`
	ProductionStatus     *StatusValue `json:"productionStatus"`
	StagingDate          *string      `json:"stagingDate"`
	ProductionDate       *string      `json:"productionDate"`
	IsVersionLocked      bool         `json:"isVersionLocked"`
	AvailableActions     []string     `json:"availableActions"`
	VersionNumber        int          `json:"versionNumber"`
	LockVersion          int          `json:"lockVersion"`
}

type VersionSummary struct {
	Status        StatusValue `json:"status"`
	VersionNumber int         `json:"versionNumber"`
}

type StatusValue string

const (
	StatusPending     StatusValue = "PENDING"
	StatusActive      StatusValue = "ACTIVE"
	StatusDeactivated StatusValue = "DEACTIVATED"
	StatusFailed      StatusValue = "FAILED"
)
