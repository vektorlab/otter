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
	Metadata() *state.Metadata // Return the Metadata of the state we are executing for
}

type Result struct {
	Consistent bool
	Metadata   *state.Metadata
	Result     error
}

type Executioner struct { // TODO: Tone down all the "execution"...?
	Executors []Executor
	Files     []*File
	Packages  []*Package
	Services  []*Service
	Results   []Result
}

func (executioner *Executioner) Run() error {
	for _, executor := range executioner.Executors {
		consistent, err := executor.Consistent() // TODO: Differentiate between results and errors
		meta := executor.Metadata()
		result := Result{
			Consistent: consistent,
			Metadata:   meta,
			Result:     err,
		}
		executioner.Results = append(executioner.Results, result)
	}

	return nil
}

func FromStateLoader(stateLoader *state.Loader) (*Executioner, error) {
	executioner := Executioner{
		Executors: make([]Executor, 0),
		Results:   make([]Result, 0),
	}
	for _, entry := range stateLoader.Files {
		executor, err := FileFromState(entry)
		if err != nil {
			return nil, err
		}
		executioner.Files = append(executioner.Files, executor)
		executioner.Executors = append(executioner.Executors, executor)
	}
	for _, entry := range stateLoader.Packages {
		executor, err := PackageFromState(entry)
		if err != nil {
			return nil, err
		}
		executioner.Packages = append(executioner.Packages, executor)
		executioner.Executors = append(executioner.Executors, executor)
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
