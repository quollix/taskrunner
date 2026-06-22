package taskrunner

import (
	"fmt"
	"runtime"
	"strings"
)

type logger interface {
	Info(string, ...any)
	Error(string, ...any)
	TaskDescription(string, ...any)
}

const (
	green      = "\033[32m"
	red        = "\033[31m"
	blueBold   = "\033[1;34m"
	cyan       = "\033[36m"
	yellow     = "\033[33m"
	blue       = "\033[34m"
	magenta    = "\033[35m"
	brightCyan = "\033[96m"
	brightBlue = "\033[94m"
	brightYel  = "\033[93m"
	white      = "\033[97m"
	reset      = "\033[0m"
)

const logPrefix = "taskrunner"

type consoleLogger struct{}

func (consoleLogger) Info(f string, a ...any) {
	fmt.Printf(green+"[%s] "+f+reset+"\n", append([]any{logPrefix}, a...)...)
}

func (consoleLogger) Error(f string, a ...any) {
	message := fmt.Sprintf(f, a...)
	fmt.Printf(red+"[%s] %s"+reset+"\n", logPrefix, message)
	stack := stackTrace(3)
	if stack == "" {
		return
	}
	fmt.Printf(red+"[%s] Stack trace:\n%s"+reset, logPrefix, stack)
}

func (consoleLogger) TaskDescription(f string, a ...any) {
	title := fmt.Sprintf(f, a...)
	fmt.Printf("\n"+blueBold+"==== %s ===="+reset+"\n\n", title)
}

func stackTrace(skip int) string {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(skip, pcs)
	if n == 0 {
		return ""
	}

	var b strings.Builder
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&b, "  %s\n    %s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}

	return b.String()
}
