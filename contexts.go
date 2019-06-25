package fs

import (
	"io/ioutil"
	"os"
)

// WithTempDir creates a temporary directory and deletes it when f returns.
func WithTempDir(prefix string, f func(tmpDir string) error) error {
	tmpDir, err := ioutil.TempDir("", prefix)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	return f(tmpDir)
}

// WithTempFile creates a temporary file and deletes it when f returns.
func WithTempFile(pattern string, f func(tmpFile string) error) error {
	tmpFile, err := ioutil.TempFile("", pattern)
	if err != nil {
		return err
	}
	tmpFileName := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpFileName)

	return f(tmpFileName)
}
