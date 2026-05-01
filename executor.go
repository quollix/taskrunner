package taskrunner

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode"
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

func (c *Command) AllowFail() *Command {
	c.allowFail = true
	return c
}

func (c *Command) Output() string {
	return c.lastOutput
}

func (c *Command) Run(format string, args ...any) *Command {
	commandStr := fmt.Sprintf(format, args...)
	cmd := c.buildCommand(commandStr)

	if c.asDaemon {
		c.startDaemon(cmd)
		return c
	}

	c.runForeground(cmd)
	return c
}

func (c *Command) buildCommand(commandStr string) *exec.Cmd {
	parts, err := splitCommandArgs(commandStr)
	if err != nil {
		c.taskRunner.Log.Error("invalid command '%s': %v", commandStr, err)
		c.taskRunner.ExitWithError()
		return nil
	}
	if len(parts) == 0 {
		c.taskRunner.Log.Error("invalid command '%s': empty command", commandStr)
		c.taskRunner.ExitWithError()
		return nil
	}

	cmd := exec.Command(parts[0], parts[1:]...) // #nosec G204 -- taskrunner intentionally executes caller-provided commands
	cmd.Dir = c.dir
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, c.envs...)
	return cmd
}

func splitCommandArgs(commandStr string) ([]string, error) {
	var parts []string
	var current []rune
	inSingleQuotes := false

	flushCurrent := func() {
		if len(current) == 0 {
			return
		}
		parts = append(parts, string(current))
		current = nil
	}

	for _, r := range commandStr {
		switch {
		case r == '\'':
			inSingleQuotes = !inSingleQuotes
		case unicode.IsSpace(r) && !inSingleQuotes:
			flushCurrent()
		default:
			current = append(current, r)
		}
	}

	if inSingleQuotes {
		return nil, errors.New("unmatched single quote")
	}

	flushCurrent()
	return parts, nil
}

func (c *Command) runForeground(cmd *exec.Cmd) {
	shortDir := strings.ReplaceAll(c.dir, c.taskRunner.Config.parentDir, "")
	commandStr := formatCommand(cmd)
	c.taskRunner.Log.Info("in directory '%s', executing '%s'", shortDir, commandStr)

	var outputBuf bytes.Buffer
	c.attachOutput(cmd, logPrefix, green, &lockedWriter{writer: &outputBuf})

	startTime := time.Now()
	err := cmd.Run()
	elapsed := time.Since(startTime)
	elapsedStr := fmt.Sprintf("%.3f", elapsed.Seconds())
	c.lastOutput = outputBuf.String()

	elapsedTimeSummary := fmt.Sprintf("Time taken: %s seconds.", elapsedStr)
	if err != nil {
		if c.allowFail {
			c.taskRunner.Log.Info(" => Command failed in directory '%s' running '%s' but continuing because AllowFail was set. %s. Error: %v", shortDir, commandStr, elapsedTimeSummary, err)
			return
		}
		c.taskRunner.Log.Error(" => Command failed in directory '%s' running '%s'. %s. Error: %v", shortDir, commandStr, elapsedTimeSummary, err)
		c.taskRunner.ExitWithError()
	} else {
		c.taskRunner.Log.Info(" => Command successful. %s", elapsedTimeSummary)
	}
}

func (c *Command) attachOutput(cmd *exec.Cmd, name, color string, capture io.Writer) {
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

	go streamOutput(stdoutPipe, os.Stdout, capture, name, color)
	go streamOutput(stderrPipe, os.Stderr, capture, name, color)
}

func streamOutput(reader io.Reader, console, capture io.Writer, name, color string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		_, _ = fmt.Fprintf(console, "%s[%s] %s%s\n", color, name, line, reset)
		if capture != nil {
			_, _ = fmt.Fprintln(capture, line)
		}
	}
}

type lockedWriter struct {
	mu     sync.Mutex
	writer io.Writer
}

func (w *lockedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.writer.Write(p)
}
