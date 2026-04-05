package tui

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

// st is set once in Run() before the program starts.
var st styles

type themeColors struct {
	Accent    string `yaml:"accent"`
	User      string `yaml:"user"`
	Assistant string `yaml:"assistant"`
	System    string `yaml:"system"`
	HeaderBg  string `yaml:"header_bg"`
	HeaderFg  string `yaml:"header_fg"`
}

type themeConfig struct {
	Colors themeColors `yaml:"colors"`
}

func defaultTheme() themeConfig {
	return themeConfig{
		Colors: themeColors{
			Accent:    "#A78BFA",
			User:      "#C4B5FD",
			Assistant: "#7C3AED",
			System:    "#9CA3AF",
			HeaderBg:  "#3B0764",
			HeaderFg:  "#EDE9FE",
		},
	}
}

func loadTheme(emmDir string) themeConfig {
	cfg := defaultTheme()
	data, err := os.ReadFile(filepath.Join(emmDir, "tui", "theme.yaml"))
	if err != nil {
		return cfg
	}
	_ = yaml.Unmarshal(data, &cfg)
	return cfg
}

type styles struct {
	user      lipgloss.Style
	assistant lipgloss.Style
	system    lipgloss.Style
	header    lipgloss.Style
	border    lipgloss.Style
}

func buildStyles(t themeConfig) styles {
	c := t.Colors
	return styles{
		user:      lipgloss.NewStyle().Foreground(lipgloss.Color(c.User)).Bold(true),
		assistant: lipgloss.NewStyle().Foreground(lipgloss.Color(c.Assistant)).Bold(true),
		system:    lipgloss.NewStyle().Foreground(lipgloss.Color(c.System)).Italic(true),
		header:    lipgloss.NewStyle().Background(lipgloss.Color(c.HeaderBg)).Foreground(lipgloss.Color(c.HeaderFg)).Bold(true),
		border:    lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(c.Accent)),
	}
}
