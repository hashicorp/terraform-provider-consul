resource "consul_config_entry_service_resolver" "web" {
  name            = "web"
  default_subset  = "v1"
  connect_timeout = "15s"

  subsets {
    name   = "v1"
    filter = "Service.Meta.version == v1"
  }

  subsets {
    name   = "v2"
    Filter = "Service.Meta.version == v2"
  }

  redirect {
    service    = "web"
    datacenter = "dc2"
  }

  failover {
    subset_name = "v2"
    datacenters = ["dc2"]
  }

  failover {
    subset_name = "*"
    datacenters = ["dc3", "dc4"]
  }

}
