#!/bin/sh

# ------------------------------------
# Purpose:
# - Build binaries for supported target systems.
#
# Releases:
# - v1.0.0 - 2025/02/20: initial release
# ------------------------------------

set -v -o verbose

# recreate directory
rm -r ./binaries
mkdir ./binaries

# renew vendor content
go mod vendor

# lint
golangci-lint run --no-config --enable gocritic

# check for known vulnerabilities
govulncheck ./...

# show compiler version
go version

# compile 'darwin' (macOS)
env GOOS=darwin GOARCH=arm64 go build -v -o binaries/darwin-arm64/gemini-prompt
env GOOS=darwin GOARCH=amd64 go build -v -o binaries/darwin-amd64/gemini-prompt

# compile 'linux'
env GOOS=linux GOARCH=amd64 go build -v -o binaries/linux-amd64/gemini-prompt
env GOOS=linux GOARCH=arm64 go build -v -o binaries/linux-arm64/gemini-prompt

# compile 'windows'
env GOOS=windows GOARCH=amd64 go build -v -o binaries/windows-amd64/gemini-prompt.exe
env GOOS=windows GOARCH=arm go build -v -o binaries/windows-arm/gemini-prompt.exe

# compile 'freebsd'
env GOOS=freebsd GOARCH=amd64 go build -v -o binaries/freebsd-amd64/gemini-prompt
env GOOS=freebsd GOARCH=arm64 go build -v -o binaries/freebsd-arm64/gemini-prompt

# compile 'openbsd'
env GOOS=openbsd GOARCH=amd64 go build -v -o binaries/openbsd-amd64/gemini-prompt
env GOOS=openbsd GOARCH=arm64 go build -v -o binaries/openbsd-arm64/gemini-prompt

# compile 'netbsd'
env GOOS=netbsd GOARCH=amd64 go build -v -o binaries/netbsd-amd64/gemini-prompt
env GOOS=netbsd GOARCH=arm64 go build -v -o binaries/netbsd-arm64/gemini-prompt

