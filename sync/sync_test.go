package sync

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFilesDiffer_SizeDifferent(t *testing.T) {
	tmpFile1, err := os.CreateTemp("", "sync_src_1")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile1.Name())

	tmpFile2, err := os.CreateTemp("", "sync_tgt_1")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile2.Name())

	// Write some data
	_, _ = tmpFile1.Write([]byte("some data"))
	_, _ = tmpFile2.Write([]byte("different size data"))

	fi1, _ := os.Stat(tmpFile1.Name())
	fi2, _ := os.Stat(tmpFile2.Name())

	//then
	assert.Truef(t, filesDiffer(fi1, fi2), "files should differ based on size", nil)
}

func TestFilesDiffer_TimeDifferent(t *testing.T) {
	//given
	tmpFile1, err := os.CreateTemp("", "sync_src_2")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile1.Name())

	tmpFile2, err := os.CreateTemp("", "sync_tgt_2")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile2.Name())

	data := []byte("identical data")
	_, _ = tmpFile1.Write(data)
	_, _ = tmpFile2.Write(data)

	os.Chtimes(tmpFile1.Name(), time.Now(), time.Now().Add(-10*time.Hour))

	fi1, _ := os.Stat(tmpFile1.Name())
	fi2, _ := os.Stat(tmpFile2.Name())

	//then
	assert.Truef(t, filesDiffer(fi1, fi2), "files should differ based on time of editing", nil)
}

func TestCopyFile(t *testing.T) {
	//given
	tmpFile1, err := os.CreateTemp("", "sync_src_3")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile1.Name())

	srcData := []byte("content to copy")
	_, _ = tmpFile1.Write(srcData)
	tmpFile1.Close()

	dstFile, err := os.CreateTemp("", "sync_dst_3")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(dstFile.Name())
	dstFile.Close()

	//when
	if err := copyFile(tmpFile1.Name(), dstFile.Name()); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	copiedData, err := os.ReadFile(dstFile.Name())
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	//then
	assert.Equalf(t, string(copiedData), string(srcData), "Expected copied data %s, got %s", srcData, copiedData)
}

func TestCopyFileOrCreateDirs(t *testing.T) {
	// given
	srcDir, err := os.MkdirTemp("", "sync_test_src_4")
	if err != nil {
		t.Fatalf("Failed to create temp source dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	srcFilePath := filepath.Join(srcDir, "testfile.txt")
	err = os.WriteFile(srcFilePath, []byte("file content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write src file: %v", err)
	}

	targetDir, err := os.MkdirTemp("", "sync_test_tgt_4")
	if err != nil {
		t.Fatalf("Failed to create temp target dir: %v", err)
	}
	defer os.RemoveAll(targetDir)

	dstFilePath := filepath.Join(targetDir, "sub", "folder", "testfile.txt")

	// when
	if err := copyFileOrCreateDirs(srcFilePath, dstFilePath); err != nil {
		t.Fatalf("Failed to copy file with directory creation: %v", err)
	}

	// then
	info, err := os.Stat(dstFilePath)
	assert.NoError(t, err, "copied file should exist")
	assert.NotZerof(t, info.Size(), "Expected the file to have content, got size=0")
}
