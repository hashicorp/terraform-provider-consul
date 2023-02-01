# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Create a peering between the EU and US Consul clusters

provider "consul" {
  alias   = "eu"
  address = "eu-cluster:8500"
}

provider "consul" {
  alias   = "us"
  address = "us-cluster:8500"
}

resource "consul_peering_token" "eu-us" {
  provider  = consul.us
  peer_name = "eu-cluster"
}

resource "consul_peering" "eu-us" {
  provider = consul.eu

  peer_name     = "eu-cluster"
  peering_token = consul_peering_token.token.peering_token

  meta = {
    hello = "world"
  }
}
