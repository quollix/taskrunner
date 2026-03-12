package taskrunner

import (
	"log"
	"os"
	"path/filepath"
)

func GetTaskRunner() *TaskRunner {
	return &TaskRunner{
		Config: &Config{
			CleanupOnFailure:            true,
			CleanupFunc:                 nil,
			DefaultEnvironmentVariables: []string{},
			idsOfDaemonProcessesCreated: []int{},
			parentDir:                   getParentDir(),
		},
		Log: consoleLogger{},
	}
}

type TaskRunner struct {
	Config *Config
	Log    logger
}

type Config struct {
	CleanupOnFailure            bool
	CleanupFunc                 func()
	DefaultEnvironmentVariables []string
	idsOfDaemonProcessesCreated []int
	parentDir                   string
}

func getParentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current dir: %v", err)
	}
	return filepath.Dir(currentDir)
}
