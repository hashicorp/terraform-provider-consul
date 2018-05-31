package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/schoology/terraform-provider-consul/consul"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: consul.Provider})
}
