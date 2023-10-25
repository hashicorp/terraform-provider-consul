data "consul_acl_role" "test" {
  name = "example-role"
}

output "consul_acl_role" {
  value = data.consul_acl_role.test.id
}
