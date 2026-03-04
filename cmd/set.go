package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aneokin12/vouch/internal/vault"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Store a secret in the vault",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

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

		// Load or create a new vault
		v, err := vault.LoadVault(password, vaultPath)
		if err != nil {
			if err == vault.ErrVaultNotFound {
				v = make(vault.Vault)
			} else {
				fmt.Println("Error loading vault:", err)
				os.Exit(1)
			}
		}

		// Store as a Secret CRDT including the updated timestamp
		v[key] = vault.Secret{
			Value:     value,
			UpdatedAt: time.Now().Unix(),
			Deleted:   false,
		}

		if err := vault.SaveVault(v, password, vaultPath); err != nil {
			fmt.Println("Error saving vault:", err)
			os.Exit(1)
		}

		fmt.Printf("Secret %s stored securely.\n", key)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
