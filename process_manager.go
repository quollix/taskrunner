package taskrunner

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/quollix/taskrunner/platform"
)

func (c *Command) startDaemon(cmd *exec.Cmd, commandStr string) {
	platform.SetProcessGroup(cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if cmd.Process == nil {
		c.taskRunner.Log.Error("error - the process was not able to start properly.")
		c.taskRunner.ExitWithError()
		return
	}

	c.taskRunner.Config.idsOfDaemonProcessesCreated = append(c.taskRunner.Config.idsOfDaemonProcessesCreated, cmd.Process.Pid)

	if err != nil {
		c.taskRunner.Log.Error("Command: '%s' -> failed with error: %v", commandStr, err)
		c.taskRunner.ExitWithError()
		return
	}

	c.taskRunner.Log.Info("started daemon with ID '%v' using command '%s'", cmd.Process.Pid, commandStr)

	go func() {
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
	t.ResetCursor()
}

func (t *TaskRunner) ExitWithError() {
	if t.Config.CleanupOnFailure && t.Config.CleanupFunc != nil {
		t.Cleanup()
	} else {
		t.killDaemonProcessesCreateDuringThisRun()
	}
	t.ResetCursor()
	os.Exit(1)
}

func (t *TaskRunner) ResetCursor() {
	fmt.Print("\x1b[?25h") // Shows the terminal cursor again if it was hidden.
	fmt.Print("\x1b[0m")   // Resets all terminal text attributes (color, bold, underline) back to default.
}

func (t *TaskRunner) killDaemonProcessesCreateDuringThisRun() {
	if len(t.Config.idsOfDaemonProcessesCreated) == 0 {
		return
	}
	t.Log.Info("Killing daemon processes")
	for _, pid := range t.Config.idsOfDaemonProcessesCreated {
		t.Log.Info("  Killing process with ID '%v'", pid)
		if err := platform.KillProcessGroup(pid); err != nil {
			t.Log.Error("Failed to kill process with ID '%v' because of error: %v", pid, err)
		}
	}
	t.Config.idsOfDaemonProcessesCreated = nil
}

func appendEnvsToCommand(cmd *exec.Cmd, envs []string) {
	envsWithLogLevel := append(envs, DefaultEnvs...)
	cmd.Env = append(os.Environ(), envsWithLogLevel...)
}

var DefaultEnvs []string

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
