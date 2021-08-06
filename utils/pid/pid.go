package pid

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// File is a file used to store the process ID of a running process.
type File struct {
	path string
}

// New creates a file using the specified path.
func New(path string) *File {
	return &File{path: path}
}

// Init will write current process pid to file
func (file *File) Init() error {
	// Note MkdirAll returns nil if a directory already exists
	err := os.MkdirAll(filepath.Dir(file.path), os.FileMode(0755))
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(file.path, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Remove removes the PIDFile. Do not use it without Init executed
func (file *File) Remove() error {
	return os.Remove(file.path)
}

func (file *File) Get() (int, error) {
	pidByte, err := ioutil.ReadFile(file.path)
	if err != nil {
		return 0, err
	}

	pidString := strings.TrimSpace(string(pidByte))
	pid, err := strconv.Atoi(pidString)
	if err != nil {
		return 0, err
	}

	return pid, nil
}
