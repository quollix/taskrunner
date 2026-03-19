package taskrunner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func (t *TaskRunner) KillProcesses(processSubStrings []string) {
	for _, process := range processSubStrings {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("taskkill", "/F", "/IM", process+".exe", "/T")
		} else {
			cmd = exec.Command("pkill", "-f", process)
		}
		_ = cmd.Run()
	}
}

func (c *Command) Dir(dir string) *Command {
	c.dir = dir
	return c
}

func (c *Command) Env(key, value string) *Command {
	c.envs = append(c.envs, fmt.Sprintf("%s=%s", key, value))
	return c
}

func (c *Command) AsDaemon() *Command {
	c.asDaemon = true
	return c
}

func (c *Command) Run(format string, args ...any) {
	commandStr := fmt.Sprintf(format, args...)
	cmd := c.buildCommand(commandStr)

	if c.asDaemon {
		c.startDaemon(cmd)
		return
	}

	c.runForeground(cmd)
}

func (c *Command) buildCommand(commandStr string) *exec.Cmd {
	parts := strings.Fields(commandStr)
	if len(parts) == 0 {
		c.taskRunner.Log.Error("invalid command '%s': empty command", commandStr)
		c.taskRunner.ExitWithError()
		return nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = c.dir
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, c.envs...)
	return cmd
}

func (c *Command) runForeground(cmd *exec.Cmd) {
	shortDir := strings.ReplaceAll(c.dir, c.taskRunner.Config.parentDir, "")
	commandStr := formatCommand(cmd)
	c.taskRunner.Log.Info("in directory '%s', executing '%s'", shortDir, commandStr)

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutMulti := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderrMulti := io.MultiWriter(os.Stderr, &stderrBuf)
	cmd.Stdout = stdoutMulti
	cmd.Stderr = stderrMulti

	startTime := time.Now()
	err := cmd.Run()
	elapsed := time.Since(startTime)
	elapsedStr := fmt.Sprintf("%.3f", elapsed.Seconds())

	elapsedTimeSummary := fmt.Sprintf("Time taken: %s seconds.", elapsedStr)
	if err != nil {
		c.taskRunner.Log.Error(" => Command failed in directory '%s' running '%s'. %s. Error: %v", shortDir, commandStr, elapsedTimeSummary, err)
		c.taskRunner.ExitWithError()
	} else {
		c.taskRunner.Log.Info(" => Command successful. %s", elapsedTimeSummary)
	}
}
