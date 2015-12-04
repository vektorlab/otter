package executors

import (
	"encoding/json"
	"github.com/vektorlab/otter/state"
	"os"
)

type File struct {
	state *state.File
}

func (file *File) Consistent() (bool, error) {
	var err error
	if err != nil {
		return false, err
	}

	f, err := os.Stat(file.state.Path)

	if f == nil {
		return false, err
	}

	return f != nil, nil
}

func (file *File) Execute() (Result, error) {
	return Result{}, nil
}

func FileFromJson(data []byte) (*File, error) {
	file := File{
		state: &state.File{},
	}
	err := json.Unmarshal(data, &file.state)
	return &file, err
}

func FileFromState(state *state.File) (*File, error) {
	file := File{
		state: state,
	}
	return &file, nil
}
