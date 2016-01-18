/*
A File represents a single file on an operating system.
States -
  absent: A file is not present on the operating system // TODO
  linked: A symbolic link is created from one destination to another // TODO
  rendered: A file is copied from a another source and rendered // TODO
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

type File struct {
	Mode     string   `json:"mode"`   // File should be set to this octal mode
	Path     string   `json:"path"`   // File destination
	Source   string   `json:"source"` // File source
	Metadata Metadata `json:"metadata"`
}

func (file *File) Meta() Metadata {
	return file.Metadata
}

func (file *File) State() *Result {
	result := &Result{
		Metadata:   &file.Metadata,
		Consistent: false,
	}
	switch file.Metadata.State {
	case "absent":
		if _, err := os.Stat(file.Path); err == nil {
			result.Message = fmt.Sprintf("File: %s exists", file.Path)
			return result
		}
	case "linked":
		f, err := os.Stat(file.Path)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		if f.Mode() != os.ModeSymlink {
			result.Message = fmt.Sprintf("File %s is not a symlink", file.Path)
			return result
		}
	case "rendered":
		_, err := os.Stat(file.Path) // TODO: Compare contents of rendered remote file to local file.
		if err != nil {
			result.Message = err.Error()
			return result
		}
	}
	result.Consistent = true
	return result
}

func (file *File) Apply() *Result {
	result := file.State()
	if result.Consistent == true {
		return result
	}
	switch file.Metadata.State {
	case "absent":
		err := os.Remove(file.Path)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Message = "File removed"
		result.Consistent = true
	case "linked":
		err := os.Symlink(file.Source, file.Path)
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Message = "Symlink created"
		result.Consistent = true
	case "rendered":
		err := file.writeFile()
		if err != nil {
			result.Message = err.Error()
			return result
		}
		result.Message = "File rendered"
		result.Consistent = true
	}
	return result
}

/*
Create and validate a new File State
*/
func newFile(metadata Metadata, data []byte) (*File, error) {
	file := &File{}
	err := json.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}
	file.Metadata = metadata
	switch metadata.State {
	case "absent":
	case "linked":
	case "rendered":
	default:
		return nil, fmt.Errorf("Invalid file state: %s", metadata.State)
	}
	if file.Path == "" {
		file.Path = metadata.Name
	}
	return file, nil

}

/*
Retrieve a file from a remote source
*/
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

/*
Write a file to local disk
*/
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
