#!/bin/bash

# 获取 agent 版本号
# ver=$("./agent" -version | head -n 1 | awk '{print $1}')
ver=${1:-"0.0.1"}  # 从第一个命令行参数获取版本号，如果没有提供则默认为"0.0.1"
echo "Detected version: $ver"

make arm

go-selfupdate  -platform linux-arm64 -o update pcdnagent $ver

ssh root@101.37.182.58 "rm -rf /opt/pcdn/update/*"
scp -r update root@101.37.182.58:/opt/pcdn/