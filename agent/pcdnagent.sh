#!/bin/bash

# crontab 每分钟执行一次
# crontab -e 
# * * * * * /opt/pcdnagent/pcdnagent.sh

# 工作目录
WORK_DIR="/opt/pcdnagent/"
# 定义进程名称
PROCESS_NAME="pcdnagent"

# 进入工作目录
cd "$WORK_DIR" || exit 1

# 定义检查间隔时间（秒）
CHECK_INTERVAL=30

# bash ./sn.sh &>/dev/null
sn=`cat /proc/cmdline | grep -o 'snum=[^ ]*' |cut -d= -f2`

# 使用pgrep检查进程是否存在
if ! pgrep -x "$PROCESS_NAME" > /dev/null; then
    echo "$(date): Process $PROCESS_NAME is not running. Starting it..." > ./nohup.log
    setcap 'cap_net_admin=+ep' pcdnagent
    nohup ./$PROCESS_NAME -sn ${sn} &
fi
