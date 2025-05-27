#!/bin/bash

sn=`cat /proc/cmdline | grep -o 'snum=[^ ]*' |cut -d= -f2`
if [ -z "$sn" ]; then
    echo "-" >> /etc/pcdnsn
    exit 1
else
    deviceID=$sn
    echo $deviceID > /etc/pcdnsn
    current_time=$(date +%Y-%m-%d\ %H:%M:%S.%3N)
    echo $current_time >> /var/log/start.log
    #cat /etc/pcdnsn | qrencode -m 2 -t ansi > /etc/issue
fi