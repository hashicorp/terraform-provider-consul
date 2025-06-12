// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccConsulACLAuthMethod_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,

		CheckDestroy: testAuthMethodDestroy(client),
		Steps: []resource.TestStep{
			{
				Config:      testResourceACLAuthMethodConfigBasic_NoConfig,
				ExpectError: regexp.MustCompile("one of 'config' or 'config_json' must be set"),
			},
			{
				Config: testResourceACLAuthMethodConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "name", "minikube"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "display_name", "Minikube Auth Method"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "type", "kubernetes"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "description", "dev minikube cluster"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "max_token_ttl", "2m0s"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "token_locality", "global"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.%", "3"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.Host", "https://192.0.2.42:8443"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.CACert", testCert+"\n"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.ServiceAccountJWT", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"),
					testAuthMethodCACert(client, "minikube", testCert),
				),
			},
			{
				Config: testResourceACLAuthMethodConfigBasic_Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "name", "auth_method"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "display_name", ""),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "type", "kubernetes"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "max_token_ttl", "0s"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "token_locality", ""),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "description", ""),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.%", "3"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.Host", "https://localhost:8443"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.CACert", testCert2+"\n"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.ServiceAccountJWT", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0IiwibmFtZSI6InRlc3QiLCJpYXQiOjE1MTYyMzkwMjJ9.uOnQsCs6ZAqj2F1VMA09tdgRZyFT1GQH2DwIC4TTn-A"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config_json", "{\"CACert\":\"-----BEGIN CERTIFICATE-----\\nMIIBsTCCARoCCQCOgZn2+rDWSDANBgkqhkiG9w0BAQsFADAdMQswCQYDVQQGEwJG\\nUjEOMAwGA1UECAwFUGFyaXMwHhcNMTkwNjI4MTA1NzA4WhcNMjAwNjI3MTA1NzA4\\nWjAdMQswCQYDVQQGEwJGUjEOMAwGA1UECAwFUGFyaXMwgZ8wDQYJKoZIhvcNAQEB\\nBQADgY0AMIGJAoGBAMMBf+kSoZYon8fGBWqoyY7QzPXbg3GWMt2bxVxc6EmV/tcN\\nPIWGFFlycjnzDWwaGqzdqWkUrfi/o1VdlQobnzr4i+qcZpxlrZi2oa7FmkJMimsX\\nVmjXaeqpZA4JXLUzGHi+oCl2zX8wVGaUf7avcUxI3FVLCiibjWofpOf2pyUTAgMB\\nAAEwDQYJKoZIhvcNAQELBQADgYEAMddaDm4csxGnT47sths8CDxtzNdBhIXVIOLy\\njfvmBQ0aqC46gaUEoqNSzBPTTKJQGHxlGrF6fcnoUyjMcgHYZDrVySgmQpcfL9Uo\\nh61wQqlvkoFb/qPC/gvxdoQKUcddd7IhEujJjaddo9TV0w4nYX4Cq2Ybd5N3hgED\\n8GuzduY=\\n-----END CERTIFICATE-----\\n\\n\",\"Host\":\"https://localhost:8443\",\"ServiceAccountJWT\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0IiwibmFtZSI6InRlc3QiLCJpYXQiOjE1MTYyMzkwMjJ9.uOnQsCs6ZAqj2F1VMA09tdgRZyFT1GQH2DwIC4TTn-A\"}"),
					testAuthMethodCACert(client, "auth_method", testCert2),
				),
			},
			{
				Config: testResourceACLAuthMethodConfigBasic_ConfigJSON,
				Check:  testAuthMethodCACert(client, "auth_method", testCert),
			},
			{
				Config: testResourceACLAuthMethodConfigBasic_ConfigJSONUpdate,
				Check:  testAuthMethodCACert(client, "auth_method", testCert2),
			},
			{
				Config: testResourceACLAuthMethodConfigBasicConfigJSON,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "name", "auth_method"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "type", "jwt"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.%", "0"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config_json", "{\"BoundIssuer\":\"corp-issuer\",\"ClaimMappings\":{\"http://example.com/first_name\":\"first_name\",\"http://example.com/last_name\":\"last_name\"},\"JWTValidationPubKeys\":[\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAryQICCl6NZ5gDKrnSztO\\n3Hy8PEUcuyvg/ikC+VcIo2SFFSf18a3IMYldIugqqqZCs4/4uVW3sbdLs/6PfgdX\\n7O9D22ZiFWHPYA2k2N744MNiCD1UE+tJyllUhSblK48bn+v1oZHCM0nYQ2NqUkvS\\nj+hwUU3RiWl7x3D2s9wSdNt7XUtW05a/FXehsPSiJfKvHJJnGOX0BgTvkLnkAOTd\\nOrUZ/wK69Dzu4IvrN4vs9Nes8vbwPa/ddZEzGR0cQMt0JBkhk9kU/qwqUseP1QRJ\\n5I1jR4g8aYPL/ke9K35PxZWuDp3U0UPAZ3PjFAh+5T+fc7gzCs9dPzSHloruU+gl\\nFQIDAQAB\\n-----END PUBLIC KEY-----\"],\"ListClaimMappings\":{\"http://example.com/groups\":\"groups\"}}"),
				),
			},
			{
				Config: testResourceACLAuthMethod_interpolation,
			},
		},
	})
}

func TestAccConsulACLAuthMethod_namespaceCE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testResourceACLAuthMethodNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccConsulACLAuthMethod_namespaceEE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLAuthMethodNamespaceEE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "name", "minikube"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "namespace", "test-auth-method"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "namespace_rule.#", "0"),
				),
			},
			{
				Config: testResourceACLAuthMethodNamespaceEE_namespaceRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "name", "minikube"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "namespace", "default"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "namespace_rule.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "namespace_rule.0.selector", "serviceaccount.namespace==default and serviceaccount.name!=vault"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "namespace_rule.0.bind_namespace", "prefixed-${serviceaccount.name}"),
				),
			},
		},
	})
}

func testAuthMethodCACert(client *consulapi.Client, name, v string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		ACL := client.ACL()

		authMethod, _, err := ACL.AuthMethodRead(name, nil)
		if err != nil {
			return err
		}
		if authMethod == nil {
			return fmt.Errorf("Auth method %q does not exists", name)
		}

		cert := strings.TrimSpace(authMethod.Config["CACert"].(string))
		v = strings.TrimSpace(v)
		if cert != v {
			return fmt.Errorf("Wrong value for CACert: %q != %q", cert, v)
		}

		return nil
	}
}

func testAuthMethodDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		ACL := client.ACL()
		qOpts := &consulapi.QueryOptions{}

		role, _, err := ACL.AuthMethodRead("minikube2", qOpts)
		if err != nil {
			return err
		}

		if role != nil {
			return fmt.Errorf("Auth method 'minikube2' still exists")
		}

		return nil
	}
}

const testResourceACLAuthMethodConfigBasic_NoConfig = `
resource "consul_acl_auth_method" "test" {
	name        = "auth_method"
    type        = "kubernetes"
}`

const testResourceACLAuthMethodConfigBasic = `
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
}`

const testResourceACLAuthMethodConfigBasic_Update = `
resource "consul_acl_auth_method" "test" {
	name        = "auth_method"
    type        = "kubernetes"

	config_json  = jsonencode({
        Host = "https://localhost:8443"
		CACert = <<-EOF
` + testCert2 + `
		EOF
        ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0IiwibmFtZSI6InRlc3QiLCJpYXQiOjE1MTYyMzkwMjJ9.uOnQsCs6ZAqj2F1VMA09tdgRZyFT1GQH2DwIC4TTn-A"
    }
	)
}`

const testResourceACLAuthMethodConfigBasic_ConfigJSON = `
resource "consul_acl_auth_method" "test" {
	name        = "auth_method"
    type        = "kubernetes"

	config_json = jsonencode({
        Host = "https://localhost:8443"
		CACert = <<-EOF
` + testCert + `
		EOF
        ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0IiwibmFtZSI6InRlc3QiLCJpYXQiOjE1MTYyMzkwMjJ9.uOnQsCs6ZAqj2F1VMA09tdgRZyFT1GQH2DwIC4TTn-A"
    })
}`

const testResourceACLAuthMethodConfigBasic_ConfigJSONUpdate = `
resource "consul_acl_auth_method" "test" {
	name        = "auth_method"
    type        = "kubernetes"

	config_json = jsonencode({
        Host = "https://localhost:8443"
		CACert = <<-EOF
` + testCert2 + `
		EOF
        ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0IiwibmFtZSI6InRlc3QiLCJpYXQiOjE1MTYyMzkwMjJ9.uOnQsCs6ZAqj2F1VMA09tdgRZyFT1GQH2DwIC4TTn-A"
    })
}`

const testResourceACLAuthMethodConfigBasicConfigJSON = `
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
}`

const testCert = `-----BEGIN CERTIFICATE-----
MIIBsTCCARoCCQCaNE5FiX2XdjANBgkqhkiG9w0BAQsFADAdMQswCQYDVQQGEwJG
UjEOMAwGA1UECAwFUGFyaXMwHhcNMTkwNjI4MTA0ODUzWhcNMjAwNjI3MTA0ODUz
WjAdMQswCQYDVQQGEwJGUjEOMAwGA1UECAwFUGFyaXMwgZ8wDQYJKoZIhvcNAQEB
BQADgY0AMIGJAoGBAK4fNg9Hzq7Q87an4wgKcHWP97clnRTlozrUuV/WLQyKzS47
ISHM0x1Iy9b8VuIFidjS7cz9YB9nAUrV4rrzeBe08hDOGPAUsSUDMGFH7g2E7YYZ
SfLJdoTo/qzCpU5lPG7iD11iplHgHYxSp+N3y8oL7wAiWfSd6GOoBiDNIb37AgMB
AAEwDQYJKoZIhvcNAQELBQADgYEAReAfqISFnIJCbLzUVmZDHQUQqL4mck4nnJ8v
gjFdDL52hG0jduSKll0qDdj54nnPKJBXv6/Q0HgY5UTa/YhqJmL2D4J2TgRL7au8
sNDoqAR38rv33fExReu+VEaz9nrMIwnrPKm/4A3cViAkp7t9r1FjYAkBqakGy1S2
/CfvsqQ=
-----END CERTIFICATE-----
`

const testCert2 = `-----BEGIN CERTIFICATE-----
MIIBsTCCARoCCQCOgZn2+rDWSDANBgkqhkiG9w0BAQsFADAdMQswCQYDVQQGEwJG
UjEOMAwGA1UECAwFUGFyaXMwHhcNMTkwNjI4MTA1NzA4WhcNMjAwNjI3MTA1NzA4
WjAdMQswCQYDVQQGEwJGUjEOMAwGA1UECAwFUGFyaXMwgZ8wDQYJKoZIhvcNAQEB
BQADgY0AMIGJAoGBAMMBf+kSoZYon8fGBWqoyY7QzPXbg3GWMt2bxVxc6EmV/tcN
PIWGFFlycjnzDWwaGqzdqWkUrfi/o1VdlQobnzr4i+qcZpxlrZi2oa7FmkJMimsX
VmjXaeqpZA4JXLUzGHi+oCl2zX8wVGaUf7avcUxI3FVLCiibjWofpOf2pyUTAgMB
AAEwDQYJKoZIhvcNAQELBQADgYEAMddaDm4csxGnT47sths8CDxtzNdBhIXVIOLy
jfvmBQ0aqC46gaUEoqNSzBPTTKJQGHxlGrF6fcnoUyjMcgHYZDrVySgmQpcfL9Uo
h61wQqlvkoFb/qPC/gvxdoQKUcddd7IhEujJjaddo9TV0w4nYX4Cq2Ybd5N3hgED
8GuzduY=
-----END CERTIFICATE-----
`

const testResourceACLAuthMethodNamespaceCE = `
resource "consul_acl_auth_method" "test" {
	name        = "minikube"
    type        = "kubernetes"
	description = "dev minikube cluster"
	namespace   = "test"

	config = {
        Host = "https://192.0.2.42:8443"
		CACert = <<-EOF
` + testCert + `
		EOF
        ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
    }
}`

const testResourceACLAuthMethodNamespaceEE = `
resource "consul_namespace" "test" {
  name = "test-auth-method"
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
}`

const testResourceACLAuthMethodNamespaceEE_namespaceRule = `
resource "consul_namespace" "test" {
	name = "test-auth-method"
}

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
}`

const testResourceACLAuthMethod_interpolation = `
resource "consul_keys" "test" {
  datacenter = "dc1"

  key {
		path   = "test/service_account"
		value  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		delete = true
	}
}

data "consul_keys" "test" {
  datacenter = consul_keys.test.datacenter

  key {
		name = "service_account"
		path = "test/service_account"
	}
}

resource "consul_acl_auth_method" "auth_method" {
	name = "test"
	type = "kubernetes"

	config_json = jsonencode({
	    Host = "host"
	    CACert = <<-EOF
` + testCert + `
		EOF
	    ServiceAccountJWT = data.consul_keys.test.var.service_account
	})
  }
`
