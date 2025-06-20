#!/bin/bash

# 获取 agent 版本号
ver=${1:-"0.0.1"}  # 从第一个命令行参数获取版本号，如果没有提供则默认为"0.0.1"
echo "Detected version: $ver"

make arm

./upx --best --lzma pcdnagent

# 计算压缩后文件的MD5
md5sum=$(md5sum pcdnagent | awk '{print $1}')

# 生成更新包
cat <<EOF > update.json
{
  "version": "$ver",
  "md5": "$md5sum",
  "url": "http://update.intelliflyt.com/upgrade/pcdnagent"
}
EOF

# go-selfupdate  -platform linux-arm64 -o update pcdnagent $ver
# ssh root@101.37.182.58 "rm -rf /opt/pcdn/update/*"
# scp -r update root@101.37.182.58:/opt/pcdn/

ssh root@101.37.182.58 "rm -rf /opt/pcdn/upgrade/*"
scp update.json pcdnagent root@101.37.182.58:/opt/pcdn/upgrade/

ssh root@101.37.182.58 "chown -R www:www /opt/pcdn/"
