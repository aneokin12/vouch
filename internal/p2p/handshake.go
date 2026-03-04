package p2p

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/schollz/pake/v3"
)

var words = []string{
	"apple", "tree", "ghost", "river", "mountain", "cloud", "fire", "stone", "wind", "shadow",
	"ocean", "eagle", "tiger", "moon", "star", "sun", "sand", "snow", "ice", "blade",
}

// GenerateMagicCode creates a memorable 3-word passphrase
func GenerateMagicCode() (string, error) {
	var code string
	for i := 0; i < 3; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
		if err != nil {
			return "", err
		}
		if i > 0 {
			code += "-"
		}
		code += words[n.Int64()]
	}
	return code, nil
}

// NewHostPake initializes the SPAKE2 handshake state for the Host (Role: 0)
func NewHostPake(magicCode string) (*pake.Pake, error) {
	return pake.InitCurve([]byte(magicCode), 0, "ed25519") // 0 for host
}

// NewClientPake initializes the SPAKE2 handshake state for the Client (Role: 1)
func NewClientPake(magicCode string) (*pake.Pake, error) {
	return pake.InitCurve([]byte(magicCode), 1, "ed25519") // 1 for client
}

// RunHandshake exchanges the PAKE payloads over the stream and derives the symmetric session key
func RunHandshake(stream network.Stream, pk *pake.Pake) ([]byte, error) {
	// 1. Send our public material
	myBytes := pk.Bytes()
	if err := writeMsg(stream, myBytes); err != nil {
		return nil, fmt.Errorf("failed to send PAKE payload: %w", err)
	}

	// 2. Read peer's public material
	peerBytes, err := readMsg(stream)
	if err != nil {
		return nil, fmt.Errorf("failed to read peer PAKE payload: %w", err)
	}

	// 3. Update the PAKE state machine
	if err := pk.Update(peerBytes); err != nil {
		return nil, fmt.Errorf("failed PAKE update (possible wrong magic code): %w", err)
	}

	// 4. Extract the shared session key
	sessionKey, err := pk.SessionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to extract session key: %w", err)
	}

	return sessionKey, nil
}

// writeMsg writes a length-prefixed message to the stream
func writeMsg(w io.Writer, msg []byte) error {
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(msg)))
	if _, err := w.Write(lenBuf); err != nil {
		return err
	}
	_, err := w.Write(msg)
	return err
}

// readMsg reads a length-prefixed message from the stream
func readMsg(r io.Reader) ([]byte, error) {
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lenBuf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf)
	if length > 1024*1024*10 { // 10MB sanity check
		return nil, fmt.Errorf("message too large: %d bytes", length)
	}

	msg := make([]byte, length)
	if _, err := io.ReadFull(r, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
