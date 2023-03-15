#!/usr/bin/env bash
protoc -I=api --go_out=. --go-grpc_out=. proto/book.proto