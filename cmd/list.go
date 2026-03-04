package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aneokin12/vouch/internal/tui"
	"github.com/aneokin12/vouch/internal/vault"
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
		vaultPath := filepath.Join(home, ".vouch", "vault.enc")

		v, err := vault.LoadVault(password, vaultPath)
		if err != nil {
			if err == vault.ErrVaultNotFound {
				fmt.Println("Vault is empty. Use `vouch set` to add secrets.")
				os.Exit(0)
			}
			fmt.Println("Error loading vault:", err)
			os.Exit(1)
		}

		p := tea.NewProgram(tui.NewListModel(v))
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running TUI:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
