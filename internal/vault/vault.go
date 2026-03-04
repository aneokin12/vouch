package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
)

// Secret represents a single key-value entry with CRDT LWW metadata
type Secret struct {
	Value     string `json:"value"`
	UpdatedAt int64  `json:"updatedAt"`
	Deleted   bool   `json:"deleted"`
}

// Vault represents the stored map of Secrets
type Vault map[string]Secret

var (
	ErrIncorrectPassword = errors.New("incorrect password or corrupted vault")
	ErrVaultNotFound     = errors.New("vault not found")
)

const (
	saltSize  = 16
	nonceSize = 12
	keySize   = 32
)

// deriveKey generates a 32-byte AES key using Argon2id
func deriveKey(password []byte, salt []byte) []byte {
	return argon2.IDKey(password, salt, 1, 64*1024, 4, keySize)
}

// Encrypt payload and save to disk
func SaveVault(v Vault, password string, path string) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	key := deriveKey([]byte(password), salt)

	block, err := aes.NewCipher(key)
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

	// Final format: [SALT (16 bytes)] [NONCE (12 bytes)] [CIPHERTEXT]
	finalPayload := append(salt, nonce...)
	finalPayload = append(finalPayload, ciphertext...)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	return os.WriteFile(path, finalPayload, 0600)
}

// Decrypt vault from disk
func LoadVault(password string, path string) (Vault, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrVaultNotFound
		}
		return nil, err
	}

	if len(data) < saltSize+nonceSize {
		return nil, ErrIncorrectPassword
	}

	salt := data[:saltSize]
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	key := deriveKey([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrIncorrectPassword
	}

	var v Vault
	if err := json.Unmarshal(plaintext, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// MergeVaults deterministically merges a remote vault into a local vault using CRDT LWW rules
func MergeVaults(local, remote Vault) Vault {
	merged := make(Vault)

	// Start with everything in local
	for k, v := range local {
		merged[k] = v
	}

	// Apply remote updates
	for k, rv := range remote {
		if lv, exists := merged[k]; exists {
			// CRDT LWW logic
			if rv.UpdatedAt > lv.UpdatedAt {
				merged[k] = rv
			} else if rv.UpdatedAt == lv.UpdatedAt {
				// Deterministic tie-breaker (e.g. string comparison value)
				if rv.Value > lv.Value {
					merged[k] = rv
				}
			}
		} else {
			// New key from remote
			merged[k] = rv
		}
	}

	return merged
}
