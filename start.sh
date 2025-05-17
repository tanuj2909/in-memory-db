#!/bin/bash

set -e
go build -o "./bin/app" app/*.go
exec ./bin/app "$@"
