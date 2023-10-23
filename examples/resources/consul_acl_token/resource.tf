# Basic usage

resource "consul_acl_policy" "agent" {
  name  = "agent"
  rules = <<-RULE
    node_prefix "" {
      policy = "read"
    }
    RULE
}

resource "consul_acl_token" "test" {
  description = "my test token"
  policies    = [consul_acl_policy.agent.name]
  local       = true
}

# Explicitly set the `accessor_id`

resource "random_uuid" "test" {}

resource "consul_acl_token" "test_predefined_id" {
  accessor_id = random_uuid.test_uuid.result
  description = "my test uuid token"
  policies    = [consul_acl_policy.agent.name]
  local       = true
}
