#!/bin/sh

env GOOS=linux go build -ldflags="-extldflags=-static" -tags sqlite_omit_load_extension main.go
mv main storage/mgmt/
