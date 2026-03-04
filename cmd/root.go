package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vouch",
	Short: "Vouch is a zero-trust P2P secret orchestrator",
	Long:  `Vouch safely injects secrets into your processes and securely syncs them with your peers.`,
}

var namespace string

func init() {
	rootCmd.PersistentFlags().StringVarP(&namespace, "env", "e", "personal", "Namespace (environment) to load secrets from")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
