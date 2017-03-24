#!/bin/bash
OUTPUT_BIN="jcsystool"


function release() {
  VERSION="$1"

  if [ -z "$VERSION" ]; then
    echo "Version number is required to release"
    exit 1
  else
    git tag -a $VERSION -m "Release version $VERSION"

    LINUX_BIN="$OUTPUT_BIN-linux-amd64-$VERSION"
    MAC_BIN="$OUTPUT_BIN-darwin-amd64-$VERSION"

    LINUX_OUT="./build/$LINUX_BIN.tar.gz"
    MAC_OUT="./build/$MAC_BIN.tar.gz"

    mkdir -p ./build

    env GOOS=linux GOARCH=amd64  go build -o $LINUX_BIN *.go
    tar -cvf $LINUX_OUT $LINUX_BIN
    rm $LINUX_BIN


    env GOOS=darwin GOARCH=amd64 go build -o $MAC_BIN *.go
    tar -cvf $MAC_OUT $MAC_BIN
    rm $MAC_BIN

    LINUX_SHA=$(shasum -a 256 $LINUX_OUT)
    MAC_SHA=$(shasum -a 256 $MAC_OUT)

    echo "Built Linux: $LINUX_OUT"
    echo $LINUX_SHA
    echo && echo
    echo "Built Mac: $MAC_OUT"
    echo $MAC_SHA
  fi
}

release $1
