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
	GOBIN="$$(pwd)/bin" go install honnef.co/go/tools/cmd/staticcheck
	./bin/staticcheck ./pkg/...
