package cmd

import (
	"masterr/pkg/aur"
	"masterr/pkg/pacman"
	"masterr/pkg/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <package>",
	Short: "Show package info",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		pkg, err := pacman.Info(name)
		if err != nil || pkg.Name == "" {
			pkg, _ = aur.Info(name)
		}

		p := tea.NewProgram(ui.NewInfoModel(pkg))
		_, err = p.Run()
		return err
	},
}
