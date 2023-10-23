resource "consul_acl_policy" "read-policy" {
  name        = "read-policy"
  rules       = "node \"\" { policy = \"read\" }"
  datacenters = ["dc1"]
}

resource "consul_acl_role" "read" {
  name        = "foo"
  description = "bar"

  policies = [
    consul_acl_policy.read-policy.id
  ]

  service_identities {
    service_name = "foo"
  }
}
