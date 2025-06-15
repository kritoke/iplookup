#!/usr/bin/env sh
GOOS=windows GOARCH=amd64 go build -o bin/iplookup-amd64.exe main.go
GOOS=darwin GOARCH=arm64 go build -o bin/iplookup-arm64-darwin main.go
GOOS=linux GOARCH=amd64 go build -o bin/iplookup-amd64-linux main.go