package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"time"
)

type (
	// The ExportConfiguration interface supports exporting comprehensive details about a security
	// configuration version. This operation returns more data than Get configuration version details,
	// including rate and security policies, rules, hostnames, and numerous additional settings.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#export
	ExportConfiguration interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getconfigurationversionexport
		GetExportConfigurations(ctx context.Context, params GetExportConfigurationsRequest) (*GetExportConfigurationsResponse, error)
	}

	// ConditionsValue is a slice of strings that describe conditions.
	ConditionsValue []string

	// GetExportConfigurationsRequest is used to call GetExportConfigurations.
	GetExportConfigurationsRequest struct {
		ConfigID int `json:"configId"`
		Version  int `json:"version"`
	}

	// GetExportConfigurationsResponse is returned from a call to GetExportConfigurations.
	GetExportConfigurationsResponse struct {
		ConfigID   int    `json:"configId"`
		ConfigName string `json:"configName"`
		Version    int    `json:"version"`
		BasedOn    int    `json:"basedOn"`
		Staging    struct {
			Status string `json:"status"`
		} `json:"staging"`
		Production struct {
			Status string `json:"status"`
		} `json:"production"`
		CreateDate      time.Time `json:"-"`
		CreatedBy       string    `json:"createdBy"`
		SelectedHosts   []string  `json:"selectedHosts"`
		SelectableHosts []string  `json:"selectableHosts"`
		RatePolicies    []struct {
			AdditionalMatchOptions []struct {
				PositiveMatch bool     `json:"positiveMatch"`
				Type          string   `json:"type"`
				Values        []string `json:"values"`
			} `json:"additionalMatchOptions"`
			AllTraffic            bool                         `json:"allTraffic,omitempty"`
			AverageThreshold      int                          `json:"averageThreshold"`
			BurstThreshold        int                          `json:"burstThreshold"`
			ClientIdentifier      string                       `json:"clientIdentifier"`
			CreateDate            time.Time                    `json:"-"`
			Description           string                       `json:"description,omitempty"`
			FileExtensions        *RatePolicyFileExtensions    `json:"fileExtensions,omitempty"`
			Hostnames             []string                     `json:"hostnames,omitempty"`
			ID                    int                          `json:"id"`
			MatchType             string                       `json:"matchType"`
			Name                  string                       `json:"name"`
			Path                  *RatePoliciesPath            `json:"path,omitempty"`
			PathMatchType         string                       `json:"pathMatchType,omitempty"`
			PathURIPositiveMatch  bool                         `json:"pathUriPositiveMatch"`
			QueryParameters       *RatePoliciesQueryParameters `json:"queryParameters,omitempty"`
			RequestType           string                       `json:"requestType"`
			SameActionOnIpv6      bool                         `json:"sameActionOnIpv6"`
			Type                  string                       `json:"type"`
			UpdateDate            time.Time                    `json:"-"`
			UseXForwardForHeaders bool                         `json:"useXForwardForHeaders"`
			Used                  bool                         `json:"-"`
		} `json:"ratePolicies"`
		ReputationProfiles []struct {
			Condition        *ConditionReputationProfile `json:"condition,omitempty"`
			Context          string                      `json:"context,omitempty"`
			ContextReadable  string                      `json:"-"`
			Enabled          bool                        `json:"-"`
			ID               int                         `json:"id"`
			Name             string                      `json:"name"`
			SharedIPHandling string                      `json:"sharedIpHandling"`
			Threshold        int                         `json:"threshold"`
		} `json:"reputationProfiles"`
		CustomRules []struct {
			Conditions    *ConditionsExp `json:"conditions,omitempty"`
			Description   string         `json:"description,omitempty"`
			ID            int            `json:"id"`
			Name          string         `json:"name"`
			RuleActivated bool           `json:"-"`
			Structured    bool           `json:"-"`
			Tag           []string       `json:"tag"`
			Version       int            `json:"-"`
		} `json:"customRules"`
		Rulesets []struct {
			ID               int            `json:"id"`
			RulesetVersionID int            `json:"rulesetVersionId"`
			Type             string         `json:"type"`
			ReleaseDate      time.Time      `json:"releaseDate"`
			Rules            *RulesetsRules `json:"rules,omitempty"`
			AttackGroups     []struct {
				Group     string `json:"group"`
				GroupName string `json:"groupName"`
				Threshold int    `json:"threshold,omitempty"`
			} `json:"attackGroups,omitempty"`
		} `json:"rulesets"`
		MatchTargets struct {
			APITargets []struct {
				Sequence int    `json:"-"`
				ID       int    `json:"id,omitempty"`
				Type     string `json:"type,omitempty"`
				Apis     []struct {
					ID   int    `json:"id,omitempty"`
					Name string `json:"name,omitempty"`
				} `json:"apis,omitempty"`
				SecurityPolicy struct {
					PolicyID string `json:"policyId,omitempty"`
				} `json:"securityPolicy,omitempty"`
				BypassNetworkLists []struct {
					Name string `json:"name,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"bypassNetworkLists,omitempty"`
			} `json:"apiTargets,omitempty"`
			WebsiteTargets []struct {
				Type               string `json:"type"`
				BypassNetworkLists []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"bypassNetworkLists,omitempty"`
				DefaultFile                  string   `json:"defaultFile"`
				FilePaths                    []string `json:"filePaths,omitempty"`
				FileExtensions               []string `json:"fileExtensions,omitempty"`
				Hostnames                    []string `json:"hostnames,omitempty"`
				ID                           int      `json:"id"`
				IsNegativeFileExtensionMatch bool     `json:"isNegativeFileExtensionMatch"`
				IsNegativePathMatch          bool     `json:"isNegativePathMatch"`
				SecurityPolicy               struct {
					PolicyID string `json:"policyId"`
				} `json:"securityPolicy"`
				Sequence int `json:"-"`
			} `json:"websiteTargets"`
		} `json:"matchTargets"`
		SecurityPolicies []struct {
			ID                      string `json:"id"`
			Name                    string `json:"name"`
			HasRatePolicyWithAPIKey bool   `json:"hasRatePolicyWithApiKey"`
			SecurityControls        struct {
				ApplyAPIConstraints           bool `json:"applyApiConstraints"`
				ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls"`
				ApplyBotmanControls           bool `json:"applyBotmanControls"`
				ApplyNetworkLayerControls     bool `json:"applyNetworkLayerControls"`
				ApplyRateControls             bool `json:"applyRateControls"`
				ApplyReputationControls       bool `json:"applyReputationControls"`
				ApplySlowPostControls         bool `json:"applySlowPostControls"`
			} `json:"securityControls"`
			WebApplicationFirewall struct {
				RuleActions []struct {
					Action           string          `json:"action"`
					ID               int             `json:"id"`
					RulesetVersionID int             `json:"rulesetVersionId"`
					Conditions       *RuleConditions `json:"conditions,omitempty"`
					Exception        *RuleException  `json:"exception,omitempty"`
				} `json:"ruleActions,omitempty"`
				AttackGroupActions []struct {
					Action                 string                         `json:"action"`
					Group                  string                         `json:"group"`
					RulesetVersionID       int                            `json:"rulesetVersionId"`
					AdvancedExceptionsList *AttackGroupAdvancedExceptions `json:"advancedExceptions,omitempty"`
					Exception              *AttackGroupException          `json:"exception,omitempty"`
				} `json:"attackGroupActions,omitempty"`
				Evaluation  *WebApplicationFirewallEvaluation `json:"evaluation,omitempty"`
				ThreatIntel string                            `json:"threatIntel"`
			} `json:"webApplicationFirewall"`
			CustomRuleActions []struct {
				Action string `json:"action"`
				ID     int    `json:"id"`
			} `json:"customRuleActions,omitempty"`
			APIRequestConstraints *APIRequestConstraintsexp `json:"apiRequestConstraints,omitempty"`
			ClientReputation      struct {
				ReputationProfileActions *ClientReputationReputationProfileActions `json:"reputationProfileActions,omitempty"`
			} `json:"clientReputation"`
			RatePolicyActions *SecurityPoliciesRatePolicyActions `json:"ratePolicyActions,omitempty"`
			IPGeoFirewall     struct {
				Block       string `json:"block"`
				GeoControls struct {
					BlockedIPNetworkLists struct {
						NetworkList []string `json:"networkList,omitempty"`
					} `json:"blockedIPNetworkLists"`
				} `json:"geoControls"`
				IPControls struct {
					AllowedIPNetworkLists struct {
						NetworkList []string `json:"networkList,omitempty"`
					} `json:"allowedIPNetworkLists"`
					BlockedIPNetworkLists struct {
						NetworkList []string `json:"networkList,omitempty"`
					} `json:"blockedIPNetworkLists"`
				} `json:"ipControls"`
			} `json:"ipGeoFirewall,omitempty"`
			PenaltyBox       *SecurityPoliciesPenaltyBox        `json:"penaltyBox,omitempty"`
			SlowPost         *SlowPostexp                       `json:"slowPost,omitempty"`
			LoggingOverrides *LoggingOverridesexp               `json:"loggingOverrides,omitempty"`
			PragmaHeader     *GetAdvancedSettingsPragmaResponse `json:"pragmaHeader,omitempty"`
		} `json:"securityPolicies"`
		Siem            *Siemexp            `json:"siem,omitempty"`
		AdvancedOptions *AdvancedOptionsexp `json:"advancedOptions,omitempty"`
		CustomDenyList  *CustomDenyListexp  `json:"customDenyList,omitempty"`
		Evaluating      struct {
			SecurityPolicies []struct {
				EffectiveSecurityControls struct {
					ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
					ApplyRateControls             bool `json:"applyRateControls,omitempty"`
					ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
				}
				Hostnames        []string `json:"hostnames,omitempty"`
				SecurityPolicyID string   `json:"id"`
			}
		} `json:"evaluating"`
	}

	// RatePoliciesPath is returned as part of GetExportConfigurationsResponse.
	RatePoliciesPath struct {
		PositiveMatch bool                    `json:"positiveMatch"`
		Values        *RatePoliciesPathValues `json:"values,omitempty"`
	}

	// ReputationProfileActionsexp is returned as part of GetExportConfigurationsResponse.
	ReputationProfileActionsexp []struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}

	// RatePolicyActionsexp is returned as part of GetExportConfigurationsResponse.
	RatePolicyActionsexp []struct {
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}

	// SlowRateThresholdExp is returned as part of GetExportConfigurationsResponse.
	SlowRateThresholdExp struct {
		Period int `json:"period"`
		Rate   int `json:"rate"`
	}

	// DurationThresholdExp is returned as part of GetExportConfigurationsResponse.
	DurationThresholdExp struct {
		Timeout int `json:"timeout"`
	}

	// SlowPostexp is returned as part of GetExportConfigurationsResponse.
	SlowPostexp struct {
		Action            string                `json:"action"`
		SlowRateThreshold *SlowRateThresholdExp `json:"slowRateThreshold,omitempty"`
		DurationThreshold *DurationThresholdExp `json:"durationThreshold,omitempty"`
	}

	// AdvancedOptionsexp is returned as part of GetExportConfigurationsResponse.
	AdvancedOptionsexp struct {
		Logging  *Loggingexp `json:"logging"`
		Prefetch struct {
			AllExtensions      bool     `json:"allExtensions"`
			EnableAppLayer     bool     `json:"enableAppLayer"`
			EnableRateControls bool     `json:"enableRateControls"`
			Extensions         []string `json:"extensions"`
		} `json:"prefetch"`
		PragmaHeader *GetAdvancedSettingsPragmaResponse `json:"pragmaHeader,omitempty"`
	}

	// CustomDenyListexp is returned as part of GetExportConfigurationsResponse.
	CustomDenyListexp []struct {
		Description string `json:"description,omitempty"`
		Name        string `json:"name"`
		ID          string `json:"id"`
		Parameters  []struct {
			DisplayName string `json:"-"`
			Name        string `json:"name"`
			Value       string `json:"value"`
		} `json:"parameters"`
	}

	// CustomRuleActionsexp  is returned as part of GetExportConfigurationsResponse.
	CustomRuleActionsexp []struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}

	// Siemexp is returned as part of GetExportConfigurationsResponse.
	Siemexp struct {
		EnableForAllPolicies    bool     `json:"enableForAllPolicies,omitempty"`
		EnableSiem              bool     `json:"enableSiem"`
		EnabledBotmanSiemEvents bool     `json:"enabledBotmanSiemEvents,omitempty"`
		FirewallPolicyIds       []string `json:"firewallPolicyIds,omitempty"`
		SiemDefinitionID        int      `json:"siemDefinitionId,omitempty"`
	}

	// PenaltyBoxexp is returned as part of GetExportConfigurationsResponse.
	PenaltyBoxexp struct {
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}

	// APIRequestConstraintsexp is returned as part of GetExportConfigurationsResponse.
	APIRequestConstraintsexp struct {
		Action       string `json:"action,omitempty"`
		APIEndpoints []struct {
			Action string `json:"action"`
			ID     int    `json:"id"`
		} `json:"apiEndpoints,omitempty"`
	}

	// Evaluationexp is returned as part of GetExportConfigurationsResponse.
	Evaluationexp struct {
		AttackGroupActions []struct {
			Action string `json:"action"`
			Group  string `json:"group"`
		} `json:"attackGroupActions"`
		EvaluationID      int `json:"evaluationId"`
		EvaluationVersion int `json:"evaluationVersion"`
		RuleActions       []struct {
			Action     string          `json:"action"`
			ID         int             `json:"id"`
			Conditions *RuleConditions `json:"conditions,omitempty"`
			Exception  *RuleException  `json:"exception,omitempty"`
		} `json:"ruleActions"`
		RulesetVersionID int `json:"rulesetVersionId"`
	}

	// ConditionReputationProfile is returned as part of GetExportConfigurationsResponse.
	ConditionReputationProfile struct {
		AtomicConditions *AtomicConditionsexp `json:"atomicConditions,omitempty"`
		CanDelete        bool                 `json:"-"`
		ConfigVersionID  int                  `json:"-"`
		ID               int                  `json:"-"`
		Name             string               `json:"-"`
		PositiveMatch    *json.RawMessage     `json:"positiveMatch,omitempty"`
		UUID             string               `json:"-"`
		Version          int64                `json:"-"`
	}

	// HeaderCookieOrParamValuesattackgroup is returned as part of GetExportConfigurationsResponse.
	HeaderCookieOrParamValuesattackgroup []struct {
		Criteria []struct {
			Hostnames []string `json:"hostnames,omitempty"`
			Paths     []string `json:"paths,omitempty"`
			Values    []string `json:"values,omitempty"`
		} `json:"criteria"`
		ValueWildcard bool     `json:"valueWildcard,omitempty"`
		Values        []string `json:"values,omitempty"`
	}

	// SpecificHeaderCookieOrParamNameValueexp is returned as part of GetExportConfigurationsResponse.
	SpecificHeaderCookieOrParamNameValueexp struct {
		Name     *json.RawMessage `json:"name,omitempty"`
		Selector string           `json:"selector,omitempty"`
		Value    *json.RawMessage `json:"value,omitempty"`
	}

	// AtomicConditionsexp is returned as part of GetExportConfigurationsResponse.
	AtomicConditionsexp []struct {
		CheckIps      *json.RawMessage `json:"checkIps,omitempty"`
		ClassName     string           `json:"className,omitempty"`
		Index         int              `json:"index,omitempty"`
		PositiveMatch *json.RawMessage `json:"positiveMatch,omitempty"`
		Value         []string         `json:"value,omitempty"`
		Name          *json.RawMessage `json:"name,omitempty"`
		NameCase      bool             `json:"nameCase,omitempty"`
		NameWildcard  *json.RawMessage `json:"nameWildcard,omitempty"`
		ValueCase     bool             `json:"valueCase,omitempty"`
		ValueWildcard *json.RawMessage `json:"valueWildcard,omitempty"`
		Host          []string         `json:"host,omitempty"`
	}

	// Loggingexp is returned as part of GetExportConfigurationsResponse.
	Loggingexp struct {
		AllowSampling bool `json:"allowSampling"`
		Cookies       struct {
			Type   string   `json:"type"`
			Values []string `json:"values,omitempty"`
		} `json:"cookies"`
		CustomHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values,omitempty"`
		} `json:"customHeaders"`
		StandardHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values,omitempty"`
		} `json:"standardHeaders"`
	}

	// LoggingOverridesexp is returned as part of GetExportConfigurationsResponse.
	LoggingOverridesexp struct {
		AllowSampling bool `json:"allowSampling"`
		Cookies       struct {
			Type   string   `json:"type"`
			Values []string `json:"values,omitempty"`
		} `json:"cookies"`
		CustomHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values,omitempty"`
		} `json:"customHeaders"`
		Override        bool `json:"override"`
		StandardHeaders struct {
			Type   string   `json:"type"`
			Values []string `json:"values,omitempty"`
		} `json:"standardHeaders"`
	}

	// ConditionsExp is returned as part of GetExportConfigurationsResponse.
	ConditionsExp []struct {
		Type          string           `json:"type"`
		PositiveMatch bool             `json:"positiveMatch"`
		Name          *json.RawMessage `json:"name,omitempty"`
		NameCase      *json.RawMessage `json:"nameCase,omitempty"`
		NameWildcard  *json.RawMessage `json:"nameWildcard,omitempty"`
		Value         *json.RawMessage `json:"value,omitempty"`
		ValueCase     *json.RawMessage `json:"valueCase,omitempty"`
		ValueWildcard *json.RawMessage `json:"valueWildcard,omitempty"`
	}

	// RatePoliciesPathValues is returned as part of GetExportConfigurationsResponse.
	RatePoliciesPathValues []string

	// RatePoliciesQueryParameters is returned as part of GetExportConfigurationsResponse.
	RatePoliciesQueryParameters []struct {
		Name          string                             `json:"name"`
		PositiveMatch bool                               `json:"positiveMatch"`
		ValueInRange  bool                               `json:"valueInRange"`
		Values        *RatePoliciesQueryParametersValues `json:"values,omitempty"`
	}

	// RatePoliciesQueryParametersValues is returned as part of GetExportConfigurationsResponse.
	RatePoliciesQueryParametersValues []string

	// SecurityPoliciesPenaltyBox is returned as part of GetExportConfigurationsResponse.
	SecurityPoliciesPenaltyBox struct {
		Action               string `json:"action,omitempty"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection,omitempty"`
	}

	// WebApplicationFirewallEvaluation is returned as part of GetExportConfigurationsResponse.
	WebApplicationFirewallEvaluation struct {
		AttackGroupActions []struct {
			Action string `json:"action"`
			Group  string `json:"group"`
		} `json:"attackGroupActions,omitempty"`
		EvaluationID      int `json:"evaluationId"`
		EvaluationVersion int `json:"evaluationVersion"`
		RuleActions       []struct {
			Action     string          `json:"action"`
			ID         int             `json:"id"`
			Conditions *RuleConditions `json:"conditions,omitempty"`
			Exception  *RuleException  `json:"exception,omitempty"`
		} `json:"ruleActions,omitempty"`
		RulesetVersionID int `json:"rulesetVersionId"`
	}

	// RulesetsRules is returned as part of GetExportConfigurationsResponse.
	RulesetsRules []struct {
		ID                  int      `json:"id"`
		InspectRequestBody  bool     `json:"inspectRequestBody"`
		InspectResponseBody bool     `json:"inspectResponseBody"`
		Outdated            bool     `json:"outdated"`
		RuleVersion         int      `json:"ruleVersion"`
		Score               int      `json:"score"`
		Tag                 string   `json:"tag"`
		Title               string   `json:"title"`
		AttackGroups        []string `json:"attackGroups,omitempty"`
	}

	// ClientReputationReputationProfileActions is returned as part of GetExportConfigurationsResponse.
	ClientReputationReputationProfileActions []struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}

	// SecurityPoliciesRatePolicyActions is returned as part of GetExportConfigurationsResponse.
	SecurityPoliciesRatePolicyActions []struct {
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}
)

// UnmarshalJSON reads a ConditionsValue struct from its data argument.
func (c *ConditionsValue) UnmarshalJSON(data []byte) error {
	var nums interface{}
	err := json.Unmarshal(data, &nums)
	if err != nil {
		return err
	}

	items := reflect.ValueOf(nums)
	switch items.Kind() {
	case reflect.String:
		*c = append(*c, items.String())

	case reflect.Slice:
		*c = make(ConditionsValue, 0, items.Len())
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			switch item.Kind() {
			case reflect.String:
				*c = append(*c, item.String())
			case reflect.Interface:
				*c = append(*c, item.Interface().(string))
			}
		}
	}
	return nil
}

func (p *appsec) GetExportConfigurations(ctx context.Context, params GetExportConfigurationsRequest) (*GetExportConfigurationsResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("GetExportConfigurations")

	var rval GetExportConfigurationsResponse

	uri := fmt.Sprintf(
		"/appsec/v1/export/configs/%d/versions/%d",
		params.ConfigID,
		params.Version)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetExportConfigurations request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("GetExportConfigurations request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}
