//go:build windows
// +build windows

package platform

import (
	"fmt"
	"os/exec"
	"strconv"
	"syscall"
)

func SetProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func KillProcessGroup(processID int) error {
	c := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(processID))
	if out, err := c.CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to kill process tree for PID '%v': %v: %s", processID, err, out)
	}
	return nil
}
