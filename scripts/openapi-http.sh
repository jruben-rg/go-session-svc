#!/bin/bash
set -e

#Server
oapi-codegen --old-config-style -generate types -o server/openapi_types.gen.go -package server api/openapi/session.yml
oapi-codegen --old-config-style -generate chi-server -o server/openapi_api.gen.go -package server api/openapi/session.yml
#Client
oapi-codegen --old-config-style -generate types -o client/openapi_types.gen.go -package client api/openapi/session.yml
oapi-codegen --old-config-style -generate client -o client/openapi_client_gen.go -package client api/openapi/session.yml