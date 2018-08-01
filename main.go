package main

import (
	"github.com/hashicorp/terraform/plugin"
	"schoology/terraform-provider-consul-yaml/consulyaml"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: consulyaml.Provider})
}
