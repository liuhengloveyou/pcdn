package logics

import (
	"fmt"
	"net"
	"pcdnagent/common"
	"pcdnagent/proxy"
	"time"

	"github.com/liuhengloveyou/pcdn/protos"
	"go.uber.org/zap"
)

// 常见的路由器IP地址
var commonRouterIPs = []string{
	"192.168.1.1",
	"192.168.0.1",
	"10.0.0.1",
	"10.1.1.1",
	"192.168.2.1",
	"192.168.10.1",
	"192.168.11.1",
	"192.168.100.1",
	"192.168.101.1",
	"192.168.88.1",
}

// 常见的路由器端口
var commonRouterPorts = []int{80, 8080, 443}

// DetectRouterIP 检测路由器IP地址
func DetectRouterIP() (string, int, error) {
	// 首先尝试获取默认网关
	gatewayIP, err := getDefaultGateway()
	if err == nil && gatewayIP != "" {
		common.Logger.Info("Found default gateway", zap.String("ip", gatewayIP))

		// 检查网关是否可访问
		for _, port := range commonRouterPorts {
			if isPortOpen(gatewayIP, port) {
				common.Logger.Info("Router detected", zap.String("ip", gatewayIP), zap.Int("port", port))
				return gatewayIP, port, nil
			}
		}
	}

	// 如果默认网关不可访问，尝试常见的路由器IP
	for _, ip := range commonRouterIPs {
		for _, port := range commonRouterPorts {
			if isPortOpen(ip, port) {
				common.Logger.Info("Router detected", zap.String("ip", ip), zap.Int("port", port))
				return ip, port, nil
			}
		}
	}

	return "", 0, fmt.Errorf("no router found")
}

// getDefaultGateway 获取默认网关
func getDefaultGateway() (string, error) {
	// 这里简化处理，实际应该解析路由表
	// 在Linux系统中，可以通过解析/proc/net/route文件获取

	// 简单实现：假设默认网关是192.168.1.1
	return "192.168.1.1", nil
}

// isPortOpen 检查指定IP和端口是否开放
func isPortOpen(ip string, port int) bool {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// CreateRouterAdminProxy 创建路由器管理代理
func CreateRouterAdminProxy() (string, error) {
	// 检测路由器IP和端口
	routerIP, routerPort, err := DetectRouterIP()
	if err != nil {
		common.Logger.Error("Failed to detect router", zap.Error(err))
		return "", err
	}

	// 创建代理连接
	proxyID := "router-admin"
	err = proxy.CreateProxyConnection(proxyID, routerIP, routerPort)
	if err != nil {
		common.Logger.Error("Failed to create proxy connection", zap.Error(err))
		return "", err
	}

	// 获取代理URL
	proxyURL := proxy.GetProxyURL(proxyID)
	common.Logger.Info("Created router admin proxy",
		zap.String("router_ip", routerIP),
		zap.Int("router_port", routerPort),
		zap.String("proxy_url", proxyURL))

	return proxyURL, nil
}

// HandleRouterAdmin 处理路由器管理任务
func HandleRouterAdmin(task *protos.Task) error {
	common.Logger.Info("Handling router admin task")

	// 创建路由器管理代理
	proxyURL, err := CreateRouterAdminProxy()
	if err != nil {
		common.Logger.Error("Failed to create router admin proxy", zap.Error(err))
		task.ErrMsg = fmt.Sprintf("Failed to create router admin proxy: %v", err)
		return err
	}

	// 设置任务的URL字段
	task.Url = &proxyURL
	common.Logger.Info("Router admin task completed", zap.String("url", proxyURL))

	return nil
}
