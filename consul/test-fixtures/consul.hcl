limits = {
  http_max_conns_per_client = -1
}

acl = {
  enabled        = true
  default_policy = "allow"
  down_policy    = "extend-cache"

  tokens = {
    initial_management = "12345678-1234-1234-1234-1234567890ab"
  }
}
