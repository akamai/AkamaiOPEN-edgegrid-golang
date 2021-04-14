# EDGEGRID GOLANG RELEASE NOTES

## 2.3.1

#### BUG FIXES

* APPSEC
  * Suppress 'null' text on output of empty/false values
  * Prevent configuration drift when reapplying configuration after importing or creating resources

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
