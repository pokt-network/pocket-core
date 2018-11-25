default: build

build:
	@go build -o core ./cmd/pocket_core/main.go

test:
	@go test ./tests/*/

install:
	@dep ensure
