build:
	@go build ./cmd/pocket_core/main.go

test:
	@go test ./tests/*/

install:
	@dep ensure
