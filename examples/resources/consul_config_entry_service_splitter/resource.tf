resource "consul_config_entry" "web" {
  name = "web"
  kind = "service-defaults"

  config_json = jsonencode({
    Protocol         = "http"
    Expose           = {}
    MeshGateway      = {}
    TransparentProxy = {}
  })
}

resource "consul_config_entry_service_resolver" "service_resolver" {
  name           = "service-resolver"
  default_subset = "v1"

  subsets {
    name   = "v1"
    filter = "Service.Meta.version == v1"
  }

  subsets {
    name   = "v2"
    Filter = "Service.Meta.version == v2"
  }
}

resource "consul_config_entry_service_splitter" "foo" {
  name = consul_config_entry_service_resolver.service_resolver.name

  meta = {
    key = "value"
  }

  splits {
    weight         = 80
    service        = "web"
    service_subset = "v1"

    request_headers {
      set = {
        "x-web-version" = "from-v1"
      }
    }

    response_headers {
      set = {
        "x-web-version" = "to-v1"
      }
    }
  }

  splits {
    weight         = 10
    service        = "web"
    service_subset = "v2"

    request_headers {
      set = {
        "x-web-version" = "from-v2"
      }
    }

    response_headers {
      set = {
        "x-web-version" = "to-v2"
      }
    }
  }

  splits {
    weight         = 10
    service        = "web"
    service_subset = "v2"
  }
}
