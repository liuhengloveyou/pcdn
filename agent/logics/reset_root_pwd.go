package logics

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"pcdnagent/common"
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

	// If chpasswd fails, try using passwd command
	common.Logger.Info(fmt.Sprintf("Attempting to reset password for user %s using passwd", username))
	if err := changePasswordWithPasswd(username, password); err == nil {
		common.Logger.Info(fmt.Sprintf("Successfully changed password for user %s using passwd", username))
		return nil
	} else {
		common.Logger.Error(fmt.Sprintf("Failed to change password using passwd: %v", err))
		return fmt.Errorf("failed to change password for user %s: all methods exhausted", username)
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
	if err := cmd.Start(); err != nil {
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
