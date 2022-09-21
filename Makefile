all: build test bench race staticcheck

build:
	go build ./...

test:
	go test -v ./...

bench:
	go test -v -bench=. ./...

race:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

bin/staticcheck:
	GOBIN=$$(pwd)/bin go install honnef.co/go/tools/cmd/staticcheck@latest

staticcheck: bin/staticcheck
	./bin/staticcheck ./...
