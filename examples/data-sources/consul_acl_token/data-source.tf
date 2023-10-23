data "consul_acl_token" "test" {
  accessor_id = "00000000-0000-0000-0000-000000000002"
}

output "consul_acl_policies" {
  value = data.consul_acl_token.test.policies
}
