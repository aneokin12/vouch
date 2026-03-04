package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aneokin12/vouch/internal/vault"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [key]",
	Short: "Remove a secret from the vault",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

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
			fmt.Println("Error loading vault:", err)
			os.Exit(1)
		}

		if secret, exists := v[key]; exists && !secret.Deleted {
			v[key] = vault.Secret{
				Value:     "", // clear value
				UpdatedAt: time.Now().Unix(),
				Deleted:   true,
			}

			if err := vault.SaveVault(v, password, vaultPath); err != nil {
				fmt.Println("Error saving vault:", err)
				os.Exit(1)
			}
			fmt.Printf("Secret %s successfully removed.\n", key)
		} else {
			fmt.Printf("Secret %s not found in the vault.\n", key)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
