/*
Executors accept State objects to validate that the operating system's state is consistent with what is
described in the loaded object as well as executing that state against the OS.
*/

package executors

import (
	"github.com/vektorlab/otter/state"
)

type Executor interface {
	Consistent() (bool, error) // The operating system state is consistent with the loaded state object
	Execute() (Result, error)  // Apply the loaded state object to the operating system
}

type Result struct{}

type Executioner struct { // TODO: Tone down all the "execution"...?
	Executors []Executor
	Files     []*File
	Packages  []*Package
	Services  []*Service
}

func FromStateLoader(stateLoader *state.Loader) (*Executioner, error) {
	executioner := Executioner{
		Executors: make([]Executor, stateLoader.Count()),
	}
	for _, entry := range stateLoader.Files {
		executor, err := FileFromState(entry)
		if err != nil {
			return nil, err
		}
		executioner.Files = append(executioner.Files, executor)
	}
	for _, entry := range stateLoader.Packages {
		executor, err := PackageFromState(entry)
		if err != nil {
			return nil, err
		}
		executioner.Packages = append(executioner.Packages, executor)
	}
	for _, entry := range stateLoader.Services {
		executor, err := ServiceFromState(entry)
		if err != nil {
			return nil, err
		}
		executioner.Services = append(executioner.Services, executor)
		executioner.Executors = append(executioner.Executors, executor)
	}

	return &executioner, nil
}
