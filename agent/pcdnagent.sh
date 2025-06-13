#!/bin/bash

echo `date` > /tmp/pcdnagent.sh.log

# crontab every on minuter
# crontab -e 
# * * * * * /opt/pcdnagent/pcdnagent.sh

WORK_DIR="/opt/pcdnagent/"
PROCESS_NAME="pcdnagent"


cd "$WORK_DIR" || exit 1

# bash ./sn.sh &>/dev/null
sn=`cat /proc/cmdline | grep -o 'snum=[^ ]*' |cut -d= -f2`

# 使用pgrep检查进程是否存在
if ! pgrep -x "$PROCESS_NAME" > /dev/null; then
    echo "$(date): Process $PROCESS_NAME is not running. Starting it..." > ./nohup.log
    setcap 'cap_net_admin=+ep' pcdnagent
    nohup ./$PROCESS_NAME -sn ${sn} &
fi
