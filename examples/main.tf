resource "consul-yaml_key_prefix_from_file" "test" {
  path_prefix  = "mysettings/"
  subkeys_file = "subkeys.yaml"
}
