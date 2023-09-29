export PATH := $(abspath bin/):${PATH}

.PHONY: clean
clean:
	rm -rf dist/

.PHONY: build
build: dist/dnsbench
dist/dnsbench:
	go build -o dist/dnsbench cmd/dnsbench/main.go

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
