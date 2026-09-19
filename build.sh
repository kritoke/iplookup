#!/usr/bin/env sh
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/iplookup-amd64.exe
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/iplookup-arm64-darwin
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/iplookup-amd64-linux