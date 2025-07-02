#!/bin/bash

set -e

PKG_LIST=$(go list ./... | grep -v /vendor/ | xargs) 
go test -covermode=count -coverprofile .coverage $PKG_LIST || true
uncover -min 85 .coverage
