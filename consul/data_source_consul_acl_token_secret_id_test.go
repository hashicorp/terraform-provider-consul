package consul

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataACLTokenSecretID_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataACLTokenSecretIDConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_acl_token_secret_id.read", "pgp_key", ""),
					resource.TestCheckResourceAttr("data.consul_acl_token_secret_id.read", "encrypted_secret_id", ""),
					testAccCheckTokenExistsAndValidUUID("data.consul_acl_token_secret_id.read", "secret_id"),
				),
			},
		},
	})
}

func TestAccDataACLTokenSecretID_namespaceCE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccDataACLTokenSecretIDConfigNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccDataACLTokenSecretID_namespaceEE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataACLTokenSecretIDConfigNamespaceEE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_acl_token_secret_id.read", "pgp_key", ""),
					resource.TestCheckResourceAttr("data.consul_acl_token_secret_id.read", "encrypted_secret_id", ""),
					resource.TestCheckResourceAttr("data.consul_acl_token_secret_id.read", "namespace", "test-data-token-secret"),
					testAccCheckTokenExistsAndValidUUID("data.consul_acl_token_secret_id.read", "secret_id"),
				),
			},
		},
	})
}

func TestAccDataACLTokenSecretID_partitionCE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccDataACLTokenSecretIDConfigPartitionCE,
				ExpectError: regexp.MustCompile("(?i)Consul Enterprise feature"),
			},
		},
	})
}

func TestAccDataACLTokenSecretID_partitionEE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataACLTokenSecretIDConfigPartitionEE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_acl_token_secret_id.read", "pgp_key", ""),
					resource.TestCheckResourceAttr("data.consul_acl_token_secret_id.read", "encrypted_secret_id", ""),
					resource.TestCheckResourceAttr("data.consul_acl_token_secret_id.read", "partition", "test-data-token-secret"),
					testAccCheckTokenExistsAndValidUUID("data.consul_acl_token_secret_id.read", "secret_id"),
				),
			},
		},
	})
}

func TestAccDataACLTokenSecretID_PGP(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataACLTokenSecretIDPGPConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_acl_token_secret_id.read", "pgp_key", "keybase:remilapeyre"),
					resource.TestCheckResourceAttr("data.consul_acl_token_secret_id.read", "secret_id", ""),
				),
			},
		},
	})
}

func testAccCheckTokenExistsAndValidUUID(n string, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		secretID := rs.Primary.Attributes[attr]
		r := regexp.MustCompile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}")
		if !r.MatchString(secretID) {
			return fmt.Errorf("No valid UUID format %q", secretID)
		}
		return nil
	}
}

const testAccDataACLTokenSecretIDConfig = `
resource "consul_acl_policy" "test" {
	name = "test-data-token-secret"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies = [consul_acl_policy.test.name]
	local = true
}

data "consul_acl_token_secret_id" "read" {
    accessor_id = "${consul_acl_token.test.id}"
}
`

const testAccDataACLTokenSecretIDPGPConfig = `
resource "consul_acl_policy" "test" {
	name = "test"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies = ["${consul_acl_policy.test.name}"]
	local = true
}

data "consul_acl_token_secret_id" "read" {
	accessor_id = "${consul_acl_token.test.id}"
	pgp_key     = "keybase:remilapeyre"
}
`

const testAccDataACLTokenSecretIDConfigNamespaceCE = `
data "consul_acl_token" "read" {
  accessor_id = "foo"
  namespace   = "test-data-token"
}
`

const testAccDataACLTokenSecretIDConfigNamespaceEE = `
resource "consul_namespace" "test" {
  name = "test-data-token-secret"
}

resource "consul_acl_policy" "test" {
  name        = "test-data-token-secret"
  rules       = "node \"\" { policy = \"read\" }"
  datacenters = [ "dc1" ]
  namespace   = consul_namespace.test.name
}

resource "consul_acl_token" "test" {
  description = "test"
  policies    = [consul_acl_policy.test.name]
  local       = true
  namespace   = consul_namespace.test.name
}

data "consul_acl_token_secret_id" "read" {
  accessor_id = consul_acl_token.test.id
  namespace   = consul_namespace.test.name
}
`

const testAccDataACLTokenSecretIDConfigPartitionCE = `
data "consul_acl_token" "read" {
  accessor_id = "foo"
  partition   = "test-data-token"
}
`

const testAccDataACLTokenSecretIDConfigPartitionEE = `
resource "consul_admin_partition" "test" {
  name = "test-data-token-secret"
}

resource "consul_acl_policy" "test" {
  name        = "test-data-token-secret"
  rules       = "node \"\" { policy = \"read\" }"
  datacenters = [ "dc1" ]
  partition   = consul_admin_partition.test.name
}

resource "consul_acl_token" "test" {
  description = "test"
  policies    = [consul_acl_policy.test.name]
  local       = true
  partition   = consul_admin_partition.test.name
}

data "consul_acl_token_secret_id" "read" {
  accessor_id = consul_acl_token.test.id
  partition   = consul_admin_partition.test.name
}
`
