package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aneokin12/vouch/internal/vault"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5FD2")).MarginBottom(1)
	keyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF"))
	valStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
)

type ListModel struct {
	vault vault.Vault
	keys  []string
}

func NewListModel(v vault.Vault) *ListModel {
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return &ListModel{
		vault: v,
		keys:  keys,
	}
}

func (m *ListModel) Init() tea.Cmd {
	return nil
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
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

	b.WriteString(titleStyle.Render("🤝 Vouch Secrets Vault"))
	b.WriteString("\n")

	if len(m.keys) == 0 {
		b.WriteString("No secrets found in the vault.\n")
	} else {
		for _, k := range m.keys {
			secret := m.vault[k]

			// Do not print tombstones
			if !secret.Deleted {
				masked := maskValue(secret.Value)
				b.WriteString(fmt.Sprintf("%s: %s\n", keyStyle.Render(k), valStyle.Render(masked)))
			}
		}
	}

	b.WriteString("\nPress 'q' or 'esc' to quit.\n")
	return b.String()
}
