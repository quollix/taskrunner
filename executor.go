package taskrunner

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func (c *Command) Dir(dir string) *Command {
	c.dir = dir
	return c
}

func (c *Command) Env(key, value string) *Command {
	c.envs = append(c.envs, fmt.Sprintf("%s=%s", key, value))
	return c
}

func (c *Command) AsDaemon(name string) *Command {
	c.asDaemon = true
	c.name = name
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
	c.attachOutput(cmd, logPrefix, green, &stdoutBuf, &stderrBuf)

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

func (c *Command) attachOutput(cmd *exec.Cmd, name, color string, stdoutBuf, stderrBuf io.Writer) {
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		c.taskRunner.Log.Error("failed to attach stdout pipe for '%s': %v", formatCommand(cmd), err)
		c.taskRunner.ExitWithError()
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		c.taskRunner.Log.Error("failed to attach stderr pipe for '%s': %v", formatCommand(cmd), err)
		c.taskRunner.ExitWithError()
		return
	}

	go streamOutput(stdoutPipe, io.MultiWriter(os.Stdout, stdoutBuf), name, color)
	go streamOutput(stderrPipe, io.MultiWriter(os.Stderr, stderrBuf), name, color)
}

func streamOutput(reader io.Reader, raw io.Writer, name, color string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		_, _ = fmt.Fprintf(raw, "%s[%s] %s%s\n", color, name, line, reset)
	}
}
