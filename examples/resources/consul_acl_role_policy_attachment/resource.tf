data "consul_acl_role" "my_role" {
  name = "my_role"
}

resource "consul_acl_policy" "read_policy" {
  name        = "read-policy"
  rules       = "node \"\" { policy = \"read\" }"
  datacenters = ["dc1"]
}

resource "consul_acl_role_policy_attachment" "my_role_read_policy" {
  role_id = data.consul_acl_role.test.id
  policy  = consul_acl_policy.read_policy.name
}
