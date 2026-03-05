# Vouch

```
⠀⣷⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⠀
⠀⣿⣿⣷⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⣿⣿⠀
⠀⣿⣿⣿⣿⣷⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⣿⣿⣿⣿⠀
⠀⣿⣿⣿⣿⡟⢀⣾⣿⣷⡖⠒⣀⣀⣤⣶⣾⣿⣷⣶⣤⣤⣴⣾⣆⠘⣿⣿⣿⠀    ⠰⣶⡔⠀⠀⠀⣶⠐⠀⣠⡴⠒⠒⢦⣄⠀⠐⢲⡶⠂⠀⠀⢲⠂⠀⢀⣤⠖⠒⠒⢤⡄⠂⣶⡖⠀⠀⠐⣶⡖⠀
⠀⣿⣿⣿⡿⠁⣾⣿⣿⣿⣧⣀⠙⠛⠛⠛⠋⣈⠻⢿⣿⣿⣿⣿⣿⣧⠈⢿⣿⠀    ⠀⢻⣿⠀⠀⢠⠃⠀⢰⡿⠀⠀⠀⠀⢻⣇⠀⢸⣇⠀⠀⠀⢸⠀⢠⣿⡏⠀⠀⠀⠈⠃⠀⣿⡇⠀⠀⠀⣿⡇⠀
⠀⣿⣿⣿⠁⣼⣿⠟⠻⠿⣿⣿⣿⣷⣶⣾⣿⠿⣷⣄⡉⠻⣿⣿⣿⣿⣧⠈⢿⠀    ⠀⠀⢻⣧⠀⡞⠀⠀⢸⣇⠀⠀⠀⠀⢸⡿⠀⢸⣇⠀⠀⠀⢸⠀⢘⣿⡄⠀⠀⠀⠀⠀⠀⣿⡍⠙⠉⠉⣿⡇⠀
⠀⠛⠛⠃⠀⠻⠋⣠⣶⠄⠙⠛⢻⣿⣿⣿⣷⣦⣈⠛⢿⣦⣄⠙⠻⠛⠁⠀⠀⠀    ⠀⠀⠀⢿⡾⠀⠀⠀⠈⢿⣄⠀⠀⣠⡿⠃⠀⠸⣷⡀⠀⢀⡼⠀⠀⠻⣷⣀⠀⢀⣠⠀⠀⣿⡇⠀⠀⠀⣿⡇⠀
⠀⠀⠀⠀⠀⢠⣾⠟⠁⣠⣿⠇⠈⠛⣿⣿⣈⠙⠿⣷⣦⡈⠛⠛⠀⠀⠀⠀⠀⠀    ⠀⠀⠀⠈⠁⠀⠀⠀⠀⠀⠈⠉⠉⠁⠀⠀⠀⠀⠀⠉⠉⠉⠀⠀⠀⠀⠀⠉⠉⠉⠀⠀⠉⠉⠉⠀⠀⠈⠉⠉⠁
⠀⠀⠀⠀⠀⠀⠁⢠⣾⡿⠁⣠⣾⠆⠸⠉⠻⣷⣤⡈⠙⠿⠃⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠉⠀⣾⠟⢁⣴⡶⠀⣤⡀⠙⠟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠰⣿⠟⠁⠀⠛⠟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
```

**A zero-trust, peer-to-peer secret orchestrator.**

Vouch safely injects secrets into your processes and securely syncs them with your peers — no `.env` files, no centralized managers, no plaintext secrets over Slack.

---

## Table of Contents

- [Why Vouch?](#why-vouch)
- [Architecture Overview](#architecture-overview)
- [Requirements](#requirements)
- [Installation](#installation)
  - [Building from Source](#building-from-source)
- [Getting Started](#getting-started)
- [Commands](#commands)
  - [`vouch set`](#vouch-set)
  - [`vouch list`](#vouch-list)
  - [`vouch rm`](#vouch-rm)
  - [`vouch init`](#vouch-init)
  - [`vouch invite`](#vouch-invite)
  - [`vouch join`](#vouch-join)
- [Namespaces](#namespaces)
- [Security Model](#security-model)
  - [Encryption At Rest](#encryption-at-rest)
  - [P2P Sync Security](#p2p-sync-security)
  - [CRDT Merge Algorithm](#crdt-merge-algorithm)
- [License](#license)

---

## Why Vouch?

Developers share secrets over Slack, commit `.env` files to repos, and juggle a patchwork of secret managers. Vouch eliminates this entirely:

- **No `.env` files.** Secrets never touch the filesystem in plaintext.
- **No centralized service.** Your secrets live in encrypted vaults on *your machine*.
- **No password sharing.** P2P sync uses SPAKE2 key exchange — master passwords are never transmitted.
- **In-memory injection.** `vouch init` uses `syscall.Exec` to inject secrets directly into a process's environment, replacing the Vouch process entirely from memory.

---

## Architecture Overview

Vouch is composed of three core workflows:

```
┌─────────────────────────────────────────────────────────────────────┐
│   YOUR MACHINE                                                      │
│                                                                     │
│   vouch set / list / rm ──▶ ~/.vouch/<namespace>.enc                │
│                                (AES-256-GCM + Argon2id)             │
│                                                                     │
│   vouch init ──▶ Decrypt vault ──▶ syscall.Exec (inject into env)   │
│                                                                     │
│   vouch invite ──▶ libp2p node ──▶ mDNS broadcast                  │
│                       │              ▼                               │
│                       │     SPAKE2 Handshake ◀── vouch join          │
│                       │              ▼                               │
│                       └──▶ Encrypted CRDT Vault Sync                │
└─────────────────────────────────────────────────────────────────────┘
```

**1. Local Vault Management** — `set`, `list`, `rm` manage secrets within namespaced, encrypted vaults stored at `~/.vouch/`.

**2. Process Injection** — `init` decrypts a single namespace and replaces the current process via `syscall.Exec`, passing secrets directly as environment variables. Secrets exist only in memory.

**3. P2P Encrypted Sync** — `invite` and `join` use [libp2p](https://libp2p.io/) and [SPAKE2](https://en.wikipedia.org/wiki/SPAKE2) to securely sync vaults between peers on the same local network, without ever transmitting master passwords.

---

## Requirements

| Requirement | Version |
|---|---|
| **Go** | `1.24.6` or later |
| **Operating System** | macOS, Linux (Windows support for `syscall.Exec` is limited) |
| **Network** | Local network (LAN) for P2P sync via mDNS |

### Environment Variable

Vouch requires the `VOUCH_PASSWORD` environment variable to be set before running any command. This is the master password used to derive the AES encryption key for your vaults.

```bash
export VOUCH_PASSWORD="your-secure-password"
```

On macOS, Vouch can also store and retrieve your password from the system Keychain automatically after the first interactive prompt.

---

## Installation

### Building from Source

Clone the repository and build with `go build`:

```bash
# Clone the repository
git clone https://github.com/aneokin12/vouch.git
cd vouch

# Download dependencies
go mod download

# Build the binary
go build -o vouch .

# (Optional) Move to a directory on your PATH
sudo mv vouch /usr/local/bin/
```

To verify the installation:

```bash
vouch --help
```

You can also build and run in one step during development:

```bash
go run main.go --help
```

---

## Getting Started

```bash
# 1. Set your vault password
export VOUCH_PASSWORD="my-secret-password"

# 2. Store your first secret
vouch set STRIPE_KEY sk_live_abc123

# 3. View your secrets in the TUI
vouch list

# 4. Run your app with secrets injected
vouch init -- npm start

# 5. Share secrets with a teammate over P2P
vouch invite --env=shared        # On your machine (generates a magic code)
vouch join orange-water-ghost    # On their machine (uses the magic code)
```

---

## Commands

### `vouch set`

Store a secret in an encrypted vault.

```
vouch set [key] [value]
```

| Flag | Short | Default | Description |
|---|---|---|---|
| `--env` | `-e` | `personal` | Namespace (environment) to store the secret in |

**Examples:**

```bash
# Store in the default "personal" namespace
vouch set DATABASE_URL "postgres://user:pass@localhost:5432/mydb"

# Store in a shared namespace
vouch set API_KEY "sk_live_123" --env=shared-prod
```

If the namespace vault doesn't exist yet, it will be created automatically.

---

### `vouch list`

List stored secrets securely using an interactive TUI dashboard.

```
vouch list
```

Opens a terminal UI that displays all available namespaces (vault files) in `~/.vouch/`. Select a namespace to decrypt and view the key-value pairs inside. Values are displayed securely within the TUI and never written to stdout in plaintext.

---

### `vouch rm`

Remove a secret or an entire namespace from the vault.

```
vouch rm [key]
vouch rm --all
```

| Flag | Short | Default | Description |
|---|---|---|---|
| `--env` | `-e` | `personal` | Namespace to remove secrets from |
| `--all` | `-a` | `false` | Delete the entire namespace vault file |

**Examples:**

```bash
# Soft-delete a single key (tombstoned for CRDT sync)
vouch rm OLD_API_KEY --env=shared-prod

# Permanently delete an entire namespace vault
vouch rm --all --env=staging
```

> **Note:** Single-key deletions are implemented as CRDT tombstones (`deleted: true` with a timestamp). This ensures deletions propagate correctly during P2P sync.

---

### `vouch init`

Execute a command with secrets injected into its environment.

```
vouch init [command to run]
```

| Flag | Short | Default | Description |
|---|---|---|---|
| `--env` | `-e` | `personal` | Namespace to load secrets from |

**Examples:**

```bash
# Run a Node.js app with personal secrets
vouch init -- npm start

# Run a Go application with shared production secrets
vouch init --env=shared-prod -- go run main.go

# Run any command
vouch init --env=dev -- python manage.py runserver
```

**How it works:** Vouch decrypts the specified namespace vault, builds an environment array with your secrets, and uses `syscall.Exec` to completely replace the Vouch process in memory with your target command. The target process inherits the injected environment variables. Secrets only ever exist in-memory — they are never written to disk in plaintext.

> **Important:** Only one namespace can be injected at a time. This follows the principle of least privilege and prevents key collisions.

---

### `vouch invite`

Host a secure P2P session to share a vault namespace with a peer.

```
vouch invite
```

| Flag | Short | Default | Description |
|---|---|---|---|
| `--env` | `-e` | `personal` | Namespace to share |

**Examples:**

```bash
# Share the shared-prod vault with a teammate
vouch invite --env=shared-prod
```

**What happens:**
1. Vouch generates a short **Magic Code** (e.g., `orange-water-ghost`).
2. A libp2p node starts and broadcasts via mDNS on the local network.
3. The command prints direct connection fallback addresses for networks where mDNS is blocked.
4. When a peer connects, a SPAKE2 handshake is performed using the Magic Code.
5. Vaults are exchanged, encrypted with the ephemeral session key.
6. Both vaults are merged using the CRDT LWW algorithm.

---

### `vouch join`

Join a P2P session to receive and sync a vault namespace.

```
vouch join [magic-code]
```

| Flag | Short | Default | Description |
|---|---|---|---|
| `--env` | `-e` | `personal` | Namespace to sync into |
| `--peer` | | `""` | Directly connect to a peer's Multiaddress, bypassing mDNS |

**Examples:**

```bash
# Join via mDNS discovery
vouch join orange-water-ghost --env=shared-prod

# Direct connection (when mDNS is unavailable)
vouch join orange-water-ghost --peer /ip4/192.168.1.5/tcp/9000/p2p/QmPeer...
```

**What happens:**
1. Vouch searches the local network via mDNS (or connects directly if `--peer` is provided).
2. A SPAKE2 handshake is performed — if the Magic Codes don't match, the connection drops silently.
3. The remote vault is received, decrypted with the negotiated session key.
4. Both local and remote vaults are merged using the CRDT LWW algorithm.
5. The merged result is re-encrypted with *your local password* and saved to disk.

---

## Namespaces

Namespaces are discrete, encrypted vault files that represent a single environment or project. They are stored at `~/.vouch/<namespace>.enc`.

```
~/.vouch/
├── personal.enc        # Default namespace
├── shared-prod.enc     # Shared production secrets
├── shared-dev.enc      # Shared development secrets
└── staging.enc         # Staging environment
```

Every command accepts `--env` (or `-e`) to specify which namespace to operate on. The default namespace is `personal`.

---

## Security Model

### Encryption At Rest

All vaults are encrypted using **AES-256-GCM** with keys derived via **Argon2id** (1 iteration, 64 MB memory, 4 threads).

The on-disk format is:

```
[SALT (16 bytes)] [NONCE (12 bytes)] [CIPHERTEXT...]
```

- A fresh random salt and nonce are generated for every save operation.
- Vault files are written with `0600` permissions (owner read/write only).
- Passwords are completely local and never shared or transmitted.

### P2P Sync Security

The P2P sync protocol is designed for **zero-trust** operation:

1. **SPAKE2 Handshake** — A Password-Authenticated Key Exchange (PAKE) using the human-readable Magic Code. This yields an ephemeral session key without ever transmitting the code itself. If codes don't match, the connection drops silently (MITM protection).
2. **Double Encryption** — The host decrypts their vault with their *local password*, re-encrypts the payload with the *ephemeral session key*, and transmits it. The receiver decrypts with the session key, merges, and re-encrypts with *their own password*.
3. **mDNS Discovery** — Peers discover each other on the local network. Direct connection via `--peer` multiaddress is also supported as a fallback.

### CRDT Merge Algorithm

Vouch uses a **Last-Write-Wins (LWW) Map** backed by tombstones for deterministic, conflict-free merging:

```json
{
  "STRIPE_KEY": {
    "value": "sk_live_123",
    "updatedAt": 1698765432,
    "deleted": false
  },
  "OLD_KEY": {
    "value": "",
    "updatedAt": 1698765100,
    "deleted": true
  }
}
```

**Merge rules:**
1. If a key exists in only one vault, it is included in the merged result.
2. If a key exists in both, the entry with the highest `updatedAt` timestamp wins.
3. Equal timestamps use a deterministic tie-breaker (lexicographic value comparison).
4. Tombstoned keys (`deleted: true`) with a newer timestamp propagate deletions across peers.

---

## Project Structure

```
vouch/
├── main.go                     # Entrypoint
├── cmd/                        # CLI command definitions (Cobra)
│   ├── root.go                 # Root command & global flags
│   ├── init.go                 # vouch init
│   ├── set.go                  # vouch set
│   ├── list.go                 # vouch list
│   ├── rm.go                   # vouch rm
│   ├── invite.go               # vouch invite (P2P host)
│   ├── join.go                 # vouch join  (P2P client)
│   └── password.go             # Keychain-backed password helper
├── internal/
│   ├── vault/vault.go          # Encryption, decryption, CRDT merge
│   ├── runner/runner.go        # syscall.Exec process injection
│   ├── keychain/               # macOS Keychain integration
│   ├── p2p/                    # libp2p networking
│   │   ├── host.go             # Node creation & mDNS discovery
│   │   ├── handshake.go        # SPAKE2 key exchange
│   │   └── transfer.go         # Encrypted vault transfer protocol
│   └── tui/                    # Terminal UI (Bubble Tea)
└── go.mod
```