package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgegriderr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The ExportConfiguration interface supports exporting comprehensive details about a security
	// configuration version. This operation returns more data than Get configuration version details,
	// including rate and security policies, rules, hostnames, and numerous additional settings.
	ExportConfiguration interface {
		// GetExportConfiguration returns comprehensive details about a security configuration version.
		//
		// See: https://techdocs.akamai.com/application-security/reference/get-export-config-version
		GetExportConfiguration(ctx context.Context, params GetExportConfigurationRequest) (*GetExportConfigurationResponse, error)
	}

	// ConditionsValue is a slice of strings that describe conditions.
	ConditionsValue []string

	// RapidRules is returned as part of GetExportConfigurationResponse.
	RapidRules struct {
		DefaultAction string                 `json:"defaultAction"`
		Enabled       bool                   `json:"enabled"`
		ThreatIntel   string                 `json:"threatIntel"`
		PolicyRules   []RapidRulesPolicyRule `json:"policyRules,omitempty"`
	}

	// RapidRulesPolicyRule represents an individual policy rule in RapidRules.
	RapidRulesPolicyRule struct {
		Action             string                  `json:"action"`
		Lock               bool                    `json:"lock"`
		RuleID             int                     `json:"ruleId"`
		RuleVersion        int                     `json:"ruleVersion"`
		Group              string                  `json:"group"`
		RulesetVersionID   int                     `json:"rulesetVersionId"`
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}

	// GetExportConfigurationRequest is used to call GetExportConfiguration.
	GetExportConfigurationRequest struct {
		ConfigID int64  `json:"configId"`
		Version  int64  `json:"version"`
		Source   string `json:"source,omitempty"`
	}

	// EvaluatingSecurityPolicy is returned from a call to GetExportConfiguration.
	EvaluatingSecurityPolicy struct {
		EffectiveSecurityControls EvaluatingSecurityPolicyEffectiveSecurityControls
		Hostnames                 []string `json:"hostnames,omitempty"`
		SecurityPolicyID          string   `json:"id"`
	}

	// EvaluatingSecurityPolicyEffectiveSecurityControls describes effective controls for evaluating policy.
	EvaluatingSecurityPolicyEffectiveSecurityControls struct {
		ApplyApplicationLayerControls bool `json:"applyApplicationLayerControls,omitempty"`
		ApplyRateControls             bool `json:"applyRateControls,omitempty"`
		ApplySlowPostControls         bool `json:"applySlowPostControls,omitempty"`
	}

	// GetExportConfigurationResponse is returned from a call to GetExportConfiguration.
	GetExportConfigurationResponse struct {
		ConfigID                  int64                                             `json:"configId"`
		ConfigName                string                                            `json:"configName"`
		Version                   int64                                             `json:"version"`
		BasedOn                   int                                               `json:"basedOn"`
		Staging                   GetExportConfigurationResponseStatus              `json:"staging"`
		Production                GetExportConfigurationResponseStatus              `json:"production"`
		TargetProduct             string                                            `json:"targetProduct"`
		CreateDate                time.Time                                         `json:"-"`
		CreatedBy                 string                                            `json:"createdBy"`
		SelectedHosts             []string                                          `json:"selectedHosts"`
		SelectableHosts           []string                                          `json:"selectableHosts"`
		RatePolicies              []GetExportConfigurationResponseRatePolicy        `json:"ratePolicies"`
		ReputationProfiles        []GetExportConfigurationResponseReputationProfile `json:"reputationProfiles"`
		CustomRules               []GetExportConfigurationResponseCustomRule        `json:"customRules"`
		Rulesets                  []GetExportConfigurationResponseRuleset           `json:"rulesets"`
		MatchTargets              GetExportConfigurationResponseMatchTargets        `json:"matchTargets"`
		SecurityPolicies          []GetExportConfigurationResponseSecurityPolicy    `json:"securityPolicies"`
		Siem                      *Siemexp                                          `json:"siem,omitempty"`
		AdvancedOptions           *AdvancedOptionsexp                               `json:"advancedOptions,omitempty"`
		CustomDenyList            *CustomDenyListexp                                `json:"customDenyList,omitempty"`
		Evaluating                GetExportConfigurationResponseEvaluating          `json:"evaluating,omitempty"`
		MalwarePolicies           []MalwarePolicyBody                               `json:"malwarePolicies,omitempty"`
		CustomBotCategories       []map[string]interface{}                          `json:"customBotCategories,omitempty"`
		CustomDefinedBots         []map[string]interface{}                          `json:"customDefinedBots,omitempty"`
		CustomBotCategorySequence []string                                          `json:"customBotCategorySequence,omitempty"`
		CustomClients             []map[string]interface{}                          `json:"customClients,omitempty"`
		CustomClientSequence      []string                                          `json:"customClientSequence,omitempty"`
		ResponseActions           *ResponseActions                                  `json:"responseActions,omitempty"`
		AdvancedSettings          *AdvancedSettings                                 `json:"advancedSettings,omitempty"`
		URLProtectionPolicies     []GetURLProtectionPolicyResponse                  `json:"urlProtectionPolicies,omitempty"`
	}

	// GetExportConfigurationResponseStatus represents the status of staging or production environment.
	GetExportConfigurationResponseStatus struct {
		Status string `json:"status"`
	}

	// GetExportConfigurationResponseRatePolicy represents an individual rate policy item.
	GetExportConfigurationResponseRatePolicy struct {
		AdditionalMatchOptions []AdditionalMatchOption      `json:"additionalMatchOptions,omitempty"`
		AllTraffic             bool                         `json:"allTraffic,omitempty"`
		AverageThreshold       int                          `json:"averageThreshold"`
		BurstThreshold         int                          `json:"burstThreshold"`
		BurstWindow            *int                         `json:"burstWindow,omitempty"`
		ClientIdentifiers      []string                     `json:"clientIdentifiers,omitempty"`
		Condition              *RatePolicyCondition         `json:"condition,omitempty"`
		CreateDate             time.Time                    `json:"-"`
		CounterType            string                       `json:"counterType"`
		Description            string                       `json:"description,omitempty"`
		FileExtensions         *RatePolicyFileExtensions    `json:"fileExtensions,omitempty"`
		Hosts                  *RatePoliciesHosts           `json:"hosts,omitempty"`
		Hostnames              []string                     `json:"hostnames,omitempty"`
		ID                     int                          `json:"id"`
		MatchType              string                       `json:"matchType"`
		Name                   string                       `json:"name"`
		Path                   *RatePoliciesPath            `json:"path,omitempty"`
		PathMatchType          string                       `json:"pathMatchType,omitempty"`
		PathURIPositiveMatch   bool                         `json:"pathUriPositiveMatch"`
		PenaltyBoxDuration     string                       `json:"penaltyBoxDuration"`
		QueryParameters        *RatePoliciesQueryParameters `json:"queryParameters,omitempty"`
		RequestType            string                       `json:"requestType"`
		SameActionOnIpv6       bool                         `json:"sameActionOnIpv6"`
		Type                   string                       `json:"type"`
		UpdateDate             time.Time                    `json:"-"`
		UseXForwardForHeaders  bool                         `json:"useXForwardForHeaders"`
		Used                   bool                         `json:"-"`
	}

	// AdditionalMatchOption represents an additional match option.
	AdditionalMatchOption struct {
		PositiveMatch bool     `json:"positiveMatch"`
		Type          string   `json:"type"`
		Values        []string `json:"values"`
	}

	// GetExportConfigurationResponseReputationProfile represents a reputation profile item.
	GetExportConfigurationResponseReputationProfile struct {
		Condition        *ConditionReputationProfile `json:"condition,omitempty"`
		Context          string                      `json:"context,omitempty"`
		ContextReadable  string                      `json:"-"`
		Enabled          bool                        `json:"-"`
		ID               int                         `json:"id"`
		Name             string                      `json:"name"`
		SharedIPHandling string                      `json:"sharedIpHandling"`
		Threshold        int                         `json:"threshold"`
	}

	// GetExportConfigurationResponseCustomRule represents an individual custom rule item.
	GetExportConfigurationResponseCustomRule struct {
		ID                  int                        `json:"id"`
		Name                string                     `json:"name"`
		Description         string                     `json:"description,omitempty"`
		Version             int                        `json:"-"`
		RuleActivated       bool                       `json:"-"`
		Structured          bool                       `json:"-"`
		Tag                 []string                   `json:"tag,omitempty"`
		Conditions          []CustomRuleCondition      `json:"conditions,omitempty"`
		EffectiveTimePeriod *CustomRuleEffectivePeriod `json:"effectiveTimePeriod,omitempty"`
		SamplingRate        int                        `json:"samplingRate,omitempty"`
		LoggingOptions      *json.RawMessage           `json:"loggingOptions,omitempty"`
		Operation           string                     `json:"operation,omitempty"`
	}

	// GetExportConfigurationResponseRuleset represents a ruleset item.
	GetExportConfigurationResponseRuleset struct {
		ID               int                                                `json:"id"`
		RulesetVersionID int                                                `json:"rulesetVersionId"`
		Type             string                                             `json:"type"`
		ReleaseDate      time.Time                                          `json:"releaseDate"`
		Rules            *RulesetsRules                                     `json:"rules,omitempty"`
		AttackGroups     []GetExportConfigurationResponseRulesetAttackGroup `json:"attackGroups,omitempty"`
	}

	// GetExportConfigurationResponseRulesetAttackGroup represents an attack group item in a ruleset.
	GetExportConfigurationResponseRulesetAttackGroup struct {
		Group     string `json:"group"`
		GroupName string `json:"groupName"`
		Threshold int    `json:"threshold,omitempty"`
	}

	// GetExportConfigurationResponseMatchTargets represents the match targets section of the export response.
	GetExportConfigurationResponseMatchTargets struct {
		APITargets     []GetExportConfigurationResponseMatchTargetsAPITarget     `json:"apiTargets,omitempty"`
		WebsiteTargets []GetExportConfigurationResponseMatchTargetsWebsiteTarget `json:"websiteTargets"`
	}

	// GetExportConfigurationResponseMatchTargetsAPITarget represents an API target.
	GetExportConfigurationResponseMatchTargetsAPITarget struct {
		Sequence           int                                                                    `json:"sequence"`
		ID                 int                                                                    `json:"id,omitempty"`
		TargetID           int                                                                    `json:"targetId"`
		Type               string                                                                 `json:"type,omitempty"`
		Apis               []GetExportConfigurationResponseMatchTargetsAPITargetAPI               `json:"apis,omitempty"`
		SecurityPolicy     GetExportConfigurationResponseMatchTargetsAPITargetSecurityPolicy      `json:"securityPolicy,omitempty"`
		BypassNetworkLists []GetExportConfigurationResponseMatchTargetsAPITargetBypassNetworkList `json:"bypassNetworkLists,omitempty"`
	}

	// GetExportConfigurationResponseMatchTargetsAPITargetAPI represents an API in an API target.
	GetExportConfigurationResponseMatchTargetsAPITargetAPI struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	// GetExportConfigurationResponseMatchTargetsAPITargetSecurityPolicy represents the security policy in an API target.
	GetExportConfigurationResponseMatchTargetsAPITargetSecurityPolicy struct {
		PolicyID string `json:"policyId,omitempty"`
	}

	// GetExportConfigurationResponseMatchTargetsAPITargetBypassNetworkList represents a bypass network list in an API target.
	GetExportConfigurationResponseMatchTargetsAPITargetBypassNetworkList struct {
		Name string `json:"name,omitempty"`
		ID   string `json:"id,omitempty"`
	}

	// GetExportConfigurationResponseMatchTargetsWebsiteTarget represents a website target.
	GetExportConfigurationResponseMatchTargetsWebsiteTarget struct {
		Type                         string                                                                     `json:"type"`
		BypassNetworkLists           []GetExportConfigurationResponseMatchTargetsWebsiteTargetBypassNetworkList `json:"bypassNetworkLists,omitempty"`
		DefaultFile                  string                                                                     `json:"defaultFile"`
		FilePaths                    []string                                                                   `json:"filePaths,omitempty"`
		FileExtensions               []string                                                                   `json:"fileExtensions,omitempty"`
		Hostnames                    []string                                                                   `json:"hostnames,omitempty"`
		ID                           int                                                                        `json:"id"`
		IsNegativeFileExtensionMatch bool                                                                       `json:"isNegativeFileExtensionMatch"`
		IsNegativePathMatch          bool                                                                       `json:"isNegativePathMatch"`
		SecurityPolicy               GetExportConfigurationResponseMatchTargetsWebsiteTargetSecurityPolicy      `json:"securityPolicy"`
		Sequence                     int                                                                        `json:"-"`
	}

	// GetExportConfigurationResponseMatchTargetsWebsiteTargetBypassNetworkList represents a bypass network list in a website target.
	GetExportConfigurationResponseMatchTargetsWebsiteTargetBypassNetworkList struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	// GetExportConfigurationResponseMatchTargetsWebsiteTargetSecurityPolicy represents the security policy in a website target.
	GetExportConfigurationResponseMatchTargetsWebsiteTargetSecurityPolicy struct {
		PolicyID string `json:"policyId"`
	}

	// GetExportConfigurationResponseSecurityPolicy represents a security policy item.
	GetExportConfigurationResponseSecurityPolicy struct {
		ID                             string                                                             `json:"id"`
		Name                           string                                                             `json:"name"`
		HasRatePolicyWithAPIKey        bool                                                               `json:"hasRatePolicyWithApiKey"`
		SecurityControls               GetExportConfigurationResponseSecurityPolicySecurityControls       `json:"securityControls"`
		WebApplicationFirewall         GetExportConfigurationResponseSecurityPolicyWebApplicationFirewall `json:"webApplicationFirewall"`
		CustomRuleActions              []ExportActionItem                                                 `json:"customRuleActions,omitempty"`
		APIRequestConstraints          *APIRequestConstraintsexp                                          `json:"apiRequestConstraints,omitempty"`
		ClientReputation               GetExportConfigurationResponseSecurityPolicyClientReputation       `json:"clientReputation"`
		RatePolicyActions              *SecurityPoliciesRatePolicyActions                                 `json:"ratePolicyActions,omitempty"`
		MalwarePolicyActions           []MalwarePolicyActionBody                                          `json:"malwarePolicyActions,omitempty"`
		URLProtectionPoliciesActions   []URLProtectionActionResp                                          `json:"urlProtectionPolicyActions,omitempty"`
		IPGeoFirewall                  *IPGeoFirewall                                                     `json:"ipGeoFirewall,omitempty"`
		RapidRules                     *RapidRules                                                        `json:"rapidRules,omitempty"`
		PenaltyBox                     *SecurityPoliciesPenaltyBox                                        `json:"penaltyBox,omitempty"`
		EvaluationPenaltyBox           *SecurityPoliciesPenaltyBox                                        `json:"evaluationPenaltyBox,omitempty"`
		PenaltyBoxConditions           *SecurityPoliciesPenaltyBoxConditions                              `json:"penaltyBoxConditions,omitempty"`
		EvaluationPenaltyBoxConditions *SecurityPoliciesPenaltyBoxConditions                              `json:"evaluationPenaltyBoxConditions,omitempty"`
		SlowPost                       *SlowPostexp                                                       `json:"slowPost,omitempty"`
		AccountProtection              *AccountProtection                                                 `json:"accountProtection,omitempty"`
		LoggingOverrides               *LoggingOverridesexp                                               `json:"loggingOverrides,omitempty"`
		AttackPayloadLoggingOverrides  *AttackPayloadLoggingOverrides                                     `json:"attackPayloadLoggingOverrides,omitempty"`
		PragmaHeader                   *GetAdvancedSettingsPragmaResponse                                 `json:"pragmaHeader,omitempty"`
		EvasivePathMatch               *EvasivePathMatchexp                                               `json:"evasivePathMatch,omitempty"`
		RequestBody                    *RequestBody                                                       `json:"requestBody,omitempty"`
		BotManagement                  *BotManagement                                                     `json:"botManagement,omitempty"`
		AIRules                        *ExportAIRulesResponse                                             `json:"aiRules,omitempty"`
	}

	// GetExportConfigurationResponseSecurityPolicySecurityControls describes security controls for a policy.
	GetExportConfigurationResponseSecurityPolicySecurityControls struct {
		ApplyAccountProtectionControls bool `json:"applyAccountProtectionControls"`
		ApplyAPIConstraints            bool `json:"applyApiConstraints"`
		ApplyApplicationLayerControls  bool `json:"applyApplicationLayerControls"`
		ApplyBotmanControls            bool `json:"applyBotmanControls"`
		ApplyNetworkLayerControls      bool `json:"applyNetworkLayerControls"`
		ApplyRateControls              bool `json:"applyRateControls"`
		ApplyReputationControls        bool `json:"applyReputationControls"`
		ApplySlowPostControls          bool `json:"applySlowPostControls"`
		ApplyMalwareControls           bool `json:"applyMalwareControls"`
		ApplyURLProtectionControls     bool `json:"applyUrlProtectionControls"`
	}

	// GetExportConfigurationResponseSecurityPolicyWebApplicationFirewall describes WAF settings for a policy.
	GetExportConfigurationResponseSecurityPolicyWebApplicationFirewall struct {
		RuleActions        []GetExportConfigurationResponseSecurityPolicyWebApplicationFirewallRuleAction        `json:"ruleActions,omitempty"`
		AttackGroupActions []GetExportConfigurationResponseSecurityPolicyWebApplicationFirewallAttackGroupAction `json:"attackGroupActions,omitempty"`
		Evaluation         *WebApplicationFirewallEvaluation                                                     `json:"evaluation,omitempty"`
		ThreatIntel        string                                                                                `json:"threatIntel"`
	}

	// GetExportConfigurationResponseSecurityPolicyWebApplicationFirewallRuleAction describes a rule action in WAF.
	GetExportConfigurationResponseSecurityPolicyWebApplicationFirewallRuleAction struct {
		Action                 string              `json:"action"`
		ID                     int                 `json:"id"`
		RulesetVersionID       int                 `json:"rulesetVersionId"`
		Conditions             *RuleConditions     `json:"conditions,omitempty"`
		AdvancedExceptionsList *AdvancedExceptions `json:"advancedExceptions,omitempty"`
		Exception              *RuleException      `json:"exception,omitempty"`
	}

	// GetExportConfigurationResponseSecurityPolicyWebApplicationFirewallAttackGroupAction describes an attack group action in WAF.
	GetExportConfigurationResponseSecurityPolicyWebApplicationFirewallAttackGroupAction struct {
		Action                 string                         `json:"action"`
		Group                  string                         `json:"group"`
		RulesetVersionID       int                            `json:"rulesetVersionId"`
		AdvancedExceptionsList *AttackGroupAdvancedExceptions `json:"advancedExceptions,omitempty"`
		Exception              *AttackGroupException          `json:"exception,omitempty"`
	}

	// GetExportConfigurationResponseSecurityPolicyClientReputation describes client reputation settings for a policy.
	GetExportConfigurationResponseSecurityPolicyClientReputation struct {
		ReputationProfileActions *ClientReputationReputationProfileActions `json:"reputationProfileActions,omitempty"`
	}

	// GetExportConfigurationResponseEvaluating describes evaluating policies in the export response.
	GetExportConfigurationResponseEvaluating struct {
		SecurityPolicies []EvaluatingSecurityPolicy `json:"securityPolicies,omitempty"`
	}

	// RatePoliciesPath is returned as part of GetExportConfigurationResponse.
	RatePoliciesPath struct {
		PositiveMatch bool                    `json:"positiveMatch"`
		Values        *RatePoliciesPathValues `json:"values,omitempty"`
	}

	// ExportActionItem describes an action and its associated identifier.
	ExportActionItem struct {
		Action string `json:"action"`
		ID     int    `json:"id"`
	}

	// ReputationProfileActionsexp is returned as part of GetExportConfigurationResponse.
	ReputationProfileActionsexp []ExportActionItem

	// RatePolicyActionItem describes a rate policy action.
	RatePolicyActionItem struct {
		ID         int    `json:"id"`
		Ipv4Action string `json:"ipv4Action"`
		Ipv6Action string `json:"ipv6Action"`
	}

	// RatePolicyActionsexp is returned as part of GetExportConfigurationResponse.
	RatePolicyActionsexp []RatePolicyActionItem

	// SlowRateThresholdExp is returned as part of GetExportConfigurationResponse.
	SlowRateThresholdExp struct {
		Period int `json:"period"`
		Rate   int `json:"rate"`
	}

	// DurationThresholdExp is returned as part of GetExportConfigurationResponse.
	DurationThresholdExp struct {
		Timeout int `json:"timeout"`
	}

	// SlowPostexp is returned as part of GetExportConfigurationResponse.
	SlowPostexp struct {
		Action            string                `json:"action"`
		SlowRateThreshold *SlowRateThresholdExp `json:"slowRateThreshold,omitempty"`
		DurationThreshold *DurationThresholdExp `json:"durationThreshold,omitempty"`
	}

	// AdvancedOptionsexp is returned as part of GetExportConfigurationResponse.
	AdvancedOptionsexp struct {
		Logging              *Loggingexp                                   `json:"logging"`
		AttackPayloadLogging *AttackPayloadLogging                         `json:"attackPayloadLogging"`
		EvasivePathMatch     *EvasivePathMatchexp                          `json:"evasivePathMatch,omitempty"`
		JA4Fingerprint       *JA4Fingerprintexp                            `json:"ja4Fingerprint,omitempty"`
		Prefetch             *Prefetch                                     `json:"prefetch"`
		PragmaHeader         *GetAdvancedSettingsPragmaResponse            `json:"pragmaHeader,omitempty"`
		AsePenaltyBox        *AsePenaltyBoxexp                             `json:"asePenaltyBox,omitempty"`
		RequestBody          *RequestBody                                  `json:"requestBody,omitempty"`
		PIILearning          *PIILearningexp                               `json:"piiLearning,omitempty"`
		URLEvasionDefense    *GetAdvancedSettingsURLEvasionDefenseResponse `json:"urlEvasion,omitempty"`
	}

	// CustomDenyListexp is returned as part of GetExportConfigurationResponse.
	CustomDenyListexp []CustomDenyListexpItem

	// CustomDenyListexpItem represents an individual custom deny list item in the export response.
	CustomDenyListexpItem struct {
		Description string                `json:"description,omitempty"`
		Name        string                `json:"name"`
		ID          string                `json:"id"`
		Parameters  []CustomDenyParameter `json:"parameters"`
	}

	// CustomRuleActionsexp is returned as part of GetExportConfigurationResponse.
	CustomRuleActionsexp []ExportActionItem

	// Siemexp is returned as part of GetExportConfigurationResponse.
	Siemexp struct {
		EnableForAllPolicies        bool     `json:"enableForAllPolicies,omitempty"`
		EnableSiem                  bool     `json:"enableSiem"`
		EnabledBotmanSiemEvents     bool     `json:"enabledBotmanSiemEvents,omitempty"`
		IncludeJA4FingerprintToSiem *bool    `json:"includeJA4FingerprintToSiem,omitempty"`
		FirewallPolicyIDs           []string `json:"firewallPolicyIds,omitempty"`
		SiemDefinitionID            int      `json:"siemDefinitionId,omitempty"`
	}

	// PenaltyBoxexp is returned as part of GetExportConfigurationResponse.
	PenaltyBoxexp struct {
		Action               string `json:"action"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection"`
	}

	// APIRequestConstraintsexp is returned as part of GetExportConfigurationResponse.
	APIRequestConstraintsexp struct {
		Action       string             `json:"action,omitempty"`
		APIEndpoints []ExportActionItem `json:"apiEndpoints,omitempty"`
	}

	// Evaluationexp is returned as part of GetExportConfigurationResponse.
	Evaluationexp struct {
		AttackGroupActions []EvaluationexpAttackGroupActionsItem `json:"attackGroupActions"`
		EvaluationID       int                                   `json:"evaluationId"`
		EvaluationVersion  int                                   `json:"evaluationVersion"`
		RuleActions        []EvaluationexpRuleActionsItem        `json:"ruleActions"`
		RulesetVersionID   int                                   `json:"rulesetVersionId"`
	}

	// EvaluationexpAttackGroupActionsItem describes an attack group action in an evaluation.
	EvaluationexpAttackGroupActionsItem struct {
		Action string `json:"action"`
		Group  string `json:"group"`
	}

	// EvaluationexpRuleActionsItem describes a rule action in an evaluation.
	EvaluationexpRuleActionsItem struct {
		Action     string          `json:"action"`
		ID         int             `json:"id"`
		Conditions *RuleConditions `json:"conditions,omitempty"`
		Exception  *RuleException  `json:"exception,omitempty"`
	}

	// ConditionReputationProfile is returned as part of GetExportConfigurationResponse.
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

	// HeaderCookieOrParamValuesattackgroup is returned as part of GetExportConfigurationResponse.
	HeaderCookieOrParamValuesattackgroup []HeaderCookieOrParamValuesattackgroupItem

	// HeaderCookieOrParamValuesattackgroupItem represents a criteria item in HeaderCookieOrParamValuesattackgroup.
	HeaderCookieOrParamValuesattackgroupItem struct {
		Criteria      []HeaderCookieOrParamValuesattackgroupItemCriteriaItem `json:"criteria"`
		ValueWildcard bool                                                   `json:"valueWildcard,omitempty"`
		Values        []string                                               `json:"values,omitempty"`
	}

	// HeaderCookieOrParamValuesattackgroupItemCriteriaItem represents hostnames, paths, and values in criteria.
	HeaderCookieOrParamValuesattackgroupItemCriteriaItem struct {
		Hostnames []string `json:"hostnames,omitempty"`
		Paths     []string `json:"paths,omitempty"`
		Values    []string `json:"values,omitempty"`
	}

	// SpecificHeaderCookieOrParamNameValueexp is returned as part of GetExportConfigurationResponse.
	SpecificHeaderCookieOrParamNameValueexp struct {
		Name     *json.RawMessage `json:"name,omitempty"`
		Selector string           `json:"selector,omitempty"`
		Value    *json.RawMessage `json:"value,omitempty"`
	}

	// AtomicConditionsexp is returned as part of GetExportConfigurationResponse.
	AtomicConditionsexp []AtomicConditionsexpItem

	// AtomicConditionsexpItem represents an atomic condition item.
	AtomicConditionsexpItem struct {
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

	// LoggingHeaderOrCookie describes header or cookie logging configuration.
	LoggingHeaderOrCookie struct {
		Type   string   `json:"type"`
		Values []string `json:"values,omitempty"`
	}

	// Loggingexp is returned as part of GetExportConfigurationResponse.
	Loggingexp struct {
		AllowSampling   bool                  `json:"allowSampling"`
		Cookies         LoggingHeaderOrCookie `json:"cookies"`
		CustomHeaders   LoggingHeaderOrCookie `json:"customHeaders"`
		StandardHeaders LoggingHeaderOrCookie `json:"standardHeaders"`
	}

	// LoggingOverridesexp is returned as part of GetExportConfigurationResponse.
	LoggingOverridesexp struct {
		AllowSampling   bool                  `json:"allowSampling"`
		Cookies         LoggingHeaderOrCookie `json:"cookies"`
		CustomHeaders   LoggingHeaderOrCookie `json:"customHeaders"`
		Override        bool                  `json:"override"`
		StandardHeaders LoggingHeaderOrCookie `json:"standardHeaders"`
	}

	// AttackPayloadLoggingPayloadType describes payload logging type.
	AttackPayloadLoggingPayloadType struct {
		Type string `json:"type"`
	}

	// AttackPayloadLogging is returned as part of GetExportConfigurationResponse.
	AttackPayloadLogging struct {
		Enabled      bool                            `json:"enabled"`
		RequestBody  AttackPayloadLoggingPayloadType `json:"requestBody"`
		ResponseBody AttackPayloadLoggingPayloadType `json:"responseBody"`
	}

	// AttackPayloadLoggingOverrides is returned as part of GetExportConfigurationResponse.
	AttackPayloadLoggingOverrides struct {
		Enabled      bool                            `json:"enabled"`
		RequestBody  AttackPayloadLoggingPayloadType `json:"requestBody"`
		ResponseBody AttackPayloadLoggingPayloadType `json:"responseBody"`
		Override     bool                            `json:"override"`
	}

	// EvasivePathMatchexp contains the EnablePathMatch setting
	EvasivePathMatchexp struct {
		EnablePathMatch bool `json:"enabled"`
	}

	// JA4Fingerprintexp contains the JA4 Client TLS Fingerprint setting
	JA4Fingerprintexp struct {
		HeaderNames []string `json:"headerNames,omitempty"`
	}

	// AsePenaltyBoxexp contains the ASE Penalty Box setting
	AsePenaltyBoxexp struct {
		RequestCount            int                                     `json:"requestCount"`
		BlockDuration           int                                     `json:"blockDuration"`
		ClientIdentifiers       []string                                `json:"clientIdentifiers"`
		AkamaiManagedExclusions AsePenaltyBoxexpAkamaiManagedExclusions `json:"akamaiManagedExclusions"`
		QualificationExclusions AsePenaltyBoxexpQualificationExclusions `json:"qualificationExclusions"`
	}

	// AsePenaltyBoxexpAkamaiManagedExclusions describes Akamai managed exclusions for ASE penalty box.
	AsePenaltyBoxexpAkamaiManagedExclusions struct {
		Rules       []int  `json:"rules"`
		LastUpdated string `json:"lastUpdated,omitempty"`
	}

	// AsePenaltyBoxexpQualificationExclusions describes qualification exclusions for ASE penalty box.
	AsePenaltyBoxexpQualificationExclusions struct {
		AttackGroups []string `json:"attackGroups"`
		Rules        []int    `json:"rules"`
	}

	// PIILearningexp contains the PIILearning setting
	PIILearningexp struct {
		EnablePIILearning bool `json:"enabled"`
	}

	// Prefetch is returned as part of AdvancedOptionsexp
	Prefetch struct {
		AllExtensions      bool     `json:"allExtensions"`
		EnableAppLayer     bool     `json:"enableAppLayer"`
		EnableRateControls bool     `json:"enableRateControls"`
		Extensions         []string `json:"extensions,omitempty"`
	}

	// RequestBody is returned as part of GetExportConfigurationResponse.
	RequestBody struct {
		RequestBodyInspectionLimitInKB     string `json:"requestBodyInspectionLimitInKB"`
		RequestBodyInspectionLimitOverride bool   `json:"override"`
	}

	// ConditionsExp is returned as part of GetExportConfigurationResponse.
	ConditionsExp []ConditionsExpItem

	// ConditionsExpItem represents an individual condition item in ConditionsExp.
	ConditionsExpItem struct {
		Type          string           `json:"type"`
		PositiveMatch bool             `json:"positiveMatch"`
		Name          *json.RawMessage `json:"name,omitempty"`
		NameCase      *json.RawMessage `json:"nameCase,omitempty"`
		NameWildcard  *json.RawMessage `json:"nameWildcard,omitempty"`
		Value         *json.RawMessage `json:"value,omitempty"`
		ValueCase     *json.RawMessage `json:"valueCase,omitempty"`
		ValueWildcard *json.RawMessage `json:"valueWildcard,omitempty"`
	}

	// RatePoliciesPathValues is returned as part of GetExportConfigurationResponse.
	RatePoliciesPathValues []string

	// RatePoliciesQueryParameters is returned as part of GetExportConfigurationResponse.
	RatePoliciesQueryParameters []RatePoliciesQueryParametersItem

	// RatePoliciesQueryParametersItem represents a query parameter in a rate policy.
	RatePoliciesQueryParametersItem struct {
		Name          string                             `json:"name"`
		PositiveMatch bool                               `json:"positiveMatch"`
		ValueInRange  bool                               `json:"valueInRange"`
		Values        *RatePoliciesQueryParametersValues `json:"values,omitempty"`
	}

	// RatePoliciesQueryParametersValues is returned as part of GetExportConfigurationResponse.
	RatePoliciesQueryParametersValues []string

	// SecurityPoliciesPenaltyBox is returned as part of GetExportConfigurationResponse.
	SecurityPoliciesPenaltyBox struct {
		Action               string `json:"action,omitempty"`
		PenaltyBoxProtection bool   `json:"penaltyBoxProtection,omitempty"`
	}

	// SecurityPoliciesPenaltyBoxConditions is returned as part of GetExportConfigurationResponse.
	SecurityPoliciesPenaltyBoxConditions struct {
		ConditionOperator string          `json:"conditionOperator,omitempty"`
		Conditions        *RuleConditions `json:"conditions,omitempty"`
	}

	// WebApplicationFirewallEvaluation is returned as part of GetExportConfigurationResponse.
	WebApplicationFirewallEvaluation struct {
		AttackGroupActions []WebApplicationFirewallEvaluationAttackGroupActionsItem `json:"attackGroupActions,omitempty"`
		EvaluationID       int                                                      `json:"evaluationId"`
		EvaluationVersion  int                                                      `json:"evaluationVersion"`
		RuleActions        []WebApplicationFirewallEvaluationRuleActionsItem        `json:"ruleActions,omitempty"`
		RulesetVersionID   int                                                      `json:"rulesetVersionId"`
	}

	// WebApplicationFirewallEvaluationAttackGroupActionsItem describes an attack group action in WAF evaluation.
	WebApplicationFirewallEvaluationAttackGroupActionsItem struct {
		Action                 string              `json:"action"`
		Group                  string              `json:"group"`
		Exception              *RuleException      `json:"exception,omitempty"`
		AdvancedExceptionsList *AdvancedExceptions `json:"advancedExceptions,omitempty"`
	}

	// WebApplicationFirewallEvaluationRuleActionsItem describes a rule action in WAF evaluation.
	WebApplicationFirewallEvaluationRuleActionsItem struct {
		Action                 string              `json:"action"`
		ID                     int                 `json:"id"`
		Conditions             *RuleConditions     `json:"conditions,omitempty"`
		Exception              *RuleException      `json:"exception,omitempty"`
		AdvancedExceptionsList *AdvancedExceptions `json:"advancedExceptions,omitempty"`
	}

	// RulesetsRules is returned as part of GetExportConfigurationResponse.
	RulesetsRules []RulesetsRulesItem

	// RulesetsRulesItem represents an individual rule item in a ruleset.
	RulesetsRulesItem struct {
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

	// ClientReputationReputationProfileActions is returned as part of GetExportConfigurationResponse.
	ClientReputationReputationProfileActions []ExportActionItem

	// SecurityPoliciesRatePolicyActions is returned as part of GetExportConfigurationResponse.
	SecurityPoliciesRatePolicyActions []RatePolicyActionItem

	// AdvancedSettings is returned as part of GetExportConfigurationResponse
	AdvancedSettings struct {
		BotAnalyticsCookieSettings              map[string]interface{} `json:"botAnalyticsCookieSettings,omitempty"`
		BotAnalyticsSettings                    map[string]interface{} `json:"botAnalyticsSettings,omitempty"`
		ClientSideSecuritySettings              map[string]interface{} `json:"clientSideSecuritySettings,omitempty"`
		TransactionalEndpointProtectionSettings map[string]interface{} `json:"transactionalEndpointProtectionSettings,omitempty"`
		UserRiskResponseStrategySettings        map[string]interface{} `json:"userRiskResponseStrategySettings,omitempty"`
		UserAllowListIdSettings                 map[string]interface{} `json:"userAllowListIdSettings,omitempty"`
	}

	// ResponseActions is returned as part of GetExportConfigurationResponse
	ResponseActions struct {
		ChallengeActions           []map[string]interface{} `json:"challengeActions,omitempty"`
		ConditionalActions         []map[string]interface{} `json:"conditionalActions,omitempty"`
		CustomDenyActions          []map[string]interface{} `json:"customDenyActions,omitempty"`
		ServeAlternateActions      []map[string]interface{} `json:"serveAlternateActions,omitempty"`
		ChallengeInterceptionRules map[string]interface{}   `json:"challengeInterceptionRules,omitempty"`
		ChallengeInjectionRules    map[string]interface{}   `json:"challengeInjectionRules,omitempty"`
	}

	// BotManagement is returned as part of GetExportConfigurationResponse
	BotManagement struct {
		AkamaiBotCategoryActions                  []map[string]interface{} `json:"akamaiBotCategoryActions,omitempty"`
		BotDetectionActions                       []map[string]interface{} `json:"botDetectionActions,omitempty"`
		BotManagementSettings                     map[string]interface{}   `json:"botManagementSettings,omitempty"`
		CustomBotCategoryActions                  []map[string]interface{} `json:"customBotCategoryActions,omitempty"`
		JavascriptInjectionRules                  map[string]interface{}   `json:"javascriptInjectionRules,omitempty"`
		TransactionalEndpoints                    *TransactionalEndpoints  `json:"transactionalEndpoints,omitempty"`
		ContentProtectionRules                    []map[string]interface{} `json:"contentProtectionRules,omitempty"`
		ContentProtectionRuleSequence             []string                 `json:"contentProtectionRuleSequence,omitempty"`
		ContentProtectionJavaScriptInjectionRules []map[string]interface{} `json:"contentProtectionJavaScriptInjectionRules,omitempty"`
	}

	// AccountProtection is returned as part of GetExportConfigurationResponse
	AccountProtection struct {
		GeneralSettings        map[string]interface{}   `json:"generalSettings,omitempty"`
		TransactionalEndpoints []map[string]interface{} `json:"transactionalEndpoints,omitempty"`
	}

	// TransactionalEndpoints is returned as port of GetExportConfigurationResponse
	TransactionalEndpoints struct {
		BotProtection           []map[string]interface{} `json:"botProtection,omitempty"`
		BotProtectionExceptions map[string]interface{}   `json:"botProtectionExceptions,omitempty"`
	}

	// ExportAIRulesResponse is the aiRules object nested inside a security policy in the export response.
	ExportAIRulesResponse struct {
		AIRuleStatus string               `json:"aiRuleStatus"`
		AIRules      []ExportPolicyAIRule `json:"aiRules"`
	}

	// ExportPolicyAIRule is an AI rule as returned in the exported response.
	ExportPolicyAIRule struct {
		RuleID              int64                            `json:"ruleId"`
		RuleVersion         int64                            `json:"ruleVersion"`
		Title               string                           `json:"title"`
		RiskScoreGroup      string                           `json:"riskScoreGroup"`
		RuleDescription     string                           `json:"ruleDescription"`
		Action              string                           `json:"action"`
		ConditionExceptions []ExportAIRuleConditionException `json:"conditionExceptions"`
	}

	// ExportAIRuleConditionException is a single condition exception entry within an exported AI rule.
	ExportAIRuleConditionException struct {
		Exception *ExportAIRuleException `json:"exception"`
	}

	// ExportAIRuleException holds the exception selectors for an exported AI rule.
	ExportAIRuleException struct {
		Selectors []ExportAIRuleExceptionSelector `json:"selectors"`
	}

	// ExportAIRuleExceptionSelector is a single selector entry within an AI rule exception.
	ExportAIRuleExceptionSelector struct {
		Names          []string `json:"names"`
		Selector       string   `json:"selector"`
		Wildcard       bool     `json:"wildcard"`
		RiskScoreGroup string   `json:"riskScoreGroup"`
	}
)

// Validate validates an GetExportConfigurationRequest struct.
func (v GetExportConfigurationRequest) Validate() error {
	return edgegriderr.ParseValidationErrors(validation.Errors{
		"Source": validation.Validate(v.Source, validation.In("TF").Error(
			fmt.Sprintf("value '%s' is invalid. Must be one of: 'TF' or empty", v.Source))),
	})
}

var (
	// ErrGetExportConfiguration is returned when ErrGetExportConfiguration fails
	ErrGetExportConfiguration = errors.New("get export configuration")
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

func (p *appsec) GetExportConfiguration(ctx context.Context, params GetExportConfigurationRequest) (*GetExportConfigurationResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetExportConfiguration")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetExportConfiguration, ErrStructValidation, err)
	}

	uri, err := url.Parse(fmt.Sprintf("/appsec/v1/export/configs/%d/versions/%d", params.ConfigID, params.Version))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %s", err)
	}

	if params.Source != "" {
		q := uri.Query()
		q.Add("source", params.Source)
		uri.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetExportConfiguration request: %w", err)
	}

	var result GetExportConfigurationResponse
	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("get export configuration request failed: %w", err)
	}
	defer session.CloseResponseBody(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil
}
