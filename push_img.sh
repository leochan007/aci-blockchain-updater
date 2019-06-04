#!/bin/bash

set -x

source module_def.sh

if [ -n "$1" ]; then
  FLAG=$1
fi

if [ xstg = x"$2" ]; then
  img_name=$img_name-stg
fi

docker push $PREFIX/$img_name:v1
