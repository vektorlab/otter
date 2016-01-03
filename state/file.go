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
	"os"
	"strings"
)

type SourceType string

type FileType string

type File struct {
	Mode       int64      `json:"mode"`   // File should be set to this octal mode
	Path       string     `json:"path"`   // File destination
	Source     string     `json:"source"` // File source
	SourceType SourceType `json:"-"`      // SourceType e.g. git, filesystem, s3
	Metadata   Metadata   `json:"metadata"`
	Require    []string   `json:"require"`
}

/*
Determine the "source type" of a file's content based on it's path prefix
*/
func (file *File) SetSourceType(field string) (SourceType, error) {
	switch {
	case strings.Contains(field, "git:///"):
		return SourceType("git"), nil
	case strings.Contains(field, "file:///"):
		return SourceType("filesystem"), nil
	case strings.Contains(field, "s3:///"):
		return SourceType("s3"), nil
	default:
		return SourceType(""), fmt.Errorf("Unable to parse source type: %s", field)
	}
}

func (file *File) Initialize() error {
	var err error
	file.SourceType, err = file.SetSourceType(file.Source)
	if err != nil {
		return err
	}
	state := file.Metadata.State
	if state != "rendered" {
		return fmt.Errorf("Invalid file state: %s", state)
	}
	if file.Path == "" {
		file.Path = file.Metadata.Name
	}
	return err
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
		Metadata: &file.Metadata,
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
	return file.Consistent()
}
