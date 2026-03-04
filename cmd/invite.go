package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Invite a peer to sync secrets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("invite called")
	},
}

func init() {
	rootCmd.AddCommand(inviteCmd)
}
