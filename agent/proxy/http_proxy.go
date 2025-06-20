package proxy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"pcdnagent/common"

	"github.com/liuhengloveyou/pcdn/protos"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 处理HTTP代理请求消息
func ProcessHttpProxyReqMsg(conn net.Conn, msgByte []byte) error {
	var req protos.HttpProxyRequest
	if err := proto.Unmarshal(msgByte, &req); err != nil {
		common.Logger.Sugar().Errorf("processHttpProxyReqMsg unmarshal error: %v, data: %s", err, string(msgByte))
		return sendHttpProxyErrorResponse(conn, &req, fmt.Sprintf("解析请求失败: %v", err))
	}

	common.Logger.Info("收到HTTP代理请求",
		zap.String("session_id", req.SessionId),
		zap.String("device_sn", req.DeviceSn),
		zap.String("method", req.Method),
		zap.String("url", req.Url),
		zap.String("proxy_id", req.ProxyId),
	)

	// 根据代理ID处理不同类型的代理请求
	switch req.ProxyId {
	case "router-admin":
		return handleRouterAdminProxy(conn, &req)
	default:
		return sendHttpProxyErrorResponse(conn, &req, fmt.Sprintf("不支持的代理类型: %s", req.ProxyId))
	}
}

// 处理路由器管理代理请求
func handleRouterAdminProxy(conn net.Conn, req *protos.HttpProxyRequest) error {
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest(req.Method, req.Url, bytes.NewReader(req.Body))
	if err != nil {
		common.Logger.Error("创建HTTP请求失败", zap.Error(err))
		return sendHttpProxyErrorResponse(conn, req, fmt.Sprintf("创建HTTP请求失败: %v", err))
	}

	// 设置请求头
	for name, value := range req.Headers {
		httpReq.Header.Set(name, value)
	}

	// 发送HTTP请求
	httpResp, err := client.Do(httpReq)
	if err != nil {
		common.Logger.Error("发送HTTP请求失败", zap.Error(err))
		return sendHttpProxyErrorResponse(conn, req, fmt.Sprintf("发送HTTP请求失败: %v", err))
	}
	defer httpResp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		common.Logger.Error("读取HTTP响应失败", zap.Error(err))
		return sendHttpProxyErrorResponse(conn, req, fmt.Sprintf("读取HTTP响应失败: %v", err))
	}

	// 创建代理响应
	resp := &protos.HttpProxyResponse{
		SessionId:  req.SessionId,
		StatusCode: int32(httpResp.StatusCode),
		Headers:    make(map[string]string),
		Body:       respBody,
	}

	common.Logger.Debug(">>>>>>>>>", zap.Any("body", respBody))
	// 复制响应头
	for name, values := range httpResp.Header {
		if len(values) > 0 {
			resp.Headers[name] = values[0]
		}
	}

	// 发送响应
	return sendHttpProxyResponse(conn, resp)
}

// 发送HTTP代理响应
func sendHttpProxyResponse(conn net.Conn, resp *protos.HttpProxyResponse) error {
	// 序列化为二进制数据
	data, err := proto.Marshal(resp)
	if err != nil {
		common.Logger.Error("HTTP代理响应序列化失败", zap.Error(err))
		return err
	}

	// 构建消息头
	buf := new(bytes.Buffer)
	buf.Write([]byte("\r\n"))

	// 写入消息类型
	msgType := uint32(protos.MsgType_MSG_TYPE_HTTP_PROXY_RESP)
	binary.Write(buf, binary.LittleEndian, msgType)

	// 写入消息长度
	msgLen := uint32(len(data))
	binary.Write(buf, binary.LittleEndian, msgLen)

	// 写入消息体
	buf.Write(data)

	// 发送消息
	if conn == nil {
		return fmt.Errorf("连接未建立，无法发送HTTP代理响应")
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		common.Logger.Error("发送HTTP代理响应失败", zap.Error(err))
		return err
	}

	common.Logger.Info("发送HTTP代理响应成功",
		zap.String("session_id", resp.SessionId),
		zap.Int32("status_code", resp.StatusCode),
		zap.Int("body_size", len(resp.Body)),
	)

	return nil
}

// 发送HTTP代理错误响应
func sendHttpProxyErrorResponse(conn net.Conn, req *protos.HttpProxyRequest, errMsg string) error {
	resp := &protos.HttpProxyResponse{
		SessionId:  req.SessionId,
		StatusCode: 500,
		Headers:    make(map[string]string),
		Error:      errMsg,
	}

	// 设置内容类型为纯文本
	resp.Headers["Content-Type"] = "text/plain; charset=utf-8"

	// 设置错误消息作为响应体
	resp.Body = []byte(errMsg)

	return sendHttpProxyResponse(conn, resp)
}
