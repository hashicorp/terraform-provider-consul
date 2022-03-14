package consul

import (
	"fmt"
	"regexp"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulACLBindingRule_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testBindingRuleDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testResourceACLBindingRuleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "auth_method", "minikube"),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "description", "foobar"),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "selector", "serviceaccount.namespace==default"),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "bind_type", "service"),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "bind_name", "minikube"),
				),
			},
			{
				Config: testResourceACLBindingRuleConfigBasic_Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "auth_method", "minikube2"),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "description", ""),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "selector", "serviceaccount.namespace==default2"),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "bind_type", "role"),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "bind_name", "minikube2"),
				),
			},
			{
				Config: testResourceACLBindingRuleConfig_node,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "auth_method", "minikube2"),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "description", ""),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "selector", "serviceaccount.namespace==default2"),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "bind_type", "node"),
					resource.TestCheckResourceAttr("consul_acl_binding_rule.test", "bind_name", "minikube2"),
				),
			},
			{
				Config:      testResourceACLBindingRuleConfig_wrongType,
				ExpectError: regexp.MustCompile(`Invalid Binding Rule: unknown BindType "foobar"`),
			},
		},
	})
}

func TestAccConsulACLBindingRule_namespaceCE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testResourceACLBindingRuleConfig_namespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccConsulACLBindingRule_namespaceEE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLBindingRuleConfig_namespaceEE,
			},
		},
	})
}

func testBindingRuleDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		ACL := client.ACL()
		qOpts := &consulapi.QueryOptions{}

		rules, _, err := ACL.BindingRuleList("minikube2", qOpts)
		if err != nil {
			return err
		}

		if len(rules) != 0 {
			return fmt.Errorf("Binding rule of 'minikube2' still exists")
		}

		return nil
	}
}

const testResourceACLBindingRuleConfigBasic = `
resource "consul_acl_auth_method" "test" {
	name        = "minikube"
    type        = "kubernetes"
    description = "dev minikube cluster"

	config = {
        Host = "https://192.0.2.42:8443"
		CACert = <<-EOF
` + testCert + `
		EOF
        ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
    }
}

resource "consul_acl_binding_rule" "test" {
	auth_method = "${consul_acl_auth_method.test.name}"
	description = "foobar"
	selector    = "serviceaccount.namespace==default"
	bind_type   = "service"
	bind_name   = "minikube"
}`

const testResourceACLBindingRuleConfigBasic_Update = `
resource "consul_acl_auth_method" "test" {
	name        = "minikube2"
    type        = "kubernetes"
    description = "dev minikube cluster"

	config = {
        Host = "https://192.0.2.42:8443"
		CACert = <<-EOF
` + testCert + `
		EOF
        ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
    }
}

resource "consul_acl_binding_rule" "test" {
	auth_method = "${consul_acl_auth_method.test.name}"
	selector    = "serviceaccount.namespace==default2"
	bind_type   = "role"
	bind_name   = "minikube2"
}`

const testResourceACLBindingRuleConfig_namespaceCE = `
resource "consul_acl_auth_method" "test" {
	name        = "minikube"
    type        = "kubernetes"
    description = "dev minikube cluster"

	config = {
        Host = "https://192.0.2.42:8443"
		CACert = <<-EOF
` + testCert + `
		EOF
        ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
    }
}

resource "consul_acl_binding_rule" "test" {
	auth_method = "${consul_acl_auth_method.test.name}"
	description = "foobar"
	selector    = "serviceaccount.namespace==default"
	bind_type   = "service"
	bind_name   = "minikube"
	namespace   = "test-binding-rule"
}`

const testResourceACLBindingRuleConfig_namespaceEE = `
resource "consul_namespace" "test" {
  name = "test-binding-rule"
}

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

resource "consul_acl_binding_rule" "test" {
	auth_method = "${consul_acl_auth_method.test.name}"
	description = "foobar"
	selector    = "serviceaccount.namespace==default"
	bind_type   = "service"
	bind_name   = "minikube"
	namespace   = consul_namespace.test.name
}`

const testResourceACLBindingRuleConfig_node = `
resource "consul_acl_auth_method" "test" {
	name        = "minikube2"
    type        = "kubernetes"
    description = "dev minikube cluster"

	config = {
        Host = "https://192.0.2.42:8443"
		CACert = <<-EOF
` + testCert + `
		EOF
        ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
    }
}

resource "consul_acl_binding_rule" "test" {
	auth_method = "${consul_acl_auth_method.test.name}"
	selector    = "serviceaccount.namespace==default2"
	bind_type   = "node"
	bind_name   = "minikube2"
}`

const testResourceACLBindingRuleConfig_wrongType = `
resource "consul_acl_auth_method" "test" {
	name        = "minikube2"
    type        = "kubernetes"
    description = "dev minikube cluster"

	config = {
        Host = "https://192.0.2.42:8443"
		CACert = <<-EOF
` + testCert + `
		EOF
        ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
    }
}

resource "consul_acl_binding_rule" "test" {
	auth_method = "${consul_acl_auth_method.test.name}"
	selector    = "serviceaccount.namespace==default2"
	bind_type   = "foobar"
	bind_name   = "minikube2"
}`
