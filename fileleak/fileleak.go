//go:build !solution

package fileleak

import (
	"io/fs"
	"os"
	"reflect"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func getFileDescriptors() ([]fs.DirEntry, error) {
	return os.ReadDir("/proc/self/fd")
}

func getFileInfo(entries []fs.DirEntry) []fs.FileInfo {
	files := make([]fs.FileInfo, 0)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, info)
	}
	return files
}

func isFileLeakDetected(startFiles, endFiles []fs.FileInfo) bool {
	return !reflect.DeepEqual(startFiles, endFiles)
}

func VerifyNone(t testingT) {
	entries, err := getFileDescriptors()
	if err != nil {
		t.Errorf("failed to read file descriptors: %v", err)
		return
	}
	startFiles := getFileInfo(entries)

	t.Cleanup(func() {
		entries, err := getFileDescriptors()
		if err != nil {
			t.Errorf("failed to read file descriptors: %v", err)
			return
		}
		endFiles := getFileInfo(entries)

		if isFileLeakDetected(startFiles, endFiles) {
			t.Errorf("file leak detected")
		}
	})
}
