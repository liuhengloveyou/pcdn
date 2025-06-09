#!/bin/bash

killall -9 pcdnagent

nohup ./pcdnagent -tcp_server 127.0.0.1:10001 &
