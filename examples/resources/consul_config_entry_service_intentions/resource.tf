resource "consul_config_entry" "jwt_provider" {
  name = "okta"
  kind = "jwt-provider"

  config_json = jsonencode({
    ClockSkewSeconds = 30
    Issuer           = "test-issuer"
    JSONWebKeySet = {
      Remote = {
        URI                 = "https://127.0.0.1:9091"
        FetchAsynchronously = true
      }
    }
  })
}

resource "consul_config_entry_service_intentions" "web" {
  name = "web"

  jwt {
    providers {
      name = consul_config_entry.jwt_provider.name

      verify_claims {
        path  = ["perms", "role"]
        value = "admin"
      }
    }
  }

  sources {
    name   = "frontend-webapp"
    type   = "consul"
    action = "allow"
  }

  sources {
    name   = "nightly-cronjob"
    type   = "consul"
    action = "deny"
  }
}
