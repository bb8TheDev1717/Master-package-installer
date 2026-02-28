package ui

import (
	"fmt"
	"masterr/pkg/pacman"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ResultsModel struct {
	packages []pacman.Package
	title    string
	cursor   int
	Selected *pacman.Package
}

func NewResultsModel(title string, pkgs []pacman.Package) ResultsModel {
	return ResultsModel{title: title, packages: pkgs}
}

func (m ResultsModel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m ResultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	limit := len(m.packages)
	if limit > 20 {
		limit = 20
	}
	switch msg := msg.(type) {
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.MouseButtonWheelDown:
			if m.cursor < limit-1 {
				m.cursor++
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < limit-1 {
				m.cursor++
			}
		case "enter":
			if limit > 0 {
				pkg := m.packages[m.cursor]
				m.Selected = &pkg
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m ResultsModel) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("  master  ") + "\n")
	sb.WriteString(menuDesc.Render(fmt.Sprintf("Search: %s", m.title)) + "\n\n")

	if len(m.packages) == 0 {
		sb.WriteString(menuDesc.Render("No packages found.") + "\n")
	} else {
		sb.WriteString(Separator() + "\n")
		limit := len(m.packages)
		if limit > 20 {
			limit = 20
		}
		for i := 0; i < limit; i++ {
			sb.WriteString(RenderPackage(m.packages[i], i, i == m.cursor) + "\n")
		}
		sb.WriteString(Separator() + "\n")
	}

	sb.WriteString("\n" + hintStyle.Render("↑/↓  navigate   enter  select   esc  back"))
	return sb.String()
}
