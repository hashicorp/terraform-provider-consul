resource "consul-yaml" "test" {
  path_prefix  = "mysettings/"
  subkeys_file = "subkeys.yaml"
}
