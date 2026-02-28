package cmd

import (
	"fmt"
	"masterr/pkg/aur"
	"masterr/pkg/pacman"
	"masterr/pkg/search"
	"masterr/pkg/ui"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <package>",
	Short: "Search and install a package from pacman or AUR",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.TrimSpace(args[0])
		fmt.Printf("Searching for \"%s\"...\n", query)

		results, err := search.Search(query)
		if err != nil || len(results) == 0 {
			fmt.Println("No packages found.")
			return nil
		}

		var pkgs []pacman.Package
		for _, r := range results {
			pkgs = append(pkgs, r.Package)
		}

		p := tea.NewProgram(ui.NewResultsModel(query, pkgs), tea.WithMouseCellMotion())
		result, err := p.Run()
		if err != nil {
			return err
		}

		final, ok := result.(ui.ResultsModel)
		if !ok || final.Selected == nil {
			return nil
		}

		pkg := *final.Selected
		fmt.Printf("Installing %s from %s...\n", pkg.Name, pkg.Source)

		if pkg.Source == "aur" {
			c := aur.Install(pkg.Name)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			return c.Run()
		}
		c := pacman.Install(pkg.Name)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		return c.Run()
	},
}
