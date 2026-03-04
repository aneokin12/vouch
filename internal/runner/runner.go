package runner

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/aneokin12/vouch/internal/vault"
)

// InjectAndRun injects the vault secrets into the current environment and replaces the process
// with the provided command and arguments.
func InjectAndRun(v vault.Vault, cmdArgs []string) error {
	if len(cmdArgs) == 0 {
		return fmt.Errorf("no command provided to run")
	}

	// 1. Resolve path to the binary (e.g., "npm" -> "/usr/local/bin/npm")
	binary, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		return err
	}

	// 2. Prepare the environment
	env := os.Environ()
	for key, val := range v {
		env = append(env, fmt.Sprintf("%s=%s", key, val))
	}

	// 3. Syscall Exec replacement
	// Notes:
	// - args should include the command itself as args[0].
	// - This completely replaces the current process in memory.
	return syscall.Exec(binary, cmdArgs, env)
}
