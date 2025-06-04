#!/bin/bash

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <version-tag>"
    exit 1
fi

TAG=$1

mkdir -p bin

# linux
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/fileshare .
tar -cvf "fileshare-$TAG-linux.tar.gz" bin/fileshare

# windows
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o bin/fileshare.exe .
tar -cvf "fileshare-$TAG-windows.tar.gz" bin/fileshare.exe

unset CC
unset CXX

rm bin/fileshare bin/fileshare.exe

echo "Building docker image"

docker build -t fileshare:$TAG .

docker rmi $(docker images --filter "dangling=true" -q --no-trunc)

docker save fileshare:$TAG > fileshare-$TAG.docker.zip
