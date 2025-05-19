#!/bin/bash

# 获取 agent 版本号
ver=$("./agent" -version | head -n 1 | awk '{print $1}')
echo "Detected version: $ver"

go-selfupdate -o update agent $ver