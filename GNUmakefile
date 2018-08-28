TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=consul
CONSUL_VERSION ?= "latest"

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

test-serv: fmtcheck
	@docker pull "consul:$(CONSUL_VERSION)"
	docker run --rm -p 127.0.0.1:8500:8500 "consul:$(CONSUL_VERSION)"

test-serv-acl-with-bootstrap: fmtcheck
	@docker pull "consul:$(CONSUL_VERSION)"
	docker run --rm -p 127.0.0.1:8500:8500 \
	    -e 'CONSUL_LOCAL_CONFIG={ "acl_datacenter": "dc1", "acl_master_token": "6b0de9ab-6d95-4af8-a965-78ca35a67018" }' \
		"consul:$(CONSUL_VERSION)" agent -server -client=0.0.0.0 -bootstrap-expect=1

test-serv-acl-without-bootstrap: fmtcheck
	@docker pull "consul:$(CONSUL_VERSION)"
	docker run --rm -p 127.0.0.1:8500:8500 \
		-e 'CONSUL_LOCAL_CONFIG={ "acl_datacenter": "dc1" }' \
		"consul:$(CONSUL_VERSION)" agent -server -client=0.0.0.0 -bootstrap-expect=1

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile test-serv website website-test

