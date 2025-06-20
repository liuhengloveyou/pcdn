package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"strings"
	"sync"
	"time"

	"pcdnagent/common"

	"go.uber.org/zap"
)

// ProxyManager 管理代理连接
type ProxyManager struct {
	mu          sync.RWMutex
	connections map[string]*ProxyConnection
	server      *http.Server
	listenPort  int
}

// ProxyConnection 代表一个代理连接
type ProxyConnection struct {
	id         string
	targetHost string
	targetPort int
	createdAt  time.Time
}

var (
	manager *ProxyManager
)

// InitProxyManager 初始化代理管理器
func InitProxyManager(port int) error {
	manager = &ProxyManager{
		connections: make(map[string]*ProxyConnection),
		listenPort:  port,
	}

	// 启动HTTP代理服务器
	mux := http.NewServeMux()
	mux.HandleFunc("/", manager.handleProxy)

	manager.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		common.Logger.Info("Starting proxy server", zap.Int("port", port))
		if err := manager.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.Logger.Error("Proxy server failed", zap.Error(err))
		}
	}()

	return nil
}

// CreateProxyConnection 创建代理连接
func CreateProxyConnection(id, targetHost string, targetPort int) error {
	if manager == nil {
		return fmt.Errorf("proxy manager not initialized")
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	conn := &ProxyConnection{
		id:         id,
		targetHost: targetHost,
		targetPort: targetPort,
		createdAt:  time.Now(),
	}

	manager.connections[id] = conn
	common.Logger.Info("Created proxy connection",
		zap.String("id", id),
		zap.String("target", fmt.Sprintf("%s:%d", targetHost, targetPort)))

	return nil
}

// RemoveProxyConnection 移除代理连接
func RemoveProxyConnection(id string) {
	if manager == nil {
		return
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	delete(manager.connections, id)
	common.Logger.Info("Removed proxy connection", zap.String("id", id))
}

// handleProxy 处理代理请求
func (pm *ProxyManager) handleProxy(w http.ResponseWriter, r *http.Request) {
	// 从请求头中获取代理ID
	proxyID := r.Header.Get("X-Proxy-ID")
	
	// 如果没有代理ID，检查是否是路由器访问请求
	if proxyID == "" {
		if r.URL.Query().Get("router_access") == "true" {
			proxyID = "router-admin"
		} else {
			http.Error(w, "Missing X-Proxy-ID header", http.StatusBadRequest)
			return
		}
	}

	pm.mu.RLock()
	conn, exists := pm.connections[proxyID]
	pm.mu.RUnlock()

	if !exists {
		http.Error(w, "Proxy connection not found", http.StatusNotFound)
		return
	}

	// 构建目标URL
	targetURL := fmt.Sprintf("http://%s:%d%s", conn.targetHost, conn.targetPort, r.URL.Path)
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	common.Logger.Debug("Proxying request",
		zap.String("proxy_id", proxyID),
		zap.String("method", r.Method),
		zap.String("target_url", targetURL))

	// 创建到目标服务器的请求
	targetReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		common.Logger.Error("Failed to create target request", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 复制请求头（除了代理相关的头）
	for name, values := range r.Header {
		if !strings.HasPrefix(strings.ToLower(name), "x-proxy-") {
			for _, value := range values {
				targetReq.Header.Add(name, value)
			}
		}
	}

	// 发送请求到目标服务器
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(targetReq)
	if err != nil {
		common.Logger.Error("Failed to proxy request", zap.Error(err))
		http.Error(w, "Failed to connect to target", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// 设置状态码
	w.WriteHeader(resp.StatusCode)

	// 复制响应体
	io.Copy(w, resp.Body)
}

// HandleTCPProxy 处理TCP代理（用于非HTTP协议）
func HandleTCPProxy(proxyID string, clientConn net.Conn) {
	defer clientConn.Close()

	if manager == nil {
		common.Logger.Error("Proxy manager not initialized")
		return
	}

	manager.mu.RLock()
	conn, exists := manager.connections[proxyID]
	manager.mu.RUnlock()

	if !exists {
		common.Logger.Error("Proxy connection not found", zap.String("proxy_id", proxyID))
		return
	}

	// 连接到目标服务器
	targetAddr := fmt.Sprintf("%s:%d", conn.targetHost, conn.targetPort)
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		common.Logger.Error("Failed to connect to target",
			zap.String("target", targetAddr),
			zap.Error(err))
		return
	}
	defer targetConn.Close()

	common.Logger.Info("TCP proxy connection established",
		zap.String("proxy_id", proxyID),
		zap.String("target", targetAddr))

	// 双向数据转发
	go func() {
		io.Copy(targetConn, clientConn)
		targetConn.Close()
	}()

	io.Copy(clientConn, targetConn)
}

// GetProxyURL 获取代理访问URL
func GetProxyURL(proxyID string) string {
	if manager == nil {
		return ""
	}
	
	// 对于路由器管理代理，使用特殊的参数
	if proxyID == "router-admin" {
		return fmt.Sprintf("http://localhost:%d/?router_access=true", manager.listenPort)
	}
	
	return fmt.Sprintf("http://localhost:%d/?proxy_id=%s", manager.listenPort, proxyID)
}

// ListConnections 列出所有代理连接
func ListConnections() map[string]*ProxyConnection {
	if manager == nil {
		return nil
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	result := make(map[string]*ProxyConnection)
	for k, v := range manager.connections {
		result[k] = v
	}
	return result
}

// StopProxyManager 停止代理管理器
func StopProxyManager() error {
	if manager == nil || manager.server == nil {
		return nil
	}

	return manager.server.Close()
}
