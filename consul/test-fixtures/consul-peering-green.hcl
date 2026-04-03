# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

limits = {
  http_max_conns_per_client = -1
}

acl = {
  enabled        = true
  default_policy = "allow"
  down_policy    = "extend-cache"

  tokens = {
    initial_management = "12345678-1234-1234-1234-1234567890ab"
    replication        = "12345678-1234-1234-1234-1234567890ab"
  }
}

ports = {
  dns      = -1
  http     = 9500
  grpc     = 9502
  grpc_tls = 9503
  serf_lan = 9301
  serf_wan = 9302
  server   = 9300
}
