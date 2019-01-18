# Stinger
_Across the Great Wall, we can reach every corner in the world._

[![Build Status](https://travis-ci.org/RitterHou/stinger.svg?branch=master)](https://travis-ci.org/RitterHou/stinger)
![Golang](https://img.shields.io/badge/golang-1.10.3-blue.svg)
![GitHub](https://img.shields.io/badge/license-Apache%202.0-green.svg)

### Usage

在[下载界面](https://github.com/RitterHou/stinger/releases)下载相应平台客户端和服务端的可执行程序以及它们的配置文件。

在远程主机上启动服务程序

    ./stinger_server -c stinger_server.yaml

在本地启动代理客户端

    ./stinger_local -c stinger_local.yaml

在浏览器中进行代理设置，支持直接代理和PAC代理，建议使用PAC方式

### Introduce

目录介绍

| 文件目录 | 文件目录的作用 |
| --- | --- |
| core/codec | 对local和server之前传输的数据进行加密 |
| core/common | 一些通用的工具 |
| core/mylog | 日志相关的配置 |
| core/network | 对go语言的连接做了一些封装 |
| local/assets | 存储了静态文件，编译前会使用statik进行打包 |
| local/conf | 管理local的配置信息 |
| local/http | 管理http服务器，提供相关的http服务 |
| local/resource | 封装对statik生成的文本的调用 |
| local/socks | 与socks协议相关的数据传输与处理 |
| local/statik | 保存statik生成的文本数据 |
| local/main | local的main方法 |
| server/conf | 管理server的配置信息 |
| server/main | server的main方法 |
| stinger_local.yaml | 本地进程的配置文件 |
| stinger_server.yaml | 远程主机的配置文件 |

工作流程



### Contribute

Install glide

    curl https://glide.sh/get | sh

Install statik

    go get -d github.com/rakyll/statik
    go install github.com/rakyll/statik

Install dependencies

    glide install

Build

    ./build.sh
