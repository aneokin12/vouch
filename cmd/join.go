package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aneokin12/vouch/internal/p2p"
	"github.com/aneokin12/vouch/internal/vault"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
)

var peerAddress string

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

		// If a manual peer address is provided, connect directly and bypass mDNS
		if peerAddress != "" {
			fmt.Printf("Attempting direct connection to peer: %s\n", peerAddress)

			maddr, err := multiaddr.NewMultiaddr(peerAddress)
			if err != nil {
				fmt.Println("Invalid --peer multiaddress:", err)
				os.Exit(1)
			}

			pi, err := peer.AddrInfoFromP2pAddr(maddr)
			if err != nil {
				fmt.Println("Failed to parse peer address info:", err)
				os.Exit(1)
			}

			if err := h.Connect(ctx, *pi); err != nil {
				fmt.Println("Failed to connect directly to peer:", err)
				os.Exit(1)
			}
			fmt.Printf("Connected directly to Host: %s\n", pi.ID.String())

			// Run the sync payload
			err = syncWithHost(ctx, h, *pi, magicCode, localVault, password, vaultPath)
			if err != nil {
				fmt.Println("Sync session failed:", err)
			}
			return
		}

		fmt.Printf("Searching local network for Host using Magic Code: %s\n", magicCode)

		// 3. Start mDNS Discovery
		peerChan, err := p2p.StartMDNSDiscovery(ctx, h, magicCode)
		if err != nil {
			fmt.Println("Error starting mDNS discovery:", err)
			os.Exit(1)
		}

		// 4. Listen for discovered peers
		for {
			select {
			case pi := <-peerChan:
				fmt.Printf("Found local Host peer: %s\n", pi.ID.String())

				err = h.Connect(ctx, pi)
				if err != nil {
					fmt.Printf("Failed to connect to host: %v\n", err)
					continue
				}

				err = syncWithHost(ctx, h, pi, magicCode, localVault, password, vaultPath)
				if err != nil {
					fmt.Println("Sync session failed:", err)
				}
				return // End after sync attempt

			case <-ctx.Done():
				return
			}
		}

	},
}

// syncWithHost handles the PAKE exchange and CRDT vault merging with a connected Host
func syncWithHost(ctx context.Context, h host.Host, pi peer.AddrInfo, magicCode string, localVault vault.Vault, password string, vaultPath string) error {
	// Open a stream to the vault sync protocol
	s, err := h.NewStream(ctx, pi.ID, protocolID)
	if err != nil {
		return fmt.Errorf("failed to open stream to host: %w", err)
	}
	defer s.Close()

	fmt.Println("Stream established! Starting SPAKE2 Handshake...")

	// Initialize Client PAKE State
	pk, err := p2p.NewClientPake(magicCode)
	if err != nil {
		s.Reset()
		return fmt.Errorf("handshake initialization failed: %w", err)
	}

	// Run Handshake over stream
	sessionKey, err := p2p.RunHandshakeClient(s, pk)
	if err != nil {
		s.Reset()
		return fmt.Errorf("handshake failed: %w", err)
	}

	fmt.Println("Handshake successful! Exchanging Vaults...")

	// Receive Remote Vault first (Host transfers first)
	remoteVault, err := p2p.ReceiveVault(s, sessionKey)
	if err != nil {
		s.Reset()
		return fmt.Errorf("failed to receive remote vault payload: %w", err)
	}

	// Transfer Local Vault
	if err := p2p.TransferVault(s, sessionKey, localVault); err != nil {
		s.Reset()
		return fmt.Errorf("failed to send vault payload: %w", err)
	}

	fmt.Println("Vault exchanged successfully. Merging CRDTs...")
	// Merge Vaults deterministically
	mergedVault := vault.MergeVaults(localVault, remoteVault)

	// Save to disk
	if err := vault.SaveVault(mergedVault, password, vaultPath); err != nil {
		return fmt.Errorf("error saving merged vault to disk: %w", err)
	}

	fmt.Printf("Sync complete! %d secrets combined.\n", len(mergedVault))
	return nil
}

func init() {
	joinCmd.Flags().StringVar(&peerAddress, "peer", "", "Directly connect to a peer's Multiaddress, bypassing mDNS discovery")
	rootCmd.AddCommand(joinCmd)
}
