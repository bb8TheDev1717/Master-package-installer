package cmd

import (
	"fmt"
	"os"

	"masterr/pkg/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "master",
	Short: "Package manager wrapper for pacman + AUR",
	Long:  "master – unified CLI for pacman and paru (AUR)",
	// Running master without subcommands opens the interactive TUI menu.
	// After each action completes, the menu is shown again until the user quits.
	RunE: func(cmd *cobra.Command, args []string) error {
		for {
			m := ui.NewModel()
			p := tea.NewProgram(m)
			result, err := p.Run()
			if err != nil {
				return err
			}

			final, ok := result.(ui.Model)
			if !ok || final.Result == "quit" || final.Result == "" {
				return nil
			}

			// Build subcommand args from TUI selection
			subArgs := []string{final.Result}
			if final.Pkg != "" {
				subArgs = append(subArgs, final.Pkg)
			}
			cmd.Root().SetArgs(subArgs)
			cmd.Root().Execute() //nolint:errcheck — errors printed to stderr by cobra
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(listCmd)
}
