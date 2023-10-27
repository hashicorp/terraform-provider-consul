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
