# EDGEGRID GOLANG RELEASE NOTES

## X.X.X (X X, X)

#### BREAKING CHANGES:


























#### FEATURES/ENHANCEMENTS:
































#### BUG FIXES:






























## 8.4.0 (Aug 22, 2024)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added field `ClientLists` to `RuleConditions` and `AttackGroupConditions`
  * The `RequestBodyInspectionLimitOverride` field has been added in the following structures:
    * `GetAdvancedSettingsRequestBodyResponse`,
    * `UpdateAdvancedSettingsRequestBodyRequest`,
    * `UpdateAdvancedSettingsRequestBodyResponse`,
    * `RemoveAdvancedSettingsRequestBodyRequest`,
    * `RemoveAdvancedSettingsRequestBodyResponse`

* IAM
  * Added new methods:
    * [GetProperty](https://techdocs.akamai.com/iam-api/reference/get-property)
    * [ListProperties](https://techdocs.akamai.com/iam-api/reference/get-properties)
    * [MoveProperty](https://techdocs.akamai.com/iam-api/reference/put-property)
    * `MapPropertyIDToName` - to provide property name for given IAM property ID
    
* PAPI
  * Added new method `MapPropertyNameToID` to provide PAPI property ID for given property name

## 8.3.0 (July 09, 2024)

#### FEATURES/ENHANCEMENTS:

* General
  * Added `To` utility function in the `ptr` package that helps with creating value pointers

* BOTMAN
  * Added Content Protection APIs
    * [CreateContentProtectionRule](https://techdocs.akamai.com/content-protector/reference/post-content-protection-rule)
    * [GetContentProtectionRuleList](https://techdocs.akamai.com/content-protector/reference/get-content-protection-rules)
    * [GetContentProtectionRule](https://techdocs.akamai.com/content-protector/reference/get-content-protection-rule)
    * [UpdateContentProtectionRule](https://techdocs.akamai.com/content-protector/reference/put-content-protection-rule)
    * [RemoveContentProtectionRule](https://techdocs.akamai.com/content-protector/reference/delete-content-protection-rule)
    * [GetContentProtectionRuleSequence](https://techdocs.akamai.com/content-protector/reference/get-content-protection-rule-sequence)
    * [UpdateContentProtectionRuleSequence](https://techdocs.akamai.com/content-protector/reference/put-content-protection-rule-sequence)
    * [GetContentProtectionJavaScriptInjectionRuleList](https://techdocs.akamai.com/content-protector/reference/get-content-protection-javascript-injection-rules)
    * [GetContentProtectionJavaScriptInjectionRule](https://techdocs.akamai.com/content-protector/reference/get-content-protection-javascript-injection-rule)
    * [CreateContentProtectionJavaScriptInjectionRule](https://techdocs.akamai.com/content-protector/reference/post-content-protection-javascript-injection-rule)
    * [UpdateContentProtectionJavaScriptInjectionRule](https://techdocs.akamai.com/content-protector/reference/put-content-protection-javascript-injection-rule)
    * [RemoveContentProtectionJavaScriptInjectionRule](https://techdocs.akamai.com/content-protector/reference/delete-content-protection-javascript-injection-rule)

* Added Cloud Access Manager API support
  * Access Keys
    * [GetAccessKeyStatus](https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key-create-request)
    * [CreateAccessKey](https://techdocs.akamai.com/cloud-access-mgr/reference/post-access-key)
    * [GetAccessKey](https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key)
    * [ListAccessKeys](https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-keys)
    * [UpdateAccessKey](https://techdocs.akamai.com/cloud-access-mgr/reference/put-access-key)
    * [DeleteAccessKey](https://techdocs.akamai.com/cloud-access-mgr/reference/delete-access-key)
  * Access Key Versions
    * [GetAccessKeyVersionStatus](https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key-version-create-request)
    * [GetAccessKeyVersion](https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key-version)
    * [CreateAccessKeyVersion](https://techdocs.akamai.com/cloud-access-mgr/reference/post-access-key-version)
    * [ListAccessKeyVersions](https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key-versions)
    * [DeleteAccessKeyVersion](https://techdocs.akamai.com/cloud-access-mgr/reference/delete-access-key-version)
  * Properties using Access Key
    * [LookupProperties](https://techdocs.akamai.com/cloud-access-mgr/reference/get-access-key-version-properties)
    * [GetAsyncPropertiesLookupID](https://techdocs.akamai.com/cloud-access-mgr/reference/get-async-version-property-lookup)
    * [PerformAsyncPropertiesLookup](https://techdocs.akamai.com/cloud-access-mgr/reference/get-property-lookup)

* DNS
  * Added [GetZonesDNSSecStatus](https://techdocs.akamai.com/edge-dns/reference/post-zones-dns-sec-status) method returning the current DNSSEC status for one or more zones

#### Deprecations

* Deprecated the following functions in the `tools` package. Use `ptr.To` instead.
  * `BoolPtr`
  * `IntPtr`
  * `Int64Ptr`
  * `Float32Ptr`
  * `Float64Ptr`
  * `StringPtr`

## 8.2.0 (May 21, 2024)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added `CounterType` field to `CreateRatePolicyResponse`, `UpdateRatePolicyResponse`, `RemoveRatePolicyResponse`, `GetRatePoliciesResponse` and `GetRatePolicyResponse` structs to support managing rate policy counter type

* BOTMAN 
  * Added [GetCustomBotCategoryItemSequence](https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-category-item-sequence)  and [UpdateCustomBotCategoryItemSequence](https://techdocs.akamai.com/bot-manager/reference/put-custom-bot-category-item-sequence)

* HAPI
  * Added method to return certificate for the edge hostname 
    * [GetCertificate](https://techdocs.akamai.com/edge-hostnames/reference/get-edge-hostname-certificate)
  * Added fields to `GetEdgeHostnameResponse`: `ProductID`, `MapAlias` and `UseCases`

#### BUG FIXES:

* APPSEC
  * The `Override` field in the following structs has been updated from a pointer to a value type within the `AdvancedSettingsAttackPayloadLogging` interface:
    * `GetAdvancedSettingsAttackPayloadLoggingResponse`,
    * `UpdateAdvancedSettingsAttackPayloadLoggingResponse`,
    * `RemoveAdvancedSettingsAttackPayloadLoggingRequest`,
    * `RemoveAdvancedSettingsAttackPayloadLoggingResponse`
      This update was made to address a drift issue related to policy level settings.
  * Omit `Prefetch` within `AdvancedOptions` in `GetExportConfigurationResponse` when empty

* CLOUDLETS
  * Added validation that `ObjectMatchValue` is not supported with `MatchType` `query` in `MatchRuleER` ([#535](https://github.com/akamai/terraform-provider-akamai/issues/535))

## 8.1.0 (April 11, 2024)

#### FEATURES/ENHANCEMENTS:

* DNS
  * Modified `ParseRData` method to remove priority, weight and port from targets **only** when those values are same for all `SRV` targets.
    Otherwise, targets are returned untouched and `priority`, `weight` and `port` in the map are not populated.

* Image and Video Manager
  * Added `SmartCrop` transformation

## 8.0.0 (March 19, 2024)

#### BREAKING CHANGES:

* Migrated to go 1.21

* CPS
  * Split request and response structures for create and update enrollment operations

* DNS
  * Renamed following structs:
    * `RecordsetQueryArgs` into `RecordSetQueryArgs`
    * `Recordsets` into `RecordSets`
    * `Recordset` into `RecordSet`
    * `MetadataH` into `Metadata`
  * Renamed following fields:
    * `GroupId` into `GroupID` in `ListGroupRequest`
    * `Recordsets` into `RecordSets` in `RecordSetResponse`
    * `ContractIds` into `ContractIDs` in `TSIGQueryString`
    * `Gid` into `GID` in `TSIGQueryString` and `TSIGReportMeta`
    * `TsigKey` into `TSIGKey` in `ZoneCreate` and `ZoneResponse`
    * `VersionId` into `VersionID` in `ZoneResponse`
    * `RequestId` into `RequestID` in `BulkZonesResponse`, `BulkStatusResponse`, `BulkCreateResultResponse` and `BulkDeleteResultResponse`
  * Renamed `RecordSets` interface into `Recordsets`
  * Renamed following methods:
    * `ListTsigKeys` into `ListTSIGKeys`
    * `GetTsigKeyZones` into `GetTSIGKeyZones`
    * `GetTsigKeyAliases` into `GetTSIGKeyAliases`
    * `TsigKeyBulkUpdate` into `TSIGKeyBulkUpdate`
    * `GetTsigKey` into `GetTSIGKey`
    * `DeleteTsigKey` into `DeleteTSIGKey`
    * `UpdateTsigKey` into `UpdateTSIGKey`
    * `GetRecordsets` into `GetRecordSets`
    * `CreateRecordsets` into `CreateRecordSets`
    * `UpdateRecordsets` into `UpdateRecordSets`
  * Deleted following methods:
    * `NewAuthorityResponse`
    * `NewChangeListResponse`
    * `NewRecordBody`
    * `NewRecordSetResponse`
    * `NewTsigKey`
    * `NewTsigQueryString`
    * `NewZone`
    * `NewZoneQueryString`
    * `NewZoneResponse`
    * `RecordToMap`
  * Unexported following methods:
    * `FullIPv6`
    * `PadCoordinates`
    * `ValidateZone`

* GTM
  * Renamed following structs:
    * `AsAssignment` into `ASAssignment`
    * `AsMap` into `ASMap`
    * `AsMapList` into `ASMapList`
    * `CidrAssignment` into `CIDRAssignment`
    * `CidrMap` into `CIDRMap`
    * `CidrMapList` into `CIDRMapList`
    * `CidrMapResponse` into `CIDRMapResponse`
    * `AsMapResponse` into `ASMapResponse`
    * `HttpHeader` into `HTTPHeader`
  * Renamed following fields:
    * `AsNumbers` into `ASNumbers` in `ASAssignment`
    * `AsMapItems` into `ASMapItems` in `ASMapList`
    * `CidrMapItems` into `CIDRMapItems` in `CIDRMapList`
    * `ChangeId` into `ChangeID` in `ResponseStatus`
    * `DatacenterId` into `DatacenterID` in `DatacenterBase`, `Datacenter`, `TrafficTarget` and `ResourceInstance`
    * `AsMaps` into `ASMaps` in `Domain`
    * `DefaultSslClientPrivateKey` into `DefaultSSLClientPrivateKey` in `Domain`
    * `CnameCoalescingEnabled` into `CNameCoalescingEnabled` in `Domain`
    * `CidrMaps` into `CIDRMaps` in `Domain`
    * `DefaultSslClientCertificate` into `DefaultSSLClientCertificate` in `Domain`
    * `AcgId` into `AcgID` in `DomainItem`
    * `HttpError3xx` into `HTTPError3xx` in `LivenessTest`
    * `HttpError4xx` into `HTTPError4xx` in `LivenessTest`
    * `HttpError5xx` into `HTTPError5xx` in `LivenessTest`
    * `SslClientPrivateKey` into `SSLClientPrivateKey` in `LivenessTest`
    * `SslClientCertificate` into `SSLClientCertificate` in `LivenessTest`
    * `HttpHeaders` into `HTTPHeaders` in `LivenessTest`
    * `Ipv6` into `IPv6` in `Property`
    * `BackupIp` into `BackupIP` in `Property`
  * Renamed `CidrMaps` interface into `CIDRMaps`
  * Renamed following methods:
    * `ListAsMaps` into `ListASMaps`
    * `GetAsMap` into `GetASMap`
    * `CreateAsMap` into `CreateASMap`
    * `DeleteAsMap` into `DeleteASMap`
    * `UpdateAsMap` into `UpdateASMap`
    * `ListCidrMaps` into `ListCIDRMaps`
    * `GetCidrMap` into `GetCIDRMap`
    * `CreateCidrMap` into `CreateCIDRMap`
    * `DeleteCidrMap` into `DeleteCIDRMap`
    * `UpdateCidrMap` into `UpdateCIDRMap`
  * Deleted following methods:
    * `NewASAssignment`
    * `NewAsMap`
    * `NewCidrAssignment`
    * `NewCidrMap`
    * `NewDatacenter`
    * `NewDatacenterBase`
    * `NewDatacenterResponse`
    * `NewDefaultDatacenter`
    * `NewDomain`
    * `NewGeoAssignment`
    * `NewHttpHeader`
    * `NewGeoMap`
    * `NewLivenessTest`
    * `NewLoadObject`
    * `NewProperty`
    * `NewResource`
    * `NewResourceInstance`
    * `NewResponseStatus`
    * `NewStaticRRSet`
    * `NewTrafficTarget`

#### FEATURES/ENHANCEMENTS:

* Added default value `application/json` for `Accept` header for all requests sent to API

* Appsec
  * Added `PenaltyBoxConditions` API - read and update
  * Added `EvalPenaltyBoxConditions` API - read and update

* CPS
  * Added `ID`, `OrgID`, `ProductionSlots`, `StagingSlots` and `AssignedSlots` to the response structures of `GetEnrollment` and `ListEnrollment` operations

* GTM
  * Added new fields:
    * `SignAndServe` and `SignAndServeAlgorithm` for the `Domain`
    * `HTTPMethod`, `HTTPRequestBody`, `Pre2023SecurityPosture` and `AlternateCACertificates` for the `LivenessTest` in `Property`
    * `Precedence` for the `TrafficTarget` in `Property`
  * Enhanced error details by addition of `Errors` field in `Error` structure
  * Added support for the creation of `ranked-failover` properties

#### BUG FIXES:

* DNS
  * Removed not working `DeleteZone` method
*
* PAPI
  * Updated documentation link for `GetProperties` method

## 7.6.1 (February 14, 2024)

#### BUG FIXES:

* Edgeworkers
  * Fixed case when not providing optional `note` field in `ActivateVersion` would cause activation to fail

## 7.6.0 (February 8, 2024)

#### FEATURES/ENHANCEMENTS:

* General
  * Enhanced error handling when Error is not in standard format.

* Added Cloudlets V3 API support
  * Cloudlet Info
    * [ListCloudlets](https://techdocs.akamai.com/cloudlets/reference/get-cloudlets)
  * Policies
    * [ListPolicies](https://techdocs.akamai.com/cloudlets/reference/get-policies)
    * [CreatePolicy](https://techdocs.akamai.com/cloudlets/reference/post-policy)
    * [DeletePolicy](https://techdocs.akamai.com/cloudlets/reference/delete-policy)
    * [GetPolicy](https://techdocs.akamai.com/cloudlets/reference/get-policy)
    * [UpdatePolicy](https://techdocs.akamai.com/cloudlets/reference/put-policy)
    * [ClonePolicy](https://techdocs.akamai.com/cloudlets/reference/post-policy-clone)
  * Policy Properties
    * [ListActivePolicyProperties](https://techdocs.akamai.com/cloudlets/reference/get-policy-properties)
  * Policy Versions
    * [ListPolicyVersions](https://techdocs.akamai.com/cloudlets/reference/get-policy-versions)
    * [GetPolicyVersion](https://techdocs.akamai.com/cloudlets/reference/get-policy-version)
    * [CreatePolicyVersion](https://techdocs.akamai.com/cloudlets/reference/post-policy-version)
    * [DeletePolicyVersion](https://techdocs.akamai.com/cloudlets/reference/delete-policy-version)
    * [UpdatePolicyVersion](https://techdocs.akamai.com/cloudlets/reference/put-policy-version)
  * Policy Activations
    * [ListPolicyActivations](https://techdocs.akamai.com/cloudlets/reference/get-policy-activations)
    * [GetPolicyActivation](https://techdocs.akamai.com/cloudlets/reference/get-policy-activation)
    * [ActivatePolicy and DeactivatePolicy](https://techdocs.akamai.com/cloudlets/reference/post-policy-activations)
  * Supported cloudlet types
    * API Prioritization (AP)
    * Application Segmentation (AS)
    * Edge Redirector (ER)
    * Forward Rewrite (FR)
    * Phased Release (PR aka CD)
    * Request Control (RC aka IG)

* DNS
  * Added `ListGroups` method
    * [ListGroups](https://techdocs.akamai.com/edge-dns/reference/get-data-groups)

* Edgeworkers
  * Added `note` field to `Activation` and `ActivateVersion` structs for EdgeWorkers Activation

* GTM
  * Added new fields to `DomainItem` struct

* IVM
  * Extended `OutputImage` for support of `AllowPristineOnDownsize` and `PreferModernFormats`
  * Extended `PolicyInputImage` for support of `ServeStaleDuration`
  * Extended `RolloutInfo` for support of `ServeStaleEndTime`

#### BUG FIXES:

* APPSEC
  * Added `updateLatestNetworkStatus` query parameter in GetActivations request to resolve drift on manual changes to infrastructure

## 7.5.0 (November 28, 2023)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added `ASNControls` field to `UpdateIPGeoRequest` and `IPGeoFirewall` structs to support firewall blocking by ASN client lists

* BOTMAN
  * Added API support for Custom Code - read and update

## 7.4.0 (October 24, 2023)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Updated `GetExportConfigurationResponse` struct to export rate policy `burstWindow` and `condition` fields

* Cloudlets
  * Added MatchesAlways field to ER cloudlet

* IAM
  * Phone number is no longer required for IAM user for `CreateUser` and `UpdateUserInfo` methods

## 7.3.0 (September 19, 2023)

#### FEATURES/ENHANCEMENTS:

* ClientLists
  * Updated `GetClientListResponse` and `UpdateClientListResponse` to include `GroupID`

* GTM
  * Added custom error `ErrNotFound` that can be used to check if GTM api returned 404 not found

* HAPI
  * Added `GetChangeRequest`

* Updated `yaml.v3` dependency

## 7.2.1 (August 25, 2023)

#### BUG FIXES:

* CloudWrapper
  * Fixed build for 32-bit systems

## 7.2.0 (August 22, 2023)

#### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Added CloudWrapper API support
  * Capacities
    * [ListCapacities](https://techdocs.akamai.com/cloud-wrapper/reference/get-capacity-inventory)
  * Configurations
    * [GetConfiguration](https://techdocs.akamai.com/cloud-wrapper/reference/get-configuration)
    * [ListConfigurations](https://techdocs.akamai.com/cloud-wrapper/reference/get-configurations)
    * [CreateConfiguration](https://techdocs.akamai.com/cloud-wrapper/reference/post-configuration)
    * [UpdateConfiguration](https://techdocs.akamai.com/cloud-wrapper/reference/put-configuration)
    * [ActivateConfiguration](https://techdocs.akamai.com/cloud-wrapper/reference/post-configuration-activations)
  * Locations
    * [ListLocations](https://techdocs.akamai.com/cloud-wrapper/reference/get-locations)
  * MultiCDN
    * [ListAuthKeys](https://techdocs.akamai.com/cloud-wrapper/reference/get-auth-keys)
    * [ListCDNProviders](https://techdocs.akamai.com/cloud-wrapper/reference/get-providers)
  * Properties
    * [ListProperties](https://techdocs.akamai.com/cloud-wrapper/reference/get-properties)
    * [ListOrigins](https://techdocs.akamai.com/cloud-wrapper/reference/get-origins)

* [IMPORTANT] Added Client Lists API Support
  * ClientLists
    * [GetClientLists](https://techdocs.akamai.com/client-lists/reference/get-lists)
      * Support filter by name or type
    * [GetClientList](https://techdocs.akamai.com/client-lists/reference/get-list)
    * [UpdateClientList](https://techdocs.akamai.com/client-lists/reference/put-update-list)
    * [UpdateClientListItems](https://techdocs.akamai.com/client-lists/reference/post-update-items)
    * [CreateClientList](https://techdocs.akamai.com/client-lists/reference/post-create-list)
    * [DeleteClientList](https://techdocs.akamai.com/client-lists/reference/delete-list)
  * Activations
    * [GetActivation](https://techdocs.akamai.com/client-lists/reference/get-retrieve-activation-status)
    * [GetActivationStatus](https://techdocs.akamai.com/client-lists/reference/get-activation-status)
    * [CreateActivation](https://techdocs.akamai.com/client-lists/reference/post-activate-list)

* APPSEC
  * Added Bot Management API Support
    * Custom Client Sequence - read and update

## 7.1.0 (July 25, 2023)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added Bot Management API Support
    * Challenge Injection Rules - read, update
    * Add `CreateSecurityPolicyWithDefaultProtections` method to the `SecurityPolicy` interface to support creating a security policy with all available protections enabled.
  * Update marshaling of PII learning setting

### Deprecations

* Challenge Interceptions Rules has been deprecated
* Deprecate the following interfaces used to maintain individual policy protections:
  * `ApiConstraintsProtection`
  * `IPGeoProtection`
  * `MalwareProtection`
  * `NetworkLayerProtection`
  * `RateProtection`
  * `ReputationProtection`
  * `SlowPostProtection`
  * `WAFProtection`
* Deprecate the `CreateSecurityPolicy` method of the `SecurityPolicy` interface.

## 7.0.0 (June 20, 2023)

### BREAKING CHANGES:

* DataStream
  * Updated `connectors` details in DataStream 2 API v2.
  * Updated `GetProperties` and `GetDatasetFields` methods in DataStream 2 API v2.
  * Updated `CreateStream`, `GetStream`, `UpdateStream`, `DeleteStream` and `ListStreams` methods in DataStream 2 API v2.
  * Updated `Activate`, `Deactivate`, `ActivationHistory` and `Stream` details in DataStream 2 API v2 and also changed their corresponding response objects.

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Update Geo control to include Action for Ukraine.
  * Add `AdvancedSettingsPIILearning` interface to support reading and updating the PII learning setting.

### BUG FIXES:

* APPSEC
  * Add error handling for failed NetworkList client calls.

## 6.0.0 (May 23, 2023)

### BREAKING CHANGES:

* APPSEC
  * Update malware policy `ContentTypes` to include `EncodedContentAttributes`.
  * Malware policy's `ContentTypes` is reported as part of an individual policy but is no longer included in the bulk report of all policies.

* CLOUDLETS
  * `ActivatePolicyVersion` also returns list of triggerred activations

* PAPI
  * Fix property variables fields with empty and null values are ignored
  * Remove `ProductID` field from `GetEdgeHostname` response

### BUG FIXES:
* APPSEC
  * Omit `clientIdentifier` and `additionalMatchOptions` in `GetExportConfigurationResponse` when empty

## 5.0.0 (March 28, 2023)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Add `AdvancedSettingsRequestBody` interface to support configuring request size inspection limit

* EDGEKV
  * [ListGroupsWithinNamespace](https://techdocs.akamai.com/edgekv/reference/get-groups)

* Image and Video Manager
  * Add possible value `avif` for `forcedFormats` and `allowedFormats`

* PAPI
  * Add `complianceRecord` field to `Activation` struct for PAPI activation

#### BREAKING CHANGES:

* APPSEC
  * Remove deprecated `EvalHost` and `EvalProtectHost` interfaces. (Use the `WAPSelectedHostnames` interface instead.)
  * Remove deprecated `BypassNetworkList` interface. (Use the `WAPBypassNetworkList` interface instead.)

## 4.1.0 (Feb 27, 2023)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added following BotManager fields to GetExportConfigurationResponse
    * BotManagement
    * CustomBotCategories
    * CustomDefinedBots
    * CustomBotCategorySequence
    * CustomClients
    * ResponseActions
    * AdvancedSettings
  * Added AdvancedSettingsAttackPayloadLogging interface

#### BUG FIXES:

* Fix V4 of Edgegrid doesn't parse hostname ([#182](https://github.com/akamai/AkamaiOPEN-edgegrid-golang/pull/182))

## 4.0.0 (Jan 31, 2023)

#### BREAKING CHANGES:

* Migrate to go 1.18

* PAPI
  * Fix response structures for GetAvailableBehaviors and GetAvailableCriteria:
    * [GetAvailableCriteria](https://techdocs.akamai.com/property-mgr/reference/get-available-criteria)
    * [GetAvailableBehaviors](https://techdocs.akamai.com/property-mgr/reference/get-available-behaviors)

* CPS
  * Update `Accept` header to the latest schema `application/vnd.akamai.cps.enrollment.v11+json` for the following endpoints:
    * [ListEnrollments](https://techdocs.akamai.com/cps/reference/get-enrollments)
    * [GetEnrollment](https://techdocs.akamai.com/cps/reference/get-enrollment)

* APPSEC
  * Fix incorrect return type structure in `UpdateBypassNetworkListsResponse`
  * Return `RatePolicyCondition` via a pointer in response structs of `RatePolicy` APIs

#### FEATURES/ENHANCEMENTS:

* Replace obsolete APIs documentation links with new one from [https://techdocs.akamai.com](https://techdocs.akamai.com)

* APPSEC
  * Add `burstWindow` and `condition` fields to RatePolicy

* CPS
  * Add `preferredTrustChain` field to `csr` struct ([#351](https://github.com/akamai/terraform-provider-akamai/issues/351))
  * Set `utf-8 charset` in `content-type` header for requests

#### BUG FIXES:

* Fix code errors in documentation examples ([#177](https://github.com/akamai/AkamaiOPEN-edgegrid-golang/pull/177))

* IAM
  * Issue updating user information - removed validation on user update

## 3.1.0 (Dec 12, 2022)

#### FEATURES/ENHANCEMENTS:

* General
  * Add badges to readme and improve code quality

* [IMPORTANT] Added Property Include API support
  * Includes
    * [ListIncludes](https://techdocs.akamai.com/property-mgr/reference/get-includes)
    * [ListIncludeParents](https://techdocs.akamai.com/property-mgr/reference/get-include-parents)
    * [GetInclude](https://techdocs.akamai.com/property-mgr/reference/get-include)
    * [CreateInclude](https://techdocs.akamai.com/property-mgr/reference/post-includes)
    * [DeleteInclude](https://techdocs.akamai.com/property-mgr/reference/delete-include)
  * Include Activations
    * [ActivateInclude](https://techdocs.akamai.com/property-mgr/reference/post-include-activation)
    * [DeactivateInclude](https://techdocs.akamai.com/property-mgr/reference/post-include-activation)
    * [CancelIncludeActivation](https://techdocs.akamai.com/property-mgr/reference/delete-include-activation)
    * [GetIncludeActivation](https://techdocs.akamai.com/property-mgr/reference/get-include-activation)
    * [ListIncludeActivations](https://techdocs.akamai.com/property-mgr/reference/get-include-activations)
  * Include Rules
    * [GetIncludeRuleTree](https://techdocs.akamai.com/property-mgr/reference/get-include-version-rules)
    * [UpdateIncludeRuleTree](https://techdocs.akamai.com/property-mgr/reference/patch-include-version-rules)
  * Include Versions
    * [CreateIncludeVersion](https://techdocs.akamai.com/property-mgr/reference/post-include-versions)
    * [GetIncludeVersion](https://techdocs.akamai.com/property-mgr/reference/get-include-version)
    * [ListIncludeVersions](https://techdocs.akamai.com/property-mgr/reference/get-include-versions)
    * [ListIncludeVersionAvailableCriteria](https://techdocs.akamai.com/property-mgr/reference/get-include-available-criteria)
    * [ListIncludeVersionAvailableBehaviors](https://techdocs.akamai.com/property-mgr/reference/get-include-available-behaviors)

#### BREAKING CHANGES:

* APPSEC
  * Factor out `PolicySecurityControls` struct

## 3.0.0 (November 28, 2022)

### Deprecations

* CPS
  * UpdateChange() function has been deprecated

#### FEATURES/ENHANCEMENTS:

* CPS
  * ChangeManagementInfo - get or acknowledge change management info, get change deployment info
  * Deployments - list deployments, get production deployment, get staging deployment
  * DeploymentSchedules - get deployment schedule, update deployment schedule
  * History - get DV history, get certificate history, get change history
  * PostVerification - get or acknowledge post verification warnings
  * ThirdPartyCSR - get third-party CSR, upload certificate

#### BREAKING CHANGES:

* Rename package `configdns` to `dns`
* Rename package `configgtm` to `gtm`
* CPS
  * Renamed structs: Challenges and ValidationRecords to Challenge and ValidationRecord accordingly
  * Type change: `NotAfter` and `NotBefore` fields in `DeploymentSchedule` struct used in response for `GetChangeStatus` are `*string` instead of `string`

## 2.17.0 (October 24, 2022)

#### FEATURES/ENHANCEMENTS:

* Datastream
  * Add ListStreams
  * Add new connectors: Elasticsearch, NewRelic and Loggly
  * Extend Splunk and Custom HTTPS connectors mTLS certificates configuration
  * Extend SumoLogic, Splunk and Custom HTTPS connectors with ability to specify custom HTTP headers

#### BUG FIXES:

* APPSEC
  * Fix incorrect JSON sent when applying appsec_ip_geo resource in allow mode

## 2.16.0 (September 26, 2022)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Add interfaces to support file malware scanning (FMS):
    * MalwareContentTypes
    * MalwarePolicy
    * MalwarePolicyAction
    * MalwareProtection
  * Add GetRuleRecommendations method to TuningRecommendations interface
  * Add deprecation notes for the following:
    * methods:
      * GetIPGeoProtections
      * GetNetworkLayerProtections
      * GetRateProtections
      * GetReputationProtections
      * GetSlowPostProtectionSetting
      * GetSlowPostProtections
      * GetWAFProtections
      * RemoveNetworkLayerProtection
      * RemovePolicyProtections
      * RemoveReputationProtection
    * structs:
      * GetIPGeoProtectionsRequest
      * GetNetworkLayerProtectionsRequest
      * GetRateProtectionsRequest
      * GetReputationProtectionsRequest
      * GetSlowPostProtectionSettingRequest
      * GetSlowPostProtectionSettingResponse
      * GetSlowPostProtectionsRequest
      * GetWAFProtectionsRequest
      * RemoveNetworkLayerProtectionRequest
      * RemovePolicyProtectionsRequest
      * RemoveReputationProtectionRequest

* [IMPORTANT] Added Bot Management API Support
    * Akamai Bot Category - read
    * Akamai Bot Category Action - read, update
    * Akamai Defined Bot - read
    * Bot Analytics Cookie - read, update
    * Bot Analytics Cookie Values - read
    * Bot Category Exception - read, update
    * Bot Detection - read
    * Bot Detection Action - read, update
    * Bot Endpoint Coverage Report - read
    * Bot Management Setting - read, update
    * Challenge Action - create, read, update, delete
    * Challenge Interception Rules - read, update
    * Client Side Security - read, update
    * Conditional Action - create, read, update, delete
    * Custom Bot Category - create, read, update, delete
    * Custom Bot Category Action - read, update
    * Custom Bot Category Sequence - read, update
    * Custom Client - create, read, update, delete
    * Custom Defined Bot - create, read, update, delete
    * Custom Deny Action - create, read, update, delete
    * Javascript Injection - read, update
    * Recategorized Akamai Defined Bot - create, read, update, delete
    * Response Action - read
    * Serve Alternate Action - create, read, update, delete
    * Transactional Endpoint - create, read, update, delete
    * Transactional Endpoint Protection - read, update

## 2.15.0 (August 22, 2022)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Add xff field to custom rule conditions
  * Add NotificationEmails to Activation struct

* GTM
  * Improved error messages

* CPS
  * Add cps ListEnrollments
  * Extend CreateEnrollment with AllowDuplicateCN option

## 2.14.1 (July 26, 2022)

#### BUG FIXES:

* IAM
  * Change IAM GroupID type to int64

## 2.14.0 (June 28, 2022)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added penalty box support for security policy in evaluation mode

* HAPI
  * EdgeHostname - update

* IAM
  * Blocked properties - read, update
  * Group - create, read, update, delete
  * Role - create, read, update, delete
  * User - lock, unlock, TFA, set password, reset password

#### BUG FIXES:
* APPSEC
  * Fixed incorrect error message on activation failure
  * The `EffectiveTimePeriod`, `SamplingRate`, `LoggingOptions`, and `Operation` fields of the various `CustomRule` response structs are now marshalled correctly

## 2.13.0 (May 31, 2022)

#### FEATURES/ENHANCEMENTS:

* Image and Video Manager
  * Add new `ImQuery` transformation
  * New `PostBreakPointTransformationType`

#### BUG FIXES:

* Image and Video Manager
  * `default_value` field on variable in image policy should not be required
  * Change all primitive optional parameters to pointers
  * Correct `Anchor` field in `RectangleShapeType`
  * Value field for `NumberVariableInline` should be defined as `float64`
  * Rename `PointShapeType.True` to `PointShapeType.Y`, to match the OpenAPI definition
  * Add `Composite` transformation to `PostBreakpointTransformations`
  * Fix `PostBreakpointTransformations.PolicyInputImage`

## 2.12.0 (Apr. 25, 2022)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Add WAPBypassNetworkLists interface, to be used in preference to deprecated BypassNetworkLists interface

* Support for account switch keys from environment ([#149](https://github.com/akamai/AkamaiOPEN-edgegrid-golang/pull/149))

## 2.11.0 (March 24, 2022)

#### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Added Image and Video Manager API support
  * Policy Set - create, read, update, delete
  * Policy - create, read, update, delete, rollback to previous version, view policy history

* CLOUDLETS
  * Support for RC cloudlet type (Request Control)

* PAPI
  * CP code - read, update

## 2.10.0 (Feb. 28, 2022)

#### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Added EdgeWorkers and EdgeKV API support
  * EDGEWORKERS
    * Ids - create, read, update, delete, clone
    * Versions - create, read, delete, validate version bundle
    * Activations - create, read, delete
    * Deactivations - read, delete
    * Resource tiers - read
    * Reports - read
    * Secure token - create
    * Permission groups - read
    * Properties - read
    * Contracts - read
  * EDGEKV
    * Items - create, read, update, delete
    * Namespaces - create, read, update
    * Initialization - create, read
    * Access token - create, read, delete

* APPSEC
  * Source for evasive path match interface updated with links to documentation

* CLOUDLETS
  * Support for AS cloudlet type (Audience Segmentation)

## 2.9.1 (Feb. 7, 2022)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Remove deprecation notes for individual policy protection methods

#### BUG FIXES:

* CLOUDLETS
  * Fixed validation for ALB version DataCenter percent

## 2.9.0 (Jan. 24, 2022)

#### FEATURES/ENHANCEMENTS:

* CLOUDLETS
  * Support for VP cloudlet type (Visitor Prioritization)
  * Support for CD cloudlet type (Continuous Deployment / Phased Release)
  * Support for FR cloudlet type (Forward Rewrite)
  * Support for AP cloudlet type (API Prioritization)

* APPSEC
  * Add support for Evasive Path Match feature
  * Deprecate individual policy protection interface methods

* NETWORK LISTS
  * Include ContractID and GroupID in GetNetworkListResponse

## 2.8.1 (Nov. 30, 2021)

#### FEATURES/ENHANCEMENTS:
* [IMPORTANT] Added Cloudlets API support
  * Policy (Application Load Balancer) -  create, read, update, delete policy
  * Policy (Edge Redirector) -  create, read, update, delete policy
  * Policy activation - create, read
  * Application Load Balancer configuration - create, update, read
  * Activation for Application Load Balancer configuration - create, read

* APPSEC
  * Add support for advanced exceptions in ASE rules
  * Update bypass-network-list datasource and resource for multi-policy WAP

## 2.7.0 (Oct 19, 2021)

#### FEATURES/ENHANCEMENTS:
* [IMPORTANT] Added DataStream API support
  * Stream operations
  * Stream activation operations
  * Read access to various DataStream properties
* Added HAPI v1 support
  * Delete edge hostname

## 2.6.0 (Aug 16, 2021)

#### BUG FIXES:
* APPSEC
  * Fix incorrect comments/URL references in inline documentation

#### FEATURES/ENHANCEMENTS
* APPSEC
  * Get an evaluation attack group's or risk score group's action

* NETWORK LISTS
  * Support contract_id and group_id for network list create/update

## 2.5.0 (Jun 15, 2021)

#### BREAKING CHANGES:
* APPSEC
  * The following have been removed, togther with their unit tests and test data:
    * pkg/appsec/attack_group_action.go
    * pkg/appsec/attack_group_condition_exception.go
    * pkg/appsec/eval_rule_action.go
    * pkg/appsec/eval_rule_condition_exception.go
    * pkg/appsec/rule_action.go
    * pkg/appsec/rule_condition_exception.go

#### BUG FIXES:
* DNSv2
    * Fixed parsing SVCB, HTTPS rdata.

#### FEATURES/ENHANCEMENTS:
* [IMPORTANT] CPS - Added Certificate Provisioning API support
  * Enrollments - create, read, update, delete enrollments
  * Change status API - get change status, cancel change
  * DV certificate API - get and acknowledge DV challenges
  * Pre verification warnings - get and acknowledge pre verification warnings

* APPSEC
  * The following have been added, together with their unit tests and test data:
    * pkg/appsec/api_constraints_protection.go
    * pkg/appsec/advanced_settings_pragma_header.go
    * pkg/appsec/attack_group.go
    * pkg/appsec/eval_rule.go
    * pkg/appsec/rule.go
    * pkg/appsec/ip_geo_protection.go

## 2.4.1 (Apr 19, 2021)

#### BUG FIXES:

* APPSEC
  * Suppress 'null' text on output of empty/false values
  * Prevent configuration drift when reapplying configuration after importing or creating resources

## 2.4.0 (Mar 29, 2021) PAPI - Secure by default

* PAPI
   * Support to provision default certs as part of hostnames request
   * New cert status object in hostnames response if it exists

## 2.3.0 (Mar 15, 2021) Network Lists

Add support for the following operations in the Network Lists API v2:

* Create a network list
* Update an existing network list
* Get the existing network lists, including optional filtering by name or type
* Subscribe to a network list
* Activate a network list

## 2.2.1 (Mar 3, 2021)
* PAPI - Fixed issue with rules causing advanced locked behaviors to fail

## 2.2.0 (Feb 23, 2021) APPSEC - Extended list of supported list endpoints from APPSEC API

#### BUG FIXES:
* PAPI
    * Fixed issue with version and rule comments being dropped
    * Fixed client side validation to allow certain PAPI errors to passthrough

#### FEATURES/ENHANCEMENTS:
* APPSEC
    * Custom Deny
    * SIEM Setting
    * Advanced Options Settings
    * API Match Target
    * API Request Constraint
    * Create/Delete/Rename Security Policy
    * Host Coverage / Edit Version Notes
    * All WAP Features / WAP Hostname Evaluation
    * Create Security Configuration
    * Rename Security Configuration Version
    * Delete Security Configuration Version
    * Clone Security Configuration
    * Import tool for adding existing resources to Terraform state
* DNS
    * Add support for HTTPS, SVCB records to ParseRData

## 2.1.1 (Feb 3, 2021)
* PAPI - Fixed validation on empty rule behaviors causing some properties with nested behaviors to fail

## 2.1.0 (Jan 13, 2021)
* [IMPORTANT] IAM - New Identity and Access Management API Support

## 2.0.4 (Dec 23, 2020)
* APPSEC - Extended list of supported endpoints from APPSEC API:
  * DDoS Protection -- Rate Policy & Action
  * DDoS Protection -- Slowpost setting & Action
  * Application Layer Protection -- Rule Action, Exceptions & Conditions
  * Application Layer Protection -- Rule Evaluation Action, Exceptions & Conditions
  * Application Layer Protection -- Attack Group Action, Exceptions & Conditions
  * Application Layer Protection -- Rule Upgrade & Change Mode for Rule Eval
  * Reputation Profile & Action
  * Network Layer Control -- IP & GEO setting

## 2.0.3 (Dec 7, 2020)
* PAPI - Property hostname validation fix for missing hostnames.
* PAPI - fix minor typo in rules error messages

## 2.0.2 (Nov 19, 2020)
* [IMPORTANT] APPSEC - Added Application Security API
* [ENHANCEMENT] DNS - Bulk Api endpoints added
* ALL - Re-enabled global account switch key support in edgerc files
* PAPI - Edgehostname IPV6 support fix.  Added enums with allowed values.
* PAPI - Edgehostname blank cname or egdehostname id fix
* PAPI - propertyversion blank etag field fix

## 2.0.1 (Oct 15, 2020)
* [IMPORTANT] Breaking changes from earlier clients. Project updated to use v2 directory structure.
* [ENHANCEMENT] PAPI - Api error return to the user when an activation or validation error occurs.
* [NOTE] Project re-organized to prepare for additional APIs to be included in future versions of this library.

## 1.0.0 (Oct 15, 2020)
* Official release for the EdgeGrid Golang library
* DNSv2 - Zone create signature to pass blank instead of nil
* PAPI - Return nil instead of error if no cp code was found
* GTM - Datacenter API requires blank instead of nil

## 0.9.18 (Jul 13, 2020)
* [AT-40][Add] Preliminary Logging CorrelationID

## 0.9.17 (Jun 9, 2020)
* Corrected AKAMAICDN target parsing
* Added endpoints for list zones, creating and updating multiple recordsets
* Refactored recordsets into separate source file

## 0.9.16 (May 29, 2020)
* Client-v1, Papi-v1 Updates
* Add lock around http request creation.
* papi - add logging to papi endpoints.

## 0.9.15 (May 15, 2020)
* DNSv2 - Added CERT, TSLA Record parsing. Removed MX Record parsing

## 0.9.14 (May 12, 2020)
* DNSv2 - Enhance RecordError functions

## 0.9.13 (Apr 26, 2020)
* DNSv2 - filterZoneCreate check upper case Type

## 0.9.12 (Apr 21, 2020)
* DNSv2 - Added optional arg to bypass dns record lock for create, update and delete functions. default preserves prior behavior

## 0.9.11 (Apr 13 , 2020)
* DNSv2 Updates
  * Add additional fields, including TSIG, to zone
  * Support alias zone types
  * Add utility functions for Rdata parsing and process.
  * Add GetRecord, GetRecordSet functions
  * Add additional Recordset metadata
* Add http request/response logging

## 0.9.10 (Mar 5, 2020)
* Add support for caching Edgehostnames and Products
* Support for cache in papi library for edgehostnames and products to minimize round trips to fetch repeated common data to avoid
  WAF deny rule IPBLOCK-BURST4-54013 issue

## 0.9.9 (Feb 29, 2020)
* Add support for caching Contract, Groups, and Cp Codes
* cache to minimize round trips on repeated common data fetches to avoid
  WAF deny rule IPBLOCK-BURST4-54013 issue

## 0.9.0 (Aug 6, 2019)
* Added support for GTM
