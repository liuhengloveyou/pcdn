#!/bin/bash
# sn --prefix="paopaoarm-" >> /var/log/start.log 2>&1
# file_path="/etc/pcdnuuid"
# cat "$file_path" | qrencode -m 2 -t ansi > /etc/issue
sn=`cat /proc/cmdline | grep -o 'snum=[^ ]*' |cut -d= -f2`
if [ -z "$sn" ]; then
    echo "sn is empty" >> /var/log/start.log
    exit 1
else
    deviceID=$sn
    echo $deviceID > /etc/pcdnuuid
    current_time=$(date +%Y-%m-%d\ %H:%M:%S.%3N)
    cat /etc/pcdnuuid | qrencode -m 2 -t ansi > /etc/issue
    echo $current_time >> /var/log/start.log
fi