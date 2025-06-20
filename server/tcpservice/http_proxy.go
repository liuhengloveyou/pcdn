package tcpservice

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"pcdn-server/common"

	"github.com/google/uuid"
	"github.com/liuhengloveyou/pcdn/protos"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 存储HTTP代理会话信息
type HttpProxySession struct {
	SessionID  string
	DeviceSN   string
	ProxyID    string
	ResponseCh chan *protos.HttpProxyResponse
	CreatedAt  time.Time
}

var (
	// 会话管理
	httpProxySessions      = make(map[string]*HttpProxySession)
	httpProxySessionsMutex sync.RWMutex
)

// 发送HTTP代理请求到设备
func SendHttpProxyRequest(deviceSN, method, url, proxyID string, headers map[string]string, body []byte) (*protos.HttpProxyResponse, error) {
	deviceSN = strings.ToUpper(deviceSN)

	// 查找设备连接
	httpProxySessionsMutex.RLock()
	tmpAgent, ok := AgentMap[deviceSN]
	httpProxySessionsMutex.RUnlock()

	if !ok || tmpAgent.ClientTcpConn == nil {
		return nil, fmt.Errorf("设备 %s 不在线", deviceSN)
	}

	// 创建会话ID
	sessionID := uuid.New().String()

	// 创建响应通道
	respCh := make(chan *protos.HttpProxyResponse, 1)

	// 保存会话信息
	session := &HttpProxySession{
		SessionID:  sessionID,
		DeviceSN:   deviceSN,
		ProxyID:    proxyID,
		ResponseCh: respCh,
		CreatedAt:  time.Now(),
	}

	httpProxySessionsMutex.Lock()
	httpProxySessions[sessionID] = session
	httpProxySessionsMutex.Unlock()

	// 创建HTTP代理请求消息
	request := &protos.HttpProxyRequest{
		SessionId: sessionID,
		DeviceSn:  deviceSN,
		Method:    method,
		Url:       url,
		Headers:   headers,
		Body:      body,
		ProxyId:   proxyID,
	}

	// 序列化请求
	msgByte, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	// 发送请求
	buff := bytes.NewBuffer([]byte("\r\n"))
	if err := binary.Write(buff, binary.LittleEndian, uint32(protos.MsgType_MSG_TYPE_HTTP_PROXY_REQ)); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.LittleEndian, uint32(len(msgByte))); err != nil {
		return nil, err
	}
	if n, err := buff.Write(msgByte); n != len(msgByte) || err != nil {
		return nil, err
	}

	tmpAgent.MU.Lock()
	if n, err := tmpAgent.ClientTcpConn.Write(buff.Bytes()); n != buff.Len() || err != nil {
		tmpAgent.MU.Unlock()
		return nil, err
	}
	tmpAgent.MU.Unlock()

	// 等待响应，设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	select {
	case resp := <-respCh:
		// 清理会话
		httpProxySessionsMutex.Lock()
		delete(httpProxySessions, sessionID)
		httpProxySessionsMutex.Unlock()
		return resp, nil
	case <-ctx.Done():
		// 超时，清理会话
		httpProxySessionsMutex.Lock()
		delete(httpProxySessions, sessionID)
		httpProxySessionsMutex.Unlock()
		return nil, fmt.Errorf("请求超时")
	}
}

// 处理来自设备的HTTP代理响应
func processHttpProxyRespMsg(conn net.Conn, msgByte []byte) error {
	var response protos.HttpProxyResponse
	if err := proto.Unmarshal(msgByte, &response); err != nil {
		common.Logger.Error("解析HTTP代理响应失败", zap.Error(err))
		return err
	}

	// 查找会话
	httpProxySessionsMutex.RLock()
	session, exists := httpProxySessions[response.SessionId]
	httpProxySessionsMutex.RUnlock()

	if !exists {
		common.Logger.Error("找不到HTTP代理会话", zap.String("session_id", response.SessionId))
		return fmt.Errorf("会话不存在: %s", response.SessionId)
	}

	// 发送响应到通道
	select {
	case session.ResponseCh <- &response:
		// 成功发送
	default:
		// 通道已满或已关闭
		common.Logger.Error("无法发送HTTP代理响应到通道", zap.String("session_id", response.SessionId))
	}

	return nil
}

// 定期清理过期的会话
func startHttpProxySessionCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			cleanupExpiredHttpProxySessions()
		}
	}()
}

func cleanupExpiredHttpProxySessions() {
	now := time.Now()
	expiredSessions := []string{}

	// 查找过期会话
	httpProxySessionsMutex.RLock()
	for id, session := range httpProxySessions {
		if now.Sub(session.CreatedAt) > 30*time.Minute {
			expiredSessions = append(expiredSessions, id)
		}
	}
	httpProxySessionsMutex.RUnlock()

	// 删除过期会话
	if len(expiredSessions) > 0 {
		httpProxySessionsMutex.Lock()
		for _, id := range expiredSessions {
			delete(httpProxySessions, id)
		}
		httpProxySessionsMutex.Unlock()

		common.Logger.Info("已清理过期HTTP代理会话", zap.Int("count", len(expiredSessions)))
	}
}
