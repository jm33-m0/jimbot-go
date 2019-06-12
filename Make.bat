@echo off
echo "Building for Linux x86_64 under Windows..."
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build