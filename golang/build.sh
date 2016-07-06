#!/bin/sh

# global variables
BUILD_DIR=build
CMDS="emanate_udp_receiver"
TARGETS="darwin linux windows"

# create the build directory
mkdir -p $BUILD_DIR

# iterate through each command executable to build
echo ""
for cmd in $CMDS; do
   # iterate through each target
   for t in $TARGETS; do
      echo "Building '${cmd}' executable for '${t}' target";
      GOOS=${t} GOARCH=amd64 go build -o build/${cmd}_${t}-x64 cmd/${cmd}/main.go
   done
done

echo "DONE!"
echo ""
