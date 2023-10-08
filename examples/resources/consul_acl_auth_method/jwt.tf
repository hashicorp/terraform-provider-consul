resource "consul_acl_auth_method" "oidc" {
  name          = "auth0"
  type          = "oidc"
  max_token_ttl = "5m"

  config_json = jsonencode({
    AllowedRedirectURIs = [
      "http://localhost:8550/oidc/callback",
      "http://localhost:8500/ui/oidc/callback"
    ]
    BoundAudiences = [
      "V1RPi2MYptMV1RPi2MYptMV1RPi2MYpt"
    ]
    ClaimMappings = {
      "http://example.com/first_name" = "first_name"
      "http://example.com/last_name"  = "last_name"
    }
    ListClaimMappings = {
      "http://consul.com/groups" = "groups"
    }
    OIDCClientID     = "V1RPi2MYptMV1RPi2MYptMV1RPi2MYpt"
    OIDCClientSecret = "...(omitted)..."
    OIDCDiscoveryURL = "https://my-corp-app-name.auth0.com/"
  })
}
