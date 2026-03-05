package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/aneokin12/vouch/internal/keychain"
	"golang.org/x/term"
)

// getPassword retrieves the vault password. It first checks the macOS Keychain,
// and if not found, prompts the user interactively and saves it to the Keychain.
func getPassword() string {
	// 1. Try the Keychain first
	pw, err := keychain.GetPassword()
	if err == nil && pw != "" {
		return pw
	}

	// 2. Not in Keychain — prompt interactively
	fmt.Print("Enter vault password: ")
	raw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // newline after hidden input
	if err != nil {
		fmt.Println("Error reading password:", err)
		os.Exit(1)
	}

	pw = string(raw)
	if pw == "" {
		fmt.Println("Error: password cannot be empty")
		os.Exit(1)
	}

	// 3. Store in Keychain for next time
	if err := keychain.SetPassword(pw); err != nil {
		fmt.Println("Warning: could not save password to Keychain:", err)
	}

	return pw
}
