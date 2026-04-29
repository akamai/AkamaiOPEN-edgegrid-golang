package mtlstruststore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/internal/texts"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// CreateCASetRequest holds request body for CreateCASet.
	CreateCASetRequest struct {
		// CASetName is a descriptive name for the set.
		CASetName string `json:"caSetName"`

		// Description is an optional description for the set.
		Description *string `json:"description,omitempty"`
	}

	// CASetResponse contains response from create, list operations.
	CASetResponse struct {
		// AccountID is the ID of the account under which this CA set was created.
		AccountID string `json:"accountId"`

		// CASetID is a unique identifier representing the CA set.
		CASetID string `json:"caSetId"`

		// CASetLink is the hypermedia link to the CA set.
		CASetLink string `json:"caSetLink"`

		// CASetName is the name of the CA set.
		CASetName string `json:"caSetName"`

		// CASetStatus is the status of the CA set - "NOT_DELETED", "DELETING", "DELETED".
		CASetStatus string `json:"caSetStatus"`

		// Description is a description of the CA set.
		Description *string `json:"description"`

		// LatestVersionLink is the hypermedia link for newly created / cloned version in a CA set.
		LatestVersionLink *string `json:"latestVersionLink"`

		// LatestVersion is the version number for newly created / cloned version in a CA set.
		LatestVersion *int64 `json:"latestVersion"`

		// StagingVersionLink is the hypermedia link for the version of the CA set that is active on staging.
		// This field could be `nil` if no version of the CA set was activated on staging.
		StagingVersionLink *string `json:"stagingVersionLink"`

		// StagingVersion is the version number of the CA set that is active on staging.
		StagingVersion *int64 `json:"stagingVersion"`

		// ProductionVersionLink is the hypermedia link for the version of the CA set that is active on production.
		// This field could be `nil` if no version of the set was activated on production.
		ProductionVersionLink *string `json:"productionVersionLink"`

		// ProductionVersion is the version number of the CA set that is active on production.
		ProductionVersion *int64 `json:"productionVersion"`

		// VersionsLink is the hypermedia link to the list of versions in the CA set.
		VersionsLink string `json:"versionsLink"`

		// CreatedDate is the date the set was created.
		CreatedDate time.Time `json:"createdDate"`

		// CreatedBy is the user who created the set.
		CreatedBy string `json:"createdBy"`

		// DeletedDate is the date the CA set was deleted if the CA set has been archived, `nil` otherwise.
		DeletedDate *time.Time `json:"deletedDate"`

		// DeletedBy is the user who deleted the CA set if the CA set has been deleted, `nil` otherwise.
		DeletedBy *string `json:"deletedBy"`
	}

	// CreateCASetResponse contains response from CreateCASet.
	CreateCASetResponse CASetResponse

	// Network represents the network type when filtering list of ca sets.
	Network string

	// ListCASetsRequest holds request body for ListCASets.
	ListCASetsRequest struct {
		// CASetNamePrefix is the name prefix to filter out matching CA sets.
		CASetNamePrefix string

		// ActivatedOn is the network type to filter out matching CA sets.
		// A CA set is included in the response if any version of it is active on that network.
		// The values that could be provided are `INACTIVE`, `STAGING`, `PRODUCTION`, `STAGING+PRODUCTION`, `PRODUCTION+STAGING`, `STAGING,PRODUCTION`, `PRODUCTION,STAGING`.
		// A CA set will not be included if it was created but none of its versions was ever activated.
		ActivatedOn Network

		// CASetStatuses filters CA sets by status values. Allowed values: `NOT_DELETED`, `DELETING`, `DELETED`. Defaults to `NOT_DELETED`.
		CASetStatuses []string
	}

	// ListCASetsResponse contains response from ListCASets.
	ListCASetsResponse struct {
		// CASets is a list of CA set objects with each object representing the details of one CA set.
		CASets []CASetResponse `json:"caSets"`
	}

	// GetCASetRequest holds request body for GetCASet.
	GetCASetRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string `json:"caSetId"`
	}

	// GetCASetResponse contains response from GetCASet.
	GetCASetResponse CASetResponse

	// DeleteCASetRequest holds request body for DeleteCASet.
	DeleteCASetRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string `json:"caSetId"`
	}

	// AssociationType is type of associations to be returned.
	AssociationType string

	// ListCASetAssociationsRequest holds request for ListCASetAssociations.
	ListCASetAssociationsRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string

		// AssociationType is optional type of associations to be returned. The values that could be provided are `enrollments` or `properties`.
		// If not provided, both types of associations are returned.
		AssociationType AssociationType `json:"associationType,omitempty"`
	}

	// ListCASetAssociationsResponse holds response for ListCASetAssociations.
	ListCASetAssociationsResponse struct {
		// Associations contains list of properties and enrollments.
		Associations Associations `json:"associations"`
	}

	// Associations represents associations data for given CASetID.
	Associations struct {
		// Properties contains details about properties with associated hostnames.
		Properties []AssociationProperty `json:"properties"`

		// Enrollments contains details about enrollments associated with this CA Set.
		Enrollments []AssociationEnrollment `json:"enrollments"`
	}

	// AssociationProperty represents one property associations data for given CASetID.
	AssociationProperty struct {
		// PropertyID is a unique identifier for the property.
		PropertyID string `json:"propertyId"`

		// PropertyName is a unique, descriptive name for the property.
		PropertyName *string `json:"propertyName"`

		// PropertyLink is link to the property.
		PropertyLink string `json:"propertyLink"`

		// AssetID is an alternative identifier for the property.
		AssetID *int64 `json:"assetId"`

		// GroupID identifies the group to which the property is assigned.
		GroupID *int64 `json:"groupId"`

		// Hostnames contains details about associated hostnames.
		Hostnames []AssociationHostname `json:"hostnames"`
	}

	// AssociationHostname represents one hostname associations data for given CASetID.
	AssociationHostname struct {
		// Hostname is name of device.
		Hostname string `json:"hostName"`

		// Network indicates the network on which CA set to hostname association is formed/removed/in progress. The values for this are "STAGING", "PRODUCTION".
		Network string `json:"network"`

		// Status indicates the status of CA set to hostname association. The values for it are - "ATTACHING", "DETACHING", "ATTACHED".
		Status string `json:"status"`
	}

	// AssociationEnrollment represents one enrollment associations data for given CASetID.
	AssociationEnrollment struct {
		// EnrollmentID is unique identifier for the enrollment.
		EnrollmentID int64 `json:"enrollmentId"`

		// EnrollmentLink is link to CPS enrollment.
		EnrollmentLink string `json:"enrollmentLink"`

		// StagingSlots are slots where the certificate is deployed on the staging network.
		StagingSlots []int64 `json:"stagingSlots"`

		// ProductionSlots are slots where the certificate is deployed on the production network.
		ProductionSlots []int64 `json:"productionSlots"`

		// CN is the domain name to use for the certificate, also known as the common name.
		CN string `json:"cn"`
	}

	// CloneCASetRequest holds request for CloneCASet.
	CloneCASetRequest struct {
		// CloneFromSetID is a CA set ID which should be used to create new CA set from.
		CloneFromSetID string `json:"-"`

		// CloneFromVersion is an optional version of CA set which should be used to create new CA set from.
		// If not provided, the latest version of the CA set is used by default.
		CloneFromVersion int64 `json:"-"`

		// NewCASetName is descriptive name for the set.
		NewCASetName string `json:"caSetName"`

		// NewDescription is optional description for the set.
		NewDescription *string `json:"description"`
	}

	// CloneCASetResponse holds response body for CloneCASet.
	CloneCASetResponse CASetResponse

	// GetCASetDeletionStatusRequest holds request for GetCASetDeleteStatus.
	GetCASetDeletionStatusRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string
	}

	// GetCASetDeletionStatusResponse holds response for GetCASetDeleteStatus.
	GetCASetDeletionStatusResponse struct {
		// Status is current status of CA set deletion, either "IN_PROGRESS", "COMPLETE", or "FAILED".
		Status string `json:"status"`

		// StatusLink is the hypermedia link to the deletion status.
		StatusLink string `json:"statusLink"`

		// CASetLink is the hypermedia link to the CA set.
		CASetLink string `json:"caSetLink"`

		// ResourceMethod is the REST API method of the actual operation. Always has a value "delete".
		ResourceMethod *string `json:"resourceMethod"`

		// CASetID is a unique identifier representing the CA set.
		CASetID string `json:"caSetId"`

		// CASetName is descriptive name for the set.
		CASetName string `json:"caSetName"`

		// EstimatedEndTime is approximated time when CA set should be deleted.
		EstimatedEndTime *time.Time `json:"estimatedEndTime"`

		// EndTime is time when CA set was deleted.
		EndTime *time.Time `json:"endTime"`

		// FailureReason is a reason why CA set was not deleted.
		FailureReason *string `json:"failureReason"`

		// StartTime is a time when CA set was asked to be deleted.
		StartTime time.Time `json:"startTime"`

		// Deletions is a list of status for each network about CA set deletion.
		Deletions []CASetNetworkDeleteStatus `json:"deletions"`

		// RetryAfter is a time when CA set deletion status can be checked again.
		// Usually 300 seconds after the deletion request was made.
		// This header value is returned only if the CA set deletion status is "IN_PROGRESS".
		RetryAfter time.Time
	}

	// CASetNetworkDeleteStatus holds information about one network for GetCASetDeleteStatus response.
	CASetNetworkDeleteStatus struct {
		// Network represents network type for this delete status.
		Network string `json:"network"`

		// PercentComplete represents progress of CA set deletion of given network.
		PercentComplete int `json:"percentComplete"`

		// Status represents CA set deletion status on given network, either "IN_PROGRESS", "COMPLETE", or "FAILED".
		Status string `json:"status"`

		// FailureReason is a reason why CA set was not deleted on given network.
		FailureReason *string `json:"failureReason"`
	}

	// ListCASetActivitiesRequest holds request for ListCASetActivities.
	ListCASetActivitiesRequest struct {
		// CASetID is a unique identifier representing the CA set.
		CASetID string

		// Start is optional Date in ISO-8601 format with fractional seconds. Include an action object in the result data if the actionDate is greater than or equal to start date.
		Start time.Time

		// End is optional Date in ISO-8601 format with fractional seconds. Include an action object in the result data if the actionDate is less than or equal to end date.
		End time.Time
	}

	// ListCASetActivitiesResponse holds response for ListCASetActivities.
	ListCASetActivitiesResponse struct {
		// CASetID is unique identifier of the set.
		CASetID string `json:"caSetID"`

		// CASetLink is hypermedia link to the CA set resource.
		CASetLink string `json:"caSetLink"`

		// CASetName is name of the CA set.
		CASetName string `json:"caSetName"`

		// CreatedDate is date the CA set was created in ISO-8601 format.
		CreatedDate time.Time `json:"createdDate"`

		// CreatedBy is user who created the CA set.
		CreatedBy string `json:"createdBy"`

		//CASetStatus indicates the status of the CA set. Could be one of: "NOT_DELETED", "DELETED", "DELETING".
		CASetStatus string `json:"caSetStatus"`

		// DeletedDate is date the CA set was deleted in ISO-8601 format. null if the CA set is not deleted.
		DeletedDate *time.Time `json:"deletedDate"`

		// DeletedBy is user who deleted the CA set if its deleted, null otherwise.
		DeletedBy *string `json:"deletedBy"`

		// Activities are list of activities on the CA set with each element in the list representing one activity object.
		Activities []CASetActivity `json:"activities"`
	}

	// CASetActivity holds one activity entry from ListCASetActivities response.
	CASetActivity struct {
		// Type is activity. Could be one of: "CREATE_CA_SET", "CREATE_CA_SET_VERSION", "ACTIVATE_CA_SET_VERSION", "DEACTIVATE_CA_SET_VERSION", "DELETE_CA_SET".
		Type string `json:"type"`

		// Network is the network for any activation-related activities, either "STAGING" or "PRODUCTION".
		Network *string `json:"network"`

		// Version on which user acted on.
		Version *int64 `json:"version"`

		// ActivityDate is the date associated with the activity in ISO-8601 format.
		ActivityDate time.Time `json:"activityDate"`

		// ActivityBy is user who submitted the request.
		ActivityBy string `json:"activityBy"`
	}
)

const (
	// NetworkInactive represents staging network.
	NetworkInactive Network = "INACTIVE"
	// NetworkStaging represents staging network.
	NetworkStaging Network = "STAGING"
	// NetworkProduction represents production network.
	NetworkProduction Network = "PRODUCTION"
	// NetworkStagingAndProduction represents staging and production networks.
	NetworkStagingAndProduction Network = "STAGING+PRODUCTION"
	// NetworkProductionAndStaging also represents staging and production networks.
	NetworkProductionAndStaging Network = "PRODUCTION+STAGING"
	// NetworkStagingOrProduction represents staging or production networks.
	NetworkStagingOrProduction Network = "STAGING,PRODUCTION"
	// NetworkProductionOrStaging also represents staging or production networks.
	NetworkProductionOrStaging Network = "PRODUCTION,STAGING"

	// CASetNamePattern is the regex pattern for CA set name.
	CASetNamePattern string = `^[%.a-zA-Z0-9_-]+$`

	// CASetNameDescription describes allowed characters for CA set name.
	CASetNameDescription = "allowed characters are alphanumerics (a-z, A-Z, 0-9), underscore (_), hyphen (-), percent (%) and period (.)"

	// DeletionStatusInProgress represents CA set deletion status in progress.
	DeletionStatusInProgress string = "IN_PROGRESS"
	// DeletionStatusComplete represents CA set deletion status complete.
	DeletionStatusComplete string = "COMPLETE"
	// DeletionStatusFailed represents CA set deletion status failed.
	DeletionStatusFailed string = "FAILED"

	// CASetStatusNotDeleted represents CA set status not deleted.
	CASetStatusNotDeleted string = "NOT_DELETED"
	// CASetStatusDeleted represents CA set status deleted.
	CASetStatusDeleted string = "DELETED"
	// CASetStatusDeleting represents CA set status deleting.
	CASetStatusDeleting string = "DELETING"

	// AssociationTypeEnrollments represents enrollments associations type.
	AssociationTypeEnrollments AssociationType = "enrollments"
	// AssociationTypeProperties represents properties associations type.
	AssociationTypeProperties AssociationType = "properties"
)

var (
	// CASetNameRegex is compiled regex for CA set name.
	CASetNameRegex = regexp.MustCompile(CASetNamePattern)

	// ErrCreateCASet is returned when the request to create a CA set fails.
	ErrCreateCASet = errors.New("create ca set failed")
	// ErrGetCASet is returned when the request to get a CA set fails.
	ErrGetCASet = errors.New("get ca set failed")
	// ErrListCASets is returned when the request to list CA sets fails.
	ErrListCASets = errors.New("list ca sets failed")
	// ErrDeleteCASet is returned when the request to delete a CA set fails.
	ErrDeleteCASet = errors.New("delete ca set failed")
	// ErrListCASetAssociations is returned when the request to get a CA set associations fails.
	ErrListCASetAssociations = errors.New("list ca sets associations failed")
	// ErrCloneCASet is returned when the request to clone a CA set fails.
	ErrCloneCASet = errors.New("clone ca set failed")
	// ErrGetCASetDeletionStatus is returned when the request to get a CA set deletion status fails.
	ErrGetCASetDeletionStatus = errors.New("list ca set deletion status failed")
	// ErrListCASetActivities is returned when the request to get a CA set activities fails.
	ErrListCASetActivities = errors.New("get ca set activities failed")
	// ErrValidateCertificates is returned when the request to validate certificates fails or some provided certificates are not valid.
	ErrValidateCertificates = errors.New("validate certificates failed")
)

// Validate validates CreateCASetRequest.
// Allowed characters are alphanumerics (a-z, A-Z, 0-9), underscore (_), hyphen (-), percent (%) and period (.) with no three consecutive periods (…).
// Length must be between 3 and 64 characters.
func (r CreateCASetRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetName": validation.Validate(r.CASetName,
			validation.Required,
			validation.Length(3, 64),
			validation.Match(CASetNameRegex).Error(CASetNameDescription),
			validateCASetName()),
		"Description": validation.Validate(r.Description, validation.NilOrNotEmpty, validation.Length(1, 255)),
	})
}

// Validate validates GetCASetRequest.
func (r GetCASetRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(r.CASetID, validation.Required, validation.Length(1, 0)),
	})
}

// Validate validates ListCASetsRequest.
func (r ListCASetsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetNamePrefix": validation.Validate(r.CASetNamePrefix,
			validation.Length(0, 64),
			validation.Match(CASetNameRegex).Error(CASetNameDescription),
			validateCASetName()),
		"ActivatedOn":   validation.Validate(r.ActivatedOn, r.ActivatedOn.Validate()),
		"CASetStatuses": validation.Validate(r.CASetStatuses, validation.By(caSetStatusRule)),
	})
}

// AllCASetStatuses returns list of all allowed CA set status values.
func AllCASetStatuses() []string {
	return []string{CASetStatusNotDeleted, CASetStatusDeleted, CASetStatusDeleting}
}

func caSetStatusRule(value any) error {
	statuses, ok := value.([]string)
	if !ok {
		return fmt.Errorf("expected []string, got %T", value)
	}

	allCASetStatuses := AllCASetStatuses()
	for _, s := range statuses {
		if !slices.Contains(allCASetStatuses, s) {
			return fmt.Errorf("list element '%s' is invalid. Each element must be one of: %s",
				s, "'"+strings.Join(allCASetStatuses, "', '")+"'")
		}
	}

	return nil
}

// Validate validates DeleteCASetsRequest.
func (r DeleteCASetRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(r.CASetID, validation.Required, validation.Length(1, 0)),
	})
}

// Validate validates ActivationNetwork.
func (n Network) Validate() validation.InRule {
	return validation.In(NetworkInactive, NetworkStaging, NetworkProduction, NetworkStagingAndProduction, NetworkProductionAndStaging,
		NetworkStagingOrProduction, NetworkProductionOrStaging).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s', '%s', '%s', '%s', '%s', '%s' or '%s'.",
			n, NetworkInactive, NetworkStaging, NetworkProduction, NetworkStagingAndProduction, NetworkProductionAndStaging,
			NetworkStagingOrProduction, NetworkProductionOrStaging))
}

func validateCASetName() validation.StringRule {
	return validation.NewStringRule(func(s string) bool {
		return !strings.Contains(s, "...")
	}, "cannot contain three consecutive periods (...)")
}

// Validate validates ListCASetAssociationsRequest.
func (r ListCASetAssociationsRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID":         validation.Validate(r.CASetID, validation.Required, validation.Length(1, 0)),
		"AssociationType": validation.Validate(r.AssociationType, r.AssociationType.Validate()),
	})
}

// Validate validates AssociationType.
func (t AssociationType) Validate() validation.InRule {
	return validation.In(AssociationTypeEnrollments, AssociationTypeProperties).
		Error(fmt.Sprintf("value '%s' is invalid. Must be one of: '%s' or '%s'.",
			t, AssociationTypeEnrollments, AssociationTypeProperties))
}

// Validate validates CloneCASetRequest.
func (r CloneCASetRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CloneFromSetID":   validation.Validate(r.CloneFromSetID, validation.Required),
		"CloneFromVersion": validation.Validate(r.CloneFromVersion, validation.Min(1)),
		"NewCASetName": validation.Validate(r.NewCASetName,
			validation.Required,
			validation.Length(3, 64),
			validation.Match(CASetNameRegex).Error(CASetNameDescription),
			validateCASetName()),
		"NewDescription": validation.Validate(r.NewDescription, validation.NilOrNotEmpty, validation.Length(1, 255)),
	})
}

// Validate validates GetCASetDeletionStatusRequest.
func (r GetCASetDeletionStatusRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(r.CASetID, validation.Required, validation.Length(1, 0)),
	})
}

// Validate validates ListCASetActivitiesRequest.
func (r ListCASetActivitiesRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"CASetID": validation.Validate(r.CASetID, validation.Required, validation.Length(1, 0)),
	})
}

func (m *mtlstruststore) CreateCASet(ctx context.Context, params CreateCASetRequest) (*CreateCASetResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("CreateCASet")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCreateCASet, ErrStructValidation, err)
	}

	uri, err := url.Parse("/mtls-edge-truststore/v2/ca-sets")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCreateCASet, err)
	}

	var result CreateCASetResponse
	resp, err := m.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCreateCASet, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) GetCASet(ctx context.Context, params GetCASetRequest) (*GetCASetResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("GetCASet")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCASet, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s", params.CASetID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrGetCASet, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCASet, err)
	}

	var result GetCASetResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCASet, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) ListCASets(ctx context.Context, params ListCASetsRequest) (*ListCASetsResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("ListCASets")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListCASets, ErrStructValidation, err)
	}

	uri, err := url.Parse("/mtls-edge-truststore/v2/ca-sets")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListCASets, err)
	}
	q := uri.Query()
	if params.CASetNamePrefix != "" {
		q.Add("caSetNamePrefix", params.CASetNamePrefix)
	}
	if params.ActivatedOn != "" {
		q.Add("activatedOn", string(params.ActivatedOn))
	}
	if len(params.CASetStatuses) > 0 {
		q.Add("caSetStatus", texts.JoinStringBased(params.CASetStatuses, ","))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCASets, err)
	}

	var result ListCASetsResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCASets, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) DeleteCASet(ctx context.Context, params DeleteCASetRequest) error {
	logger := m.Log(ctx)
	logger.Debug("DeleteCASet")

	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: %w: %s", ErrDeleteCASet, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s", params.CASetID))
	if err != nil {
		return fmt.Errorf("%w: failed to parse url: %s", ErrDeleteCASet, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return fmt.Errorf("%w: failed to create request: %s", ErrDeleteCASet, err)
	}

	resp, err := m.Exec(req, nil)
	if err != nil {
		return fmt.Errorf("%w: request failed: %s", ErrDeleteCASet, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusAccepted {
		return m.Error(resp)
	}

	return nil
}

func (m *mtlstruststore) ListCASetAssociations(ctx context.Context, params ListCASetAssociationsRequest) (*ListCASetAssociationsResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("ListCASetAssociations")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListCASetAssociations, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s/associations", params.CASetID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListCASets, err)
	}
	q := uri.Query()
	if params.AssociationType != "" {
		q.Add("associationType", string(params.AssociationType))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCASetAssociations, err)
	}

	var result ListCASetAssociationsResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCASetAssociations, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) CloneCASet(ctx context.Context, params CloneCASetRequest) (*CloneCASetResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("CloneCASet")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrCloneCASet, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s/clone", params.CloneFromSetID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrCloneCASet, err)
	}
	q := uri.Query()
	if params.CloneFromVersion != 0 {
		q.Add("cloneFromVersion", strconv.FormatInt(params.CloneFromVersion, 10))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrCloneCASet, err)
	}

	var result CloneCASetResponse
	resp, err := m.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrCloneCASet, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusCreated {
		return nil, m.Error(resp)
	}

	return &result, nil
}

func (m *mtlstruststore) GetCASetDeletionStatus(ctx context.Context, params GetCASetDeletionStatusRequest) (*GetCASetDeletionStatusResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("GetCASetDeletionStatus")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetCASetDeletionStatus, ErrStructValidation, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s/status/delete", params.CASetID), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetCASetDeletionStatus, err)
	}

	var result GetCASetDeletionStatusResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetCASetDeletionStatus, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusAccepted &&
		resp.StatusCode != http.StatusMultiStatus {
		return nil, m.Error(resp)
	}

	if resp.Header.Get("Retry-After") != "" {
		after, err := time.Parse(time.RFC1123, resp.Header.Get("Retry-After"))
		if err != nil {
			return nil, fmt.Errorf("%w: failed to parse Retry-After header: %s",
				ErrGetCASetDeletionStatus, err)
		}
		result.RetryAfter = after
	}

	return &result, nil
}

func (m *mtlstruststore) ListCASetActivities(ctx context.Context, params ListCASetActivitiesRequest) (*ListCASetActivitiesResponse, error) {
	logger := m.Log(ctx)
	logger.Debug("ListCASetActivities")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrListCASetActivities, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/mtls-edge-truststore/v2/ca-sets/%s/activities", params.CASetID))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse url: %s", ErrListCASetActivities, err)
	}
	q := uri.Query()
	if !params.Start.IsZero() {
		q.Add("start", params.Start.Format(time.RFC3339Nano))
	}
	if !params.End.IsZero() {
		q.Add("end", params.End.Format(time.RFC3339Nano))
	}
	uri.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListCASetActivities, err)
	}

	var result ListCASetActivitiesResponse
	resp, err := m.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListCASetActivities, err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, m.Error(resp)
	}

	return &result, nil
}
