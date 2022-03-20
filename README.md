<h1>
  <img src="./assets/logo.svg" align="left" height="46px" alt="Consul logo"/>
  <span>Consul Terraform Provider</span>
</h1>

[![Discuss](https://img.shields.io/badge/discuss-consul?logo=consul)](https://discuss.hashicorp.com/c/consul) [![Gitter chat](https://badges.gitter.im/hashicorp-consul/Lobby.png)](https://gitter.im/hashicorp-consul/Lobby)

- [Terraforrm Website](https://www.terraform.io/)
- [Consul Docs](https://www.consul.io/docs/intro)
- [Consul Terraform Provider Docs](https://www.terraform.io/docs/providers/consul/)

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
| 2.15.0                     | >= 1.11.0      |
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

Running the tests
-----------------

In order to test the provider, you can simply run `make test`.

```sh
make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.
This should be performed before merging or opening pull requests.

The acceptance test will automatically look for a `consul` binary in the computer
to start a development server. The binary used by the acceptance tests can be
specified by setting the `CONSUL_TEST_BINARY` environment variable.

```sh
make testacc
```

Documentation
-------------

Full, comprehensive documentation is available on the Terraform Registry:

<https://registry.terraform.io/providers/hashicorp/consul/latest/docs>

If you wish to contribute to the documentation, the source can be found in this
repository under website/docs/. To preview documentation changes prior to
submitting a pull request, please use the Terraform Registry's
[doc preview](https://registry.terraform.io/tools/doc-preview) tool.
