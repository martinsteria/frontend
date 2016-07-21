#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
LIB_DIR="library"
LIB_GIT_URL="https://github.com/martinsteria/library"
USR_DIR="users"


rm -rf ${DIR}/${USR_DIR}
if [ ! -d "${DIR}/${USR_DIR}" ]; then
    mkdir ${DIR}/${USR_DIR}
fi

rm -rf ${DIR}/${LIB_DIR}
if [ ! -d "${DIR}/${LIB_DIR}" ]; then
    git clone $LIB_GIT_URL ${DIR}/${LIB_DIR}
fi
