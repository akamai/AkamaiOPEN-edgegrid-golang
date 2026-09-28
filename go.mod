module github.com/akamai/AkamaiOPEN-edgegrid-golang/v14

go 1.26.0

toolchain go1.26.8

require (
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/google/uuid v1.6.0
	github.com/spf13/cast v1.10.0
	github.com/stretchr/testify v1.12.1
	github.com/wk8/go-ordered-map/v2 v2.1.8
	go.uber.org/ratelimit v0.3.1
	golang.org/x/net v0.59.0
	gopkg.in/ini.v1 v1.67.3
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/buger/jsonparser v1.6.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8
	github.com/mailru/easyjson v0.9.2 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/stretchr/testify v1.4.0 => github.com/stretchr/testify v1.10.0 // Fix security vulnerability
