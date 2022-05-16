#!/bin/env bash

APP_NAME="template-gen"
RELEASE_VER="0.1.1"
BUILD_PATH="build"
BIN_PATH="${BUILD_PATH}/bin"

if [[ ! -d "$BUILD_PATH" ]]; then
    mkdir $BUILD_PATH
    if [[ ! -d "$BIN_PATH" ]]; then
        mkdir $BIN_PATH
    fi
fi

# compile the application for linux and archive it
GOOS=linux GOARCH=amd64 go build -o "${BIN_PATH}"/$APP_NAME app/main.go

tar -czf build/"${APP_NAME}"_"${RELEASE_VER}"_linux.tar.gz "${BIN_PATH}/${APP_NAME}"


# compile the application for windows and archive it
GOOS=windows GOARCH=amd64 go build -o "${BIN_PATH}"/$APP_NAME.exe app/main.go

tar -czf build/"${APP_NAME}"_"${RELEASE_VER}"_windows.tar.gz "${BIN_PATH}/${APP_NAME}".exe
