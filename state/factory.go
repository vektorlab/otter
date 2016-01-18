package state

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

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
