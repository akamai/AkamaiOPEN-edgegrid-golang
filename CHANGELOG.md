# EDGEGRID GOLANG RELEASE NOTES

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
