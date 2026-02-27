package cmd

import (
	"fmt"
	"masterr/pkg/pacman"
	"masterr/pkg/search"
	"masterr/pkg/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <package>",
	Short: "Search for a package in pacman and AUR",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		results, err := search.Search(query)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Search failed:", err)
			return err
		}

		var pkgs []pacman.Package
		for _, r := range results {
			pkgs = append(pkgs, r.Package)
		}

		p := tea.NewProgram(ui.NewResultsModel(query, pkgs))
		_, err = p.Run()
		return err
	},
}
