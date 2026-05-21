package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/muesli/reflow/wrap"

	"github.com/bubnyukab/go-todo-cli/store"
)

func (m model) View() tea.View {
	str := m.headerView() + "\n" + m.input.View() + "\n"
	var wrappedText string
	i, ok := m.list.SelectedItem().(store.Todo)
	if ok && m.width > 0 {
		wrappedText = wrap.String(i.Body, m.width)
	}
	str += m.list.View() + "\n" + wrappedText + "\n" + m.footerView()
	v := tea.NewView(str)
	v.AltScreen = true
	return v
}

func (m model) headerView() string { return "What will you do today?\n" }

func (m model) footerView() string {
	return "\n(ctrl+c to quit) (ctrl+e to edit) (ctrl+d to delete) (space to mark done)"
}
