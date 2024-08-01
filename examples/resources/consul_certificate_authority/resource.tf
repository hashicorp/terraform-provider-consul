# Using the built-in CA with specific TTL
resource "consul_certificate_authority" "connect" {
  connect_provider = "consul"

  config_json = jsonencode({
    LeafCertTTL         = "24h"
    RotationPeriod      = "2160h"
    IntermediateCertTTL = "8760h"
  })
}


# Using Vault to manage and sign certificates
resource "consul_certificate_authority" "connect" {
  connect_provider = "vault"

  config_json = jsonencode({
    Address             = "http://localhost:8200"
    Token               = "..."
    RootPKIPath         = "connect-root"
    IntermediatePKIPath = "connect-intermediate"
  })
}


# Using the AWS Certificate Manager Private Certificate Authority
#  * https://aws.amazon.com/certificate-manager/private-certificate-authority/
resource "consul_certificate_authority" "connect" {
  connect_provider = "aws-pca"

  config_json = jsonencode({
    ExistingARN = "arn:aws:acm-pca:region:account:certificate-authority/12345678-1234-1234-123456789012"
  })
}
