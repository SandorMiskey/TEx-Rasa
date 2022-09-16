#!/bin/sh

# TODO: get agents directori from cli/command parameter...

go run main.go \
	-dbAddr ./storage/engine/mgmt.db \
	-httpStaticRoot ./storage/engine/public
