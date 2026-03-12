#!/bin/bash

set -e

go get -u ./...
go mod tidy

bash test.sh