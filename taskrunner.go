package taskrunner

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func GetTaskRunner() *TaskRunner {
	tr := &TaskRunner{
		Config: &Config{
			CleanupOnFailure:            true,
			CleanupFunc:                 nil,
			DefaultEnvironmentVariables: []string{},
			parentDir:                   getParentDir(),
			DefaultWaitTimeout:          60,
		},
		Log: consoleLogger{},
	}
	tr.File = fileOps{taskRunner: tr}
	return tr
}

type TaskRunner struct {
	Config          *Config
	File            fileOps
	Log             logger
	daemons         []*daemonProcess
	nextDaemonColor int
}

type Command struct {
	taskRunner *TaskRunner
	dir        string
	envs       []string
	asDaemon   bool
	name       string
	allowFail  bool
	lastOutput string
}

type daemonProcess struct {
	cmd     *exec.Cmd
	command string
	name    string
	color   string
}

type fileOps struct {
	taskRunner *TaskRunner
}

type pendingFileTarget struct {
	taskRunner *TaskRunner
	srcPath    string
	action     string
}

type pendingRenameTarget struct {
	taskRunner *TaskRunner
	srcPath    string
}

type Config struct {
	CleanupOnFailure            bool
	CleanupFunc                 func()
	DefaultEnvironmentVariables []string
	parentDir                   string
	DefaultWaitTimeout          int
}

func (t *TaskRunner) Cmd() *Command {
	return &Command{
		taskRunner: t,
		dir:        ".",
		envs:       []string{},
	}
}

func getParentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current dir: %v", err)
	}
	return filepath.Dir(currentDir)
}
