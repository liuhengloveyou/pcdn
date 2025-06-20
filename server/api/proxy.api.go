package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"pcdn-server/common"
	"pcdn-server/tcpservice"

	gocommon "github.com/liuhengloveyou/go-common"
	"go.uber.org/zap"
)

func init() {
	// 注册HTTP代理API
	Apis["/proxy/router-admin"] = ApiStruct{
		Handler:    handleRouterAdminProxy,
		Method:     "", // 支持所有HTTP方法
		NeedLogin:  true,
		NeedAccess: true,
	}
}

// 处理路由器管理代理请求
func handleRouterAdminProxy(w http.ResponseWriter, r *http.Request) {
	// 获取设备序列号
	deviceSN := r.URL.Query().Get("sn")
	if deviceSN == "" {
		gocommon.HttpJsonErr(w, http.StatusBadRequest, common.ErrParam)
		return
	}

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		common.Logger.Error("读取请求体失败", zap.Error(err))
		gocommon.HttpJsonErr(w, http.StatusInternalServerError, common.ErrService)
		return
	}

	// 获取目标URL
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		targetURL = "http://192.168.1.1:8080/" // 默认路由器地址
	}

	// 转换请求头为map
	headers := make(map[string]string)
	for name, values := range r.Header {
		// 跳过一些特殊的头部
		if strings.ToLower(name) == "host" ||
			strings.ToLower(name) == "content-length" ||
			strings.ToLower(name) == "connection" ||
			strings.ToLower(name) == "authorization" {
			continue
		}

		if len(values) > 0 {
			headers[name] = values[0]
		}
	}

	// 发送HTTP代理请求到设备
	resp, err := tcpservice.SendHttpProxyRequest(
		deviceSN,
		r.Method,
		targetURL,
		"router-admin", // 代理ID
		headers,
		body,
	)

	if err != nil {
		common.Logger.Error("发送HTTP代理请求失败", zap.Error(err))
		gocommon.HttpJsonErr(w, http.StatusInternalServerError, fmt.Errorf("代理请求失败: %v", err))
		return
	}

	// 检查是否有错误
	if resp.Error != "" {
		common.Logger.Error("HTTP代理请求返回错误", zap.String("error", resp.Error))
		gocommon.HttpJsonErr(w, http.StatusInternalServerError, fmt.Errorf("代理请求错误: %s", resp.Error))
		return
	}

	// 设置响应头
	for name, value := range resp.Headers {
		w.Header().Set(name, value)
	}

	// 设置状态码
	w.WriteHeader(int(resp.StatusCode))

	// 写入响应体
	if _, err := w.Write(resp.Body); err != nil {
		common.Logger.Error("写入响应体失败", zap.Error(err))
	}
}
