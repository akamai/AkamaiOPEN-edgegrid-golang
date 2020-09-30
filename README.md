# Akamai OPEN EdgeGrid for GoLang v2

This module is presently in development and provides Akamai REST API support for the Akamai Terraform Provider. It does **not** yet implement the full scope of Akamai endpoints. It is recommended to continue to use the the v1 API.

## Backward Compatibility
This module is not backward compatible with the previous version.

## Concurrent Usage
The packages of library can be imported alongside the v1 library versions without conflict, for example:

```
import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	papiv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
)
```
