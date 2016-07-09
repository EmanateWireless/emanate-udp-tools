#!/bin/sh

# global variables
BUILD_DIR=build
CMDS="emanate_udp_sender emanate_udp_receiver"

# create the build directory
mkdir -p $BUILD_DIR

# iterate through each command executable to build
echo ""
for cmd in $CMDS; do
   echo "Building '${cmd}' executable for OSX target";
   GOOS=darwin GOARCH=386 go build -o build/${cmd}_osx cmd/${cmd}/main.go

   echo "Building '${cmd}' executable for Windows target";
   GOOS=windows GOARCH=386 go build -o build/${cmd}_win.exe cmd/${cmd}/main.go

   echo "Building '${cmd}' executable for Linux x86 target";
   GOOS=linux GOARCH=386 go build -o build/${cmd}_linux_x86 cmd/${cmd}/main.go
done

echo "DONE!"
echo ""
