#!/bin/bash

# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"
make

# /d/upx.exe ./server

scp ./pcdn-server root@101.37.182.58:/opt/pcdn/server/pcdn-server
# scp ./server root@67.223.117.146:/opt/pcdn-server/pcdn-server/server
# rm ./server

