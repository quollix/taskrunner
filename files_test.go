//go:build integration

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
	tr.MakeDir(dir)
	assertExists(t, dir)
	p := filepath.Join(dir, name)
	createFile(t, p)
	assertExists(t, p)
	return p
}

func TestDirCopy(t *testing.T) {
	defer tr.Remove(tmpDir, tmpDir2)
	setupDirWithFile(t, tmpDir, "test.txt")
	tr.MakeDir(tmpDir2)

	tr.Copy(tmpDir, tmpDir2)

	assertExists(t, tmpDir)
	assertExists(t, filepath.Join(tmpDir2, "test.txt"))
	assert.Equal(t, mustMode(t, tmpDir), mustMode(t, tmpDir2))
	assert.Equal(t, mustMode(t, filepath.Join(tmpDir, "test.txt")), mustMode(t, filepath.Join(tmpDir2, "test.txt")))
}

func TestDirMove(t *testing.T) {
	defer tr.Remove(tmpDir, tmpDir2)
	setupDirWithFile(t, tmpDir, "test.txt")

	srcDirMode := mustMode(t, tmpDir)
	srcFileMode := mustMode(t, filepath.Join(tmpDir, "test.txt"))

	tr.Move(tmpDir, tmpDir2)

	assertNotExists(t, tmpDir)
	assertExists(t, tmpDir2)
	assertExists(t, filepath.Join(tmpDir2, "test.txt"))
	assert.Equal(t, srcDirMode, mustMode(t, tmpDir2))
	assert.Equal(t, srcFileMode, mustMode(t, filepath.Join(tmpDir2, "test.txt")))
}

func TestRenameDir(t *testing.T) {
	defer tr.Remove(tmpDir)
	setupDirWithFile(t, tmpDir, "x.txt")

	sub := filepath.Join(tmpDir, "subdir")
	tr.MakeDir(sub)
	createFile(t, filepath.Join(sub, "y.txt"))
	assertExists(t, filepath.Join(sub, "y.txt"))

	newSub := filepath.Join(tmpDir, "renamed")
	tr.Rename(sub, newSub)

	assertNotExists(t, sub)
	assertExists(t, newSub)
	assertExists(t, filepath.Join(newSub, "y.txt"))
}

func TestFileCopy(t *testing.T) {
	defer tr.Remove(tmpDir, tmpDir2)
	src := setupDirWithFile(t, tmpDir, "a.txt")
	dst := filepath.Join(tmpDir2, "b.txt")
	tr.MakeDir(tmpDir2)

	tr.Copy(src, dst)

	assertExists(t, src)
	assertExists(t, dst)
	assert.Equal(t, mustMode(t, src), mustMode(t, dst))
}

func TestFileMove(t *testing.T) {
	defer tr.Remove(tmpDir, tmpDir2)
	src := setupDirWithFile(t, tmpDir, "a.txt")
	dst := filepath.Join(tmpDir2, "b.txt")
	srcMode := mustMode(t, src)
	tr.MakeDir(tmpDir2)

	tr.Move(src, dst)

	assertNotExists(t, src)
	assertExists(t, dst)
	assert.Equal(t, srcMode, mustMode(t, dst))
}

func TestFileRename(t *testing.T) {
	defer tr.Remove(tmpDir)
	src := setupDirWithFile(t, tmpDir, "a.txt")
	dst := filepath.Join(tmpDir, "b.txt")
	srcMode := mustMode(t, src)

	tr.Rename(src, dst)

	assertNotExists(t, src)
	assertExists(t, dst)
	assert.Equal(t, srcMode, mustMode(t, dst))
}
