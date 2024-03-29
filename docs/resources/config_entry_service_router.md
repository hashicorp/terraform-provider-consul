---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "consul_config_entry_service_router Resource - terraform-provider-consul"
subcategory: ""
description: |-
  The consul_config_entry_service_router resource configures a service router https://developer.hashicorp.com/consul/docs/connect/config-entries/service-router to redirect a traffic request for a service to one or more specific service instances.
---

# consul_config_entry_service_router (Resource)

The `consul_config_entry_service_router` resource configures a [service router](https://developer.hashicorp.com/consul/docs/connect/config-entries/service-router) to redirect a traffic request for a service to one or more specific service instances.

## Example Usage

```terraform
resource "consul_config_entry_service_defaults" "admin_service_defaults" {
  name     = "web"
  protocol = "http"
}

resource "consul_config_entry_service_defaults" "admin_service_defaults" {
  name     = "dashboard"
  protocol = "http"
}


resource "consul_config_entry_service_router" "foo" {
  name = consul_config_entry.web.name

  routes {
    match {
      http {
        path_prefix = "/admin"
      }
    }

    destination {
      service = consul_config_entry.admin_service.name
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Specifies a name for the configuration entry.

### Optional

- `meta` (Map of String) Specifies key-value pairs to add to the KV store.
- `namespace` (String) Specifies the namespace to apply the configuration entry.
- `partition` (String) Specifies the admin partition to apply the configuration entry.
- `routes` (Block List) Defines the possible routes for L7 requests. (see [below for nested schema](#nestedblock--routes))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--routes"></a>
### Nested Schema for `routes`

Optional:

- `destination` (Block List, Max: 1) Specifies the target service to route matching requests to, as well as behavior for the request to follow when routed. (see [below for nested schema](#nestedblock--routes--destination))
- `match` (Block List, Max: 1) Describes a set of criteria that Consul compares incoming L7 traffic with. (see [below for nested schema](#nestedblock--routes--match))

<a id="nestedblock--routes--destination"></a>
### Nested Schema for `routes.destination`

Optional:

- `idle_timeout` (String) Specifies the total amount of time permitted for the request stream to be idle.
- `namespace` (String) Specifies the Consul namespace to resolve the service from instead of the current namespace.
- `num_retries` (Number) Specifies the number of times to retry the request when a retry condition occurs.
- `partition` (String) Specifies the Consul admin partition to resolve the service from instead of the current partition.
- `prefix_rewrite` (String) Specifies rewrites to the HTTP request path before proxying it to its final destination.
- `request_headers` (Block List, Max: 1) Specifies a set of HTTP-specific header modification rules applied to requests routed with the service router. (see [below for nested schema](#nestedblock--routes--destination--request_headers))
- `request_timeout` (String) Specifies the total amount of time permitted for the entire downstream request to be processed, including retry attempts.
- `response_headers` (Block List, Max: 1) Specifies a set of HTTP-specific header modification rules applied to responses routed with the service router. (see [below for nested schema](#nestedblock--routes--destination--response_headers))
- `retry_on` (List of String) Specifies a list of conditions for Consul to retry requests based on the response from an upstream service.
- `retry_on_connect_failure` (Boolean) Specifies that connection failure errors that trigger a retry request.
- `retry_on_status_codes` (List of Number) Specifies a list of integers for HTTP response status codes that trigger a retry request.
- `service` (String) Specifies the name of the service to resolve.
- `service_subset` (String) Specifies a named subset of the given service to resolve instead of the one defined as that service's `default_subset` in the service resolver configuration entry.

<a id="nestedblock--routes--destination--request_headers"></a>
### Nested Schema for `routes.destination.request_headers`

Optional:

- `add` (Map of String) Defines a set of key-value pairs to add to the header. Use header names as the keys.
- `remove` (List of String) Defines a list of headers to remove.
- `set` (Map of String) Defines a set of key-value pairs to add to the request header or to replace existing header values with.


<a id="nestedblock--routes--destination--response_headers"></a>
### Nested Schema for `routes.destination.response_headers`

Optional:

- `add` (Map of String) Defines a set of key-value pairs to add to the header. Use header names as the keys
- `remove` (List of String) Defines a list of headers to remove.
- `set` (Map of String) Defines a set of key-value pairs to add to the response header or to replace existing header values with



<a id="nestedblock--routes--match"></a>
### Nested Schema for `routes.match`

Optional:

- `http` (Block List, Max: 1) Specifies a set of HTTP criteria used to evaluate incoming L7 traffic for matches. (see [below for nested schema](#nestedblock--routes--match--http))

<a id="nestedblock--routes--match--http"></a>
### Nested Schema for `routes.match.http`

Optional:

- `header` (Block List) Specifies information in the HTTP request header to match with. (see [below for nested schema](#nestedblock--routes--match--http--header))
- `methods` (List of String) Specifies HTTP methods that the match applies to.
- `path_exact` (String) Specifies the exact path to match on the HTTP request path.
- `path_prefix` (String) Specifies the path prefix to match on the HTTP request path.
- `path_regex` (String) Specifies a regular expression to match on the HTTP request path.
- `query_param` (Block List) Specifies information to match to on HTTP query parameters. (see [below for nested schema](#nestedblock--routes--match--http--query_param))

<a id="nestedblock--routes--match--http--header"></a>
### Nested Schema for `routes.match.http.header`

Optional:

- `exact` (String) Specifies that a request matches when the header with the given name is this exact value.
- `invert` (Boolean) Specifies that the logic for the HTTP header match should be inverted.
- `name` (String) Specifies the name of the HTTP header to match.
- `prefix` (String) Specifies that a request matches when the header with the given name has this prefix.
- `present` (Boolean) Specifies that a request matches when the value in the `name` argument is present anywhere in the HTTP header.
- `regex` (String) Specifies that a request matches when the header with the given name matches this regular expression.
- `suffix` (String) Specifies that a request matches when the header with the given name has this suffix.


<a id="nestedblock--routes--match--http--query_param"></a>
### Nested Schema for `routes.match.http.query_param`

Optional:

- `exact` (String) Specifies that a request matches when the query parameter with the given name is this exact value.
- `name` (String) Specifies the name of the HTTP query parameter to match.
- `present` (Boolean) Specifies that a request matches when the value in the `name` argument is present anywhere in the HTTP query parameter.
- `regex` (String) Specifies that a request matches when the query parameter with the given name matches this regular expression.
