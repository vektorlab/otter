package state

import (
	"encoding/json"
	"fmt"
)

func StateFactory(metadata Metadata, data []byte) (State, error) {
	switch metadata.Type {
	case "file":
		file := &File{Metadata: metadata}
		err := json.Unmarshal(data, &file)
		if err != nil {
			return file, err
		}
		switch metadata.State {
		case "absent":
		case "linked":
		case "rendered":
		default:
			return file, fmt.Errorf("Invalid file state: %s", metadata.State)
		}
		if file.Path == "" {
			file.Path = metadata.Name
		}
		return file, nil
	case "package":
		pkg := &Package{Metadata: metadata}
		err := json.Unmarshal(data, &pkg)
		if err != nil {
			return pkg, err
		}
		switch metadata.State {
		case "installed":
		case "removed":
		default:
			return pkg, fmt.Errorf("Invalid package state: %s", metadata.State)
		}
		if pkg.Name == "" {
			pkg.Name = metadata.Name
		}
		return pkg, nil
	case "service":
		service := &Service{Metadata: metadata}
		err := json.Unmarshal(data, &service)
		if err != nil {
			return service, err
		}
		switch metadata.State {
		case "running":
		case "stopped":
		default:
			return service, fmt.Errorf("Invalid service state: %s", metadata.State)
		}
		if service.Name == "" {
			service.Name = metadata.Name
		}
		return service, nil
	default:
		panic(fmt.Errorf("Unknown state keyword: %s", metadata.Type))
	}
}
