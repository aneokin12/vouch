package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aneokin12/vouch/internal/p2p"
	"github.com/aneokin12/vouch/internal/vault"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/spf13/cobra"
)

const protocolID = "/vouch/sync/1.0.0"

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Host a secure P2P session to share a namespace",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

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

		// 1. Load Local Vault
		localVault, err := vault.LoadVault(password, vaultPath)
		if err != nil {
			if err == vault.ErrVaultNotFound {
				localVault = make(vault.Vault)
			} else {
				fmt.Println("Error loading local vault:", err)
				os.Exit(1)
			}
		}

		// 2. Generate Magic Code
		magicCode, err := p2p.GenerateMagicCode()
		if err != nil {
			fmt.Println("Error generating magic code:", err)
			os.Exit(1)
		}

		fmt.Printf("Starting Vouch Host for namespace: '%s'...\n\n", namespace)
		fmt.Printf("Magic Code: %s\n\n", magicCode)
		fmt.Println("Tell the receiver to run: `vouch join " + magicCode + "`")

		// 3. Create LibP2P Host
		h, err := p2p.NewNode(ctx)
		if err != nil {
			fmt.Println("Error starting P2P host:", err)
			os.Exit(1)
		}
		defer h.Close()

		// 4. Set up the Stream Handler for incoming sync requests
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
			sessionKey, err := p2p.RunHandshakeHost(s, pk)
			if err != nil {
				fmt.Println("Handshake failed:", err)
				s.Reset()
				return
			}

			fmt.Println("Handshake successful! Exchanging Vaults...")

			// Transfer Local Vault first
			if err := p2p.TransferVault(s, sessionKey, localVault); err != nil {
				fmt.Println("Failed to send vault payload:", err)
				s.Reset()
				return
			}

			// Receive Remote Vault
			remoteVault, err := p2p.ReceiveVault(s, sessionKey)
			if err != nil {
				fmt.Println("Failed to receive remote vault payload:", err)
				s.Reset()
				return
			}

			fmt.Println("Vault exchanged successfully. Merging CRDTs...")
			// Merge Vaults deterministically
			mergedVault := vault.MergeVaults(localVault, remoteVault)

			// Save to disk
			if err := vault.SaveVault(mergedVault, password, vaultPath); err != nil {
				fmt.Println("Error saving merged vault to disk:", err)
				return
			}

			fmt.Printf("Sync complete! %d secrets combined. You can now close the connection (Ctrl+C).\n", len(mergedVault))
		})

		// 5. Start mDNS Discovery
		peerChan, err := p2p.StartMDNSDiscovery(ctx, h, magicCode)
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
