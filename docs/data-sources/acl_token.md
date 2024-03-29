---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "consul_acl_token Data Source - terraform-provider-consul"
subcategory: ""
description: |-
  The consul_acl_token data source returns the information related to the consul_acl_token resource with the exception of its secret ID.
  If you want to get the secret ID associated with a token, use the consul_acl_token_secret_id data source.
---

# consul_acl_token (Data Source)

The `consul_acl_token` data source returns the information related to the `consul_acl_token` resource with the exception of its secret ID.

If you want to get the secret ID associated with a token, use the [`consul_acl_token_secret_id` data source](/docs/providers/consul/d/acl_token_secret_id.html).

## Example Usage

```terraform
data "consul_acl_token" "test" {
  accessor_id = "00000000-0000-0000-0000-000000000002"
}

output "consul_acl_policies" {
  value = data.consul_acl_token.test.policies
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `accessor_id` (String) The accessor ID of the ACL token.

### Optional

- `namespace` (String) The namespace to lookup the ACL token.
- `partition` (String) The partition to lookup the ACL token.

### Read-Only

- `description` (String) The description of the ACL token.
- `expiration_time` (String) If set this represents the point after which a token should be considered revoked and is eligible for destruction.
- `id` (String) The ID of this resource.
- `local` (Boolean) Whether the ACL token is local to the datacenter it was created within.
- `node_identities` (List of Object) The list of node identities attached to the token. (see [below for nested schema](#nestedatt--node_identities))
- `policies` (List of Object) A list of policies associated with the ACL token. (see [below for nested schema](#nestedatt--policies))
- `roles` (List of Object) List of roles linked to the token (see [below for nested schema](#nestedatt--roles))
- `service_identities` (List of Object) The list of service identities attached to the token. (see [below for nested schema](#nestedatt--service_identities))
- `templated_policies` (List of Object) The list of templated policies that should be applied to the token. (see [below for nested schema](#nestedatt--templated_policies))

<a id="nestedatt--node_identities"></a>
### Nested Schema for `node_identities`

Read-Only:

- `datacenter` (String)
- `node_name` (String)


<a id="nestedatt--policies"></a>
### Nested Schema for `policies`

Read-Only:

- `id` (String)
- `name` (String)


<a id="nestedatt--roles"></a>
### Nested Schema for `roles`

Read-Only:

- `id` (String)
- `name` (String)


<a id="nestedatt--service_identities"></a>
### Nested Schema for `service_identities`

Read-Only:

- `datacenters` (List of String)
- `service_name` (String)


<a id="nestedatt--templated_policies"></a>
### Nested Schema for `templated_policies`

Read-Only:

- `datacenters` (List of String)
- `template_name` (String)
- `template_variables` (List of Object) (see [below for nested schema](#nestedobjatt--templated_policies--template_variables))

<a id="nestedobjatt--templated_policies--template_variables"></a>
### Nested Schema for `templated_policies.template_variables`

Read-Only:

- `name` (String)
