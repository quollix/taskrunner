//go:build linux
// +build linux

package platform

import (
	"fmt"
	"os/exec"
	"syscall"
)

func SetProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}

func KillProcessGroup(processID int) error {
	processGroupID, err := syscall.Getpgid(processID)
	if err != nil {
		return fmt.Errorf("Failed to get process group ID of process ID '%v' because of error: %v", processID, err)
	}
	if err := syscall.Kill(-processGroupID, syscall.SIGKILL); err != nil {
		return fmt.Errorf("Failed to kill process group ID '%v' because of error: %v", processID, err)
	}
	return nil
}
