package taskrunner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createFile(t *testing.T, path string) {
	f, err := os.Create(path)
	assert.NoError(t, err)
	assert.NoError(t, f.Close())
}

func mustMode(t *testing.T, path string) os.FileMode {
	i, err := os.Stat(path)
	assert.NoError(t, err)
	return i.Mode()
}

func assertExists(t *testing.T, path string) {
	assert.True(t, checkIfExists(path))
}

func assertNotExists(t *testing.T, path string) {
	assert.False(t, checkIfExists(path))
}

func setupDirWithFile(t *testing.T, dir, name string) string {
	tr.File.Remove("%s", dir)
	tr.File.MakeDir("%s", dir)
	assertExists(t, dir)
	p := filepath.Join(dir, name)
	createFile(t, p)
	assertExists(t, p)
	return p
}

func TestDirCopy(t *testing.T) {
	tr.File.Remove("%s", tmpDir)
	tr.File.Remove("%s", tmpDir2)
	defer tr.File.Remove("%s", tmpDir)
	defer tr.File.Remove("%s", tmpDir2)
	setupDirWithFile(t, tmpDir, "test.txt")
	tr.File.MakeDir("%s", tmpDir2)

	tr.File.Copy("%s", tmpDir).To("%s", tmpDir2)

	assertExists(t, tmpDir)
	assertExists(t, filepath.Join(tmpDir2, "test.txt"))
	assert.Equal(t, mustMode(t, tmpDir), mustMode(t, tmpDir2))
	assert.Equal(t, mustMode(t, filepath.Join(tmpDir, "test.txt")), mustMode(t, filepath.Join(tmpDir2, "test.txt")))
}

func TestDirMove(t *testing.T) {
	tr.File.Remove("%s", tmpDir)
	tr.File.Remove("%s", tmpDir2)
	defer tr.File.Remove("%s", tmpDir)
	defer tr.File.Remove("%s", tmpDir2)
	setupDirWithFile(t, tmpDir, "test.txt")

	srcDirMode := mustMode(t, tmpDir)
	srcFileMode := mustMode(t, filepath.Join(tmpDir, "test.txt"))

	tr.File.Move("%s", tmpDir).To("%s", tmpDir2)

	assertNotExists(t, tmpDir)
	assertExists(t, tmpDir2)
	assertExists(t, filepath.Join(tmpDir2, "test.txt"))
	assert.Equal(t, srcDirMode, mustMode(t, tmpDir2))
	assert.Equal(t, srcFileMode, mustMode(t, filepath.Join(tmpDir2, "test.txt")))
}

func TestRenameDir(t *testing.T) {
	tr.File.Remove("%s", tmpDir)
	defer tr.File.Remove("%s", tmpDir)
	setupDirWithFile(t, tmpDir, "x.txt")

	sub := filepath.Join(tmpDir, "subdir")
	tr.File.MakeDir("%s", sub)
	createFile(t, filepath.Join(sub, "y.txt"))
	assertExists(t, filepath.Join(sub, "y.txt"))

	newSub := filepath.Join(tmpDir, "renamed")
	tr.File.Rename("%s", sub).To("%s", "renamed")

	assertNotExists(t, sub)
	assertExists(t, newSub)
	assertExists(t, filepath.Join(newSub, "y.txt"))
}

func TestFileCopy(t *testing.T) {
	tr.File.Remove("%s", tmpDir)
	tr.File.Remove("%s", tmpDir2)
	defer tr.File.Remove("%s", tmpDir)
	defer tr.File.Remove("%s", tmpDir2)
	src := setupDirWithFile(t, tmpDir, "a.txt")
	dst := filepath.Join(tmpDir2, "b.txt")
	tr.File.MakeDir("%s", tmpDir2)

	tr.File.Copy("%s", src).To("%s", dst)

	assertExists(t, src)
	assertExists(t, dst)
	assert.Equal(t, mustMode(t, src), mustMode(t, dst))
}

func TestFileMove(t *testing.T) {
	tr.File.Remove("%s", tmpDir)
	tr.File.Remove("%s", tmpDir2)
	defer tr.File.Remove("%s", tmpDir)
	defer tr.File.Remove("%s", tmpDir2)
	src := setupDirWithFile(t, tmpDir, "a.txt")
	dst := filepath.Join(tmpDir2, "b.txt")
	srcMode := mustMode(t, src)
	tr.File.MakeDir("%s", tmpDir2)

	tr.File.Move("%s", src).To("%s", dst)

	assertNotExists(t, src)
	assertExists(t, dst)
	assert.Equal(t, srcMode, mustMode(t, dst))
}

func TestFileRename(t *testing.T) {
	tr.File.Remove("%s", tmpDir)
	defer tr.File.Remove("%s", tmpDir)
	src := setupDirWithFile(t, tmpDir, "a.txt")
	dst := filepath.Join(tmpDir, "b.txt")
	srcMode := mustMode(t, src)

	tr.File.Rename("%s", src).To("%s", "b.txt")

	assertNotExists(t, src)
	assertExists(t, dst)
	assert.Equal(t, srcMode, mustMode(t, dst))
}

func TestFileRenameRejectsPath(t *testing.T) {
	tr.File.Remove("%s", tmpDir)
	defer tr.File.Remove("%s", tmpDir)
	src := setupDirWithFile(t, tmpDir, "a.txt")

	tr.File.Rename("%s", src).To("%s", filepath.Join("nested", "b.txt"))

	assertExists(t, src)
	assertNotExists(t, filepath.Join(tmpDir, "nested", "b.txt"))
}
