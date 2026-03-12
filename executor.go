package taskrunner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ocelot-cloud/taskrunner/platform"
)

func (t *TaskRunner) KillProcesses(processSubStrings []string) {
	platform.KillProcesses(processSubStrings)
}

func (t *TaskRunner) ExecuteInDir(dir string, commandStr string, envs ...string) {
	shortDir := strings.Replace(dir, t.Config.parentDir, "", -1)
	t.Log.Info("in directory '%s', executing '%s'", shortDir, commandStr)

	cmd := platform.BuildCommand(dir, commandStr)
	appendEnvsToCommand(cmd, envs)

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
		t.Log.Error(" => Command failed. %s. Error: %v", elapsedTimeSummary, err)
		t.ExitWithError()
	} else {
		t.Log.Info(" => Command successful. %s", elapsedTimeSummary)
	}
}

func (t *TaskRunner) Execute(commandStr string, envs ...string) {
	t.ExecuteInDir(".", commandStr, envs...)
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
