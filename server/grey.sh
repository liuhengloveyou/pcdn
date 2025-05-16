#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 

scp ./club ubuntu@114.132.226.76:/opt/club_grey/server/club_grey

rm ./club
