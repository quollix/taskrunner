package taskrunner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/quollix/taskrunner/platform"
)

func (t *TaskRunner) KillProcesses(processSubStrings []string) {
	platform.KillProcesses(processSubStrings)
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
		c.startDaemon(cmd, commandStr)
		return
	}

	c.runForeground(cmd, commandStr)
}

func (c *Command) buildCommand(commandStr string) *exec.Cmd {
	cmd := platform.BuildCommand(c.dir, commandStr)
	appendEnvsToCommand(cmd, c.envs)
	return cmd
}

func (c *Command) runForeground(cmd *exec.Cmd, commandStr string) {
	shortDir := strings.Replace(c.dir, c.taskRunner.Config.parentDir, "", -1)
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

func (t *TaskRunner) PromptForContinuation(prompt string) {
	fmt.Printf("%s (y/N): ", prompt)
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		fmt.Println("Command aborted.")
		os.Exit(0)
	}
}
