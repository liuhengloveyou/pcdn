#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"

/d/upx.exe ./arbitrage

scp ./arbitrage root@47.239.247.87:/opt/arbitrage/server/arbitrage
# scp ./arbitrage root@67.223.117.146:/opt/arbitrage/server/arbitrage
rm ./arbitrage