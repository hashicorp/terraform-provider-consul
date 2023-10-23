resource "consul_acl_auth_method" "minikube" {
  name        = "minikube"
  type        = "kubernetes"
  description = "dev minikube cluster"

  config = {
    Host              = "https://192.0.2.42:8443"
    CACert            = "-----BEGIN CERTIFICATE-----\n...-----END CERTIFICATE-----\n"
    ServiceAccountJWT = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9..."
  }
}

resource "consul_acl_binding_rule" "test" {
  auth_method = consul_acl_auth_method.minikube.name
  description = "foobar"
  selector    = "serviceaccount.namespace==default"
  bind_type   = "service"
  bind_name   = "minikube"
}
