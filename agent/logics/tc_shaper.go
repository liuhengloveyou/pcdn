package logics

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"pcdnagent/common"
	"strings"
	"time"
)

const (
	limitRate = "4000kbit"
)

// GetDefaultInterface retrieves the default network interface name.
func GetDefaultInterface() (string, error) {
	cmd := exec.Command("ip", "route")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'ip route': %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "default via") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				return fields[4], nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning 'ip route' output: %w", err)
	}

	return "", fmt.Errorf("default interface not found")
}

func runTCCommand(args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/sbin/tc", args...)

	// 捕获 stdout 和 stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 记录执行的命令（带参数）
	common.Logger.Sugar().Infof("Executing TC command: tc %s", strings.Join(args, " "))

	// 执行命令
	err := cmd.Run()

	// 打印标准输出和标准错误（无论是否出错）
	if stdout.Len() > 0 {
		common.Logger.Sugar().Infof("TC stdout: %s", stdout.String())
	}
	if stderr.Len() > 0 {
		common.Logger.Sugar().Infof("TC stderr: %s", stderr.String())
	}

	// 处理错误（忽略可接受的错误）
	if err != nil {
		common.Logger.Sugar().Infof("TC warning (ignored): %v - %s", err, stderr.String())

		return fmt.Errorf("tc command 'tc %s' failed: %v, stdout: %q, stderr: %q",
			strings.Join(args, " "), err, stdout.String(), stderr.String())
	}

	return nil
}

// ClearTCRules removes all tc qdisc rules from the specified interface.
func ClearTCRules(ifaceName string) error {
	common.Logger.Sugar().Infoln("Clearing TC rules on interface %s", ifaceName)
	return runTCCommand("qdisc", "del", "dev", ifaceName, "root")
}

// ApplyTCRules applies HTB-based rate limiting to the specified interface.
func ApplyTCRules(ifaceName string, rate string) error {
	common.Logger.Sugar().Infoln("Applying TC rate limit %s on interface %s", rate, ifaceName)
	// First, try to delete any existing qdisc to avoid errors if it already exists with a different handle.
	_ = ClearTCRules(ifaceName) // Ignore error, as it might not exist

	err := runTCCommand("qdisc", "add", "dev", ifaceName, "root", "handle", "1:", "htb", "default", "10")
	if err != nil {
		return fmt.Errorf("failed to add root htb qdisc: %w", err)
	}

	err = runTCCommand("class", "add", "dev", ifaceName, "parent", "1:", "classid", "1:10", "htb", "rate", rate)
	if err != nil {
		// If class add fails, try to clean up the qdisc we just added.
		_ = ClearTCRules(ifaceName)
		return fmt.Errorf("failed to add htb class: %w", err)
	}
	common.Logger.Sugar().Infoln("Successfully applied TC rate limit %s on interface %s", rate, ifaceName)
	return nil
}
