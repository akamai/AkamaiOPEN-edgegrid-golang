# RELEASE NOTES

## X.X.X (X X, X)

### BREAKING CHANGES:









* Logging
  * Changed logger from apex to the custom interface
    * Logger method differences:
      *  A new method `With` has been added that condenses the apex methods (`WithError`, `WithField`, `WithFields`, `WithDuration`) into one.
          * `WithError`, `WithField`, `WithFields`, `WithDuration` methods are not included in the new logger.
          * The `Entry` ([documentation](https://pkg.go.dev/github.com/apex/log#Entry)) type no longer exists in the new logger, `With` instead of `Entry` returns a new logger instance with new fields.
      * Logging methods (`Fatal`, `Error`, `Warn`, `Info`, `Debug`) can accept key-value pairs in addition to a message,
          * The attribute arguments are processed as follows: If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair. Otherwise, the argument is treated as a value with key "!BADKEY".
          * formatted logging methods (`Fatalf`, `Errorf`, `Warnf`, `Infof`, `Debugf`) remain unchanged.
  * By default `slog` with custom handler is used.
  * `log.Interface` allows users to define default logger with `SetLogger` method and provides option to use different logger backend.
    * Instructions on using different logger backends can be found in `pkg/log/README.md` file
  * Log output structure have changed slightly,
    * Time format was adjusted, logger will use 24-hour clock with milliseconds instead of 12-hour clock used previously.
















### FEATURES/ENHANCEMENTS:




* Migrated to go 1.22



* Improved formatting of validation errors
* Added ability to return an error to `session.Option`
* PAPI
  * Added the `OriginalInput` parameter in the `GetRuleTreeRequest` to allow returning upgraded content of rules. When omitted it is equal to true, meaning that returned rules are exactly as sent.


* Updated vulnerable dependencies
* Improved code by resolving issues reported by linter 1.58.1
























### BUG FIXES:










* DNS
  * Fixed incorrect URL for `ListGroups` method.




























## 9.1.0 (Nov 14, 2024)

### FEATURES/ENHANCEMENTS:

* General
  * Added a configurable `WithRetries` option for creating new sessions with global GET retries. It can be configured with these parameters:
    * `retryMax`  - The maximum number of API request retries.
    * `retryWaitMin` - The minimum wait time in `time.Duration` between API requests retries.
    * `retryWaitMax` - The maximum wait time in `time.Duration` between API requests retries.
    * `excludedEndpoints` - The list of path expressions defining endpoints which should be excluded from the retry feature.
  * Added logic responsible for closing the response body in each method.

* APPSEC
  * Added following content protection fields to `GetExportConfigurationResponse` under `BotManagement` section
    * `ContentProtectionRules`
    * `ContentProtectionRuleSequence`
    * `ContentProtectionJavaScriptInjectionRules`
  * Changed `EnabledBotmanSiemEvents` to `*bool` and omitted from following structs when empty
    * `GetSiemSettingResponse`
    * `RemoveSiemSettingsRequest`
    * `RemoveSiemSettingsResponse`
    * `UpdateSiemSettingsRequest`
    * `UpdateSiemSettingsResponse`

* DNS
  * Added support for `OutboundZoneTransfer` field in requests and responses for these methods:
    * `CreateBulkZones`
    * `CreateZone`
    * `GetZone`
    * `ListZones`
    * `UpdateZone`

### BUG FIXES:

* APPSEC
  * Fixed SIEM exception validation for the `Exceptions` field.

* Cloud Access
  * Added custom error `ErrAccessKeyNotFound` to easier verify if provided access key does not exist.

## 9.0.0 (Oct 3, 2024)

### BREAKING CHANGES:

* General
  * Consolidated multiple sub-interfaces into a single interface for each sub-provider.
  * Renamed the `NTWRKLISTS` interface to `NetworkList` for the `networklists` provider.
  * Removed the `tools` package in favor of the `ptr` package.

* Cloudaccess
  * Changed naming of the request body fields for the `CreateAccessKeyVersionRequest` structure:
    * From `BodyParams` to `Body`.
    * From `CreateAccessKeyVersionBodyParams` to `CreateAccessKeyVersionRequestBody`.

* Cloudlets
  * Changed naming of the request body fields for the `UpdatePolicyRequest` structure:
    * From `BodyParams` to `Body`.
    * From `UpdatePolicyBodyParams` to `UpdatePolicyRequestBody`.
  * Changed naming of the request body fields for the `ClonePolicyRequest` structure:
    * From `BodyParams` to `Body`.
    * From `ClonePolicyBodyParams` to `ClonePolicyRequestBody`.

* Cloudwrapper
  * Changed naming of the request body field for the `CreateConfigurationRequest` structure:
    * From`CreateConfigurationBody` to `CreateConfigurationRequestBody`.
  * Changed naming of the request body field for the `UpdateConfigurationRequest` structure:
    * From `UpdateConfigurationBody` to `UpdateConfigurationRequestBody`.

* DNS
  * Refactored parameters in these methods:
    * `GetAuthorities` - from (context.Context, string) into (context.Context, `GetAuthoritiesRequest`)
    * `GetNameServerRecordList` - from (context.Context, string) into (context.Context, `GetNameServerRecordListRequest`)
    * `GetRecord` - from (context.Context, string, string, string) into (context.Context, `GetRecordRequest`)
    * `GetRecordList` - from (context.Context, string, string, string) into (context.Context, `GetRecordListRequest`)
    * `CreateRecord` - from (context.Context, *RecordBody, string, ...bool) into (context.Context, `CreateRecordRequest`)
    * `UpdateRecord` - from (context.Context, *RecordBody, string, ...bool) into (context.Context, `UpdateRecordRequest`)
    * `DeleteRecord` - from (context.Context, *RecordBody, string, ...bool) into (context.Context, `DeleteRecordRequest`)
    * `GetRecordSets` - from (context.Context, string, ...RecordSetQueryArgs) into (context.Context, `GetRecordSetsRequest`)
    * `CreateRecordSets` - from (context.Context, *RecordSets, string, ...bool) into (context.Context, `CreateRecordSetsRequest`)
    * `UpdateRecordSets` - from (context.Context, *RecordSets, string, ...bool) into (context.Context, `UpdateRecordSetsRequest`)
    * `ListTSIGKeys` - from (context.Context, *TSIGQueryString) into (context.Context, `ListTSIGKeysRequest`)
    * `GetTSIGKeyZones` - from (context.Context, *TSIGKey) into (context.Context, `GetTSIGKeyZonesRequest`)
    * `GetTSIGKeyAliases` - from (context.Context, string) into (context.Context, `GetTSIGKeyAliasesRequest`)
    * `UpdateTSIGKeyBulk` - from (context.Context, *TSIGKeyBulkPost) into (context.Context, `UpdateTSIGKeyBulkRequest`)
    * `GetTSIGKey` - from (context.Context, string) into (context.Context, `GetTSIGKeyRequest`)
    * `DeleteTSIGKey` - from (context.Context, string) into (context.Context, `DeleteTSIGKeyRequest`)
    * `UpdateTSIGKey` - from (context.Context, *TSIGKey, string) into (context.Context, `UpdateTSIGKeyRequest`)
    * `ListZones` - from (context.Context, ...ZoneListQueryArgs) into (context.Context, `ListZonesRequest`)
    * `GetZone` - from (context.Context, string) into (context.Context, `GetZoneRequest`)
    * `GetChangeList` - from (context.Context, string) into (context.Context, `GetChangeListRequest`)
    * `GetMasterZoneFile` - from (context.Context, string) into (context.Context, `GetMasterZoneFileRequest`)
    * `PostMasterZoneFile` - from (context.Context, string, string) into (context.Context, `PostMasterZoneFileRequest`)
    * `CreateZone` - from (context.Context, *ZoneCreate, ZoneQueryString, ...bool) into (context.Context, `CreateZoneRequest`)
    * `SaveChangeList` - from (context.Context, *ZoneCreate) into (context.Context, `SaveChangeListRequest`)
    * `SubmitChangeList` - from (context.Context, *ZoneCreate) into (context.Context, `SubmitChangeListRequest`)
    * `UpdateZone` - from (context.Context, *ZoneCreate) into (context.Context, `UpdateZoneRequest`)
    * `GetZoneNames` - from (context.Context, string) into (context.Context, `GetZoneNamesRequest`)
    * `GetZoneNameTypes` - from (context.Context, string, string) into (context.Context, `GetZoneNameTypesRequest`)
    * `GetBulkZoneCreateStatus` - from (context.Context, string) into (context.Context, `GetBulkZoneCreateStatusRequest`)
    * `GetBulkZoneDeleteStatus` - from (context.Context, string) into (context.Context, `GetBulkZoneDeleteStatusRequest`)
    * `GetBulkZoneCreateResult` - from (context.Context, string) into (context.Context, `GetBulkZoneCreateResultRequest`)
    * `GetBulkZoneDeleteResult` - from (context.Context, string) into (context.Context, `GetBulkZoneDeleteResultRequest`)
    * `CreateBulkZones` - from (context.Context, *BulkZonesCreate, ZoneQueryString) into (context.Context, `CreateBulkZonesRequest`)
    * `DeleteBulkZones` - from (context.Context, *ZoneNameListResponse, ...bool) into (context.Context, `DeleteBulkZonesRequest`)
    * `GetRdata` - from (context.Context, string, string, string) into (context.Context, `GetRdataRequest`)
  * Refactored the responses in these methods:
    * `GetAuthorities` - from `*AuthorityResponse` into `*GetAuthoritiesResponse`
    * `GetRecord` - `*RecordBody` into `*GetRecordResponse`
    * `GetRecordList` - from `*RecordSetResponse` into `*GetRecordListResponse`
    * `GetRecordSets` - from `*RecordSetResponse` into `*GetRecordSetsResponse`
    * `GetTSIGKey` - from `*TSIGKeyResponse` into `*GetTSIGKeyResponse`
    * `ListTSIGKeys` - from `*TSIGReportResponse` into `*ListTSIGKeysResponse`
    * `GetTSIGKeyZones` - from **`*ZoneNameListResponse` into `*GetTSIGKeyZonesResponse`
    * `GetTSIGKeyAliases` - from `*ZoneNameListResponse` into `*GetTSIGKeyAliasesResponse`
    * `GetZone` - from `*ZoneResponse` into `*GetZoneResponse`
    * `GetChangeList` - from `*ChangeListResponse` into `*GetChangeListResponse`
    * `GetZoneNames` - from `*ZoneNamesResponse` into `*GetZoneNamesResponse`
    * `GetZoneNameTypes` - from `*ZoneNameTypesResponse` into `*GetZoneNameTypesResponse`
    * `GetBulkZoneCreateStatus` - from `*BulkStatusResponse` into `*GetBulkZoneCreateStatusResponse`
    * `GetBulkZoneDeleteStatus` - from `*BulkStatusResponse` into `*GetBulkZoneDeleteStatusResponse`
    * `GetBulkZoneCreateResult` - from `*BulkCreateResultResponse` into `*GetBulkZoneCreateResultResponse`
    * `GetBulkZoneDeleteResult` - from `*BulkDeleteResultResponse` into `*GetBulkZoneDeleteResultResponse`
    * `CreateBulkZones` - from `*BulkZonesResponse` into `*CreateBulkZonesResponse`
    * `DeleteBulkZones` - from `*BulkZonesResponse` into `*DeleteBulkZonesResponse`
  * Removed these interfaces:
    * `Authorities`
    * `Data`
    * `Records`
    * `Recordsets`
    * `TSIGKeys`
    * `Zones`
  * Renamed these methods:
    * From `SaveChangelist` into `SaveChangeList`
    * From `SubmitChangelist` into `SubmitChangeList`
    * From `TSIGKeyBulkUpdate` into `UpdateTSIGKeyBulk`

* EdgeKV
  * For the `CreateEdgeKVAccessTokenRequest` structure, removed the `Expiry` field and added the `RestrictToEdgeWorkerIDs` field.
  * For the `CreateEdgeKVAccessTokenResponse` structure, removed the `Expiry` and `Value` fields, and added these fields:
    * `AllowOnProduction`
    * `AllowOnStaging`
    * `CPCode`
    * `IssueDate`
    * `LatestRefreshDate`
    * `NamespacePermissions`
    * `NextScheduledRefreshDate`
    * `RestrictToEdgeWorkerIDs`
    * `TokenActivationStatus`
  * Added these fields to the `EdgeKVAccessToken` structure:
    * `TokenActivationStatus`
    * `IssueDate`
    * `LatestRefreshDate`
    * `NextScheduledRefreshDate`

* Edgeworkers
  * Changed naming of request body field for this structure:
    * From `EdgeWorkerIDBodyRequest` to `EdgeWorkerIDRequestBody`.

* GTM
  * Refactored parameters in these methods:
    * `ListASMaps` - from (context.Context, string) into (context.Context, `ListASMapsRequest`)
    * `GetASMap` - from (context.Context, string, string) into (context.Context, `GetASMapRequests`)
    * `CreateASMap` - from (context.Context, *ASMap, string) into (context.Context, `CreateASMapRequest`)
    * `UpdateASMap` - from (context.Context, *ASMap, string) into (context.Context, `UpdateASMapRequest`)
    * `DeleteASMap` - from (context.Context, *ASMap, string) into (context.Context, `DeleteASMapRequest`)
    * `ListCIDRMaps` - from (context.Context, string) into (context.Context, `ListCIDRMapsRequest`)
    * `GetCIDRMap` - from (context.Context, string, string) into (context.Context, `GetCIDRMapRequest`)
    * `CreateCIDRMap` - from (context.Context, *CIDRMap, string) into (context.Context, `CreateCIDRMapRequest`)
    * `UpdateCIDRMap` - from (context.Context, *CIDRMap, string) into (context.Context, `UpdateCIDRMapRequest`)
    * `DeleteCIDRMap` - from (context.Context, *CIDRMap, string) into (context.Context, `DeleteCIDRMapRequest`)
    * `ListDatacenters` - from (context.Context, string) into (context.Context, `ListDatacentersRequest`)
    * `GetDatacenter` - from (context.Context, int, string) into (context.Context, `GetDatacenterRequest`)
    * `CreateDatacenter` - from (context.Context, *Datacenter, string) into (context.Context, `CreateDatacenterRequest`)
    * `UpdateDatacenter` - from (context.Context, *Datacenter, string) into (context.Context, `UpdateDatacenterRequest`)
    * `DeleteDatacenter` - from (context.Context, *Datacenter, string) into (context.Context, `DeleteDatacenterRequest`)
    * `GetDomainStatus` - from (context.Context, string) into (context.Context, `GetDomainStatusRequest`)
    * `GetDomain` - from (context.Context, string) into (context.Context, `GetDomainRequest`)
    * `CreateDomain` - from (context.Context, *Domain, map[string]string) into (context.Context, `CreateDomainRequest`)
    * `UpdateDomain` - from (context.Context, *Domain, map[string]string) into (context.Context, `UpdateDomainRequest`)
    * `DeleteDomain` - from (context.Context, *Domain) into (context.Context, `DeleteDomainRequest`)
    * `ListGeoMaps` - from (context.Context, string) into (context.Context, `ListGeoMapsRequest`)
    * `GetGeoMap` - from (context.Context, string, string) into (context.Context, `GetGeoMapRequest`)
    * `CreateGeoMap` - from (context.Context, *GeoMap, string) into (context.Context, `CreateGeoMapRequest`)
    * `UpdateGeoMap` - from (context.Context, *GeoMap, string) into (context.Context, `UpdateGeoMapRequest`)
    * `DeleteGeoMap` - from (context.Context, *GeoMap, string) into (context.Context, `DeleteGeoMapRequest`)
    * `ListProperties` - from (context.Context, string) into (context.Context, `ListPropertiesRequest`)
    * `GetProperty` - from (context.Context, string, string) into (context.Context, `GetPropertyRequest`)
    * `CreateProperty` - from (context.Context, *Property, string) into (context.Context, `CreatePropertyRequest`)
    * `UpdateProperty` - from (context.Context, *Property, string) into (context.Context, `UpdatePropertyRequest`)
    * `DeleteProperty` - from (context.Context, *Property, string) into (context.Context, `DeletePropertyRequest`)
    * `ListResources` - from (context.Context, string) into (context.Context, `ListResourcesRequest`)
    * `GetResource` - from (context.Context, string, string) into (context.Context, `GetResourceRequest`)
    * `CreateResource` - from (context.Context, *Resource, string) into (context.Context, `CreateResourceRequest`)
    * `UpdateResource` - from (context.Context, *Resource, string) into (context.Context, `UpdateResourceRequest`)
    * `DeleteResource` - from (context.Context, *Resource, string) into (context.Context, `DeleteResourceRequest`)
  * Refactored the responses in these methods:
    * `ListASMaps` - from `[]*ASMap` into `[]ASMap`
    * `GetASMap` - from`*ASMap` into `*GetASMapResponse`
    * `CreateASMap` - from `*ASMapResponse` into `*CreateASMapResponse`
    * `UpdateASMap` - from `*ResponseStatus` into `*UpdateASMapResponse`
    * `DeleteASMap` - from`*ResponseStatus` into `*DeleteASMapResponse`
    * `ListCIDRMaps` - from `[]*CIDRMap` into `[]CIDRMap`
    * `GetCIDRMap` - from `*CIDRMap` into `*GetCIDRMapResponse`
    * `CreateCIDRMap` - from `*CIDRMapResponse` into `*CreateCIDRMapResponse`
    * `UpdateCIDRMap` - from `*ResponseStatus` into `*UpdateCIDRMapResponse`
    * `DeleteCIDRMap` - from `*ResponseStatus` into `*DeleteCIDRMapResponse`
    * `ListDatacenters` - from `[]*Datacenter` into `[]Datacenter`
    * `CreateDatacenter` - from `*DatacenterResponse` into `*CreateDatacenterResponse`
    * `UpdateDatacenter` - from `*ResponseStatus` into `*UpdateDatacenterResponse`
    * `DeleteDatacenter` - from `*ResponseStatus` into `*DeleteDatacenterResponse`
    * `ListDomains` - from `[]*DomainItem` into `[]DomainItem`
    * `GetDomain` - `*Domain` into `*GetDomainResponse`
    * `CreateDomain` - from `*DomainResponse` into `*CreateDomainResponse`
    * `UpdateDomain` - from `*ResponseStatus` into `*UpdateDomainResponse`
    * `DeleteDomain` - `*ResponseStatus` into `*DeleteDomainResponse`
    * `GetDomainStatus` - from `*ResponseStatus` into `*GetDomainStatusResponse`
    * `ListGeoMaps` - from `[]*GeoMap` into `[]GeoMap`
    * `GetGeoMap` - from `*GeoMap` into `*GetGeoMapResponse`
    * `CreateGeoMap` - from `*GeoMapResponse` into `*CreateGeoMapResponse`
    * `UpdateGeoMap` - from `*ResponseStatus` into `*UpdateGeoMapResponse`
    * `DeleteGeoMap` - from `*ResponseStatus` into `*DeleteGeoMapResponse`
    * `ListProperties` - from `[]*Property` into `[]Property`
    * `GetProperty` - from `*Property` into `*GetPropertyResponse`
    * `CreateProperty` - from `*PropertyResponse` into `*CreatePropertyResponse`
    * `UpdateProperty` - from `*ResponseStatus` into `*UpdatePropertyResponse`
    * `DeleteProperty` - from `*ResponseStatus` into `*DeletePropertyResponse`
    * `ListResources` - from `[]*Resource` into `[]Resource`
    * `GetResource` - from `*Resource` into `*GetResourceResponse`
    * `CreateResource` - from `*ResourceResponse` into `*CreateResourceResponse`
    * `UpdateResource` - from `*ResponseStatus` into `*UpdateResourceResponse`
    * `DeleteResource` - from `*ResponseStatus` into `*DeleteResourceResponse`
  * Extended the response for these methods - previously only the status was returned, now the status and resource are returned:
    * `UpdateASMap`
    * `DeleteASMap`
    * `UpdateCIDRMap`
    * `DeleteCIDRMap`
    * `UpdateDatacenter`
    * `DeleteDatacenter`
    * `UpdateDomain`
    * `UpdateGeoMap`
    * `DeleteGeoMap`
    * `UpdateProperty`
    * `DeleteProperty`
    * `UpdateResource`
    * `DeleteResource`
  * Removed these interfaces:
    * `ASMaps`
    * `CIDRMaps`
    * `Datacenters`
    * `Domains`
    * `GeoMaps`
    * `Properties`
    * `Resources`

* IAM
  * Migrated V2 endpoints to V3.
  * Improved date handling to use `time.Time` instead of `string`.
    * Changed field types in these structures:
      * `Users`
        * `LastLoginDate`. Changed the field data type from `string` to `time.Time`.
        * `PasswordExpiryDate`. Changed the field data type from `string` to `time.Time`.
      * `UserListItem`
        * `LastLoginDate`. Changed the field data type from `string` to `time.Time`.
      * `Role`
        * `CreatedDate`. Changed the field data type from `string` to `time.Time`.
        * `ModifiedDate`. Changed the field data type from `string` to `time.Time`.
      * `RoleUser`
        * `LastLoginDate`. Changed the field data type from `string` to `time.Time`.
      * `GroupUser`
        * `LastLoginDate`. Changed the field data type from `string` to `time.Time`.
  * Changed the `Notifications` field to a pointer type in these structures:
    * `CreateUserRequest`
    * `UpdateUserNotificationsRequest`
  * Added the required `AdditionalAuthentication` field to the `CreateUserRequest` method.
  * Made the `Notifications` field required in the `UpdateUserNotifications` method.

* PAPI
  * Removed the `rule_format` and `product_id` fields from the `Property` structure, as this information is populated in the `GetPropertyVersion` method.

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added the `Exceptions` field to these structures:
    * `GetSiemSettingsResponse`
    * `GetSiemSettingResponse`
    * `UpdateSiemSettingsRequest`
    * `UpdateSiemSettingsResponse`
  * Added the `Source` field to the `GetExportConfigurationRequest` structure and the `TargetProduct` field to the `GetExportConfigurationResponse` structure.

* IAM
  * Updated these structures:
    * `User` with the `AdditionalAuthenticationConfigured` and `Actions` parameters.
    * `UserListItem` with the `AdditionalAuthenticationConfigured` and `AdditionalAuthentication` parameters.
    * `UserBasicInfo` with the `AdditionalAuthentication` parameter.
    * `UserActions` with the `CanGenerateBypassCode` parameter.
    * `UserNotificationOptions` with the `APIClientCredentialExpiry` parameter.
  * Added new methods:
    * [UpdateMFA](https://techdocs.akamai.com/iam-api/reference/put-user-profile-additional-authentication)
    * [ResetMFA](https://techdocs.akamai.com/iam-api/reference/put-ui-identity-reset-additional-authentication)
  * Added API Client Credentials methods:
    * [CreateYourCredential](https://techdocs.akamai.com/iam-api/reference/post-self-credentials) and [CreateCredential](https://techdocs.akamai.com/iam-api/reference/post-client-credentials)
    * [GetYourCredential](https://techdocs.akamai.com/iam-api/reference/get-self-credential) and [GetCredential](https://techdocs.akamai.com/iam-api/reference/get-client-credential)
    * [UpdateYourCredential](https://techdocs.akamai.com/iam-api/reference/put-self-credential) and [UpdateCredential](https://techdocs.akamai.com/iam-api/reference/put-client-credential)
    * [DeleteYourCredential](https://techdocs.akamai.com/iam-api/reference/delete-self-credential) and [DeleteCredential](https://techdocs.akamai.com/iam-api/reference/delete-client-credential)
    * [ListYourCredentials](https://techdocs.akamai.com/iam-api/reference/get-self-credentials) and [ListCredentials](https://techdocs.akamai.com/iam-api/reference/get-client-credentials)
    * [DeactivateYourCredential](https://techdocs.akamai.com/iam-api/reference/post-self-credential-deactivate) and [DeactivateCredential](https://techdocs.akamai.com/iam-api/reference/post-client-credential-deactivate)
    * [DeactivateYourCredentials](https://techdocs.akamai.com/iam-api/reference/post-self-credentials-deactivate) and [DeactivateCredentials](https://techdocs.akamai.com/iam-api/reference/post-client-credentials-deactivate)
  * Added the `UserStatus` and `AccountID` parameters to the `User` structure.
  * Added the [GetPasswordPolicy](https://techdocs.akamai.com/iam-api/reference/get-common-password-policy) method to get a password policy for an account.
  * Added Helper APIs:
    * [ListAllowedCPCodes](https://techdocs.akamai.com/iam-api/reference/post-api-clients-users-allowed-cpcodes)
    * [ListAuthorizedUsers](https://techdocs.akamai.com/iam-api/reference/get-api-clients-users)
    * [ListAllowedAPIs](https://techdocs.akamai.com/iam-api/reference/get-api-clients-users-allowed-apis)
    * [ListAccessibleGroups](https://techdocs.akamai.com/iam-api/reference/get-api-clients-users-group-access)
  * Added new methods:
    * [ListUsersForProperty](https://techdocs.akamai.com/iam-api/reference/get-property-users)
    * [BlockUsers](https://techdocs.akamai.com/iam-api/reference/put-property-users-block)
    * [DisableIPAllowlist](https://techdocs.akamai.com/iam-api/reference/post-allowlist-disable)
    * [EnableIPAllowlist](https://techdocs.akamai.com/iam-api/reference/post-allowlist-enable)
    * [GetIPAllowlistStatus](https://techdocs.akamai.com/iam-api/reference/get-allowlist-status)
    * `ListAccountSwitchKeys` based on [ListAccountSwitchKeys](https://techdocs.akamai.com/iam-api/reference/get-client-account-switch-keys) and [ListYourAccountSwitchKeys](https://techdocs.akamai.com/iam-api/reference/get-self-account-switch-keys)
    * `LockAPIClient` based on [LockAPIClient](https://techdocs.akamai.com/iam-api/reference/put-lock-api-client) and [LockYourAPIClient](https://techdocs.akamai.com/iam-api/reference/put-lock-api-client-self)
    * [UnlockAPIClient](https://techdocs.akamai.com/iam-api/reference/put-unlock-api-client)
    * [ListAPIClients](https://techdocs.akamai.com/iam-api/reference/get-api-clients)
    * [CreateAPIClient](https://techdocs.akamai.com/iam-api/reference/post-api-clients)
    * `GetAPIClient` based on [GetAPIClient](https://techdocs.akamai.com/iam-api/reference/get-api-client) and [GetYourAPIClient](https://techdocs.akamai.com/iam-api/reference/get-api-client-self)
    * `UpdateAPIClient` based on [UpdateAPIClient](https://techdocs.akamai.com/iam-api/reference/put-api-clients) and [UpdateYourAPIClient](https://techdocs.akamai.com/iam-api/reference/put-api-clients-self)
    * `DeleteAPIClient` based on [DeleteAPIClient](https://techdocs.akamai.com/iam-api/reference/delete-api-client) and [DeleteYourAPIClient](https://techdocs.akamai.com/iam-api/reference/delete-api-client-self)
    * [ListCIDRBlocks](https://techdocs.akamai.com/iam-api/reference/get-allowlist)
    * [CreateCIDRBlock](https://techdocs.akamai.com/iam-api/reference/post-allowlist)
    * [GetCIDRBlock](https://techdocs.akamai.com/iam-api/reference/get-allowlist-cidrblockid)
    * [UpdateCIDRBlock](https://techdocs.akamai.com/iam-api/reference/put-allowlist-cidrblockid)
    * [DeleteCIDRBlock](https://techdocs.akamai.com/iam-api/reference/delete-allowlist-cidrblockid)
    * [ValidateCIDRBlock](https://techdocs.akamai.com/iam-api/reference/get-allowlist-validate)

## 8.4.0 (Aug 22, 2024)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added the `ClientLists` field to the `RuleConditions` and `AttackGroupConditions` structures.
  * Added the `RequestBodyInspectionLimitOverride` field to these structures:
    * `GetAdvancedSettingsRequestBodyResponse`
    * `UpdateAdvancedSettingsRequestBodyRequest`
    * `UpdateAdvancedSettingsRequestBodyResponse`
    * `RemoveAdvancedSettingsRequestBodyRequest`
    * `RemoveAdvancedSettingsRequestBodyResponse`

* IAM
  * Added new methods:
    * [GetProperty](https://techdocs.akamai.com/iam-api/reference/get-property)
    * [ListProperties](https://techdocs.akamai.com/iam-api/reference/get-properties)
    * [MoveProperty](https://techdocs.akamai.com/iam-api/reference/put-property)
    * `MapPropertyIDToName` - to provide a property name for a given IAM property ID

* PAPI
  * Added a new method `MapPropertyNameToID` to provide a PAPI property ID for a given property name.

## 8.3.0 (July 09, 2024)

### FEATURES/ENHANCEMENTS:

* General
  * Added the `To` utility function in the `ptr` package to facilitate creating value pointers.

* BOTMAN
  * Added Content Protection APIs:
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

* Added Cloud Access Manager API support:
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
  * Added the [GetZonesDNSSecStatus](https://techdocs.akamai.com/edge-dns/reference/post-zones-dns-sec-status) method returning the current DNSSEC status for one or more zones.

### Deprecations

* Deprecated these functions in the `tools` package (use `ptr.To` instead):
  * `BoolPtr`
  * `IntPtr`
  * `Int64Ptr`
  * `Float32Ptr`
  * `Float64Ptr`
  * `StringPtr`

## 8.2.0 (May 21, 2024)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added the `CounterType` field to the `CreateRatePolicyResponse`, `UpdateRatePolicyResponse`, `RemoveRatePolicyResponse`, `GetRatePoliciesResponse`, and `GetRatePolicyResponse` structures to support managing the rate policy counter type.

* BOTMAN
  * Added the [GetCustomBotCategoryItemSequence](https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-category-item-sequence) and [UpdateCustomBotCategoryItemSequence](https://techdocs.akamai.com/bot-manager/reference/put-custom-bot-category-item-sequence) methods.

* HAPI
  * Added a new method to return a certificate for an edge hostname.
    * [GetCertificate](https://techdocs.akamai.com/edge-hostnames/reference/get-edge-hostname-certificate)
  * Added the `ProductID`, `MapAlias`, and `UseCases` fields to the `GetEdgeHostnameResponse` structure.

### BUG FIXES:

* APPSEC
  * Updated the `Override` field in these structures from a pointer to a value type within the `AdvancedSettingsAttackPayloadLogging` interface:
    * `GetAdvancedSettingsAttackPayloadLoggingResponse`
    * `UpdateAdvancedSettingsAttackPayloadLoggingResponse`
    * `RemoveAdvancedSettingsAttackPayloadLoggingRequest`
    * `RemoveAdvancedSettingsAttackPayloadLoggingResponse`
      This update was made to address a drift issue related to the policy level settings.
  * Omitted `Prefetch` within `AdvancedOptions` in the `GetExportConfigurationResponse` structure when empty.

* CLOUDLETS
  * Added validation that `ObjectMatchValue` is not supported with `MatchType` `query` in `MatchRuleER` ([#535](https://github.com/akamai/terraform-provider-akamai/issues/535)).

## 8.1.0 (April 11, 2024)

### FEATURES/ENHANCEMENTS:

* DNS
  * Modified the `ParseRData` method to remove priority, weight, and port from targets **only** when those values are same for all `SRV` targets.
    Otherwise, targets are returned untouched and `priority`, `weight`, and `port` in the map are not populated.

* Image and Video Manager
  * Added `SmartCrop` transformation.

## 8.0.0 (March 19, 2024)

### BREAKING CHANGES:

* Migrated to go 1.21.

* CPS
  * Split the request and response structures for create and update enrollment operations.

* DNS
  * Renamed these structures:
    * `RecordsetQueryArgs` into `RecordSetQueryArgs`
    * `Recordsets` into `RecordSets`
    * `Recordset` into `RecordSet`
    * `MetadataH` into `Metadata`
  * Renamed these fields:
    * `GroupId` into `GroupID` in `ListGroupRequest`
    * `Recordsets` into `RecordSets` in `RecordSetResponse`
    * `ContractIds` into `ContractIDs` in `TSIGQueryString`
    * `Gid` into `GID` in `TSIGQueryString` and `TSIGReportMeta`
    * `TsigKey` into `TSIGKey` in `ZoneCreate` and `ZoneResponse`
    * `VersionId` into `VersionID` in `ZoneResponse`
    * `RequestId` into `RequestID` in `BulkZonesResponse`, `BulkStatusResponse`, `BulkCreateResultResponse`, and `BulkDeleteResultResponse`
  * Renamed the `RecordSets` interface into `Recordsets`.
  * Renamed these methods:
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
  * Deleted these methods:
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
  * Unexported these methods:
    * `FullIPv6`
    * `PadCoordinates`
    * `ValidateZone`

* GTM
  * Renamed these structures:
    * `AsAssignment` into `ASAssignment`
    * `AsMap` into `ASMap`
    * `AsMapList` into `ASMapList`
    * `CidrAssignment` into `CIDRAssignment`
    * `CidrMap` into `CIDRMap`
    * `CidrMapList` into `CIDRMapList`
    * `CidrMapResponse` into `CIDRMapResponse`
    * `AsMapResponse` into `ASMapResponse`
    * `HttpHeader` into `HTTPHeader`
  * Renamed these fields:
    * `AsNumbers` into `ASNumbers` in `ASAssignment`
    * `AsMapItems` into `ASMapItems` in `ASMapList`
    * `CidrMapItems` into `CIDRMapItems` in `CIDRMapList`
    * `ChangeId` into `ChangeID` in `ResponseStatus`
    * `DatacenterId` into `DatacenterID` in `DatacenterBase`, `Datacenter`, `TrafficTarget`, and `ResourceInstance`
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
  * Renamed the `CidrMaps` interface into `CIDRMaps`.
  * Renamed these methods:
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
  * Deleted these methods:
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

### FEATURES/ENHANCEMENTS:

* Added the default value of `application/json` for the `Accept` header for all requests sent to API.

* APPSEC
  * Added the `PenaltyBoxConditions` API - read and update.
  * Added the `EvalPenaltyBoxConditions` API - read and update.

* CPS
  * Added the `ID`, `OrgID`, `ProductionSlots`, `StagingSlots`, and `AssignedSlots` fields to the response structures of the `GetEnrollment` and `ListEnrollment` operations.

* GTM
  * Added new fields:
    * `SignAndServe` and `SignAndServeAlgorithm` for the `Domain`
    * `HTTPMethod`, `HTTPRequestBody`, `Pre2023SecurityPosture`, and `AlternateCACertificates` for the `LivenessTest` in `Property`
    * `Precedence` for `TrafficTarget` in `Property`
  * Enhanced error details by adding the `Errors` field in the `Error` structure.
  * Added support for the creation of the `ranked-failover` properties.

### BUG FIXES:

* DNS
  * Removed the `DeleteZone` method that was not working.
*
* PAPI
  * Updated the documentation link for the `GetProperties` method.

## 7.6.1 (February 14, 2024)

### BUG FIXES:

* Edgeworkers
  * Fixed the case when not providing an optional `note` field in the `ActivateVersion` method would cause activation to fail.

## 7.6.0 (February 8, 2024)

### FEATURES/ENHANCEMENTS:

* General
  * Enhanced error handling when Error is not in standard format.

* Added Cloudlets V3 API support.
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
  * Added the `ListGroups` method.
    * [ListGroups](https://techdocs.akamai.com/edge-dns/reference/get-data-groups)

* Edgeworkers
  * Added the `note` field to the `Activation` and `ActivateVersion` structures for EdgeWorkers Activation.

* GTM
  * Added new fields to the `DomainItem` structure.

* IVM
  * Extended `OutputImage` for support of `AllowPristineOnDownsize` and `PreferModernFormats`.
  * Extended `PolicyInputImage` for support of `ServeStaleDuration`.
  * Extended `RolloutInfo` for support of `ServeStaleEndTime`.

### BUG FIXES:

* APPSEC
  * Added the `updateLatestNetworkStatus` query parameter in the `GetActivations` request to resolve drift on manual changes to infrastructure.

## 7.5.0 (November 28, 2023)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added the `ASNControls` field to the UpdateIPGeoRequest` and `IPGeoFirewall` structures to support firewall blocking by ASN client lists.

* BOTMAN
  * Added the API support for Custom Code - read and update.

## 7.4.0 (October 24, 2023)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Updated the `GetExportConfigurationResponse` structure to export the rate policy `burstWindow` and `condition` fields.

* Cloudlets
  * Added the `MatchesAlways` field to the ER cloudlet.

* IAM
  * Phone number is no longer required for IAM user for the `CreateUser` and `UpdateUserInfo` methods.

## 7.3.0 (September 19, 2023)

### FEATURES/ENHANCEMENTS:

* ClientLists
  * Updated the `GetClientListResponse` and `UpdateClientListResponse` structures to include the `GroupID` field.

* GTM
  * Added a custom error `ErrNotFound` that can be used to check if GTM API returned a 404 not found.

* HAPI
  * Added `GetChangeRequest`.

* Updated the `yaml.v3` dependency.

## 7.2.1 (August 25, 2023)

### BUG FIXES:

* CloudWrapper
  * Fixed the build for 32-bit systems.

## 7.2.0 (August 22, 2023)

### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Added CloudWrapper API support:
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
  * Added Bot Management API Support:
    * Custom Client Sequence - read and update.

## 7.1.0 (July 25, 2023)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added Bot Management API Support:
    * Challenge Injection Rules - read and update.
    * Added the `CreateSecurityPolicyWithDefaultProtections` method to the `SecurityPolicy` interface to support creating a security policy with all available protections enabled.
  * Updated marshaling of the PII learning setting.

### Deprecations

* Deprecated the Challenge Interceptions Rules.
* Deprecated these interfaces used to maintain individual policy protections:
  * `ApiConstraintsProtection`
  * `IPGeoProtection`
  * `MalwareProtection`
  * `NetworkLayerProtection`
  * `RateProtection`
  * `ReputationProtection`
  * `SlowPostProtection`
  * `WAFProtection`
* Deprecated the `CreateSecurityPolicy` method of the `SecurityPolicy` interface.

## 7.0.0 (June 20, 2023)

### BREAKING CHANGES:

* DataStream
  * Updated the `connectors` details in the DataStream 2 API v2.
  * Updated the `GetProperties` and `GetDatasetFields` methods in the DataStream 2 API v2.
  * Updated the `CreateStream`, `GetStream`, `UpdateStream`, `DeleteStream`, and `ListStreams` methods in the DataStream 2 API v2.
  * Updated the `Activate`, `Deactivate`, `ActivationHistory`, and `Stream` details in the DataStream 2 API v2 and changed their corresponding response objects.

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Updated the Geo control to include Action for Ukraine.
  * Added the `AdvancedSettingsPIILearning` interface to support reading and updating of the PII learning setting.

### BUG FIXES:

* APPSEC
  * Added error handling for the failed NetworkList client calls.

## 6.0.0 (May 23, 2023)

### BREAKING CHANGES:

* APPSEC
  * Updated the malware policy `ContentTypes` to include `EncodedContentAttributes`.
  * Malware policy's `ContentTypes` is reported as part of an individual policy but is no longer included in the bulk report of all policies.

* CLOUDLETS
  * Updated `ActivatePolicyVersion` to also return list of triggered activations.

* PAPI
  * Fixed the property variables fields – empty and null values are ignored.
  * Removed the `ProductID` field from the `GetEdgeHostname` response.

### BUG FIXES:

* APPSEC
  * Omitted the `clientIdentifier` and `additionalMatchOptions` fields in `GetExportConfigurationResponse` when empty.

## 5.0.0 (March 28, 2023)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added the `AdvancedSettingsRequestBody` interface to support configuring the request size inspection limit.

* EDGEKV
  * [ListGroupsWithinNamespace](https://techdocs.akamai.com/edgekv/reference/get-groups)

* Image and Video Manager
  * Added a possible value of `avif` for the `forcedFormats` and `allowedFormats` fields.

* PAPI
  * Added the `complianceRecord` field to the `Activation` structure for PAPI activation.

### BREAKING CHANGES:

* APPSEC
  * Removed the deprecated `EvalHost` and `EvalProtectHost` interfaces. (Use the `WAPSelectedHostnames` interface instead.)
  * Removed the deprecated `BypassNetworkList` interface. (Use the `WAPBypassNetworkList` interface instead.)

## 4.1.0 (Feb 27, 2023)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added these BotManager fields to `GetExportConfigurationResponse`:
    * `BotManagement`
    * `CustomBotCategories`
    * `CustomDefinedBots`
    * `CustomBotCategorySequence`
    * `CustomClients`
    * `ResponseActions`
    * `AdvancedSettings`
  * Added the `AdvancedSettingsAttackPayloadLogging` interface.

### BUG FIXES:

* Fixed an issue in Edgegrid v4 with parsing a hostname ([#182](https://github.com/akamai/AkamaiOPEN-edgegrid-golang/pull/182)).

## 4.0.0 (Jan 31, 2023)

### BREAKING CHANGES:

* Migrated to go 1.18.

* PAPI
  * Fixed the response structures for `GetAvailableBehaviors` and `GetAvailableCriteria`:
    * [GetAvailableCriteria](https://techdocs.akamai.com/property-mgr/reference/get-available-criteria)
    * [GetAvailableBehaviors](https://techdocs.akamai.com/property-mgr/reference/get-available-behaviors)

* CPS
  * Updated the `Accept` header to the latest schema `application/vnd.akamai.cps.enrollment.v11+json` for these endpoints:
    * [ListEnrollments](https://techdocs.akamai.com/cps/reference/get-enrollments)
    * [GetEnrollment](https://techdocs.akamai.com/cps/reference/get-enrollment)

* APPSEC
  * Fixed an incorrect return type structure in `UpdateBypassNetworkListsResponse`.
  * Returned `RatePolicyCondition` via a pointer in the response structures of the `RatePolicy` APIs.

### FEATURES/ENHANCEMENTS:

* Replaced obsolete APIs documentation links with the new ones from [https://techdocs.akamai.com](https://techdocs.akamai.com).

* APPSEC
  * Added the `burstWindow` and `condition` fields to `RatePolicy`.

* CPS
  * Added the `preferredTrustChain` field to the `csr` structure ([#351](https://github.com/akamai/terraform-provider-akamai/issues/351)).
  * Set `utf-8 charset` in the `content-type` header for requests.

### BUG FIXES:

* Fixed code errors in documentation examples ([#177](https://github.com/akamai/AkamaiOPEN-edgegrid-golang/pull/177)).

* IAM
  * Issued updating user information – removed validation on user update.

## 3.1.0 (Dec 12, 2022)

### FEATURES/ENHANCEMENTS:

* General
  * Added badges to readme and improved code quality.

* [IMPORTANT] Added Property Include API support:
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

### BREAKING CHANGES:

* APPSEC
  * Factored out the `PolicySecurityControls` structure.

## 3.0.0 (November 28, 2022)

### Deprecations

* CPS
  * Deprecated the `UpdateChange()` function.

### FEATURES/ENHANCEMENTS:

* CPS
  * `ChangeManagementInfo` - get or acknowledge the change management information, get the change deployment information.
  * `Deployments` - list deployments, get the production deployment, get the staging deployment.
  * `DeploymentSchedules` - get the deployment schedule, update the deployment schedule.
  * `History` - get the DV history, get the certificate history, get the change history.
  * `PostVerification` - get or acknowledge the post verification warnings.
  * `ThirdPartyCSR` - get the third-party CSR, upload a certificate.

### BREAKING CHANGES:

* Renamed the `configdns` package to `dns`.
* Rename the `configgtm` package to `gtm`.
* CPS
  * Renamed these structures: `Challenges` to `Challenge` and `ValidationRecords` to `ValidationRecord`.
  * Changed the fields' type: `NotAfter` and `NotBefore` fields in the `DeploymentSchedule` structure used in the response for `GetChangeStatus` are `*string` instead of `string`.

## 2.17.0 (October 24, 2022)

### FEATURES/ENHANCEMENTS:

* Datastream
  * Added the `ListStreams` method.
  * Added new connectors: `Elasticsearch`, `NewRelic`, and `Loggly`.
  * Extended the Splunk and Custom HTTPS connectors mTLS certificates configuration.
  * Extended the SumoLogic, Splunk, and Custom HTTPS connectors with the ability to specify custom HTTP headers.

### BUG FIXES:

* APPSEC
  * Fixed an incorrect JSON sent when applying the `appsec_ip_geo` resource in allow mode.

## 2.16.0 (September 26, 2022)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added new interfaces to support file malware scanning (FMS):
    * `MalwareContentTypes`
    * `MalwarePolicy`
    * `MalwarePolicyAction`
    * `MalwareProtection`
  * Added the `GetRuleRecommendations` method to the `TuningRecommendations` interface.
  * Added the deprecation notes for these:
    * methods:
      * `GetIPGeoProtections`
      * `GetNetworkLayerProtections`
      * `GetRateProtections`
      * `GetReputationProtections`
      * `GetSlowPostProtectionSetting`
      * `GetSlowPostProtections`
      * `GetWAFProtections`
      * `RemoveNetworkLayerProtection`
      * `RemovePolicyProtections`
      * `RemoveReputationProtection`
    * structures:
      * `GetIPGeoProtectionsRequest`
      * `GetNetworkLayerProtectionsRequest`
      * `GetRateProtectionsRequest`
      * `GetReputationProtectionsRequest`
      * `GetSlowPostProtectionSettingRequest`
      * `GetSlowPostProtectionSettingResponse`
      * `GetSlowPostProtectionsRequest`
      * `GetWAFProtectionsRequest`
      * `RemoveNetworkLayerProtectionRequest`
      * `RemovePolicyProtectionsRequest`
      * `RemoveReputationProtectionRequest`

* [IMPORTANT] Added Bot Management API Support:
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

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added the `xff` field to the custom rule conditions.
  * Added the `NotificationEmails` field to the `Activation` structure.

* GTM
  * Improved error messages.

* CPS
  * Added cps ListEnrollments.
  * Extended `CreateEnrollment` with the `AllowDuplicateCN` option.

## 2.14.1 (July 26, 2022)

### BUG FIXES:

* IAM
  * Changed the IAM `GroupID` type to `int64`.

## 2.14.0 (June 28, 2022)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added the penalty box support for the security policy in evaluation mode.

* HAPI
  * EdgeHostname - update

* IAM
  * Blocked properties - read, update
  * Group - create, read, update, delete
  * Role - create, read, update, delete
  * User - lock, unlock, TFA, set password, reset password

### BUG FIXES:

* APPSEC
  * Fixed an incorrect error message on the activation failure.
  * The `EffectiveTimePeriod`, `SamplingRate`, `LoggingOptions`, and `Operation` fields of the various `CustomRule` response structures are now marshaled correctly.

## 2.13.0 (May 31, 2022)

### FEATURES/ENHANCEMENTS:

* Image and Video Manager
  * Added the new `ImQuery` transformation.
  * Added the new `PostBreakPointTransformationType`.

### BUG FIXES:

* Image and Video Manager
  * The `default_value` field on variable in image policy should not be required.
  * Changed all primitive optional parameters to pointers.
  * Corrected the `Anchor` field in `RectangleShapeType`.
  * Value field for `NumberVariableInline` should be defined as `float64`.
  * Renamed `PointShapeType.True` to `PointShapeType.Y`, to match the OpenAPI definition.
  * Added the `Composite` transformation to `PostBreakpointTransformations`.
  * Fixed `PostBreakpointTransformations.PolicyInputImage`.

## 2.12.0 (Apr. 25, 2022)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added the `WAPBypassNetworkLists` interface, to be used in preference to the deprecated `BypassNetworkLists` interface.

* Added support for the account switch keys from environment ([#149](https://github.com/akamai/AkamaiOPEN-edgegrid-golang/pull/149)).

## 2.11.0 (March 24, 2022)

### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Added Image and Video Manager API support
  * Policy Set - create, read, update, delete
  * Policy - create, read, update, delete, rollback to previous version, view policy history

* CLOUDLETS
  * Support for RC cloudlet type (Request Control)

* PAPI
  * CP code - read, update

## 2.10.0 (Feb. 28, 2022)

### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Added EdgeWorkers and EdgeKV API support:
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
  * Updated the source for the evasive path match interface with links to documentation.

* CLOUDLETS
  * Added support for the AS (Audience Segmentation) cloudlet type.

## 2.9.1 (Feb. 7, 2022)

### FEATURES/ENHANCEMENTS:

* APPSEC
  * Removed the deprecation notes for individual policy protection methods.

### BUG FIXES:

* CLOUDLETS
  * Fixed validation for the ALB version DataCenter percent.

## 2.9.0 (Jan. 24, 2022)

### FEATURES/ENHANCEMENTS:

* CLOUDLETS
  * Added support for VP cloudlet type (Visitor Prioritization).
  * Added support for CD cloudlet type (Continuous Deployment / Phased Release).
  * Added support for FR cloudlet type (Forward Rewrite).
  * Added support for AP cloudlet type (API Prioritization).

* APPSEC
  * Added support for Evasive Path Match feature.
  * Deprecated the individual policy protection interface methods.

* NETWORK LISTS
  * Included `ContractID` and `GroupID` in `GetNetworkListResponse`.

## 2.8.1 (Nov. 30, 2021)

### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Added Cloudlets API support:
  * Policy (Application Load Balancer) -  create, read, update, delete policy
  * Policy (Edge Redirector) -  create, read, update, delete policy
  * Policy activation - create, read
  * Application Load Balancer configuration - create, update, read
  * Activation for Application Load Balancer configuration - create, read

* APPSEC
  * Added support for advanced exceptions in ASE rules.
  * Updated the `bypass-network-list` data source and resource for the multi-policy WAP.

## 2.7.0 (Oct 19, 2021)

### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Added the DataStream API support:
  * Stream operations
  * Stream activation operations
  * Read access to various DataStream properties
* Added the HAPI v1 support:
  * Delete edge hostname

## 2.6.0 (Aug 16, 2021)

### BUG FIXES:

* APPSEC
  * Fixed incorrect comments/URL references in inline documentation.

### FEATURES/ENHANCEMENTS

* APPSEC
  * Got an evaluation attack group's or risk score group's action.

* NETWORK LISTS
  * Added support for `contract_id` and `group_id` for network list create/update.

## 2.5.0 (Jun 15, 2021)

### BREAKING CHANGES:

* APPSEC
  * Removed these packages along with their unit tests and test data:
    * `pkg/appsec/attack_group_action.go`
    * `pkg/appsec/attack_group_condition_exception.go`
    * `pkg/appsec/eval_rule_action.go`
    * `pkg/appsec/eval_rule_condition_exception.go`
    * `pkg/appsec/rule_action.go`
    * `pkg/appsec/rule_condition_exception.go`

### BUG FIXES:

* DNSv2
    * Fixed parsing SVCB, HTTPS rdata.

### FEATURES/ENHANCEMENTS:

* [IMPORTANT] CPS - Added Certificate Provisioning API support:
  * Enrollments - create, read, update, delete enrollments
  * Change status API - get change status, cancel change
  * DV certificate API - get and acknowledge DV challenges
  * Pre verification warnings - get and acknowledge pre verification warnings

* APPSEC
  * Removed these packages along with their unit tests and test data:
    * `pkg/appsec/api_constraints_protection.go`
    * `pkg/appsec/advanced_settings_pragma_header.go`
    * `pkg/appsec/attack_group.go`
    * `pkg/appsec/eval_rule.go`
    * `pkg/appsec/rule.go`
    * `pkg/appsec/ip_geo_protection.go`

## 2.4.1 (Apr 19, 2021)

### BUG FIXES:

* APPSEC
  * Suppressed the 'null' text on output of empty/false values.
  * Prevented configuration drift when reapplying configuration after importing or creating resources.

## 2.4.0 (Mar 29, 2021) PAPI - Secure by default

* PAPI
   * Added support to provision default certs as part of the hostname request.
   * Added a new cert status object in the hostname response if it exists.

## 2.3.0 (Mar 15, 2021) Network Lists

Added support for the these operations in the Network Lists API v2:

* Create a network list.
* Update an existing network list.
* Get the existing network lists, including optional filtering by name or type.
* Subscribe to a network list.
* Activate a network list.

## 2.2.1 (Mar 3, 2021)

* PAPI - Fixed an issue with rules causing advanced locked behaviors to fail.

## 2.2.0 (Feb 23, 2021) APPSEC - Extended list of supported list endpoints from APPSEC API

### BUG FIXES:

* PAPI
    * Fixed an issue with the version and rule comments being dropped.
    * Fixed client side validation to allow certain PAPI errors to pass through.

### FEATURES/ENHANCEMENTS:

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
    * Added support for HTTPS, SVCB records to ParseRData.

## 2.1.1 (Feb 3, 2021)

* PAPI - Fixed validation on empty rule behaviors causing some properties with nested behaviors to fail.

## 2.1.0 (Jan 13, 2021)

* [IMPORTANT] IAM - New Identity and Access Management API Support.

## 2.0.4 (Dec 23, 2020)

* APPSEC - Extended list of supported endpoints from the APPSEC API:
  * DDoS Protection - Rate Policy & Action
  * DDoS Protection - Slowpost setting & Action
  * Application Layer Protection - Rule Action, Exceptions & Conditions
  * Application Layer Protection - Rule Evaluation Action, Exceptions & Conditions
  * Application Layer Protection - Attack Group Action, Exceptions & Conditions
  * Application Layer Protection - Rule Upgrade & Change Mode for Rule Eval
  * Reputation Profile & Action
  * Network Layer Control - IP & GEO setting

## 2.0.3 (Dec 7, 2020)

* PAPI - fixed property hostname validation for missing hostnames.
* PAPI - fixed minor typos in rules error messages.

## 2.0.2 (Nov 19, 2020)

* [IMPORTANT] APPSEC - added the Application Security API.
* [ENHANCEMENT] DNS - added the Bulk API endpoints.
* ALL - re-enabled global account switch key support in the `.edgerc` files.
* PAPI - Edgehostname IPV6 support fix.  Added enums with allowed values.
* PAPI - Edgehostname blank cname or egdehostname id fix.
* PAPI - propertyversion blank etag field fix.

## 2.0.1 (Oct 15, 2020)

* [IMPORTANT] Breaking changes from earlier clients. Updated the library to use the v2 directory structure.
* [ENHANCEMENT] PAPI - API error returns to the user when an activation or validation error occurs.
* [NOTE] Reorganized the library to prepare for additional APIs to be included in future versions.

## 1.0.0 (Oct 15, 2020)

* Official release for the EdgeGrid Golang library.
* DNSv2 - Zone create signature to pass blank instead of nil.
* PAPI - Return nil instead of error if no cp code was found.
* GTM - Datacenter API requires blank instead of nil.

## 0.9.18 (Jul 13, 2020)

* [AT-40][Add] Preliminary Logging CorrelationID.

## 0.9.17 (Jun 9, 2020)

* Corrected AKAMAICDN target parsing.
* Addeded endpoints for listing zones, creating, and updating multiple recordsets.
* Refactored recordsets into a separate source file.

## 0.9.16 (May 29, 2020)

* Added updates to Client-v1, Papi-v1.
* Added a lock around the http request creation.
* PAPI - added logging to PAPI endpoints.

## 0.9.15 (May 15, 2020)

* DNSv2 - Added CERT and TSLA record parsing. Removed MX record parsing.

## 0.9.14 (May 12, 2020)

* DNSv2 - Enhanced the RecordError functions.

## 0.9.13 (Apr 26, 2020)

* DNSv2 - filterZoneCreate check upper case Type.

## 0.9.12 (Apr 21, 2020)

* DNSv2 - Added an optional `arg` to bypass dns record lock for the create, update, and delete functions. The default preserves the prior behavior.

## 0.9.11 (Apr 13 , 2020)

* DNSv2 Updates:
  * Added additional fields, including TSIG, to a zone.
  * Added support for alias zone types.
  * Added the utility functions for rdata parsing and process.
  * Added the `GetRecord` and `GetRecordSet` functions.
  * Add an additional recordset metadata.
* Added http request/response logging.

## 0.9.10 (Mar 5, 2020)

* Added support for caching Edgehostnames and Products.
* Added support for cache in PAPI library for edgehostnames and products to minimize round trips to fetch repeated common data to avoid the WAF deny rule IPBLOCK-BURST4-54013 issue.

## 0.9.9 (Feb 29, 2020)

* Added support for caching Contract, Groups, and Cp Codes.
* cache to minimize round trips on repeated common data fetches to avoid the WAF deny rule IPBLOCK-BURST4-54013 issue.

## 0.9.0 (Aug 6, 2019)

* Added support for GTM.
