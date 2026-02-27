package aur

import (
	"masterr/pkg/pacman"
	"os/exec"
	"strings"
)

func Search(query string) ([]pacman.Package, error) {
	out, err := exec.Command("paru", "-Ss", "--aur", query).Output()
	if err != nil {
		return nil, nil
	}
	return parseSearchOutput(string(out)), nil
}

func Info(name string) (pacman.Package, error) {
	out, err := exec.Command("paru", "-Si", "--aur", name).Output()
	if err != nil {
		return pacman.Package{}, err
	}
	return parseInfoOutput(string(out)), nil
}

func Install(name string) *exec.Cmd {
	return exec.Command("paru", "-S", "--noconfirm", name)
}

func Remove(name string) *exec.Cmd {
	return exec.Command("sudo", "pacman", "-Rns", "--noconfirm", name)
}

func Update() *exec.Cmd {
	return exec.Command("paru", "-Syu", "--aur", "--noconfirm")
}

func parseSearchOutput(output string) []pacman.Package {
	var pkgs []pacman.Package
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
		pkgs = append(pkgs, pacman.Package{Name: name, Version: version, Source: "aur", Desc: desc})
	}
	return pkgs
}

func parseInfoOutput(output string) pacman.Package {
	var pkg pacman.Package
	pkg.Source = "aur"
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
