package main

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/muesli/reflow/wrap"

	"github.com/bubnyukab/go-todo-cli/store"
)

func (m model) View() tea.View {
	m.list.SetDelegate(itemDelegate{styles: m.styles, state: m.state})

	str := m.headerView() + "\n"
	str += m.inputView() + "\n"
	str += m.statusView() + "\n"
	str += m.listView() + "\n"
	str += m.previewView() + "\n"
	str += m.helpView()

	v := tea.NewView(str)
	v.AltScreen = true
	return v
}

func (m model) headerView() string {
	return m.styles.Header.Width(m.width).Render("\nWhat will you do today?\n")
}

func (m model) inputView() string {
	style := m.styles.InputFocused
	if m.state == listView {
		style = m.styles.InputBlurred
	}
	return style.Width(m.width).Render(m.input.View())
}

func (m model) statusView() string {
	items := m.list.Items()
	total := len(items)
	done := 0
	for _, item := range items {
		if t, ok := item.(store.Todo); ok && t.Done {
			done++
		}
	}
	if total == 0 {
		return ""
	}
	pct := done * 100 / total
	text := fmt.Sprintf("%d/%d done (%d%%)", done, total, pct)
	return m.styles.StatusBar.Padding(0, 0, 1, 2).Render(text)
}

func (m model) listView() string {
	listView := m.list.View()
	pagerView := m.list.Paginator.View()

	styledList := lipgloss.NewStyle().PaddingLeft(1).Render(listView)

	if len(m.list.Items()) > 10 {
		styledPager := lipgloss.NewStyle().PaddingLeft(1).Render(pagerView)

		return lipgloss.JoinVertical(lipgloss.Left, styledList, styledPager)
	}

	return styledList
}

func (m model) previewView() string {
	borderColor := lipgloss.Color("240")
	label := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("no item selected")

	i, ok := m.list.SelectedItem().(store.Todo)
	if ok && i.Body != "" && m.state == listView {
		borderColor = lipgloss.Color("99")
		label = wrap.String(i.Body, m.width-4)
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(m.width).
		Render(label)
}

func (m model) helpView() string {
	var keys []key.Binding
	if m.state == listView {
		keys = []key.Binding{
			key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "input")),
			key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
			key.NewBinding(
				key.WithKeys("\u2191", "\u2193"),
				key.WithHelp("\u2191/\u2193", "navigate"),
			),
			key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "toggle done")),
			key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "edit")),
			key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "delete")),
		}
	} else {
		keys = []key.Binding{
			key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "list")),
			key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "add/save")),
		}
	}

	sep := " \u2022 "
	maxWidth := m.width + 8
	var lines []string
	current := ""
	for _, k := range keys {
		h := k.Help()
		pair := h.Key + " " + h.Desc
		if current == "" {
			current = pair
		} else if len(current+sep+pair) <= maxWidth {
			current += sep + pair
		} else {
			lines = append(lines, current)
			current = pair
		}
	}
	if current != "" {
		lines = append(lines, current)
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		PaddingLeft(1).Render(strings.Join(lines, "\n"))
}
