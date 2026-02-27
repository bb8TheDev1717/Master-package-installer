package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Action struct {
	Name string
	Desc string
}

var Actions = []Action{
	{"install", "Install a package"},
	{"search", "Search for a package"},
	{"remove", "Remove a package"},
	{"update", "Update all packages"},
	{"list", "List & browse installed packages"},
}

type State int

const (
	StateMenu State = iota
	StateInput
)

type Model struct {
	cursor int
	state  State
	action Action
	input  string
	Result string
	Pkg    string
	width  int
	height int
}

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cba6f7")).Padding(0, 4)
	menuSelected = lipgloss.NewStyle().Foreground(lipgloss.Color("#1e1e2e")).Background(lipgloss.Color("#cba6f7")).Padding(0, 2)
	menuNormal   = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")).Padding(0, 2)
	menuDesc     = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
	inputStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#fab387")).Bold(true)
	hintStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#45475a"))
	boxStyle     = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#cba6f7")).
			Padding(1, 4)
)

func NewModel() Model {
	return Model{state: StateMenu}
}

func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch m.state {
		case StateMenu:
			switch msg.String() {
			case "esc", "ctrl+c":
				m.Result = "quit"
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(Actions)-1 {
					m.cursor++
				}
			case "enter":
				m.action = Actions[m.cursor]
				if m.action.Name == "update" || m.action.Name == "list" {
					m.Result = m.action.Name
					return m, tea.Quit
				}
				m.state = StateInput
			}
		case StateInput:
			switch msg.String() {
			case "ctrl+c":
				m.Result = "quit"
				return m, tea.Quit
			case "esc":
				m.state = StateMenu
				m.input = ""
			case "enter":
				if m.input != "" {
					m.Pkg = m.input
					m.Result = m.action.Name
					return m, tea.Quit
				}
			case "backspace":
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			default:
				if len(msg.String()) == 1 {
					m.input += msg.String()
				}
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var content strings.Builder

	logo := lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7")).Bold(true).Render(
		"███╗   ███╗ █████╗ ███████╗████████╗███████╗██████╗\n" +
			"████╗ ████║██╔══██╗██╔════╝╚══██╔══╝██╔════╝██╔══██╗\n" +
			"██╔████╔██║███████║███████╗   ██║   █████╗  ██████╔╝\n" +
			"██║╚██╔╝██║██╔══██║╚════██║   ██║   ██╔══╝  ██╔══██╗\n" +
			"██║ ╚═╝ ██║██║  ██║███████║   ██║   ███████╗██║  ██║\n" +
			"╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝",
	)
	content.WriteString(logo + "\n\n")

	switch m.state {
	case StateMenu:
		for i, a := range Actions {
			name := fmt.Sprintf("%-10s", a.Name)
			desc := menuDesc.Render(a.Desc)
			if i == m.cursor {
				content.WriteString(menuSelected.Render("▸ "+name) + "   " + desc + "\n\n")
			} else {
				content.WriteString(menuNormal.Render("  "+name) + "   " + desc + "\n\n")
			}
		}
		content.WriteString("\n" + hintStyle.Render("↑/↓  navigate   enter  select   esc  quit"))

	case StateInput:
		content.WriteString(menuDesc.Render("Action: ") + inputStyle.Render(m.action.Name) + "\n\n")
		content.WriteString(inputStyle.Render("▸ ") + m.input + "█\n\n")
		content.WriteString(hintStyle.Render("enter  confirm   esc  back"))
	}

	box := boxStyle.Render(content.String())

	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}
	return box
}
