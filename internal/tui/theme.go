package tui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

// st and cfg are set once in Run() before the program starts.
var st styles
var cfg themeConfig

type themeColors struct {
	Accent    string `yaml:"accent"`
	User      string `yaml:"user"`
	Assistant string `yaml:"assistant"`
	System    string `yaml:"system"`
	Tool      string `yaml:"tool"`
	HeaderBg  string `yaml:"header_bg"`
	HeaderFg  string `yaml:"header_fg"`
}

type layoutConfig struct {
	ShowHeader bool `yaml:"show_header"`
	ShowStatus bool `yaml:"show_status"`
}

type inputConfig struct {
	Placeholder string `yaml:"placeholder"`
}

type themeConfig struct {
	Colors themeColors  `yaml:"colors"`
	Layout layoutConfig `yaml:"layout"`
	Input  inputConfig  `yaml:"input"`
}

func defaultTheme() themeConfig {
	return themeConfig{
		Colors: themeColors{
			Accent:    "#7C3AED",
			User:      "#EDE9FE",
			Assistant: "#A78BFA",
			System:    "#6D28D9",
			Tool:      "#4ADE80",
			HeaderBg:  "",
			HeaderFg:  "#8B5CF6",
		},
		Layout: layoutConfig{
			ShowHeader: false,
			ShowStatus: true,
		},
		Input: inputConfig{
			Placeholder: "message  (/help)",
		},
	}
}

func loadTheme(emmDir string) (themeConfig, error) {
	cfg := defaultTheme()
	data, err := os.ReadFile(filepath.Join(emmDir, "tui", "theme.yaml"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing theme.yaml: %w", err)
	}
	return cfg, nil
}

type styles struct {
	user       lipgloss.Style
	assistant  lipgloss.Style
	system     lipgloss.Style
	toolHeader lipgloss.Style
	toolBody   lipgloss.Style
	header     lipgloss.Style
	border     lipgloss.Style
	msg        lipgloss.Style
	dim        lipgloss.Style
}

func buildStyles(t themeConfig) styles {
	c := t.Colors

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.HeaderFg)).
		Padding(0, 1)
	if c.HeaderBg != "" {
		header = header.Background(lipgloss.Color(c.HeaderBg))
	}

	return styles{
		user:       lipgloss.NewStyle().Foreground(lipgloss.Color(c.User)).Bold(true),
		assistant:  lipgloss.NewStyle().Foreground(lipgloss.Color(c.Assistant)),
		system:     lipgloss.NewStyle().Foreground(lipgloss.Color(c.System)),
		toolHeader: lipgloss.NewStyle().Foreground(lipgloss.Color(c.Tool)).Bold(true),
		toolBody:   lipgloss.NewStyle().Foreground(lipgloss.Color(c.Tool)).PaddingLeft(2),
		header:     header,
		border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(c.Accent)).
			Padding(0, 2),
		msg: lipgloss.NewStyle().PaddingLeft(2),
		dim: lipgloss.NewStyle().Foreground(lipgloss.Color(c.System)),
	}
}
