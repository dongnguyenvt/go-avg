all: build test bench race staticcheck

build:
	go build ./pkg/...

test:
	go test -v ./pkg/...

bench:
	go test -v -bench=. ./pkg/...

race:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./pkg/...

staticcheck:
	go get honnef.co/go/tools/cmd/staticcheck
	staticcheck ./pkg/...
