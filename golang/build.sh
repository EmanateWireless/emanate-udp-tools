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
   GOOS=${t} GOARCH=amd64 go build -o build/${cmd}_osx cmd/${cmd}/main.go

   echo "Building '${cmd}' executable for Windows 64-bit target";
   GOOS=${t} GOARCH=amd64 go build -o build/${cmd}_win_x64.exe cmd/${cmd}/main.go

   echo "Building '${cmd}' executable for Windows 32-bit target";
   GOOS=${t} GOARCH=386 go build -o build/${cmd}_win_386.exe cmd/${cmd}/main.go

   echo "Building '${cmd}' executable for Linux target";
   GOOS=${t} GOARCH=386 go build -o build/${cmd}_linux cmd/${cmd}/main.go
done

echo "DONE!"
echo ""
