ui = true
datacenter = "dc2"
primary_datacenter = "dc1"

limits = {
    http_max_conns_per_client = -1
}

acl = {
    enabled = true
    default_policy = "allow"
    down_policy = "extend-cache"

    tokens = {
        replication = "master-token"
    }
}

ports = {
    dns = -1
    grpc = -1
    http = 8501
    server = 8305
    serf_lan = 8306
    serf_wan = 8307
}
