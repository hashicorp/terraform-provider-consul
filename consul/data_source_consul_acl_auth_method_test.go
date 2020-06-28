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
