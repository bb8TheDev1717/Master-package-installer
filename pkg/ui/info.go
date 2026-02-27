package ui

import (
	"masterr/pkg/pacman"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type InfoModel struct {
	pkg pacman.Package
}

func NewInfoModel(pkg pacman.Package) InfoModel {
	return InfoModel{pkg: pkg}
}

func (m InfoModel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m InfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m InfoModel) View() string {
	var sb strings.Builder

	p := m.pkg
	sb.WriteString(titleStyle.Render("  master  ") + "\n\n")

	if p.Name == "" {
		sb.WriteString(menuDesc.Render("Package not found.") + "\n")
	} else {
		src := repoStyle.Render("[pacman]")
		if p.Source == "aur" {
			src = aurStyle.Render("[aur]")
		}
		sb.WriteString(Separator() + "\n")
		sb.WriteString(nameStyle.Render(p.Name) + "  " + src + "\n\n")
		sb.WriteString(inputStyle.Render("Version: ") + p.Version + "\n")
		sb.WriteString(inputStyle.Render("Desc:    ") + descStyle.Render(p.Desc) + "\n")
		sb.WriteString(Separator() + "\n")
	}

	sb.WriteString("\n" + hintStyle.Render("esc  back to menu"))
	return sb.String()
}
