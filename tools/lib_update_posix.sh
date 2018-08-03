#!/usr/bin/env bash

set -e
script=$(readlink -f "$0")
route=$(dirname "$script")

cd ${route}/../
mkdir vendor
git clone https://github.com/echopairs/vendor.git vendor
go get -a github.com/golang/protobuf/protoc-gen-go