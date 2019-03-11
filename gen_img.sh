#!/bin/bash

set -x

PREFIX=nexus.alphacario.com:8089
FLAG=testnet_stg

if [ -n "$1" ]; then
  FLAG=$1
fi

VER=`git rev-parse HEAD`

echo 'VER:'$VER

img_name=aci-blockchain-updater

if [ "testnet" != "$FLAG" ]; then
  img_name=$img_name-stg
fi

echo y | docker system prune

docker rmi $img_name

docker rmi $PREFIX/$img_name:v1

 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o updater core/main.go

docker build --no-cache -t $img_name .

docker tag $img_name:latest $PREFIX/$img_name:v1
