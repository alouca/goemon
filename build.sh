#!/bin/bash

# Setup the environment
export GOPATH=$GOPATH:$PWD

BIN_DIR=$PWD/bin
ROOT_DIR=$PWD

echo " --> Building GoEmon Server"
cd $ROOT_DIR/src/ 
FILES=`ls *.go`
go build -o $BIN_DIR/goemon $FILES
