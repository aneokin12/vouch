package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/aneokin12/vouch/internal/p2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/spf13/cobra"
)

const protocolID = "/vouch/sync/1.0.0"

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Host a secure P2P session to share a namespace",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// 1. Generate Magic Code
		magicCode, err := p2p.GenerateMagicCode()
		if err != nil {
			fmt.Println("Error generating magic code:", err)
			os.Exit(1)
		}

		fmt.Printf("Starting Vouch Host for namespace: '%s'...\n\n", namespace)
		fmt.Printf("Magic Code: %s\n\n", magicCode)
		fmt.Println("Tell the receiver to run: `vouch join " + magicCode + "`")

		// 2. Create LibP2P Host
		h, err := p2p.NewNode(ctx)
		if err != nil {
			fmt.Println("Error starting P2P host:", err)
			os.Exit(1)
		}
		defer h.Close()

		// 3. Set up the Stream Handler for incoming sync requests
		h.SetStreamHandler(protocolID, func(s network.Stream) {
			defer s.Close()
			fmt.Println("\nIncoming connection! Starting SPAKE2 Handshake...")

			// Initialize Host PAKE State
			pk, err := p2p.NewHostPake(magicCode)
			if err != nil {
				fmt.Println("Handshake initialization failed:", err)
				return
			}

			// Run Handshake over stream
			sessionKey, err := p2p.RunHandshake(s, pk)
			if err != nil {
				fmt.Println("Handshake failed:", err)
				s.Reset()
				return
			}

			fmt.Printf("Handshake successful! Derived Session Key: %s\n", hex.EncodeToString(sessionKey))
			// TODO: Encrypt Vault with Session Key and Send over stream
			fmt.Println("Closing stream (Data transfer not yet implemented)")
		})

		// 4. Start mDNS Discovery
		peerChan, err := p2p.StartMDNSDiscovery(ctx, h)
		if err != nil {
			fmt.Println("Error starting mDNS discovery:", err)
			os.Exit(1)
		}

		fmt.Println("\nWaiting for peers to discover us on the local network...")

		// 5. Listen for discovered peers
		for {
			select {
			case pi := <-peerChan:
				fmt.Printf("Discovered peer on network: %s (Waiting for connection...)\n", pi.ID.String())
			case <-ctx.Done():
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(inviteCmd)
}
