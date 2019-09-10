module github.com/terraform-providers/terraform-provider-consul

require (
	github.com/hashicorp/consul/api v1.1.0
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/terraform v0.12.8
	github.com/mitchellh/mapstructure v1.1.2
)

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.0
