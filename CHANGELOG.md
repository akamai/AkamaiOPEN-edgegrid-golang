# EDGEGRID GOLANG RELEASE NOTES

## 1.0.1 (Jan 7, 2021)
* CPSv2 - Fixed several issues with listing enrollments

## 2.0.0 (Oct 15, 2020)
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
