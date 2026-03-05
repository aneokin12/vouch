package p2p

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// discoveryNotifee gets notified when we find a new peer via mDNS
type discoveryNotifee struct {
	Peerchan chan peer.AddrInfo
}

// HandlePeerFound connects to peers discovered via mDNS
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.Peerchan <- pi
}

// NewNode creates a new libp2p Host that listens on a random local port
func NewNode(ctx context.Context) (host.Host, error) {
	// listen on local ipv4 and ipv6, let the OS pick the port
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0", "/ip6/::/tcp/0"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p host: %w", err)
	}

	return h, nil
}

// StartMDNSDiscovery starts the mDNS discovery service on the Host using a unique rendezvous string based on the magic code
func StartMDNSDiscovery(ctx context.Context, h host.Host, magicCode string) (<-chan peer.AddrInfo, error) {
	// Create a channel for peer discovery notifications
	peerChan := make(chan peer.AddrInfo)
	notifee := &discoveryNotifee{Peerchan: peerChan}

	// Setting up the mDNS service to broadcast our presence and listen for others
	// We append the magicCode to ensure we only discover the specific host we want
	rendezvous := "vouch-sync-" + magicCode
	ser := mdns.NewMdnsService(h, rendezvous, notifee)

	if err := ser.Start(); err != nil {
		return nil, fmt.Errorf("failed to start mDNS service: %w", err)
	}

	return peerChan, nil
}
