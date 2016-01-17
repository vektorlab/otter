/*
A File represents a single file on an operating system.
States -
  absent: A file is not present on the operating system // TODO
  linked: A symbolic link is created from one destination to another // TODO
  rendered: A file is copied from a another source and rendered // TODO

Source Types -
  git: A file is copied from a Git repository // TODO
  s3: A file is copied from an Amazon S3 repo // TODO
  file: A file is copied from the local file system // TODO
*/

package state

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type FileType string

type File struct {
	Mode     string   `json:"mode"`   // File should be set to this octal mode
	Path     string   `json:"path"`   // File destination
	Source   string   `json:"source"` // File source
	Metadata Metadata `json:"metadata"`
	Require  []string `json:"require"`
}

func (file *File) retrieveFile() ([]byte, error) {
	var body []byte
	switch {
	case strings.Contains(file.Source, "http://"):
		log.Printf("Calling HTTP GET: %s", file.Source)
		resp, err := http.Get(file.Source)
		if err != nil {
			return body, err
		}
		defer resp.Body.Close() // TODO: Buffer and error if a maximum file size is exceeded.
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return body, err
		}
		return body, nil
	default:
		return body, fmt.Errorf("Unable to parse source type: %s", file.Source)
	}
}

func (file *File) writeFile() error {
	data, err := file.retrieveFile()
	if err != nil {
		return err
	}
	u, err := strconv.ParseUint(file.Mode, 8, 32)
	if err != nil {
		return err
	}
	mode := os.FileMode(u)
	log.Printf("Writing to file [%s] %s", mode, file.Path)
	return ioutil.WriteFile(file.Path, data, mode)
}

func (file *File) renderFile() error {
	return nil
}

func (file *File) Requirements() []string {
	return file.Require
}

func (file *File) Meta() Metadata {
	return file.Metadata
}

func (file *File) State() *Result {
	result := &Result{
		Metadata:   &file.Metadata,
		Consistent: false,
	}
	switch {
	case file.Metadata.State == "absent":
		if _, err := os.Stat(file.Path); err == nil {
			result.Message = fmt.Sprintf("File: %s exists", file.Path)
			return result
		}
	case file.Metadata.State == "linked":
		f, err := os.Stat(file.Path)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		if f.Mode() != os.ModeSymlink {
			result.Message = fmt.Sprintf("File %s is not a symlink", file.Path)
			return result
		}
	case file.Metadata.State == "rendered":
		_, err := os.Stat(file.Path) // TODO: Compare contents of rendered remote file to local file.
		if err != nil {
			result.Message = err.Error()
			return result
		}
	}
	result.Consistent = true
	result.Message = "Success"
	return result
}

func (file *File) Apply() *Result {
	result := file.State()
	switch {
	case result.Consistent == true: // File is in the correct state
	case file.Metadata.State == "absent":
		err := os.Remove(file.Path)
		if err != nil {
			result.Message = err.Error()
			return result
		}
	case file.Metadata.State == "linked":
		err := os.Symlink(file.Source, file.Path)
		if err != nil {
			result.Message = err.Error()
			return result
		}
	case file.Metadata.State == "rendered":
		err := file.writeFile()
		if err != nil {
			result.Message = err.Error()
			return result
		}
	}
	result.Message = "Success"
	return result
}
