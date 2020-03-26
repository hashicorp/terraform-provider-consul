module github.com/terraform-providers/terraform-provider-consul

require (
	github.com/hashicorp/consul/api v1.4.0
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/terraform v0.12.9 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
)

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.0

go 1.12
