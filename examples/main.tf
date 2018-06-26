resource "consul_yaml" "test" {
  path_prefix  = "mysettings/"
  subkeys_file = "subkeys.yaml"
}
