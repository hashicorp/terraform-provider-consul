ui = true

limits = {
    http_max_conns_per_client = -1
}

acl = {
    enabled = true
    default_policy = "allow"
    down_policy = "extend-cache"

    tokens = {
        master = "master-token"
    }
}
