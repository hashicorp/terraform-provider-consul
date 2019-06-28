package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccConsulACLRole_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,

		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAuthMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLAuthMethodConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "name", "minikube"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "type", "kubernetes"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "description", "dev minikube cluster"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.%", "3"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.Host", "https://192.0.2.42:8443"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.CACert", testCert+"\n"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.ServiceAccountJWT", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"),
				),
			},
			{
				Config: testResourceACLAuthMethodConfigBasic_Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "name", "minikube2"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "type", "kubernetes"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "description", ""),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.%", "3"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.Host", "https://localhost:8443"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.CACert", testCert2+"\n"),
					resource.TestCheckResourceAttr("consul_acl_auth_method.test", "config.ServiceAccountJWT", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0IiwibmFtZSI6InRlc3QiLCJpYXQiOjE1MTYyMzkwMjJ9.uOnQsCs6ZAqj2F1VMA09tdgRZyFT1GQH2DwIC4TTn-A"),
				),
			},
		},
	})
}

func testAuthMethodDestroy(s *terraform.State) error {
	ACL := getClient(testAccProvider.Meta()).ACL()
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

const testResourceACLAuthMethodConfigBasic = `
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
}`

const testResourceACLAuthMethodConfigBasic_Update = `
resource "consul_acl_auth_method" "test" {
	name        = "minikube2"
    type        = "kubernetes"

	config = {
        Host = "https://localhost:8443"
		CACert = <<-EOF
` + testCert2 + `
		EOF
        ServiceAccountJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0IiwibmFtZSI6InRlc3QiLCJpYXQiOjE1MTYyMzkwMjJ9.uOnQsCs6ZAqj2F1VMA09tdgRZyFT1GQH2DwIC4TTn-A"
    }
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
