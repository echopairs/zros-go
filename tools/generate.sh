#!/usr/bin/env bash

script=$(readlink -f "$0")
route=$(dirname "$script")

echo "generating rpc ..."
mkdir -p ${route}/../zros_rpc
mkdir -p ${route}/../zros_example

protoc -I ${route}/../zros_proto/ ${route}/../zros_proto/zros.proto --go_out=plugins=grpc:${route}/../zros_rpc
echo "generating examples ..."
protoc -I ${route}/../zros_proto/ ${route}/../zros_proto/test_message.proto --go_out=plugins=grpc:${route}/../zros_example
protoc -I ${route}/../zros_proto/ ${route}/../zros_proto/test_service.proto --go_out=plugins=grpc:${route}/../zros_example

echo "Done!"