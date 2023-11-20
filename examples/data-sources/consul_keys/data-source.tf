data "consul_keys" "app" {
  datacenter = "nyc1"

  # Read the launch AMI from Consul
  key {
    name    = "ami"
    path    = "service/app/launch_ami"
    default = "ami-1234"
  }
}

# Start our instance with the dynamic ami value
resource "aws_instance" "app" {
  ami = data.consul_keys.app.var.ami

  # ...
}
