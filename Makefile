.PHONY: build test e2e

test:
	go test -v ./...

build: dist/dnsbench
dist/dnsbench:
	go build -o dist/dnsbench dnsbench.go

e2e: dist/dnsbench
	./dist/functional-tests.sh
