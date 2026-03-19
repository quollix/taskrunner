package taskrunner

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func GetTaskRunner() *TaskRunner {
	tr := &TaskRunner{
		Config: &Config{
			CleanupOnFailure:            true,
			CleanupFunc:                 nil,
			DefaultEnvironmentVariables: []string{},
			parentDir:                   getParentDir(),
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
	daemonMu        sync.Mutex
}

type Command struct {
	taskRunner *TaskRunner
	dir        string
	envs       []string
	asDaemon   bool
	name       string
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
