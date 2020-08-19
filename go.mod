module github.com/hashicorp/terraform-provider-consul

require (
	github.com/apparentlymart/go-dump v0.0.0-20190214190832-042adf3cf4a0 // indirect
	github.com/aws/aws-sdk-go v1.22.0 // indirect
	github.com/hashicorp/consul/api v1.5.0
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-msgpack v0.5.4 // indirect
	github.com/hashicorp/hcl v0.0.0-20180906183839-65a6292f0157 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.0.0
	github.com/keybase/go-crypto v0.0.0-20180614160407-5114a9a81e1b // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/vmihailenco/msgpack v4.0.1+incompatible // indirect
)

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.0

go 1.12
