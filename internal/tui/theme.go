package tui

import (
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
	HeaderBg  string `yaml:"header_bg"`
	HeaderFg  string `yaml:"header_fg"`
}

type layoutConfig struct {
	InputHeight int  `yaml:"input_height"` // textarea height in lines
	ShowHeader  bool `yaml:"show_header"`  // show the top header bar
	ShowStatus  bool `yaml:"show_status"`  // show the scroll % / pause indicator
}

type inputConfig struct {
	Placeholder string `yaml:"placeholder"`
}

type themeConfig struct {
	Colors themeColors `yaml:"colors"`
	Layout layoutConfig `yaml:"layout"`
	Input  inputConfig  `yaml:"input"`
}

func defaultTheme() themeConfig {
	return themeConfig{
		Colors: themeColors{
			Accent:    "#475569",
			User:      "#E2E8F0",
			Assistant: "#94A3B8",
			System:    "#475569",
			HeaderBg:  "",
			HeaderFg:  "#64748B",
		},
		Layout: layoutConfig{
			InputHeight: 3,
			ShowHeader:  true,
			ShowStatus:  true,
		},
		Input: inputConfig{
			Placeholder: "message  (/help)",
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
	user         lipgloss.Style
	assistant    lipgloss.Style
	system       lipgloss.Style
	header       lipgloss.Style
	msgUser      lipgloss.Style
	msgAssistant lipgloss.Style
	dim          lipgloss.Style
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
		user:         lipgloss.NewStyle().Foreground(lipgloss.Color(c.User)).Bold(true),
		assistant:    lipgloss.NewStyle().Foreground(lipgloss.Color(c.Assistant)),
		system:       lipgloss.NewStyle().Foreground(lipgloss.Color(c.System)),
		header:       header,
		msgUser:      lipgloss.NewStyle().PaddingLeft(2),
		msgAssistant: lipgloss.NewStyle().PaddingLeft(2),
		dim:          lipgloss.NewStyle().Foreground(lipgloss.Color(c.Accent)),
	}
}
