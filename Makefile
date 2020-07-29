all: build test bench

build:
	go build ./pkg/...

test:
	go test -v ./pkg/...

bench:
	go test -v -bench=. ./pkg/...
