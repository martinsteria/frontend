#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
LIB_DIR="library"
LIB_GIT_URL="https://github.com/martinsteria/library"
USR_DIR="users"

export GOPATH=$DIR
mkdir ${DIR}/users
if [ ! -d "${DIR}/library"]; then
    git clone $LIB_GIT_URL ${DIR}/library
fi

sudo -E go run main.go 80
