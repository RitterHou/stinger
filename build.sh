#!/usr/bin/env bash

build () {
    windows_name="stinger_$1_windows.exe"
    echo "Building $windows_name"
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${windows_name}
    mv ${windows_name} ../target

    linux_name="stinger_$1_linux"
    echo "Building $linux_name"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${linux_name}
    mv ${linux_name} ../target

    mac_name="stinger_$1_mac"
    echo "Building $mac_name"
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ${mac_name}
    mv ${mac_name} ../target
}

directoryName="target"

rm -rf ${directoryName}

if [ ! -d ${directoryName} ]; then
    echo "Creating directory $directoryName"
    mkdir ${directoryName}
fi

# 移动配置文件
mkdir -p target/local/pac
cp local/pac/pac.js target/local/pac/pac.js

# 编译源代码
cd ./server
build "server"

cd ../local
build "local"
