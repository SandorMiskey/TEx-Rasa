#!/bin/sh

cd $ROOT
env GOOS=linux go build main.go
mv main storage/engine/
