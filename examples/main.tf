provider "consul" {
  address    = "127.0.0.1:8300"
  datacenter = "dc1"
}


resource "consul-yaml" "app" {
	datacenter = "dc1"

	path_prefix = "test/"
	subkeys_file = "subkeys.yaml"
}
