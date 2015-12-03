package helpers

import (
	"fmt"
	"gopkg.in/ini.v1"
	"runtime"
)

func GetDistribution() (string, error) {
	if runtime.GOOS == "linux" {
		i, err := ini.Load([]byte(""), "/etc/os-release")
		section, err := i.Section("").GetKey("ID_LIKE")
		if err != nil {
			return "", err
		}
		switch section.String() {
		case "debian":
			return "debian", nil
		case "rhel fedora":
			return "centos", nil
		default:
			return "", fmt.Errorf("Distribution %s is not supported by Otter at this time", section.String())
		}
		return section.String(), nil
	}
	return "", fmt.Errorf("%s is not supported by Otter at this time", runtime.GOOS)
}
