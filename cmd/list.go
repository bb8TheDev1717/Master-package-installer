package cmd

import (
	"masterr/pkg/aur"
	"masterr/pkg/pacman"
	"masterr/pkg/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Browse installed packages; press Enter on a package to view info",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgs, err := pacman.List()
		if err != nil {
			return err
		}

		for {
			m := ui.NewListModel(pkgs)
			p := tea.NewProgram(m, tea.WithMouseCellMotion())
			result, err := p.Run()
			if err != nil {
				return err
			}

			final, ok := result.(ui.ListModel)
			if !ok || final.Info == nil {
				return nil
			}

			// Fetch full package info: try pacman first, fall back to AUR
			pkg := *final.Info
			full, err := pacman.Info(pkg.Name)
			if err != nil || full.Name == "" {
				full, _ = aur.Info(pkg.Name)
			}
			if full.Name == "" {
				full = pkg // use what we already have if both fail
			}

			p2 := tea.NewProgram(ui.NewInfoModel(full), tea.WithMouseCellMotion())
			if _, err := p2.Run(); err != nil {
				return err
			}
			// Loop back to list after closing info
		}
	},
}
