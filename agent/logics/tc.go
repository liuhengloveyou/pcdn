package logics

import (
	"fmt"
	"log"
	"os/exec"
	"pcdnagent/common"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var (
	tcMu sync.Mutex
)

func ApplyLimitUploadBandwidthRules(faceName, rate, targetIP string) error {
	if rate == "" {
		return fmt.Errorf("rate or targetIP is nil")
	}

	tcMu.Lock()
	defer tcMu.Unlock()

	if faceName == "" {
		interfaces, err := getPhysicalInterfaces()
		if err != nil {
			return fmt.Errorf("获取网卡失败: %w", err)
		}

		errMsg := ""
		for _, iface := range interfaces {
			if err := clearInterfaceRules(iface); err != nil {
				log.Printf("清除网卡 %s 规则失败: %v", iface, err)
			}

			if err := setupInterface(iface, rate, targetIP); err != nil {
				errMsg += fmt.Sprintf("设置网卡 %s 失败: %w\n", iface, err)
			}
		}
		if errMsg != "" {
			return fmt.Errorf("%s", errMsg)
		}
	} else {
		if err := clearInterfaceRules(faceName); err != nil {
			log.Printf("清除网卡 %s 规则失败: %v", faceName, err)
		}

		if err := setupInterface(faceName, rate, targetIP); err != nil {
			return fmt.Errorf("设置网卡 %s 失败: %w", faceName, err)
		}
	}
	return nil
}

func getPhysicalInterfaces() ([]string, error) {
	cmd := exec.Command("sh", "-c", "ip link show | awk -F': ' '/^[0-9]+: e/ {print $2}' | grep -v lo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	interfaces := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(interfaces) == 0 {
		return nil, fmt.Errorf("未找到物理网卡")
	}

	return interfaces, nil
}

func clearInterfaceRules(iface string) error {
	cmd := exec.Command("/sbin/tc", "qdisc", "del", "dev", iface, "root")
	cmd.Run() // 忽略错误，可能没有规则
	return nil
}

func setupInterface(iface, rate, targetIP string) error {
	// 创建根HTB队列
	cmds := []*exec.Cmd{
		exec.Command("/sbin/tc", "qdisc", "add", "dev", iface, "root", "handle", "1:", "htb", "default", "20"),
		exec.Command("/sbin/tc", "class", "add", "dev", iface, "parent", "1:", "classid", "1:1", "htb", "rate", rate, "ceil", rate),
		exec.Command("/sbin/tc", "class", "add", "dev", iface, "parent", "1:1", "classid", "1:10", "htb", "rate", "10000mbit", "ceil", "10000mbit"),
		exec.Command("/sbin/tc", "class", "add", "dev", iface, "parent", "1:1", "classid", "1:20", "htb", "rate", rate, "ceil", rate),
		exec.Command("/sbin/tc", "filter", "add", "dev", iface, "protocol", "ip", "parent", "1:0", "prio", "1", "u32",
			"match", "ip", "dst", targetIP, "flowid", "1:10"),
		exec.Command("/sbin/tc", "filter", "add", "dev", iface, "protocol", "ip", "parent", "1:0", "prio", "2", "u32",
			"match", "ip", "dst", "0.0.0.0/0", "flowid", "1:20"),
	}

	for _, cmd := range cmds {
		output, err := cmd.CombinedOutput()
		if err != nil {
			common.Logger.Error("cmd ERR: ", zap.Any("cmd", cmd), zap.Error(err))
			return fmt.Errorf("命令执行失败: %s\n错误输出: %s", cmd.String(), string(output))
		}
	}

	return nil
}

func ClearAllLimitUploadBandwidthRules() {
	interfaces, err := getPhysicalInterfaces()
	if err != nil {
		log.Printf("获取网卡失败: %v", err)
		return
	}

	for _, iface := range interfaces {
		clearInterfaceRules(iface)
	}
}

// getTCStatus 使用tc命令获取当前的限速规则状态
func GetTCStatus(ifaceName string) (string, string, error) {
	if ifaceName == "" {
		return "", "", fmt.Errorf("ifaceName is empty")
	}

	// 检查第一个网卡的tc规则
	cmd := exec.Command("/sbin/tc", "qdisc", "show", "dev", ifaceName)
	qdiscOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("获取qdisc规则失败: %s; %w", string(qdiscOutput), err)
	}

	// 如果没有找到HTB规则，则认为限速已禁用
	if !strings.Contains(string(qdiscOutput), "htb") {
		return "末设置", "", nil
	}

	// 获取类规则
	cmd = exec.Command("/sbin/tc", "class", "show", "dev", ifaceName)
	classOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("获取class规则失败: %w", err)
	}

	// 获取过滤器规则
	cmd = exec.Command("/sbin/tc", "filter", "show", "dev", ifaceName)
	filterOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("获取filter规则失败: %w", err)
	}

	// 解析速率和目标IP
	otherRateInfo := "未知"
	// 解析不同类别的速率
	for _, line := range strings.Split(string(classOutput), "\n") {
		if strings.Contains(line, "1:10") && strings.Contains(line, "rate") {
			// // 针对特定目标IP的流量类别
			// parts := strings.Fields(line)
			// for i, part := range parts {
			// 	if part == "rate" && i+1 < len(parts) {
			// 		targetRateInfo = parts[i+1]
			// 		break
			// 	}
			// }
		} else if strings.Contains(line, "1:20") && strings.Contains(line, "rate") {
			// 针对其他所有流量的类别
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "rate" && i+1 < len(parts) {
					otherRateInfo = parts[i+1]
					break
				}
			}
		}
	}

	return otherRateInfo, fmt.Sprintf("Qdisc:\n%s\nClass:\n%s\nFilter:\n%s\n", string(qdiscOutput), string(classOutput), string(filterOutput)), nil
}
