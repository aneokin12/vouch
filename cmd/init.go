package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aneokin12/vouch/internal/runner"
	"github.com/aneokin12/vouch/internal/vault"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [command to run]",
	Short: "Execute a command with secrets injected into its environment",
	Args:  cobra.MinimumNArgs(1),
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
		vaultPath := filepath.Join(home, ".vouch", namespace+".enc")

		v, err := vault.LoadVault(password, vaultPath)
		if err != nil {
			fmt.Println("Error loading vault:", err)
			os.Exit(1)
		}

		err = runner.InjectAndRun(v, args)
		if err != nil {
			fmt.Println("Error running command:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
