Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by the Consul team at [HashiCorp](https://www.hashicorp.com/).

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.16 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/hashicorp/terraform-provider-consul`

```sh
$ mkdir -p $GOPATH/src/github.com/hashicorp; cd $GOPATH/src/github.com/hashicorp
$ git clone git@github.com:hashicorp/terraform-provider-consul
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/hashicorp/terraform-provider-consul
$ make build
```

Developing the Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.16+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-consul
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
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
$ make test-serv
```

By default, this will use the latest version of Consul based on the latest
image in the Docker repository. You can specify a version with the following:

```sh
$ CONSUL_VERSION=1.0.1 make test-serv
```

This command will run attached and will stop Consul when
interrupted. Images will be cached locally by Docker so it is quicky to
restart the server as necessary. This will expose Consul on the default
adddress.

Nightly acceptance tests are run against the `latest` tag of the Consul
Docker image. To run the acceptance tests against a development
version of Consul, you can [compile it](https://github.com/hashicorp/consul#developing-consul)
locally and then run it in development mode:

```
$ consul agent -dev
```
