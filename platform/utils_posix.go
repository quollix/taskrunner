//go:build !windows

package platform

import (
	"fmt"
	"os"
	"os/exec"
)

func BuildCommand(dir, commandStr string) *exec.Cmd {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}
	cmd := exec.Command(shell, "-c", commandStr)
	cmd.Dir = dir
	return cmd
}

func KillProcesses(processes []string) {
	for _, p := range processes {
		cmd := exec.Command("bash", "-c", fmt.Sprintf("pgrep -f %s | xargs -r kill -9", p))
		_ = cmd.Run()
	}
}
