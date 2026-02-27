package ui

import (
	"fmt"
	"masterr/pkg/pacman"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	repoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa")) // Catppuccin blue
	aurStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")) // Catppuccin green
	nameStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cdd6f4"))
	versionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#fab387"))
	descStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#1e1e2e")).Background(lipgloss.Color("#cba6f7"))
)

func RenderPackage(p pacman.Package, index int, selected bool) string {
	src := repoStyle.Render("[pacman]")
	if p.Source == "aur" {
		src = aurStyle.Render("[aur]")
	}
	line := fmt.Sprintf("%2d. %s %s %s  %s",
		index+1,
		src,
		nameStyle.Render(p.Name),
		versionStyle.Render(p.Version),
		descStyle.Render(p.Desc),
	)
	if selected {
		return selectedStyle.Render(line)
	}
	return line
}

func PrintPackages(pkgs []pacman.Package) {
	for i, p := range pkgs {
		fmt.Println(RenderPackage(p, i, false))
	}
}

func Separator() string {
	return descStyle.Render(strings.Repeat("─", 60))
}
