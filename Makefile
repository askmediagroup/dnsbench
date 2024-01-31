export PATH := $(abspath bin/):${PATH}

.PHONY: clean
clean:
	rm -rf dist/

LD_FLAGS=-ldflags " \
    -X github.com/askmediagroup/dnsbench/pkg/cmd.dnsbenchVersion=$(shell git describe --tags --dirty --broken) \
    -X github.com/askmediagroup/dnsbench/pkg/cmd.gitCommit=$(shell git rev-parse HEAD) \
    -X github.com/askmediagroup/dnsbench/pkg/cmd.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
	"
.PHONY: build
build: dist/dnsbench
dist/dnsbench:
	go build $(LD_FLAGS) -o dist/dnsbench cmd/dnsbench/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: e2e
e2e: dist/dnsbench
	"$(CURDIR)/hack/functional-tests.sh"

.PHONY: check
checks: fmt test e2e license

.PHONY: license
license: bin/licensei
	licensei cache
	licensei check
	licensei header

.PHONY: fmt
fmt:
	"$(CURDIR)/hack/gofmt.sh"

LICENSEI_VERSION = 0.9.0
bin/licensei: bin/licensei-${LICENSEI_VERSION}
	@ln -sf licensei-${LICENSEI_VERSION} bin/licensei
bin/licensei-${LICENSEI_VERSION}:
	@mkdir -p bin
	curl -sfL https://git.io/licensei | bash -s v${LICENSEI_VERSION}
	@mv bin/licensei $@
