package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// MoveDir checks for existence of the srcDir and moves it to dstDir
func MoveDir(srcDir string, dstDir string) (bool, error) {
	_, err := os.Stat(srcDir)
	// Directory exists
	if err == nil {
		if err = os.Rename(srcDir, dstDir); err != nil {
			return false, fmt.Errorf("Could not back up %q directory: %v", srcDir, err)
		}
		return true, nil
	} else if !os.IsNotExist(err) { // Directory does not exist but got some other error
		return false, fmt.Errorf("Could not determine if %q directory exists: %v", srcDir, err)
	}
	// Directory does not already exist, nothing to do
	return false, nil
}

// CopyDir copies srcDir to dstDir, doesn't matter if src is a directory or a file
func CopyDir(srcDir, dstDir string) error {
	info, err := os.Stat(srcDir)
	if err != nil {
		return err
	}
	return copyInfo(srcDir, dstDir, info)
}

// "info" must be given here, NOT nil.
func copyInfo(src, dest string, info os.FileInfo) error {
	if info.IsDir() {
		return dcopy(src, dest, info)
	}
	return fcopy(src, dest, info)
}

func fcopy(src, dest string, info os.FileInfo) error {

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	return err
}

func dcopy(src, dest string, info os.FileInfo) error {

	if err := os.MkdirAll(dest, info.Mode()); err != nil {
		return err
	}

	infos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if err := copyInfo(
			filepath.Join(src, info.Name()),
			filepath.Join(dest, info.Name()),
			info,
		); err != nil {
			return err
		}
	}

	return nil
}
