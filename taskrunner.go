package taskrunner

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func GetTaskRunner() *TaskRunner {
	return &TaskRunner{
		Config: &Config{
			CleanupOnFailure:            true,
			CleanupFunc:                 nil,
			DefaultEnvironmentVariables: []string{},
			parentDir:                   getParentDir(),
		},
		Log: consoleLogger{},
	}
}

type TaskRunner struct {
	Config          *Config
	Log             logger
	nextDaemonColor int
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
