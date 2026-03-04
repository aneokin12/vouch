package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aneokin12/vouch/internal/p2p"
	"github.com/aneokin12/vouch/internal/vault"
	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use:   "join [magic-code]",
	Short: "Join a P2P session to receive a namespace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		magicCode := args[0]
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

		fmt.Printf("Starting Vouch Client for namespace: '%s'...\n", namespace)

		// 2. Create LibP2P Host
		h, err := p2p.NewNode(ctx)
		if err != nil {
			fmt.Println("Error starting P2P host:", err)
			os.Exit(1)
		}
		defer h.Close()

		fmt.Printf("Searching local network for Host using Magic Code: %s\n", magicCode)

		// 2. Start mDNS Discovery
		peerChan, err := p2p.StartMDNSDiscovery(ctx, h)
		if err != nil {
			fmt.Println("Error starting mDNS discovery:", err)
			os.Exit(1)
		}

		// 3. Listen for discovered peers
		for {
			select {
			case pi := <-peerChan:
				fmt.Printf("Found local Host peer: %s\n", pi.ID.String())

				// Next Step: Actually connect and perform SPAKE2 handshake using the magicCode
				err = h.Connect(ctx, pi)
				if err != nil {
					fmt.Printf("Failed to connect to host: %v\n", err)
					continue
				}

				// Open a stream to the vault sync protocol
				s, err := h.NewStream(ctx, pi.ID, protocolID)
				if err != nil {
					fmt.Printf("Failed to open stream to host: %v\n", err)
					continue
				}

				fmt.Println("Stream established! Starting SPAKE2 Handshake...")

				// Initialize Client PAKE State
				pk, err := p2p.NewClientPake(magicCode)
				if err != nil {
					fmt.Println("Handshake initialization failed:", err)
					s.Reset()
					continue
				}

				// Run Handshake over stream
				sessionKey, err := p2p.RunHandshake(s, pk)
				if err != nil {
					fmt.Println("Handshake failed:", err)
					s.Reset()
					continue
				}

				fmt.Println("Handshake successful! Exchanging Vaults...")

				// Receive Remote Vault first (Host transfers first)
				remoteVault, err := p2p.ReceiveVault(s, sessionKey)
				if err != nil {
					fmt.Println("Failed to receive remote vault payload:", err)
					s.Reset()
					continue
				}

				// Transfer Local Vault
				if err := p2p.TransferVault(s, sessionKey, localVault); err != nil {
					fmt.Println("Failed to send vault payload:", err)
					s.Reset()
					continue
				}

				fmt.Println("Vault exchanged successfully. Merging CRDTs...")
				// Merge Vaults deterministically
				mergedVault := vault.MergeVaults(localVault, remoteVault)

				// Save to disk
				if err := vault.SaveVault(mergedVault, password, vaultPath); err != nil {
					fmt.Println("Error saving merged vault to disk:", err)
					continue
				}

				fmt.Printf("Sync complete! %d secrets combined.\n", len(mergedVault))

				s.Close()
				return

			case <-ctx.Done():
				return
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
