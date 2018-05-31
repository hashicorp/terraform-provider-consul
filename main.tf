resource "consul_key_prefix_from_file" "test_yaml" {
  path_prefix  = "meh"
  subkeys_file = "path/to/file"
}

resource "consul_key_prefix" "test" {
  path_prefix = "muhconfig"

  subkeys = {
    "elb_cname"      = "asdasd"
    "s3_bucket_name" = "11111"
    "database/name"  = "22222"
  }
}
