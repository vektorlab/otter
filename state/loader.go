/*

The state package provides translation for otter's declarative YAML configuration file
into State objects which are used to configure a host server.

An example YAML configuration might look like this:

/root/readme.md: <-- Name
  file.rendered: <-- {type}.{state}
    mode: 644 <-- Key/Value arguments
    source: git@github.com/vektorlab/otter/README.md
*/

package state

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type State interface {
	Initialize() error
	Dump() ([]byte, error)
	Requirements() []string
	Meta() Metadata
}

type Metadata struct {
	Name  string // Unique name to associate with a state
	Type  string // The type of state "package", "file", etc.
	State string // The desired state "installed", "rendered", etc.
}

type Loader struct {
	stateRaw map[string]map[string]interface{}
	State    map[string][]State
}

/*
Translate each section in the otter.yml file to a State object
*/
func (loader *Loader) sectionToState(name, keyword string, data interface{}) error {
	split := strings.Split(keyword, ".")
	metadata := Metadata{name, split[0], split[1]}
	switch split[0] {
	case "file":
		var err error
		file := File{
			Metadata: metadata,
		}
		err = mapstructure.Decode(data, &file)
		if err != nil {
			return err
		}
		err = file.Initialize()
		if err != nil {
			return err
		}
		loader.State[name] = append(loader.State[name], &file)
		//loader.state[name] = &file
		return nil
	case "package":
		var err error
		pkg := Package{
			Metadata: metadata,
		}
		err = mapstructure.Decode(data, &pkg)
		if err != nil {
			return err
		}
		err = pkg.Initialize()
		if err != nil {
			return err
		}
		loader.State[name] = append(loader.State[name], &pkg)
		return nil
	case "service":
		var err error
		service := Service{
			Metadata: metadata,
		}
		err = mapstructure.Decode(data, &service)
		if err != nil {
			return err
		}
		err = service.Initialize()
		if err != nil {
			return err
		}
		loader.State[name] = append(loader.State[name], &service)
		return nil
	default:
		return fmt.Errorf("Unknown keyword %s", keyword)
	}
}

/*
Partially Unmarshal the YAML document and then call "sectionToState" to evaluate
the keywords and decode the rest with Mapstructure.
*/
func (loader *Loader) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	err = unmarshal(&loader.stateRaw)

	if err != nil {
		return err
	}

	for name, value := range loader.stateRaw {
		for keyword, data := range value {
			loader.sectionToState(name, keyword, data)
		}
	}
	return err
}

/*
Dump the entire state out to JSON
*/
func (loader *Loader) Dump() ([]byte, error) {
	return json.Marshal(loader.State)
}

func (loader *Loader) Count() int {
	var count int
	for _, value := range loader.State {
		count += len(value)
	}
	return count
}

func (loader *Loader) requirements(entry []State) []string {
	reqs := make([]string, 0)
	for s := range entry {
		entryReqs := entry[s].Requirements()
		for r := range entryReqs {
			reqs = append(reqs, entryReqs[r])
		}
	}
	return reqs
}

/*
Validate that all requirements in each state exist and that there are no circular requirements
*/
func (loader *Loader) validate() error {
	for name, entry := range loader.State {
		reqs := loader.requirements(entry)
		for req := range reqs {
			other, exists := loader.State[reqs[req]]
			if !exists {
				return fmt.Errorf("Unable to find requirement: %s", reqs[req])
			}
			otherReqs := loader.requirements(other)
			for req := range otherReqs {
				if name == otherReqs[req] {
					return fmt.Errorf("Detected circular requirement: %s", name)
				}
			}
		}
	}
	return nil
}

func FromBytes(data []byte) (*Loader, error) {

	var err error

	loader := Loader{
		State: make(map[string][]State),
	}

	err = yaml.Unmarshal(data, &loader)

	if err != nil {
		return nil, err
	}

	err = loader.validate()

	if err != nil {
		return nil, err
	}

	return &loader, nil
}

func FromPath(path string) (*Loader, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FromBytes(data)
}
