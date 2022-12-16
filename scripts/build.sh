#!/usr/bin/bash

set -e

cd ../


go mod download
go build -o ./main ./cmd/app
./main