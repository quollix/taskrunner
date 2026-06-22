package taskrunner

import "testing"

func TestLoggerVisually(t *testing.T) {
	tr := GetTaskRunner()
	tr.Log.TaskDescription("Testing logger visually")
	tr.Log.Info("This is an info message")
	tr.Log.Error("This is an error message")
}
