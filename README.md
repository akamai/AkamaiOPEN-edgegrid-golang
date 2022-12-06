# Akamai OPEN EdgeGrid for GoLang v3

![Build Status](https://github.com/akamai/akamaiOPEN-edgegrid-golang/actions/workflows/checks.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/akamaiOPEN-edgegrid-golang/v3)](https://goreportcard.com/report/github.com/akamai/akamaiOPEN-edgegrid-golang/v3)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/akamai/akamaiOPEN-edgegrid-golang)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://pkg.go.dev/badge/github.com/akamai/akamaiOPEN-edgegrid-golang?utm_source=godoc)](https://pkg.go.dev/github.com/akamai/AkamaiOPEN-edgegrid-golang/v3)

This module is presently in active development and provides Akamai REST API support for the Akamai Terraform Provider.

## Backward Compatibility

This module is not backward compatible with the version `v1`.

Originally branch `master` was representing version `v1`. Now it is representing latest version `v3` and version `v1`
was moved to dedicated `v1` branch.

## Concurrent Usage

The packages of library can be imported alongside the `v1` library versions without conflict, for example:

```
import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	papiv3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
)
```
