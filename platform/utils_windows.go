//go:build windows

package platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func BuildCommand(dir, commandStr string) *exec.Cmd {
	shell := os.Getenv("COMSPEC")
	if shell == "" {
		shell = "cmd.exe"
	}
	base := strings.ToLower(filepath.Base(shell))

	var cmd *exec.Cmd
	if strings.Contains(base, "powershell") || base == "pwsh.exe" || base == "pwsh" {
		cmd = exec.Command(shell, "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", commandStr)
	} else {
		cmd = exec.Command(shell, "/d", "/s", "/c", commandStr)
	}

	cmd.Dir = dir
	return cmd
}

func KillProcesses(processes []string) {
	for _, p := range processes {
		cmd := exec.Command("cmd", "/C", fmt.Sprintf("taskkill /F /IM %s.exe /T", p))
		_ = cmd.Run()
	}
}
