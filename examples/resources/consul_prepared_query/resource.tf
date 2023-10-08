# Creates a prepared query myquery.query.consul that finds the nearest
# healthy myapp.service.consul instance that has the active tag and not
# the standby tag.
resource "consul_prepared_query" "myapp-query" {
  name         = "myquery"
  datacenter   = "us-central1"
  token        = "abcd"
  stored_token = "wxyz"
  only_passing = true
  near         = "_agent"

  service = "myapp"
  tags    = ["active", "!standby"]

  failover {
    nearest_n   = 3
    datacenters = ["us-west1", "us-east-2", "asia-east1"]
  }

  dns {
    ttl = "30s"
  }
}

# Creates a Prepared Query Template that matches *-near-self.query.consul
# and finds the nearest service that matches the glob character (e.g.
# foo-near-self.query.consul will find the nearest healthy foo.service.consul).
resource "consul_prepared_query" "service-near-self" {
  datacenter   = "nyc1"
  token        = "abcd"
  stored_token = "wxyz"
  name         = ""
  only_passing = true
  connect      = true
  near         = "_agent"

  template {
    type   = "name_prefix_match"
    regexp = "^(.*)-near-self$"
  }

  service = "$${match(1)}"

  failover {
    nearest_n   = 3
    datacenters = ["dc2", "dc3", "dc4"]
  }

  dns {
    ttl = "5m"
  }
}
