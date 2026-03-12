package taskrunner

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type logger interface {
	Info(string, ...any)
	Error(string, ...any)
	TaskDescription(string, ...any)
}

const (
	green    = "\033[32m"
	red      = "\033[31m"
	blueBold = "\033[1;34m"
	reset    = "\033[0m"
)

type consoleLogger struct{}

func (consoleLogger) Info(f string, a ...any) {
	file, line := callerInfo()
	fmt.Printf(green+"[%s:%d] "+f+reset+"\n", append([]any{file, line}, a...)...)
}

func (consoleLogger) Error(f string, a ...any) {
	file, line := callerInfo()
	fmt.Printf(red+"[%s:%d] "+f+reset+"\n", append([]any{file, line}, a...)...)
}

func (consoleLogger) TaskDescription(f string, a ...any) {
	title := fmt.Sprintf(f, a...)
	fmt.Printf("\n"+blueBold+"==== %s ===="+reset+"\n\n", title)
}

func callerInfo() (string, int) {
	// skip 2 = skip callerInfo + Info/Error → return actual user code line
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "???", 0
	}

	file = filepath.Base(file)

	return file, line
}
