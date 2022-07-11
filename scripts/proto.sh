#!/bin/bash
set -e
protoc -Iapi/protobuf --go_out=. --go_opt=module=github.com/jruben-rg/go-session-svc --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=module=github.com/jruben-rg/go-session-svc api/protobuf/session.proto