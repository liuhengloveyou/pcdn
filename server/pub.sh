#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"

/d/upx.exe ./server

scp ./server root@47.239.247.87:/opt/pcdn-server/pcdn-server/server
# scp ./server root@67.223.117.146:/opt/pcdn-server/pcdn-server/server
rm ./server