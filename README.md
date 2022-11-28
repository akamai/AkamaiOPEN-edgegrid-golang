# Akamai OPEN EdgeGrid for GoLang v3

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
