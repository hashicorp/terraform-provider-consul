Terraform Provider
==================

- Website: <https://www.terraform.io>
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by the Consul team at [HashiCorp](https://www.hashicorp.com/).

Compatibility
-------------

The Consul Terraform provider uses features of the latest version of Consul.
Some resources may not be supported by older versions of Consul.

The known compatibility between this provider and Consul is:

| Terraform provider version | Consul version |
| -------------------------- | -------------- |
| 2.14.0                     | >= 1.10.0      |
| 2.13.0                     | >= 1.10.0      |


Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.15

Building The Provider
---------------------

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```sh
go install
```

Developing the Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

In order to test the provider, you can simply run `make test`.

```sh
make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.
This should be performed before merging or opening pull requests.

```sh
$ consul agent -dev -config-file ./consul_test.hcl &
$ export CONSUL_HTTP_ADDR=localhost:8500
$ export CONSUL_HTTP_TOKEN=master-token
$ make testacc
```

Testing the resources specific to Consul Enterprise requires a running Consul
Enterprise server. It is possible to use the Consul Enterprise Docker image
which has a license valid for six hours during development:

```sh
$ docker run --rm \
             -d \
             --name consul-test \
             -v $PWD/consul_test.hcl:/consul_test.hcl:ro \
             -p 8500:8500 \
             hashicorp/consul-enterprise:latest consul agent -dev -config-file consul_test.hcl -client=0.0.0.0
$ export CONSUL_HTTP_ADDR=localhost:8500
$ export CONSUL_HTTP_TOKEN=master-token
$ make testacc
$ docker stop consul-test
```

Running the tests requires a running Consul agent locally. This provider targets
the latest version of Consul, but older versions should be compatible where
possible. In some cases, older versions of this provider will work with
older versions of Consul.

If you have [Docker](https://docs.docker.com/install/) installed, you can
run Consul with the following command:

```sh
make test-serv
```

By default, this will use the latest version of Consul based on the latest
image in the Docker repository. You can specify a version with the following:

```sh
CONSUL_VERSION=1.0.1 make test-serv
```

This command will run attached and will stop Consul when
interrupted. Images will be cached locally by Docker so it is quickly to
restart the server as necessary. This will expose Consul on the default
address.

Nightly acceptance tests are run against the `latest` tag of the Consul
Docker image. To run the acceptance tests against a development
version of Consul, you can [compile it](https://github.com/hashicorp/consul/blob/main/.github/CONTRIBUTING.md#building-consul)
locally and then run it in development mode:

```shell
consul agent -dev
```

It is also possible to run additional tests to test the provider with multiple
datacenters:

```sh
$ consul agent -dev -config-file ./consul_test_dc2.hcl &
$ consul agent -dev -config-file ./consul_test.hcl &
$ export CONSUL_HTTP_ADDR=localhost:8500
$ export CONSUL_HTTP_TOKEN=master-token
$ TEST_REMOTE_DATACENTER=1 make testacc
```

Documentation
-----------------------

Full, comprehensive documentation is available on the Terraform Registry:

<https://registry.terraform.io/providers/hashicorp/consul/latest/docs>

If you wish to contribute to the documentation, the source can be found in this
repository under website/docs/. To preview documentation changes prior to
submitting a pull request, please use the Terraform Registry's
[doc preview](https://registry.terraform.io/tools/doc-preview) tool.
