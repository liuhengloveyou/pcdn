#!/bin/bash

# 定义网卡名称,指定或自动获取
INTERFACE=$(ip route | awk '/default/ {print $5}')
#INTERFACE=eth0

# 定义速率限制,4Mbit
RATE="4000kbps"

# 定义解除限速时间段
START_TIME="17:15"
END_TIME="23:30"

# 获取当前时间
CURRENT_TIME=$(date +%s)

second_start=`date +%s -d "$START_TIME"`
second_end=`date +%s -d "$END_TIME"`

# tc qdisc show dev eth0
# tc qdisc del dev eth0 root
# tc qdisc add dev eth0 root handle 1: htb default 10
# tc class add dev eth0 parent 1: classid 1:10 htb rate 1mbps

if [[ "$CURRENT_TIME" -ge "$second_start" && "$CURRENT_TIME" -le "$second_end" ]]; then
    # 清除现有的tc规则
    sudo tc qdisc del dev $INTERFACE root 2>/dev/null
    # 添加新的tc规则
#    sudo tc qdisc add dev $INTERFACE root handle 1: htb default 10
#    sudo tc class add dev $INTERFACE parent 1: classid 1:10 htb rate $RATE
else
    # 在指定时间范围外，重新设置规则
    sudo tc qdisc del dev $INTERFACE root 2>/dev/null
    # 添加新的tc规则
    sudo tc qdisc add dev $INTERFACE root handle 1: htb default 10
    sudo tc class add dev $INTERFACE parent 1: classid 1:10 htb rate $RATE
fi
