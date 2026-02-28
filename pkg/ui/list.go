package ui

import (
	"fmt"
	"masterr/pkg/pacman"
	"masterr/pkg/util"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ListModel is the TUI model for browsing installed packages.
// Typing filters the list with fuzzy matching; Enter opens package info.
type ListModel struct {
	all      []pacman.Package
	filtered []pacman.Package
	filter   string
	cursor   int
	Info     *pacman.Package // set when user presses Enter, triggers info view
}

func NewListModel(pkgs []pacman.Package) ListModel {
	return ListModel{all: pkgs, filtered: pkgs}
}

func (m ListModel) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.MouseButtonWheelDown:
			if m.cursor < len(m.filtered)-1 {
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
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.filtered) > 0 {
				pkg := m.filtered[m.cursor]
				m.Info = &pkg
				return m, tea.Quit
			}
		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.cursor = 0
				m.applyFilter()
			}
		default:
			if len(msg.String()) == 1 {
				m.filter += msg.String()
				m.cursor = 0
				m.applyFilter()
			}
		}
	}
	return m, nil
}

// applyFilter updates the filtered list using exact → contains → Levenshtein ranking.
func (m *ListModel) applyFilter() {
	if m.filter == "" {
		m.filtered = m.all
		return
	}
	q := strings.ToLower(m.filter)

	type scored struct {
		pkg   pacman.Package
		score int
	}
	var items []scored
	for _, p := range m.all {
		name := strings.ToLower(p.Name)
		var s int
		switch {
		case name == q:
			s = 1000
		case strings.HasPrefix(name, q):
			s = 500
		case strings.Contains(name, q):
			s = 200
		default:
			dist := util.Levenshtein(q, name)
			if dist > 4 {
				continue
			}
			s = 200 - dist
		}
		items = append(items, scored{p, s})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].score > items[j].score
	})

	m.filtered = make([]pacman.Package, len(items))
	for i, s := range items {
		m.filtered[i] = s.pkg
	}
}

func (m ListModel) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("  master  ") + "\n")
	sb.WriteString(inputStyle.Render("Filter: ") + m.filter + "█\n\n")
	sb.WriteString(Separator() + "\n")

	limit := 25
	start := 0
	if m.cursor >= limit {
		start = m.cursor - limit + 1
	}
	end := start + limit
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := start; i < end; i++ {
		p := m.filtered[i]
		line := fmt.Sprintf("  %-40s %s", p.Name, p.Version)
		if i == m.cursor {
			sb.WriteString(menuSelected.Render("▸"+line) + "\n")
		} else {
			sb.WriteString(menuNormal.Render(line) + "\n")
		}
	}

	sb.WriteString(Separator() + "\n")
	sb.WriteString(menuDesc.Render(fmt.Sprintf("%d packages", len(m.filtered))) + "\n\n")
	sb.WriteString(hintStyle.Render("type to filter   ↑/↓ navigate   enter  info   esc  back"))

	return sb.String()
}
