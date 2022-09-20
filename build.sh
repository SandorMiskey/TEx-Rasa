#!/bin/sh

# env GOOS=linux CGO_ENABLED=1 go build -ldflags="-extldflags=-static" -tags sqlite_omit_load_extension,netgo main.go
env GOOS=linux go build main.go
mv main storage/engine/
