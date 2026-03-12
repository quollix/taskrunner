#!/bin/bash

set -e

echo "building for linux"
GOOS=linux GOARCH=amd64 go build
echo "successfully built for linux"

echo "building for windows"
GOOS=windows GOARCH=amd64 go build
echo "successfully built for windows"

go test -tags=integration -count=1 ./...