package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aneokin12/vouch/internal/p2p"
	"github.com/spf13/cobra"
)

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Host a secure P2P session to share a namespace",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		fmt.Printf("Starting Vouch Host for namespace: %s...\n", namespace)

		// 1. Create LibP2P Host
		h, err := p2p.NewNode(ctx)
		if err != nil {
			fmt.Println("Error starting P2P host:", err)
			os.Exit(1)
		}
		defer h.Close()

		fmt.Printf("Listening on local network with PeerID: %s\n", h.ID().String())

		// 2. Start mDNS Discovery
		peerChan, err := p2p.StartMDNSDiscovery(ctx, h)
		if err != nil {
			fmt.Println("Error starting mDNS discovery:", err)
			os.Exit(1)
		}

		fmt.Println("Waiting for peers to join...")

		// 3. Listen for discovered peers
		for {
			select {
			case pi := <-peerChan:
				fmt.Printf("Found local peer: %s\n", pi.ID.String())
				// Next Step: Handshake Negotiation
			case <-ctx.Done():
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(inviteCmd)
}
