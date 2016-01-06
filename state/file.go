/*
A File represents a single file on an operating system.
States -
  rendered: A file is copied from a another source // TODO
  deleted: A file is removed from the operating system // TODO
  linked: A symbolic link is created from one destination to another // TODO

Source Types -
  git: A file is copied from a Git repository // TODO
  s3: A file is copied from an Amazon S3 repo // TODO
  file: A file is copied from the local file system // TODO
*/

package state

import (
	"encoding/json"
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

func (file *File) Initialize() error {
	state := file.Metadata.State
	if state != "rendered" {
		return fmt.Errorf("Invalid file state: %s", state)
	}
	if file.Path == "" {
		file.Path = file.Metadata.Name
	}
	return nil
}

func (file *File) Dump() ([]byte, error) {
	return json.Marshal(file)
}

func (file *File) Requirements() []string {
	return file.Require
}

func (file *File) Meta() Metadata {
	return file.Metadata
}

func (file *File) Consistent() *Result {
	result := &Result{
		Metadata:   &file.Metadata,
		Consistent: false,
	}
	f, err := os.Stat(file.Path)
	if f == nil {
		result.Message = err.Error()
		return result
	}
	result.Consistent = f != nil
	return result
}

func (file *File) Execute() *Result {
	err := file.writeFile()
	result := &Result{
		Metadata: &file.Metadata,
	}
	if err != nil {
		result.Consistent = false
		result.Message = err.Error()
	} else {
		result.Consistent = true
	}
	return result
}
