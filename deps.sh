#!/bin/bash
go get -u github.com/goware/modvendor
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.26.0
go mod vendor
go mod download
modvendor -copy="**/*.c **/*.h" -v
