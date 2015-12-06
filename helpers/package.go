package helpers

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetDpkgPackage(name string) (string, error) {
	out, err := exec.Command("dpkg", "-l", name).Output() // TODO: Security
	if err != nil { // dpkg doesn't know about the package
		return "rc", nil
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		split := strings.Fields(line)
		if len(split) >= 1 {
			f := string(split[0][0])
			if f == "h" || f == "i" || f == "p" || f == "r" || f == "u" {
				return split[0], nil
			}
		}
	}
	return "", nil
}

func GetPackageStatus(name string) (string, error) {
	distro, err := GetDistro()
	if err != nil {
		return "", err
	}

	switch distro.Family {
	case "debian":
		status, err := GetDpkgPackage(name)
		if err != nil {
			return "", err
		}
		switch status {
		case "ii":
			return "installed", nil
		case "rc":
			return "removed", nil
		case "":
			return "purged", nil
		default:
			return "", fmt.Errorf("Package %s is in an unknown state: %s", name, status)
		}
	default:
		return "", fmt.Errorf("Unsupported operating system: %s", distro.Family)
	}
}

func installAptPackage(name string) error {
	out, err := exec.Command("apt-get", "update").CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		return err
	}
	fmt.Println(string(out))
	out, err = exec.Command("apt-get", "install", "-y", name).CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func InstallPackage(name string) error {
	distro, err := GetDistro()
	if err != nil {
		return err
	}
	switch distro.Family {
	case "debian":
		err := installAptPackage(name)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported operating system: %s", distro.Family)
	}
	return nil
}