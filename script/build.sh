#!/bin/bash

build=$(date +%FT%T%z)
version="$1"

ldflags="-s -w -X github.com/devops-metalflow/metalflow/config.Build=$build -X github.com/devops-metalflow/metalflow/config.Version=$version"
target="metalflow"

go env -w GOPROXY=https://goproxy.cn,direct

# go tool dist list
CGO_ENABLED=0 GOARCH=$(go env GOARCH) GOOS=$(go env GOOS) go build -ldflags "$ldflags" -o bin/$target main.go

upx bin/$target
