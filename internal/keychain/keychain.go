package keychain

import (
	"errors"
	"os/exec"
	"strings"
)

const serviceName = "vouch-vault"
const accountName = "vouch"

var ErrNotFound = errors.New("password not found in keychain")

// GetPassword retrieves the vault password from the macOS Keychain.
func GetPassword() (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-s", serviceName, "-a", accountName, "-w")
	out, err := cmd.Output()
	if err != nil {
		return "", ErrNotFound
	}
	return strings.TrimSpace(string(out)), nil
}

// SetPassword stores or updates the vault password in the macOS Keychain.
func SetPassword(password string) error {
	cmd := exec.Command("security", "add-generic-password", "-s", serviceName, "-a", accountName, "-w", password, "-U")
	return cmd.Run()
}
