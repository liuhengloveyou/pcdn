#!/bin/bash


# 定义进程名称
PROCESS_NAME="agent"

# 定义检查间隔时间（秒）
CHECK_INTERVAL=30

# 无限循环，检查进程是否在运行
while true; do
    # 使用pgrep检查进程是否存在
    if ! pgrep -x "$PROCESS_NAME" > /dev/null; then
        echo "$(date): Process $PROCESS_NAME is not running. Starting it..."
        # 启动进程
        nohup ./$PROCESS_NAME &
    fi
    # 等待一段时间再次检查
    sleep $CHECK_INTERVAL
done
