#!/bin/bash

# 监控 10 秒内的平均上行速率
duration=10
interface=eth0

# 获取初始流量统计
rx_start=$(cat /sys/class/net/$interface/statistics/tx_bytes)

sleep $duration

# 获取结束流量统计
rx_end=$(cat /sys/class/net/$interface/statistics/tx_bytes)

# 计算平均速率 (Mbps)
rate=$(echo "scale=2; (($rx_end - $rx_start) * 8) / ($duration * 1000000)" | bc)
echo "Average upload rate: $rate Mbps"

# 验证是否接近 1Mbps
if (( $(echo "$rate < 1.05" | bc -l) && $(echo "$rate > 0.95" | bc -l) )); then
    echo "Rate limiting is working!"
else
    echo "Rate limiting failed!"
fi

