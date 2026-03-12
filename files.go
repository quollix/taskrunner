package taskrunner

import (
	"io"
	"os"
	"path/filepath"
)

func (t *TaskRunner) Copy(srcPath, destPath string) {
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

	in, err := os.Open(src)
	if err != nil {
		t.Log.Error("error opening %s: %v", src, err)
		return
	}
	defer in.Close()

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		t.Log.Error("error creating %s: %v", dest, err)
		return
	}
	defer out.Close()

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

func (t *TaskRunner) Remove(paths ...string) {
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

func (t *TaskRunner) MakeDir(path string) {
	path = filepath.Clean(path)
	if err := os.MkdirAll(path, 0700); err != nil {
		t.Log.Error("error creating %s: %v", path, err)
	}
}

func (t *TaskRunner) Move(srcPath, destPath string) {
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

func (t *TaskRunner) Rename(srcPath, destPath string) {
	srcPath = filepath.Clean(srcPath)
	destPath = filepath.Clean(destPath)

	if err := os.MkdirAll(filepath.Dir(destPath), 0700); err != nil {
		t.Log.Error("error creating directory %s: %v", filepath.Dir(destPath), err)
		return
	}
	if err := os.Rename(srcPath, destPath); err != nil {
		t.Log.Error("error renaming %s to %s: %v", srcPath, destPath, err)
	}
}
