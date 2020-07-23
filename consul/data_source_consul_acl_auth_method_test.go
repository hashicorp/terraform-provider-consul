package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataACLAuthMethod_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceACLAuthMethodConfigNotFound,
				ExpectError: regexp.MustCompile("Could not find auth-method 'not-found'"),
			},
			{
				Config: testAccDataSourceACLAuthMethodConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "name", "minikube"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "display_name", "Minikube Auth Method"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "type", "kubernetes"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "description", "dev minikube cluster"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "max_token_ttl", "2m0s"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "token_locality", "global"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "config.%", "3"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "config_json", `{"CACert":"-----BEGIN CERTIFICATE-----\nMIIBsTCCARoCCQCaNE5FiX2XdjANBgkqhkiG9w0BAQsFADAdMQswCQYDVQQGEwJG\nUjEOMAwGA1UECAwFUGFyaXMwHhcNMTkwNjI4MTA0ODUzWhcNMjAwNjI3MTA0ODUz\nWjAdMQswCQYDVQQGEwJGUjEOMAwGA1UECAwFUGFyaXMwgZ8wDQYJKoZIhvcNAQEB\nBQADgY0AMIGJAoGBAK4fNg9Hzq7Q87an4wgKcHWP97clnRTlozrUuV/WLQyKzS47\nISHM0x1Iy9b8VuIFidjS7cz9YB9nAUrV4rrzeBe08hDOGPAUsSUDMGFH7g2E7YYZ\nSfLJdoTo/qzCpU5lPG7iD11iplHgHYxSp+N3y8oL7wAiWfSd6GOoBiDNIb37AgMB\nAAEwDQYJKoZIhvcNAQELBQADgYEAReAfqISFnIJCbLzUVmZDHQUQqL4mck4nnJ8v\ngjFdDL52hG0jduSKll0qDdj54nnPKJBXv6/Q0HgY5UTa/YhqJmL2D4J2TgRL7au8\nsNDoqAR38rv33fExReu+VEaz9nrMIwnrPKm/4A3cViAkp7t9r1FjYAkBqakGy1S2\n/CfvsqQ=\n-----END CERTIFICATE-----\n\n","Host":"https://192.0.2.42:8443","ServiceAccountJWT":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}`),
				),
			},
			{
				Config: testAccDataSourceACLAuthMethodConfigBasicConfigJSON,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "name", "auth_method"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "type", "jwt"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "description", ""),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "config.%", "0"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "config_json", `{"BoundIssuer":"corp-issuer","ClaimMappings":{"http://example.com/first_name":"first_name","http://example.com/last_name":"last_name"},"JWTValidationPubKeys":["-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAryQICCl6NZ5gDKrnSztO\n3Hy8PEUcuyvg/ikC+VcIo2SFFSf18a3IMYldIugqqqZCs4/4uVW3sbdLs/6PfgdX\n7O9D22ZiFWHPYA2k2N744MNiCD1UE+tJyllUhSblK48bn+v1oZHCM0nYQ2NqUkvS\nj+hwUU3RiWl7x3D2s9wSdNt7XUtW05a/FXehsPSiJfKvHJJnGOX0BgTvkLnkAOTd\nOrUZ/wK69Dzu4IvrN4vs9Nes8vbwPa/ddZEzGR0cQMt0JBkhk9kU/qwqUseP1QRJ\n5I1jR4g8aYPL/ke9K35PxZWuDp3U0UPAZ3PjFAh+5T+fc7gzCs9dPzSHloruU+gl\nFQIDAQAB\n-----END PUBLIC KEY-----"],"ListClaimMappings":{"http://example.com/groups":"groups"}}`),
				),
			},
		},
	})
}

func TestAccDataACLAuthMethod_namespaceCE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceACLAuthMethodConfigNamespaceCE,
				ExpectError: regexp.MustCompile("Namespaces is a Consul Enterprise feature"),
			},
		},
	})
}

func TestAccDataACLAuthMethod_namespaceEE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceACLAuthMethodConfigNamespaceEE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "namespace", "test-data-auth-method"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "namespace_rule.#", "0"),
				),
			},
			{
				Config: testAccDataSourceACLAuthMethodConfigNamespaceEE_namespaceRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "namespace", "default"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "namespace_rule.#", "1"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "namespace_rule.0.selector", "serviceaccount.namespace==default and serviceaccount.name!=vault"),
					testAccCheckDataSourceValue("data.consul_acl_auth_method.test", "namespace_rule.0.bind_namespace", "prefixed-${serviceaccount.name}"),
				),
			},
		},
	})
}

const testAccDataSourceACLAuthMethodConfigNotFound = `
data "consul_acl_auth_method" "test" {
	name = "not-found"
}
`

const testAccDataSourceACLAuthMethodConfigBasic = `
resource "consul_acl_auth_method" "test" {
	name           = "minikube"
	display_name   = "Minikube Auth Method"
    type           = "kubernetes"
	description    = "dev minikube cluster"
	max_token_ttl  = "120s"
	token_locality = "global"

	config = {
        Host = "https://192.0.2.42:8443"
		CACert = <<-EOF
` + testCert + `
		EOF
        ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
    }
}

data "consul_acl_auth_method" "test" {
	name = consul_acl_auth_method.test.name
}
`

const testAccDataSourceACLAuthMethodConfigBasicConfigJSON = `
resource "consul_acl_auth_method" "test" {
	name        = "auth_method"
    type        = "jwt"

	config_json = jsonencode({
		BoundIssuer = "corp-issuer"
		JWTValidationPubKeys = [
			"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAryQICCl6NZ5gDKrnSztO\n3Hy8PEUcuyvg/ikC+VcIo2SFFSf18a3IMYldIugqqqZCs4/4uVW3sbdLs/6PfgdX\n7O9D22ZiFWHPYA2k2N744MNiCD1UE+tJyllUhSblK48bn+v1oZHCM0nYQ2NqUkvS\nj+hwUU3RiWl7x3D2s9wSdNt7XUtW05a/FXehsPSiJfKvHJJnGOX0BgTvkLnkAOTd\nOrUZ/wK69Dzu4IvrN4vs9Nes8vbwPa/ddZEzGR0cQMt0JBkhk9kU/qwqUseP1QRJ\n5I1jR4g8aYPL/ke9K35PxZWuDp3U0UPAZ3PjFAh+5T+fc7gzCs9dPzSHloruU+gl\nFQIDAQAB\n-----END PUBLIC KEY-----"
		]
		ClaimMappings = {
			"http://example.com/first_name" = "first_name"
			"http://example.com/last_name" = "last_name"
		}
		ListClaimMappings = {
			"http://example.com/groups" = "groups"
		}
	})
}

data "consul_acl_auth_method" "test" {
	name = consul_acl_auth_method.test.name
}`

const testAccDataSourceACLAuthMethodConfigNamespaceCE = `
data "consul_acl_auth_method" "test" {
  name      = "not-found"
  namespace = "test-data-auth-method"
}
`

const testAccDataSourceACLAuthMethodConfigNamespaceEE = `
resource "consul_acl_auth_method" "test" {
  name        = "minikube"
  type        = "kubernetes"
  description = "dev minikube cluster"
  namespace   = consul_namespace.test.name

  config = {
    Host = "https://192.0.2.42:8443"
    CACert = <<-EOF
` + testCert + `
    EOF
    ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
  }
}

resource "consul_namespace" "test" {
  name = "test-data-auth-method"
}

data "consul_acl_auth_method" "test" {
  name      = consul_acl_auth_method.test.name
  namespace = consul_acl_auth_method.test.namespace
}
`

const testAccDataSourceACLAuthMethodConfigNamespaceEE_namespaceRule = `
resource "consul_acl_auth_method" "test" {
  name        = "minikube"
  type        = "kubernetes"
  description = "dev minikube cluster"
  namespace   = "default"

  namespace_rule {
  	selector       = "serviceaccount.namespace==default and serviceaccount.name!=vault"
  	bind_namespace = "prefixed-$${serviceaccount.name}"
  }

  config = {
    Host = "https://192.0.2.42:8443"
    CACert = <<-EOF
` + testCert + `
    EOF
    ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
  }
}

data "consul_acl_auth_method" "test" {
	name      = consul_acl_auth_method.test.name
	namespace = consul_acl_auth_method.test.namespace
  }
`
