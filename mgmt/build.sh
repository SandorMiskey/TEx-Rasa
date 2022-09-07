#!/bin/sh

env GOOS=linux go build main.go -ldflags="-extldflags=-static" -tags sqlite_omit_load_extension
mv main storage/mgmt/