package p2p

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aneokin12/vouch/internal/vault"
	"github.com/libp2p/go-libp2p/core/network"
)

// TransferVault encrypts the Vault payload with the Session Key and writes it to the stream
func TransferVault(stream network.Stream, sessionKey []byte, v vault.Vault) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(sessionKey)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)

	// Prepend nonce
	payload := append(nonce, ciphertext...)

	return writeMsg(stream, payload)
}

// ReceiveVault reads the encrypted Vault payload from the stream and decrypts it with the Session Key
func ReceiveVault(stream network.Stream, sessionKey []byte) (vault.Vault, error) {
	payload, err := readMsg(stream)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(sessionKey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(payload) < nonceSize {
		return nil, fmt.Errorf("payload too short")
	}

	nonce := payload[:nonceSize]
	ciphertext := payload[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	var v vault.Vault
	if err := json.Unmarshal(plaintext, &v); err != nil {
		return nil, err
	}

	return v, nil
}
