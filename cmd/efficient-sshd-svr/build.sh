#!/bin/bash -e

export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++
export CGO_ENABLED=1
export GOOS=windows

go build

mv efficient-sshd-svr.exe ~/Downloads/efficient-sshd-svr.exe

