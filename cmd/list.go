package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aneokin12/vouch/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List stored secrets securely using a TUI dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		password := os.Getenv("VOUCH_PASSWORD")
		if password == "" {
			fmt.Println("Error: VOUCH_PASSWORD environment variable is not set")
			os.Exit(1)
		}

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			os.Exit(1)
		}
		vouchDir := filepath.Join(home, ".vouch")
		entries, err := os.ReadDir(vouchDir)
		if err != nil && !os.IsNotExist(err) {
			fmt.Println("Error reading vouch directory:", err)
			os.Exit(1)
		}

		var namespaces []string
		for _, e := range entries {
			if !e.IsDir() && filepath.Ext(e.Name()) == ".enc" {
				namespaces = append(namespaces, strings.TrimSuffix(e.Name(), ".enc"))
			}
		}

		if len(namespaces) == 0 {
			fmt.Println("No Vaults found. Use `vouch set` to create your first secret.")
			os.Exit(0)
		}

		p := tea.NewProgram(tui.NewListModel(namespaces, password, vouchDir))
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running TUI:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
