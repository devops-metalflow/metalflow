#!/bin/bash

ldflags="-s -w"
target="metalflow"

go env -w GOPROXY=https://goproxy.cn,direct

# go tool dist list
GIN_MODE=release CGO_ENABLED=0 GOARCH=$(go env GOARCH) GOOS=$(go env GOOS) go build -ldflags "$ldflags" -o bin/$target main.go

upx bin/$target
