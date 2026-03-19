#!/bin/bash

set -e

go build
go test -tags=integration -count=1 ./...
