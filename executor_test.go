package taskrunner

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	sampleTestDir = "sample_tests"
	tmpDir        = "temp"
	tmpDir2       = "temp2"
)

var tr = GetTaskRunner()

func TestMain(m *testing.M) {
	tr.File.Remove("%s", tmpDir)
	tr.File.Remove("%s", tmpDir2)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCommandSuccessful(t *testing.T) {
	tr.Cmd().Dir(sampleTestDir).Run("go test success_test.go")
}

func TestDirCreationAndDeletion(t *testing.T) {
	assert.False(t, checkIfExists(tmpDir))
	defer tr.File.Remove("%s", tmpDir)
	tr.File.MakeDir("%s", tmpDir)
	assert.True(t, checkIfExists(tmpDir))

	createFile(t, tmpDir+"/test.txt")
	assert.True(t, checkIfExists(tmpDir+"/test.txt"))
	tr.File.Remove("%s", tmpDir)
	assert.False(t, checkIfExists(tmpDir))
}

func checkIfExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	tr.Log.Error("error checking if file exists: %v", err)
	return false
}

func processExists(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err == nil
}

func TestDaemon(t *testing.T) {
	assert.Equal(t, 0, len(tr.daemons))

	tr.Cmd().AsDaemon("sleepy").Run("sleep 100")
	assert.Eventually(t, func() bool {
		return len(tr.daemons) == 1 && tr.daemons[0].name == "sleepy"
	}, time.Second, 10*time.Millisecond)
	pid := tr.daemons[0].cmd.Process.Pid
	assert.True(t, processExists(pid))

	tr.Cleanup()
	assert.Eventually(t, func() bool {
		return len(tr.daemons) == 0
	}, time.Second, 10*time.Millisecond)
	assert.Eventually(t, func() bool {
		return !processExists(pid)
	}, time.Second, 10*time.Millisecond)
}

func TestCustomCleanupFunction(t *testing.T) {
	defer tr.File.Remove("%s", tmpDir)
	previousCleanupFunc := tr.Config.CleanupFunc
	defer func() {
		tr.Config.CleanupFunc = previousCleanupFunc
	}()
	tr.Config.CleanupFunc = func() {
		tr.File.MakeDir("%s", tmpDir)
	}
	assert.False(t, checkIfExists(tmpDir))
	tr.Cleanup()
	assert.True(t, checkIfExists(tmpDir))
}

func TestCommandEnvPassed(t *testing.T) {
	tr.Cmd().
		Dir(sampleTestDir).
		Env("TASKRUNNER_TEST_ENV", "expected").
		Run("go test -run TestEnvVarAvailable env_test.go")
}
