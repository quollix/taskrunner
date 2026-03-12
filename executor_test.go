//go:build integration

package taskrunner

import (
	"fmt"
	"os"
	"testing"

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
	tr.ExecuteInDir(sampleTestDir, "go test success_test.go")
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
	assert.Equal(t, 0, len(tr.Config.idsOfDaemonProcessesCreated))

	tr.StartDaemon(".", "sleep 100")
	assert.Equal(t, 1, len(tr.Config.idsOfDaemonProcessesCreated))
	processId := tr.Config.idsOfDaemonProcessesCreated[0]
	command := fmt.Sprintf("bash -c 'ps -p %d -o cmd= | grep -q sleep'", processId)
	tr.ExecuteInDir(".", command)

	tr.Cleanup()
	assert.Equal(t, 0, len(tr.Config.idsOfDaemonProcessesCreated))
	command = fmt.Sprintf("bash -c '! ps -p %d'", processId)
	tr.ExecuteInDir(".", command)
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
