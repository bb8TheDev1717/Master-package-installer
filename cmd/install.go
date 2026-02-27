package cmd

import (
	"bufio"
	"fmt"
	"masterr/pkg/aur"
	"masterr/pkg/pacman"
	"masterr/pkg/search"
	"masterr/pkg/ui"
	"os"
	"strconv"
	"strings"

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

		// Show up to 20 results
		limit := len(results)
		if limit > 20 {
			limit = 20
		}
		fmt.Println(ui.Separator())
		for i := 0; i < limit; i++ {
			fmt.Println(ui.RenderPackage(results[i].Package, i, false))
		}
		fmt.Println(ui.Separator())

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Select number to install (or q to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" || input == "" {
			return nil
		}

		n, err := strconv.Atoi(input)
		if err != nil || n < 1 || n > limit {
			fmt.Println("Invalid selection.")
			return nil
		}

		pkg := results[n-1].Package
		fmt.Printf("Installing %s from %s...\n", pkg.Name, pkg.Source)

		// Route to correct backend based on package source
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
