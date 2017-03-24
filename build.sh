#!/bin/bash
OUTPUT_BIN="jcsystool"

LINUX_BIN="$OUTPUT_BIN-linux-amd64"
MAC_BIN="$OUTPUT_BIN-darwin-amd64"

LINUX_OUT="./build/$LINUX_BIN.tar.gz"
MAC_OUT="./build/$MAC_BIN.tar.gz"

mkdir -p ./build

env GOOS=linux GOARCH=amd64  go build -v -o $LINUX_BIN *.go
tar -cvf $LINUX_OUT $LINUX_BIN
rm $LINUX_BIN


env GOOS=darwin GOARCH=amd64 go build -v -o $MAC_BIN *.go
tar -cvf $MAC_OUT $MAC_BIN
rm $MAC_BIN

echo "Built Linux: $LINUX_OUT"
echo "Built Mac: $MAC_OUT"
