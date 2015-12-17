package state

import (
	"fmt"
	"strings"
	"encoding/json"
)


func StateFactory(name, keyword string, data []byte) (State, error) {
	split := strings.Split(keyword, ".")
	metadata := Metadata{name, split[0], split[1]}
	switch split[0] {
	case "file":
		file := &File{Metadata: metadata}
		err := json.Unmarshal(data, &file)
		if err != nil {
			return file, err
		}
		return file, nil
	case "package":
		pkg := &Package{Metadata: metadata}
		err := json.Unmarshal(data, &pkg)
		if err != nil {
			return pkg, err
		}
		return pkg, nil
	case "service":
		service := &Service{Metadata: metadata}
		err := json.Unmarshal(data, &service)
		if err != nil {
			return service, err
		}
		return service, nil
	default:
		panic(fmt.Errorf("Unknown state keyword: %s", split[0]))
	}
}