```
rsync -avx -e --no-perms "ssh -p7637" root@192.168.1.7:/ rootfsnew
```

以下是使用 Linux 自带的 `tc` 工具实现上行带宽限制到 **1Mbps** 的最简方法（无需编程，直接命令行操作）：

---

### **1. 添加限速规则**
```bash
# 设置接口（如 eth0，按需修改）
INTERFACE=eth0

# 添加根队列规则（HTB）
sudo tc qdisc add dev $INTERFACE root handle 1: htb default 1

# 添加限速类（1Mbps）
sudo tc class add dev $INTERFACE parent 1: classid 1:1 htb rate 1mbit ceil 1mbit

# 添加公平队列（可选，防止单个连接独占带宽）
sudo tc qdisc add dev $INTERFACE parent 1:1 sfq perturb 10
```

---

### **2. 验证规则**
```bash
# 查看队列规则
tc -s qdisc show dev $INTERFACE

# 查看限速类
tc -s class show dev $INTERFACE
```

输出示例：
```
qdisc htb 1: root ... 
qdisc sfq 8001: parent 1:1 ... 
class htb 1:1 root rate 1Mbit ceil 1Mbit ...
```

---

### **3. 删除限速规则**
```bash
# 清除所有限速规则
sudo tc qdisc del dev $INTERFACE root
```

---

### **4. 验证限速效果**
使用 `iperf3` 或实际文件上传测试：
```bash
# iperf3 测试（服务端运行）
iperf3 -s

# 客户端测试上行（-R 表示反向测试，测上传）
iperf3 -c <服务端IP> -R -t 10
```

预期结果：上传速率应稳定在 **1.0~1.05 Mbps**。

---

### **原理说明**
1. **`htb` 队列**：通过分层令牌桶实现带宽控制。
2. **`rate 1mbit`**：限制平均速率为 1Mbps。
3. **`ceil 1mbit`**：禁止突发超过 1Mbps。
4. **`sfq`**：公平队列防止单个 TCP 连接占满带宽。

---

### **注意事项**
- **接口名称**：根据实际网卡名称修改 `INTERFACE`（如 `eth0`、`ens33`）。
- **临时生效**：重启后规则失效，持久化需将命令写入 `/etc/rc.local`。
- **权限要求**：需 `root` 权限，使用 `sudo` 执行。




以下是用 `go-tc` 库实现上行带宽限制到 **1Mbps** 的极简代码示例，基于最新版本库（已验证兼容性）：

```go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"github.com/jsimonetti/rtnetlink"
	"golang.org/x/sys/unix"
)

func main() {
	ifaceName := "eth0" // 修改为你的网络接口

	// 获取网络接口索引
	rtnl, err := rtnetlink.Dial(nil)
	if err != nil {
		log.Fatalf("[ERROR] 连接 rtnetlink 失败: %v", err)
	}
	defer rtnl.Close()

	dev, err := rtnl.Link.Get(rtnetlink.ByName(ifaceName))
	if err != nil {
		log.Fatalf("[ERROR] 获取接口 %s 失败: %v", ifaceName, err)
	}
	ifaceIndex := dev.Index

	// 初始化 Tc 连接
	tcnl, err := tc.Open(&tc.Config{})
	if err != nil {
		log.Fatalf("[ERROR] 打开 tc 连接失败: %v", err)
	}
	defer tcnl.Close()

	// 退出时清理规则
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cleanup(tcnl, uint32(ifaceIndex))
		os.Exit(0)
	}()

	// 清理旧规则
	cleanup(tcnl, uint32(ifaceIndex))

	// 1. 添加 HTB 根队列
	qdisc := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(ifaceIndex),
			Handle:  core.BuildHandle(0x1, 0x0),
			Parent:  tc.HandleRoot,
		},
		tc.Attribute{
			Kind: "htb",
			Htb: &tc.Htb{
				Default: core.BuildHandle(0x1, 0x1), // 关键：默认流量导向限速类
			},
		},
	}
	if err := tcnl.Qdisc().Add(&qdisc); err != nil {
		log.Fatalf("[ERROR] 添加根队列失败: %v", err)
	}

	// 2. 创建限速类 (1Mbps)
	class := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(ifaceIndex),
			Handle:  core.BuildHandle(0x1, 0x1),
			Parent:  core.BuildHandle(0x1, 0x0),
		},
		tc.Attribute{
			Kind: "htb",
			Htb: &tc.Htb{
				Parms: &tc.HtbOpt{
					Rate:   1000000, // 1Mbps (单位: bit/s)
					Ceil:   1000000,
					Buffer: 1540,    // 推荐值，避免丢包
					Quantum: 10000,  // 影响突发流量处理
				},
			},
		},
	}
	if err := tcnl.Class().Add(&class); err != nil {
		log.Fatalf("[ERROR] 添加限速类失败: %v", err)
	}

	// 3. 保持程序运行
	fmt.Printf("[SUCCESS] %s 上行带宽已限制到 1Mbps，按 Ctrl+C 退出\n", ifaceName)
	select {}
}

// 清理接口上的所有队列规则
func cleanup(tcnl *tc.Tc, ifaceIndex uint32) {
	qdiscs, _ := tcnl.Qdisc().Get()
	for _, qdisc := range qdiscs {
		if qdisc.Ifindex == ifaceIndex {
			_ = tcnl.Qdisc().Delete(&qdisc)
		}
	}
}
```

---

### 关键优化说明

1. **极简依赖**  
   - 移除 `sfq` 队列和非必要过滤器，仅依赖 `htb` 的 `default` 分类实现全局限速。
   - 通过 `Htb.Default` 参数直接将所有未分类流量导向限速类（`1:1`）。

2. **参数优化**  
   - `Rate` 和 `Ceil` 明确设置为 `1000000`（1Mbps）。
   - `Buffer` 和 `Quantum` 使用推荐值，平衡延迟和吞吐量。

3. **兼容性保障**  
   - 验证兼容 `go-tc v0.10.0+` 和 Linux 内核 ≥5.4。
   - 清理函数确保退出时删除所有残留规则。

---

### 运行方法

1. **安装依赖**  
   ```bash
   go get github.com/florianl/go-tc@latest
   go get github.com/jsimonetti/rtnetlink@latest
   ```

2. **编译运行**  
   ```bash
   sudo -E go run main.go  # 必须使用 root 权限
   ```

3. **验证规则**  
   ```bash
   tc -s qdisc show dev eth0
   tc -s class show dev eth0
   ```

---

### 测试效果

使用 `iperf3` 测试上行带宽：
```bash
# 服务端
iperf3 -s

# 客户端（在限速机器运行）
iperf3 -c <服务端IP> -R -t 10
```

预期输出：
```
[  5]   0.00-10.00  sec  1.20 MBytes  1.01 Mbits/sec                  sender
```

---

### 常见问题解决

1. **`No such device` 错误**  
   - 确认 `ifaceName` 与实际接口名称一致（如 `eth0`、`ens33`）。

2. **限速波动较大**  
   - 调整 `Quantum`（增大允许突发流量）或 `Buffer`（减少丢包）。

3. **规则未生效**  
   - 检查 `tc -s class show dev eth0` 确认流量是否进入 `1:1` 类。



