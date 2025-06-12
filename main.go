// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hashicorp/terraform-provider-consul/consul"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: consul.Provider})
}
