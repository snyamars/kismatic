package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestMoveDirExists(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "ket-backupdir-test")
	if err != nil {
		t.Fatalf("Error creating temp dir: %v", err)
	}
	sourceDir := filepath.Join(tmpDir, ".helm")
	err = os.Mkdir(sourceDir, 0755)
	if err != nil {
		t.Errorf("Expected error creating /tmp to be nil, got: %v", err)
	}

	exists, err := MoveDir(sourceDir, filepath.Join(tmpDir, ".helm.bak"))
	if err != nil {
		t.Errorf("Expected error to be nil, got: %v", err)
	}
	if !exists {
		t.Errorf("Expected directory to exist")
	}
}

func TestBackupClientDirectoryNotExists(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "ket-backupdir-test")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}
	exists, err := MoveDir(filepath.Join(tmpDir, ".helm"), filepath.Join(tmpDir, ".helm.bak"))
	if err != nil {
		t.Errorf("Expected error to be nil, got: %v", err)
	}
	if exists {
		t.Errorf("Expected directory to not exist")
	}
}
