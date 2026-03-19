//go:build integration

package taskrunner

import (
	"os"
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
	tr.Remove(tmpDir, tmpDir2)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCommandSuccessful(t *testing.T) {
	tr.Cmd().Dir(sampleTestDir).Run("go test success_test.go")
}

func TestDirCreationAndDeletion(t *testing.T) {
	assert.False(t, checkIfExists(tmpDir))
	defer tr.Remove(tmpDir)
	tr.MakeDir(tmpDir)
	assert.True(t, checkIfExists(tmpDir))
	
	createFile(t, tmpDir+"/test.txt")
	assert.True(t, checkIfExists(tmpDir+"/test.txt"))
	tr.Remove(tmpDir)
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

func TestDaemon(t *testing.T) {
	tr.daemonMu.Lock()
	assert.Equal(t, 0, len(tr.daemons))
	tr.daemonMu.Unlock()

	tr.Cmd().AsDaemon().Run("sleep 100")
	assert.Eventually(t, func() bool {
		tr.daemonMu.Lock()
		defer tr.daemonMu.Unlock()
		return len(tr.daemons) == 1
	}, time.Second, 10*time.Millisecond)

	tr.Cleanup()
	assert.Eventually(t, func() bool {
		tr.daemonMu.Lock()
		defer tr.daemonMu.Unlock()
		return len(tr.daemons) == 0
	}, time.Second, 10*time.Millisecond)
}

func TestCustomCleanupFunction(t *testing.T) {
	defer tr.Remove(tmpDir)
	tr.Config.CleanupFunc = func() {
		tr.MakeDir(tmpDir)
	}
	assert.False(t, checkIfExists(tmpDir))
	tr.Cleanup()
	assert.True(t, checkIfExists(tmpDir))
}
