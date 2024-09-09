#!/bin/bash
protoc --go_out=. --go_opt=module=github.com/inabinash/grpc --go-grpc_out=. --go-grpc_opt=module=github.com/inabinash/grpc maxapi/proto/maxapi.proto