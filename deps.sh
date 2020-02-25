#!/bin/bash
go get -u github.com/goware/modvendor
go mod vendor
go mod download
modvendor -copy="**/*.c **/*.h" -v
