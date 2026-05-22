package main

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

type Styles struct {
	Header           lipgloss.Style
	Key              lipgloss.Style
	SelectedBody     lipgloss.Style
	SelectedCheckbox lipgloss.Style
	Done             lipgloss.Style
	Checkbox         lipgloss.Style
	InputFocused     lipgloss.Style
	InputBlurred     lipgloss.Style
	InputPlaceholder textinput.Styles
	StatusBar        lipgloss.Style
}

func newStyles() Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("99")).
			Padding(0, 1).Align(lipgloss.Center),

		Key: lipgloss.NewStyle().Foreground(lipgloss.Color("99")),

		SelectedBody: lipgloss.NewStyle().
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("255")).
			Bold(true),

		SelectedCheckbox: lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Background(lipgloss.Color("235")),

		// the strikethrough on some terminals may be off center
		Done: lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Bold(true).Strikethrough(true),

		InputFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			Padding(0, 1),

		InputBlurred: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1),

		InputPlaceholder: newInputStyles(),

		StatusBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			PaddingLeft(1),
	}
}

func newInputStyles() textinput.Styles {
	s := textinput.DefaultDarkStyles()
	s.Focused.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	s.Blurred.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	s.Blurred.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	s.Blurred.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	return s
}
