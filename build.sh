#!/usr/bin/env bash

# 分别为Windows、Mac和Linux编译local和server应用
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

# 开始执行
start=`date +%s`
# 删除原有的目标文件
directoryName="target"
rm -rf ${directoryName}
if [ ! -d ${directoryName} ]; then
    echo "Creating directory \"$directoryName\""
    mkdir ${directoryName}
fi

# 移动配置文件
echo "Copying configuration files"
cp ./stinger_local.yaml ./target
cp ./stinger_server.yaml ./target

# 编译源代码
cd ./server
build "server"

cd ../local
echo "Packaging static files to statik/statik.go"
go generate
build "local"

end=`date +%s`
cost=$[$end-$start]
echo "Build success, cost ${cost} second(s)"
