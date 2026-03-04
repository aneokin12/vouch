package tui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aneokin12/vouch/internal/vault"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5FD2")).MarginBottom(1)
	subtitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#A8A8A8")).MarginBottom(1)
	keyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF"))
	valStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	cursorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5FD2")).Bold(true)
	namespaceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF"))
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
)

type State int

const (
	StateNamespaces State = iota
	StateSecrets
)

type ListModel struct {
	state      State
	namespaces []string
	cursor     int
	err        error

	password string
	vaultDir string

	vault    vault.Vault
	keys     []string
	nsCursor int // cursor for secrets view
}

func NewListModel(namespaces []string, password, vaultDir string) *ListModel {
	sort.Strings(namespaces)
	return &ListModel{
		state:      StateNamespaces,
		namespaces: namespaces,
		password:   password,
		vaultDir:   vaultDir,
	}
}

func (m *ListModel) Init() tea.Cmd {
	return nil
}

func (m *ListModel) loadVaultForNamespace(ns string) error {
	vaultPath := filepath.Join(m.vaultDir, ns+".enc")
	v, err := vault.LoadVault(m.password, vaultPath)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(v))
	for k, secret := range v {
		if !secret.Deleted {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	m.vault = v
	m.keys = keys
	m.nsCursor = 0
	return nil
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.state == StateNamespaces {
				if m.cursor > 0 {
					m.cursor--
				}
			} else {
				if m.nsCursor > 0 {
					m.nsCursor--
				}
			}

		case "down", "j":
			if m.state == StateNamespaces {
				if m.cursor < len(m.namespaces)-1 {
					m.cursor++
				}
			} else {
				if m.nsCursor < len(m.keys)-1 {
					m.nsCursor++
				}
			}

		case "enter", "right", "l":
			if m.state == StateNamespaces && len(m.namespaces) > 0 {
				selectedNs := m.namespaces[m.cursor]
				if err := m.loadVaultForNamespace(selectedNs); err != nil {
					m.err = err
					return m, nil
				}
				m.err = nil
				m.state = StateSecrets
			}

		case "esc", "left", "h":
			if m.state == StateSecrets {
				m.state = StateNamespaces
				m.err = nil
			}
		}
	}
	return m, nil
}

func maskValue(val string) string {
	if len(val) <= 4 {
		return strings.Repeat("*", len(val))
	}
	if len(val) <= 8 {
		return strings.Repeat("*", len(val)-2) + val[len(val)-2:]
	}
	return strings.Repeat("*", 8) + "-" + val[len(val)-4:]
}

func (m *ListModel) View() string {
	var b strings.Builder

	if m.state == StateNamespaces {
		b.WriteString(titleStyle.Render("🤝 Vouch Vault Directory"))
		b.WriteString("\n")
		b.WriteString(subtitleStyle.Render("Select a namespace to explore:"))
		b.WriteString("\n\n")

		if m.err != nil {
			b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
			b.WriteString("\n\n")
		}

		for i, ns := range m.namespaces {
			cursor := "  "
			style := namespaceStyle
			if m.cursor == i {
				cursor = cursorStyle.Render("> ")
				style = style.Copy().Bold(true).Underline(true)
			}
			b.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(ns)))
		}

		b.WriteString("\nPress Enter/Right to select, 'q' to quit.\n")

	} else {
		selectedNs := m.namespaces[m.cursor]
		b.WriteString(titleStyle.Render(fmt.Sprintf("🔒 Inside Namespace: %s", selectedNs)))
		b.WriteString("\n\n")

		if len(m.keys) == 0 {
			b.WriteString("No active secrets found.\n")
		} else {
			for i, k := range m.keys {
				secret := m.vault[k]
				masked := maskValue(secret.Value)

				cursor := "  "
				kStyle := keyStyle
				vStyle := valStyle
				if m.nsCursor == i {
					cursor = cursorStyle.Render("> ")
					kStyle = kStyle.Copy().Bold(true)
					vStyle = vStyle.Copy().Bold(true)
				}
				b.WriteString(fmt.Sprintf("%s%s: %s\n", cursor, kStyle.Render(k), vStyle.Render(masked)))
			}
		}

		b.WriteString("\nPress Esc/Left to return, 'q' to quit.\n")
	}

	return b.String()
}
