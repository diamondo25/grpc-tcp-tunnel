#!/usr/bin/env bash

set -e

rm -f *.pb.go

protoc -I "$GOPATH" --proto_path . *.proto --go_out=plugins=grpc:.
