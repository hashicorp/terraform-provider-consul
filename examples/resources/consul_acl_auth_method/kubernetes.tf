resource "consul_acl_auth_method" "minikube" {
  name        = "minikube"
  type        = "kubernetes"
  description = "dev minikube cluster"

  config_json = jsonencode({
    Host              = "https://192.0.2.42:8443"
    CACert            = "-----BEGIN CERTIFICATE-----\n...-----END CERTIFICATE-----\n"
    ServiceAccountJWT = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9..."
  })
}
