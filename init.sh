#!/bin/bash

set -e

echo "Fetching known-good versions of dependencies..."

# We need a specific version of masterzen/winrm.
echo "Fetching github.com/masterzen/winrm (54ea5d0147)..."
go get github.com/masterzen/winrm
pushd $GOPATH/src/github.com/masterzen/winrm
git checkout 54ea5d01478cfc2afccec1504bd0dfcd8c260cfa
popd

# We need a specific version of packer-community/winrmcp.
echo "Fetching github.com/packer-community/winrmcp (f1bcf36a69)..."
mkdir -p $GOPATH/src/github.com/packer-community
git clone https://github.com/packer-community/winrmcp $GOPATH/src/github.com/packer-community/winrmcp
pushd $GOPATH/src/github.com/packer-community/winrmcp
git checkout f1bcf36a69fa2945e65dd099eee11b560fbd3346
popd

# We need a specific version of packer (actually, anything less than v9.0.0 will do).
echo "Fetching github.com/mitchellh/packer (v0.8.6)..."
go get github.com/mitchellh/packer
pushd $GOPATH/src/github.com/mitchellh/packer
git checkout v0.8.6
go get -d ./... || true
popd

echo "Done."
