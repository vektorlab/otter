package state

import (
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"fmt"
)

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
	default: return nil, fmt.Errorf("Invalid file state: %s", metadata.State)
	}
	if file.Path == "" {
		file.Path = metadata.Name
	}
	return file, nil

}

/*
Create and validate a new Package State
 */
func newPackage(metadata Metadata, data []byte) (*Package, error) {
	pkg := &Package{}
	err := json.Unmarshal(data, &pkg)
	if err != nil {
		return nil, err
	}
	pkg.Metadata = metadata
	switch metadata.State {
	case "installed":
	case "removed":
	default:
		return nil, fmt.Errorf("Invalid package state: %s", metadata.State)
	}
	if pkg.Name == "" {
		pkg.Name = metadata.Name
	}
	return pkg, nil
}

/*
Create and validate a new Service State
 */
func newService(metadata Metadata, data []byte) (*Service, error) {
	service := &Service{}
	err := json.Unmarshal(data, &service)
	if err != nil {
		return nil, err
	}
	service.Metadata = metadata
	switch metadata.State {
	case "running":
	case "stopped":
	default:
		return nil, fmt.Errorf("Invalid service state: %s", metadata.State)
	}
	if service.Name == "" {
		service.Name = metadata.Name
	}
	return service, nil
}

func StateFactory(metadata Metadata, data []byte) (State, error) {
	log.Printf("Loading state %s: %s.%s", metadata.Name, metadata.Type, metadata.State)
	err := json.Unmarshal(data, &metadata) // Load requirements from each state into Metadata
	if err != nil {
		panic(err)
	}
	switch metadata.Type {
	case "file":
		return newFile(metadata, data)
	case "package":
		return newPackage(metadata, data)
	case "service":
		return newService(metadata, data)
	default:
		panic(fmt.Errorf("Unknown state keyword: %s", metadata.Type))
	}
}
