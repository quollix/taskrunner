//go:build windows
// +build windows

package platform

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func SetProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func KillProcessGroup(processID int) error {
	fmt.Printf("Killing process group for PID %d\n", processID)
	c := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(processID))
	if out, err := c.CombinedOutput(); err != nil {
		if strings.Contains(err.Error(), "exit status 128") {
			return nil
		}
		return fmt.Errorf("failed to kill process tree for PID '%v': %v: %s", processID, err, out)
	}
	return nil
}
