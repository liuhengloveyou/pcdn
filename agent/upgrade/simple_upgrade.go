package upgrade

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pcdnagent/common"

	"go.uber.org/zap"
)

// SimpleUpgrader 简单升级器
type SimpleUpgrader struct {
	ServerURL  string // 升级服务器地址
	CurrentVer string // 当前版本
	DeviceSN   string // 设备序列号
	TempFile   string // 临时文件路径
	BackupFile string // 备份文件路径
	Executable string // 可执行文件路径
	DNSServer  string // DNS服务器地址
}

// UpdateInfo 升级信息
type UpdateInfo struct {
	Version string `json:"version"` // 版本号
	MD5     string `json:"md5"`     // MD5校验和
	URL     string `json:"url"`     // 下载地址
}

// NewSimpleUpgrader 创建一个新的简单升级器
func NewSimpleUpgrader(serverURL, currentVer, deviceSN string) *SimpleUpgrader {
	executable, _ := os.Executable()

	// 从环境变量获取DNS服务器地址，如果未设置则使用默认值
	dnsServer := os.Getenv("PCDN_DNS_SERVER")
	if dnsServer == "" {
		dnsServer = "8.8.8.8:53" // 默认使用Google DNS
	}

	return &SimpleUpgrader{
		ServerURL:  serverURL,
		CurrentVer: currentVer,
		DeviceSN:   deviceSN,
		TempFile:   executable + ".new",
		BackupFile: executable + ".bak",
		Executable: executable,
		DNSServer:  dnsServer,
	}
}

// CheckAndUpgrade 检查并升级
func (u *SimpleUpgrader) CheckAndUpgrade() error {
	common.Logger.Info("Checking for updates",
		zap.String("current_version", u.CurrentVer),
		zap.String("server", u.ServerURL))

	// 获取升级信息
	updateInfo, err := u.getUpdateInfo()
	if err != nil {
		return fmt.Errorf("failed to get update info: %v", err)
	}

	// 检查版本
	if !u.needsUpgrade(updateInfo.Version) {
		common.Logger.Info("No update needed",
			zap.String("current", u.CurrentVer),
			zap.String("latest", updateInfo.Version))
		return nil
	}

	common.Logger.Info("Update available",
		zap.String("current", u.CurrentVer),
		zap.String("latest", updateInfo.Version))

	// 下载新版本
	if err = u.downloadNewVersion(updateInfo); err != nil {
		return fmt.Errorf("failed to download new version: %v", err)
	}

	// 验证MD5
	md5sum, err := calculateFileMD5(u.TempFile)
	if err != nil {
		os.Remove(u.TempFile) // 清理临时文件
		return fmt.Errorf("failed to calculate MD5: %v", err)
	}

	if md5sum != updateInfo.MD5 {
		os.Remove(u.TempFile) // 清理临时文件
		return fmt.Errorf("MD5 mismatch: expected %s, got %s", updateInfo.MD5, md5sum)
	}

	common.Logger.Info("Downloaded new version",
		zap.String("version", updateInfo.Version),
		zap.String("md5", md5sum))

	// 备份当前版本
	if err := os.Rename(u.Executable, u.BackupFile); err != nil {
		os.Remove(u.TempFile) // 清理临时文件
		return fmt.Errorf("failed to backup current version: %v", err)
	}

	// 安装新版本
	if err := os.Rename(u.TempFile, u.Executable); err != nil {
		// 恢复备份
		os.Rename(u.BackupFile, u.Executable)
		return fmt.Errorf("failed to install new version: %v", err)
	}

	// 设置执行权限
	if err := os.Chmod(u.Executable, 0755); err != nil {
		return fmt.Errorf("failed to set executable permission: %v", err)
	}

	common.Logger.Info("Update successful", zap.String("version", updateInfo.Version))

	// 重启服务
	return u.restartService()
}

// getUpdateInfo 获取升级信息
func (u *SimpleUpgrader) getUpdateInfo() (*UpdateInfo, error) {
	// 设置HTTP客户端超时
	client := &http.Client{
		Timeout: 30 * time.Second, // 30秒超时
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second, // DNS解析和连接超时
				KeepAlive: 30 * time.Second,
				Resolver: &net.Resolver{
					PreferGo: true,
					Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
						d := net.Dialer{
							Timeout: 5 * time.Second,
						}
						// 使用配置的DNS服务器
						return d.DialContext(ctx, network, u.DNSServer)
					},
				},
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 5 * time.Second,
		},
	}

	// 重试逻辑
	var lastErr error
	for retry := 0; retry < 3; retry++ {
		if retry > 0 {
			common.Logger.Info("Retrying get update info", zap.Int("attempt", retry+1))
			time.Sleep(2 * time.Second) // 重试间隔
		}

		resp, err := client.Get(u.ServerURL + "update.json")
		if err != nil {
			lastErr = fmt.Errorf("request failed: %v", err)
			common.Logger.Warn("Failed to get update info",
				zap.Error(err),
				zap.Int("attempt", retry+1),
				zap.String("url", u.ServerURL+"update.json"))
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("server returned status %d", resp.StatusCode)
			common.Logger.Warn("Server returned non-OK status",
				zap.Int("status", resp.StatusCode),
				zap.Int("attempt", retry+1))
			continue
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %v", err)
			common.Logger.Warn("Failed to read response",
				zap.Error(err),
				zap.Int("attempt", retry+1))
			continue
		}

		var info UpdateInfo
		if err := json.Unmarshal(data, &info); err != nil {
			lastErr = fmt.Errorf("failed to parse JSON: %v", err)
			common.Logger.Warn("Failed to parse JSON",
				zap.Error(err),
				zap.Int("attempt", retry+1))
			continue
		}

		return &info, nil
	}

	return nil, fmt.Errorf("failed to get update info after 3 attempts: %v", lastErr)
}

// needsUpgrade 检查是否需要升级
func (u *SimpleUpgrader) needsUpgrade(newVersion string) bool {
	// 简单版本比较，可以根据需要实现更复杂的版本比较逻辑
	return newVersion > u.CurrentVer
}

// downloadNewVersion 下载新版本
func (u *SimpleUpgrader) downloadNewVersion(info *UpdateInfo) error {
	// 清理已存在的临时文件
	if _, err := os.Stat(u.TempFile); err == nil {
		if err := os.Remove(u.TempFile); err != nil {
			return fmt.Errorf("failed to remove existing temp file: %v", err)
		}
	}

	// 创建临时文件
	file, err := os.Create(u.TempFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 设置HTTP客户端超时
	client := &http.Client{
		Timeout: 5 * time.Minute, // 5分钟超时
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // DNS解析和连接超时
				KeepAlive: 30 * time.Second,
				Resolver: &net.Resolver{
					PreferGo: true,
					Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
						d := net.Dialer{
							Timeout: 5 * time.Second,
						}
						// 使用公共DNS服务器，可以根据需要更改
						return d.DialContext(ctx, network, u.DNSServer)
					},
				},
			}).DialContext,
			TLSHandshakeTimeout:   15 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
			MaxIdleConns:          10,
			IdleConnTimeout:       90 * time.Second,
		},
	}

	// 重试下载
	var lastErr error
	for retry := 0; retry < 3; retry++ {
		if retry > 0 {
			common.Logger.Info("Retrying download", zap.Int("attempt", retry+1))
			time.Sleep(5 * time.Second) // 重试间隔
		}

		// 重置文件指针
		if _, err := file.Seek(0, 0); err != nil {
			return fmt.Errorf("failed to reset file pointer: %v", err)
		}
		if err := file.Truncate(0); err != nil {
			return fmt.Errorf("failed to truncate file: %v", err)
		}

		// 下载文件
		resp, err := client.Get(info.URL)
		if err != nil {
			lastErr = err
			common.Logger.Error("Download failed", zap.Error(err), zap.Int("attempt", retry+1))
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("server returned status %d", resp.StatusCode)
			common.Logger.Error("Download failed", zap.Error(lastErr), zap.Int("attempt", retry+1))
			continue
		}

		// 复制响应体到文件
		if _, err := io.Copy(file, resp.Body); err != nil {
			lastErr = err
			common.Logger.Error("Download failed", zap.Error(err), zap.Int("attempt", retry+1))
			continue
		}

		// 下载成功
		return nil
	}

	// 所有重试都失败了
	os.Remove(u.TempFile) // 清理临时文件
	return fmt.Errorf("failed to download after 3 attempts: %v", lastErr)
}

// restartService 重启服务
// 直接退出， 等脚本拉起
func (u *SimpleUpgrader) restartService() error {
	os.Exit(0)
	return nil
}

// GenerateUpdateJSON 生成升级信息JSON文件
func GenerateUpdateJSON(updateServer *string, version string, executablePath string) error {
	// 计算MD5
	md5sum, err := calculateFileMD5(executablePath)
	if err != nil {
		return fmt.Errorf("failed to calculate MD5: %v", err)
	}

	// 构建URL
	baseURL := *updateServer
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	// 获取可执行文件名
	executableName := filepath.Base(executablePath)

	// 创建更新信息
	updateInfo := UpdateInfo{
		Version: version,
		MD5:     md5sum,
		URL:     baseURL + executableName,
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(updateInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// 写入文件
	if err := os.WriteFile("update.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write update.json: %v", err)
	}

	return nil
}

// calculateFileMD5 计算文件的MD5值
func calculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
