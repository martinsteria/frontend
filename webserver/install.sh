#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export GOPATH=$DIR
mkdir ${DIR}/users
git clone https://github.com/martinsteria/library ${DIR}/library
