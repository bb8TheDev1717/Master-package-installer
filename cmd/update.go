package cmd

import (
	"fmt"
	"masterr/pkg/aur"
	"masterr/pkg/pacman"
	"os"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all packages (pacman + AUR)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Updating pacman packages...")
		c := pacman.Update()
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		if err := c.Run(); err != nil {
			return err
		}

		fmt.Println("Updating AUR packages...")
		c2 := aur.Update()
		c2.Stdout = os.Stdout
		c2.Stderr = os.Stderr
		c2.Stdin = os.Stdin
		return c2.Run()
	},
}
