resource "consul_config_entry_service_defaults" "dashboard" {
  name = "dashboard"

  upstream_config {
    defaults = {
      mesh_gateway = {
        mode = "local"
      }

      limits = {
        max_connections         = 512
        max_pending_requests    = 512
        max_concurrent_requests = 512
      }
    }

    overrides {
      name = "counting"

      mesh_gateway {
        mode = "remote"
      }
    }
  }

}
