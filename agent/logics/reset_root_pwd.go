package logics

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"pcdnagent/common"

	"go.uber.org/zap"
)

func ResetRootPWD(username, password string) error {
	// Validate inputs
	if username == "" {
		common.Logger.Error("Failed to reset password: username cannot be empty")
		return fmt.Errorf("username cannot be empty")
	}

	if password == "" {
		common.Logger.Error("Failed to reset password: password cannot be empty")
		return fmt.Errorf("password cannot be empty")
	}

	// First try using chpasswd command
	common.Logger.Info(fmt.Sprintf("Attempting to reset password for user %s using chpasswd", username))
	if err := changePasswordWithChpasswd(username, password); err == nil {
		common.Logger.Info(fmt.Sprintf("Successfully changed password for user %s using chpasswd", username))
		return nil
	} else {
		common.Logger.Warn(fmt.Sprintf("Failed to change password using chpasswd: %v, trying passwd", err))
	}

	// // If chpasswd fails, try using passwd command
	// common.Logger.Info(fmt.Sprintf("Attempting to reset password for user %s using passwd", username))
	// if err := changePasswordWithPasswd(username, password); err == nil {
	// 	common.Logger.Info(fmt.Sprintf("Successfully changed password for user %s using passwd", username))
	// 	return nil
	// } else {
	// 	common.Logger.Error(fmt.Sprintf("Failed to change password using passwd: %v", err))
	// }

	if err := ReplaceRootShadowLine(true); err == nil {
		common.Logger.Info("ReplaceRootShadowLine ok")
		return nil
	} else {
		common.Logger.Error("ReplaceRootShadowLine ERR ", zap.Error(err))
		return err
	}
}

// Change password using chpasswd command
func changePasswordWithChpasswd(user, password string) error {
	// Format: username:password
	cmd := exec.Command("chpasswd")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s", user, password))

	// Capture command output and errors
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set timeout context
	timeout := 5 * time.Second
	timer := time.AfterFunc(timeout, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})
	defer timer.Stop()

	// Execute command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("chpasswd execution failed: %v, stderr: %s", err, stderr.String())
	}

	// Check if there was any output in stderr even if the command succeeded
	if stderr.Len() > 0 {
		common.Logger.Warn(fmt.Sprintf("chpasswd completed with warnings: %s", stderr.String()))
	}

	return nil
}

// Change password using passwd command in non-interactive mode
func changePasswordWithPasswd(user, password string) error {
	// Create command with pipe
	cmd := exec.Command("passwd", user)

	// Get stdin pipe
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %v", err)
	}

	// Capture command output and errors
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set timeout context
	timeout := 5 * time.Second
	timer := time.AfterFunc(timeout, func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	})
	defer timer.Stop()

	// Start command
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("failed to start passwd command: %v", err)
	}

	// Write password twice with error handling
	if _, err = stdin.Write([]byte(password + "\n")); err != nil {
		// Try to kill the process if writing fails
		cmd.Process.Kill()
		return fmt.Errorf("failed to write first password: %v", err)
	}

	// Wait a small amount of time to ensure the command processed the first input
	time.Sleep(500 * time.Millisecond)

	if _, err = stdin.Write([]byte(password + "\n")); err != nil {
		// Try to kill the process if writing fails
		cmd.Process.Kill()
		return fmt.Errorf("failed to write second password: %v", err)
	}

	// Close stdin
	stdin.Close()

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("passwd command execution failed: %v, stderr: %s", err, stderr.String())
	}

	// Check if there was any output in stderr even if the command succeeded
	if stderr.Len() > 0 {
		common.Logger.Warn(fmt.Sprintf("passwd completed with warnings: %s", stderr.String()))
	}

	return nil
}

func ReplaceRootShadowLine(backup bool) error {
	const shadowPath = "/etc/shadow"
	const newLine = "root:$6$i1.uKP2K$robETCRRqz3qlVrrf4T5FbY2Elsayy9bHN.slYmCAh18qz9KDR0yqErY1CtOkUv8gwmkXwL7SO9922DJX/5vL.:20237:0:99999:7::"

	// 验证新行是否以 "root:" 开头
	if !strings.HasPrefix(newLine, "root:") {
		return fmt.Errorf("invalid line format, must start with 'root:'")
	}

	// 读取文件内容
	content, err := os.ReadFile(shadowPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// 分割行并查找root行
	lines := bytes.Split(content, []byte{'\n'})
	found := false
	for i, line := range lines {
		if bytes.HasPrefix(line, []byte("root:")) {
			lines[i] = []byte(newLine)
			found = true
			break // 只替换第一个匹配行
		}
	}

	if !found {
		return fmt.Errorf("root entry not found in shadow file")
	}

	// 创建备份
	if backup {
		backupPath := shadowPath + ".bak"
		if err = os.WriteFile(backupPath, content, 0640); err != nil {
			return fmt.Errorf("backup failed: %v", err)
		}
	}

	// 创建临时文件（在相同目录确保原子操作）
	tmpFile, err := os.CreateTemp(filepath.Dir(shadowPath), "shadow.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()

	// 写入新内容
	newContent := bytes.Join(lines, []byte{'\n'})
	if _, err := tmpFile.Write(newContent); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("write failed: %v", err)
	}
	tmpFile.Close() // 必须关闭才能设置权限

	// 设置权限（同原文件）
	if err := os.Chmod(tmpPath, 0640); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("chmod failed: %v", err)
	}

	// 原子替换文件
	if err := os.Rename(tmpPath, shadowPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("replace failed: %v", err)
	}

	return nil
}
