package taskrunner

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func (c *Command) startDaemon(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if cmd.Process == nil {
		c.taskRunner.Log.Error("error - the process was not able to start properly.")
		c.taskRunner.ExitWithError()
		return
	}

	if err != nil {
		c.taskRunner.Log.Error("Command: '%s' -> failed with error: %v", formatCommand(cmd), err)
		c.taskRunner.ExitWithError()
		return
	}

	c.taskRunner.registerDaemon(cmd)
	c.taskRunner.Log.Info("started daemon with ID '%v' using command '%s'", cmd.Process.Pid, formatCommand(cmd))

	go func() {
		commandStr := formatCommand(cmd)
		if err = cmd.Wait(); err != nil {
			if err.Error() == "signal: killed" {
				c.taskRunner.Log.Info("command: '%s' -> stopped through cleanup process killing", commandStr)
			} else {
				c.taskRunner.Log.Error("command: '%s' -> stopped with error: %v", commandStr, err)
			}
		} else {
			c.taskRunner.Log.Info("command: '%s' -> stopped through termination", commandStr)
		}
	}()
}

func (t *TaskRunner) Cleanup() {
	t.Log.Info("\ncleanup method called")
	t.killDaemonProcessesCreateDuringThisRun()
	if t.Config.CleanupFunc != nil {
		t.Log.Info("calling custom cleanup function")
		t.Config.CleanupFunc()
	}
	t.resetCursor()
}

func (t *TaskRunner) ExitWithError() {
	if t.Config.CleanupOnFailure && t.Config.CleanupFunc != nil {
		t.Cleanup()
	} else {
		t.killDaemonProcessesCreateDuringThisRun()
	}
	t.resetCursor()
	os.Exit(1)
}

func (t *TaskRunner) resetCursor() {
	fmt.Print("\x1b[?25h") // Shows the terminal cursor again if it was hidden.
	fmt.Print("\x1b[0m")   // Resets all terminal text attributes (color, bold, underline) back to default.
}

func (t *TaskRunner) killDaemonProcessesCreateDuringThisRun() {
	t.daemonMu.Lock()
	daemons := append([]*daemonProcess(nil), t.daemons...)
	t.daemons = nil
	t.daemonMu.Unlock()

	if len(daemons) == 0 {
		return
	}
	t.Log.Info("Killing daemon processes")
	for _, daemon := range daemons {
		if daemon.cmd == nil || daemon.cmd.Process == nil {
			continue
		}
		t.Log.Info("  Killing process with ID '%v'", daemon.cmd.Process.Pid)
		if err := daemon.cmd.Process.Kill(); err != nil {
			t.Log.Error("Failed to kill process with ID '%v' because of error: %v", daemon.cmd.Process.Pid, err)
		}
	}
}

func (t *TaskRunner) registerDaemon(cmd *exec.Cmd) {
	t.daemonMu.Lock()
	defer t.daemonMu.Unlock()
	t.daemons = append(t.daemons, &daemonProcess{
		cmd:     cmd,
		command: formatCommand(cmd),
	})
}

func formatCommand(cmd *exec.Cmd) string {
	return strings.Join(append([]string{cmd.Path}, cmd.Args[1:]...), " ")
}

func (t *TaskRunner) EnableAbortForKeystrokeControlPlusC() {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		t.Log.Info("Received signal: %v. Initiating graceful shutdown...", sig)
		t.Cleanup()
		os.Exit(1)
	}()
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
