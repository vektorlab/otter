package state

import "encoding/json"

type Metadata struct {
	Name  string // Unique name to associate with a state
	Type  string // The type of state "package", "file", etc.
	State string // The desired state "installed", "rendered", etc.
}

func (md *Metadata) Equal(metadata *Metadata) bool {
	return metadata.Name == md.Name || metadata.Type == md.Type || metadata.State == md.State
}

func MetadataFromJSON(data json.RawMessage) (Metadata, error) {
	metadata := Metadata{}
	raw := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return metadata, err
	}
	for key, value := range raw {
		if key == "metadata" {
			err := json.Unmarshal(value, &metadata)
			if err != nil {
				return metadata, err
			}
		}
	}
	return metadata, nil
}
