resource "consul-yaml" "app" {
	datacenter = "dc1"

	path_prefix = "prefix_test/"
	subkeys_file = "../consulyaml/test-fixtures/cheese.yam"
}
