package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aneokin12/vouch/internal/vault"
	"github.com/spf13/cobra"
)

var rmAll bool

var rmCmd = &cobra.Command{
	Use:   "rm [key]",
	Short: "Remove a secret or an entire namespace from the vault",
	Args: func(cmd *cobra.Command, args []string) error {
		if rmAll {
			if len(args) != 0 {
				return fmt.Errorf("accepts no arguments when --all is used")
			}
			return nil
		}
		if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
		}
		return nil
	},
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

		if rmAll {
			err := os.Remove(vaultPath)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("Namespace '%s' does not exist.\n", namespace)
					os.Exit(0)
				}
				fmt.Println("Error deleting namespace:", err)
				os.Exit(1)
			}
			fmt.Printf("Namespace '%s' successfully deleted.\n", namespace)
			return
		}

		key := args[0]

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
	rmCmd.Flags().BoolVarP(&rmAll, "all", "a", false, "Delete the entire namespace vault file")
	rootCmd.AddCommand(rmCmd)
}
