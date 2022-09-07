#!/bin/sh

env GOOS=linux go build -ldflags="-extldflags=-static" -tags sqlite_omit_load_extension,netgo main.go
mv main storage/mgmt/
