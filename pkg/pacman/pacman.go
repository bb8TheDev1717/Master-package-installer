package pacman

import (
	"os/exec"
	"strings"
)

type Package struct {
	Name    string
	Version string
	Source  string // "repo" or "aur"
	Desc    string
}

func Search(query string) ([]Package, error) {
	out, err := exec.Command("pacman", "-Ss", query).Output()
	if err != nil {
		return nil, nil
	}
	return parseSearchOutput(string(out), "repo"), nil
}

func Info(name string) (Package, error) {
	out, err := exec.Command("pacman", "-Si", name).Output()
	if err != nil {
		out, err = exec.Command("pacman", "-Qi", name).Output()
		if err != nil {
			return Package{}, err
		}
	}
	return parseInfoOutput(string(out), "repo"), nil
}

func Install(name string) *exec.Cmd {
	return exec.Command("sudo", "pacman", "-S", "--noconfirm", name)
}

func Remove(name string) *exec.Cmd {
	return exec.Command("sudo", "pacman", "-Rns", "--noconfirm", name)
}

func Update() *exec.Cmd {
	return exec.Command("sudo", "pacman", "-Syu", "--noconfirm")
}

func List() ([]Package, error) {
	out, err := exec.Command("pacman", "-Q").Output()
	if err != nil {
		return nil, err
	}
	var pkgs []Package
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			pkgs = append(pkgs, Package{Name: parts[0], Version: parts[1], Source: "repo"})
		}
	}
	return pkgs, nil
}

func parseSearchOutput(output, source string) []Package {
	var pkgs []Package
	lines := strings.Split(output, "\n")
	for i := 0; i < len(lines)-1; i += 2 {
		line := lines[i]
		desc := ""
		if i+1 < len(lines) {
			desc = strings.TrimSpace(lines[i+1])
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		namePart := parts[0]
		slash := strings.Index(namePart, "/")
		name := namePart
		if slash >= 0 {
			name = namePart[slash+1:]
		}
		version := parts[1]
		pkgs = append(pkgs, Package{Name: name, Version: version, Source: source, Desc: desc})
	}
	return pkgs
}

func parseInfoOutput(output, source string) Package {
	var pkg Package
	pkg.Source = source
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "Name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				pkg.Name = strings.TrimSpace(parts[1])
			}
		}
		if strings.HasPrefix(line, "Version") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				pkg.Version = strings.TrimSpace(parts[1])
			}
		}
		if strings.HasPrefix(line, "Description") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				pkg.Desc = strings.TrimSpace(parts[1])
			}
		}
	}
	return pkg
}
