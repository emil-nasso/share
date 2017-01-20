#!/bin/bash
build()
{
  GOOS=$1 GOARCH=$2 go build -ldflags="-s -w" -o release/share-$1-$2$3
  upx -q release/share-$1-$2$3
}
rm -rf release
mkdir release

for GOOS in darwin linux; do
  for GOARCH in 386 amd64; do
    build $GOOS $GOARCH ""
  done
done

for GOARCH in 386 amd64; do
  build "windows" $GOARCH ".exe"
done

find release/* -exec ls -lhls  {} \; | grep release
