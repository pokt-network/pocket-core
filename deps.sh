#!/bin/bash
go get -u github.com/goware/modvendor
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.26.0
go mod vendor
go mod download
modvendor -copy="**/*.c **/*.h" -v
