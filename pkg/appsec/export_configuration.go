package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"time"
)

// ExportConfiguration represents a collection of ExportConfiguration
//
// See: ExportConfiguration.GetExportConfiguration()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ExportConfiguration  contains operations available on ExportConfiguration  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getexportconfiguration
	ExportConfiguration interface {
		GetExportConfigurations(ctx context.Context, params GetExportConfigurationsRequest) (*GetExportConfigurationsResponse, error)
	}

	ConditionsValue []string

	GetExportConfigurationsRequest struct {
		ConfigID int `json:"configId"`
		Version  int `json:"version"`
	}

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
				Evaluation *WebApplicationFirewallEvaluation `json:"evaluation,omitempty"`
				ThreatIntel string `json:"threatIntel"`
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
			PenaltyBox       *SecurityPoliciesPenaltyBox `json:"penaltyBox,omitempty"`
			SlowPost         *SlowPostexp                `json:"slowPost,omitempty"`
			LoggingOverrides *LoggingOverridesexp        `json:"loggingOverrides,omitempty"`
			PragmaHeader     *GetAdvancedSettingsPragmaResponse `json:"pragmaHeader,omitempty"`
		} `json:"securityPolicies"`
		Siem            *Siemexp            `json:"siem,omitempty"`
		AdvancedOptions *AdvancedOptionsexp `json:"advancedOptions,omitempty"`
		CustomDenyList  *CustomDenyListexp  `json:"customDenyList,omitempty"`
	}

	RatePoliciesPath struct {
		PositiveMatch bool                    `json:"positiveMatch"`
		Values        *RatePoliciesPathValues `json:"values,omitempty"`
	}

	//TODO fails export
	ReputationProfileActionsexp []struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}
	//TODO fails export
	RatePolicyActionsexp []struct {
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}
	SlowRateThresholdExp struct {
		Period int `json:"period"`
		Rate   int `json:"rate"`
	}

	DurationThresholdExp struct {
		Timeout int `json:"timeout"`
	}

	SlowPostexp struct {
		Action            string                `json:"action"`
		SlowRateThreshold *SlowRateThresholdExp `json:"slowRateThreshold,omitempty"`
		DurationThreshold *DurationThresholdExp `json:"durationThreshold,omitempty"`
	}
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
	//TODO breaks export
	CustomRuleActionsexp []struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}
	Siemexp struct {
		EnableForAllPolicies    bool     `json:"enableForAllPolicies,omitempty"`
		EnableSiem              bool     `json:"enableSiem"`
		EnabledBotmanSiemEvents bool     `json:"enabledBotmanSiemEvents,omitempty"`
		FirewallPolicyIds       []string `json:"firewallPolicyIds,omitempty"`
		SiemDefinitionID        int      `json:"siemDefinitionId,omitempty"`
	}
	//TODO nfails export
	PenaltyBoxexp struct {
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}
	APIRequestConstraintsexp struct {
		Action       string `json:"action,omitempty"`
		APIEndpoints []struct {
			Action string `json:"action"`
			ID     int    `json:"id"`
		} `json:"apiEndpoints,omitempty"`
	}
	//TODO breaks export
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

	HeaderCookieOrParamValuesattackgroup []struct {
		Criteria []struct {
			Hostnames []string `json:"hostnames,omitempty"`
			Paths     []string `json:"paths,omitempty"`
			Values    []string `json:"values,omitempty"`
		} `json:"criteria"`
		ValueWildcard bool     `json:"valueWildcard,omitempty"`
		Values        []string `json:"values,omitempty"`
	}

	SpecificHeaderCookieOrParamNameValueexp struct {
		Name     *json.RawMessage `json:"name,omitempty"`
		Selector string           `json:"selector,omitempty"`
		Value    *json.RawMessage `json:"value,omitempty"`
	}

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

	RatePoliciesFileExtensionsValues []string

	RatePoliciesPathValues []string

	RatePoliciesQueryParameters []struct {
		Name          string                             `json:"name"`
		PositiveMatch bool                               `json:"positiveMatch"`
		ValueInRange  bool                               `json:"valueInRange"`
		Values        *RatePoliciesQueryParametersValues `json:"values,omitempty"`
	}

	RatePoliciesQueryParametersValues []string

	SecurityPoliciesPenaltyBox struct {
		Action               string `json:"action,omitempty"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection,omitempty"`
	}

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
	ClientReputationReputationProfileActions []struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}
	SecurityPoliciesRatePolicyActions []struct {
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}
)

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
		return nil, fmt.Errorf("failed to create getexportconfigurations request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getexportconfigurations request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}
