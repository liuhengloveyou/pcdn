#!/bin/bash

# 放行目标IP
TARGET_IP="101.37.182.58"
# 限速值（10Mbps）
RATE="10mbit"

# 获取所有物理网卡（排除lo回环）
INTERFACES=$(ip link show | awk -F': ' '/^[0-9]+: e/ {print $2}' | grep -v lo)

for IF in $INTERFACES; do
    # 清除现有规则
    tc qdisc del dev $IF root 2>/dev/null

    # 创建HTB队列（根节点）
    tc qdisc add dev $IF root handle 1: htb default 20

    # 设置根分类（总带宽=限速+放行带宽）
    tc class add dev $IF parent 1: classid 1:1 htb rate ${RATE} ceil ${RATE}

    # 创建子分类1：放行目标IP（不限速）
    tc class add dev $IF parent 1:1 classid 1:10 htb rate 100mbit ceil 100mbit  # 高带宽保证不限速

    # 创建子分类2：限速其他流量（10Mbps）
    tc class add dev $IF parent 1:1 classid 1:20 htb rate ${RATE} ceil ${RATE}

    # 添加过滤器：目标IP放行
    tc filter add dev $IF protocol ip parent 1:0 prio 1 u32 \
        match ip dst $TARGET_IP flowid 1:10

    # 默认过滤器：其他流量限速
    tc filter add dev $IF protocol ip parent 1:0 prio 2 u32 \
        match ip dst 0.0.0.0/0 flowid 1:20  # 匹配所有其他IP

    echo "网卡 $IF 已设置：限速10Mbps，放行 $TARGET_IP"
done

# ### 验证命令
# ```bash
# # 查看网卡的队列规则
# tc -s qdisc show dev eth0

# # 查看分类状态
# tc -s class show dev eth0

# # 查看过滤器规则
# tc -s filter show dev eth0
# ```

# ### 撤销限速
# ```bash
# # 清除所有网卡的规则
# for IF in $(ip link | awk -F': ' '/^[0-9]+: e/ {print $2}' | grep -v lo); do
#     tc qdisc del dev $IF root
# done
# ```