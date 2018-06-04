package main

import (
	"schoology/terraform-provider-consul-yaml/consulyaml"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: consulyaml.Provider})
}
