build:
	goreleaser build --snapshot --clean

test:
	go test ./...

check:
	golangci-lint run

release:
	goreleaser release --snapshot --clean

.PHONY: build test check release
