# Stinger
_Across the Great Wall, we can reach every corner in the world._

![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)

### Usage

在远程主机上启动服务程序

    ./stinger_server -c stinger_server.yaml

在本地启动代理客户端

    ./stinger_local -c stinger_local.yaml

在浏览器中进行代理设置，支持直接代理和PAC代理，建议使用PAC方式

### Contribute

Install glide

    curl https://glide.sh/get | sh

Install statik

    go get -d github.com/rakyll/statik
    go install github.com/rakyll/statik
