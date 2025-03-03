package sync

import (
	"altimi-interview/logger"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func SyncDirectories(source, target string, deleteMissingFilesFlag bool) error {
	logger.Message(fmt.Sprintf("Starting sync: source=%s, target=%s, deleteMissing=%v at %s",
		source, target, deleteMissingFilesFlag, time.Now().Format(time.RFC3339)))

	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		logger.Error(fmt.Errorf("failed resolving source path: %v", err))
		return err
	}
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		logger.Error(fmt.Errorf("failed resolving target path: %v", err))
		return err
	}

	sourceFiles := make(map[string]os.FileInfo)
	err = filepath.Walk(sourceAbs, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			logger.Error(fmt.Errorf("failed accessing path %s: %v", path, walkErr))
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		relPath, relErr := filepath.Rel(sourceAbs, path)
		if relErr != nil {
			logger.Error(fmt.Errorf("failed getting relative path for %s: %v", path, relErr))
			return nil
		}
		sourceFiles[relPath] = info
		return nil
	})
	if err != nil {
		logger.Error(fmt.Errorf("failed walking source directory: %v", err))
	}

	for relPath, info := range sourceFiles {
		sourcePath := filepath.Join(sourceAbs, relPath)
		targetPath := filepath.Join(targetAbs, relPath)

		targetInfo, tErr := os.Stat(targetPath)
		if os.IsNotExist(tErr) {
			if copyErr := copyFileOrCreateDirs(sourcePath, targetPath); copyErr != nil {
				logger.Error(fmt.Errorf("failed copying file %s to %s: %v", sourcePath, targetPath, copyErr))
			} else {
				logger.Message(fmt.Sprintf("Copied file from %s to %s", sourcePath, targetPath))
			}
		} else if tErr == nil {
			if filesDiffer(info, targetInfo) {
				// update (overwrite target)
				if copyErr := copyFileOrCreateDirs(sourcePath, targetPath); copyErr != nil {
					logger.Error(fmt.Errorf("failed updating file %s to %s: %v", sourcePath, targetPath, copyErr))
				} else {
					logger.Message(fmt.Sprintf("Updated file at %s with %s", targetPath, sourcePath))
				}
			}
		} else {
			// some other error
			logger.Error(fmt.Errorf("failed accessing target file info %s: %v", targetPath, tErr))
		}
	}

	if deleteMissingFilesFlag {
		deleteMissingFiles(targetAbs, sourceFiles)
	}

	return nil
}

func deleteMissingFiles(targetAbs string, sourceFiles map[string]os.FileInfo) {
	err := filepath.Walk(targetAbs, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			logger.Error(fmt.Errorf("failed accessing path %s: %v", path, walkErr))
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		relPath, relErr := filepath.Rel(targetAbs, path)
		if relErr != nil {
			logger.Error(fmt.Errorf("failed getting relative path for target file %s: %v", path, relErr))
			return nil
		}

		if _, found := sourceFiles[relPath]; !found {
			if delErr := os.Remove(path); delErr != nil {
				logger.Error(fmt.Errorf("failed deleting file %s: %v", path, delErr))
			} else {
				logger.Message(fmt.Sprintf("Deleted missing file from target: %s", path))
			}
		}
		return nil
	})

	if err != nil {
		logger.Error(fmt.Errorf("failed walking target directory for delete-missing: %v", err))
	}
}

func filesDiffer(info1, info2 os.FileInfo) bool {
	if info1.Size() != info2.Size() {
		return true
	}

	if info1.Mode() != info2.Mode() {
		return true
	}

	t1 := info1.ModTime().Truncate(time.Second)
	t2 := info2.ModTime().Truncate(time.Second)
	return !t1.Equal(t2)
}

func copyFileOrCreateDirs(src, dst string) error {
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directories %s: %w", dir, err)
	}
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err == nil {
		if chErr := os.Chmod(dst, srcInfo.Mode()); chErr != nil {
			logger.Message(fmt.Sprintf("Warning: unable to preserve file mode on %s: %v", dst, chErr))
		}
		if chtErr := os.Chtimes(dst, time.Now(), srcInfo.ModTime()); chtErr != nil {
			logger.Message(fmt.Sprintf("Warning: unable to preserve mod time on %s: %v", dst, chtErr))
		}
	}

	return nil
}
