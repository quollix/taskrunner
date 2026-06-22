package taskrunner

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (f fileOps) Copy(format string, args ...any) *pendingFileTarget {
	return &pendingFileTarget{
		taskRunner: f.taskRunner,
		srcPath:    resolvePath(format, args...),
		action:     "copy",
	}
}

func (f fileOps) Move(format string, args ...any) *pendingFileTarget {
	return &pendingFileTarget{
		taskRunner: f.taskRunner,
		srcPath:    resolvePath(format, args...),
		action:     "move",
	}
}

func (f fileOps) Rename(format string, args ...any) *pendingRenameTarget {
	return &pendingRenameTarget{
		taskRunner: f.taskRunner,
		srcPath:    resolvePath(format, args...),
	}
}

func (f fileOps) MakeDir(format string, args ...any) {
	path := resolvePath(format, args...)
	f.taskRunner.makeDir(path)
}

func (f fileOps) Remove(format string, args ...any) {
	path := resolvePath(format, args...)
	f.taskRunner.remove(path)
}

func (p *pendingFileTarget) To(format string, args ...any) {
	destPath := resolvePath(format, args...)
	switch p.action {
	case "copy":
		p.taskRunner.copy(p.srcPath, destPath)
	case "move":
		p.taskRunner.move(p.srcPath, destPath)
	case "rename":
		p.taskRunner.rename(p.srcPath, destPath)
	default:
		p.taskRunner.Log.Error("unknown file action '%s'", p.action)
	}
}

func (p *pendingRenameTarget) To(format string, args ...any) {
	newName := fmt.Sprintf(format, args...)
	p.taskRunner.rename(p.srcPath, newName)
}

func resolvePath(format string, args ...any) string {
	return filepath.Clean(fmt.Sprintf(format, args...))
}

func (t *TaskRunner) copy(srcPath, destPath string) {
	srcPath = filepath.Clean(srcPath)
	destPath = filepath.Clean(destPath)

	info, err := os.Stat(srcPath)
	if err != nil {
		t.Log.Error("error stating %s: %v", srcPath, err)
		return
	}

	if info.IsDir() {
		if err := os.RemoveAll(destPath); err != nil && !os.IsNotExist(err) {
			t.Log.Error("error removing existing destination %s: %v", destPath, err)
			return
		}
		t.copyDir(srcPath, destPath)
		return
	}

	t.copyFile(srcPath, destPath)
}

func (t *TaskRunner) copyFile(src, dest string) {
	info, err := os.Stat(src)
	if err != nil {
		t.Log.Error("error stating %s: %v", src, err)
		return
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0700); err != nil {
		t.Log.Error("error creating directory %s: %v", filepath.Dir(dest), err)
		return
	}

	in, err := os.Open(src) // #nosec G304 -- taskrunner intentionally copies caller-provided file paths
	if err != nil {
		t.Log.Error("error opening %s: %v", src, err)
		return
	}
	defer func() {
		if err := in.Close(); err != nil {
			t.Log.Error("error closing %s: %v", src, err)
		}
	}()

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode()) // #nosec G304 -- taskrunner intentionally copies caller-provided file paths
	if err != nil {
		t.Log.Error("error creating %s: %v", dest, err)
		return
	}
	defer func() {
		if err := out.Close(); err != nil {
			t.Log.Error("error closing %s: %v", dest, err)
		}
	}()

	if _, err := io.Copy(out, in); err != nil {
		t.Log.Error("error copying from %s to %s: %v", src, dest, err)
		return
	}

	if err := os.Chmod(dest, info.Mode()); err != nil {
		t.Log.Error("error setting mode on %s: %v", dest, err)
	}
}

func (t *TaskRunner) copyDir(srcDir, destDir string) {
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		t.Log.Error("error stating %s: %v", srcDir, err)
		return
	}

	if err := os.MkdirAll(destDir, srcInfo.Mode()); err != nil {
		t.Log.Error("error creating directory %s: %v", destDir, err)
		return
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Log.Error("error reading directory %s: %v", srcDir, err)
		return
	}

	for _, e := range entries {
		s := filepath.Join(srcDir, e.Name())
		d := filepath.Join(destDir, e.Name())
		if e.IsDir() {
			t.copyDir(s, d)
		} else {
			t.copyFile(s, d)
		}
	}

	if err := os.Chmod(destDir, srcInfo.Mode()); err != nil {
		t.Log.Error("error setting mode on directory %s: %v", destDir, err)
	}
}

func (t *TaskRunner) remove(paths ...string) {
	for _, p := range paths {
		p = filepath.Clean(p)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			continue
		}
		if err := os.RemoveAll(p); err != nil {
			t.Log.Error("error removing %s: %v", p, err)
		}
	}
}

func (t *TaskRunner) makeDir(path string) {
	path = filepath.Clean(path)
	if err := os.MkdirAll(path, 0700); err != nil {
		t.Log.Error("error creating %s: %v", path, err)
	}
}

func (t *TaskRunner) move(srcPath, destPath string) {
	srcPath = filepath.Clean(srcPath)
	destPath = filepath.Clean(destPath)

	if err := os.MkdirAll(filepath.Dir(destPath), 0700); err != nil {
		t.Log.Error("error creating directory %s: %v", filepath.Dir(destPath), err)
		return
	}
	if err := os.Rename(srcPath, destPath); err != nil {
		t.Log.Error("error moving %s to %s: %v", srcPath, destPath, err)
	}
}

func (t *TaskRunner) rename(srcPath, newName string) {
	srcPath = filepath.Clean(srcPath)
	newName = filepath.Clean(newName)
	if newName == "." || newName == ".." || filepath.Base(newName) != newName {
		t.Log.Error("invalid rename target %q: expected a name without path separators", newName)
		return
	}
	destPath := filepath.Join(filepath.Dir(srcPath), newName)
	if err := os.Rename(srcPath, destPath); err != nil {
		t.Log.Error("error renaming %s to %s: %v", srcPath, destPath, err)
	}
}
