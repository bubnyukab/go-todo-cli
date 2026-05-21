package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/muesli/reflow/wrap"
)

func (m mainModel) View() tea.View {
	str := m.headerView() + "\n" + m.input.View() + "\n"
	var wrapped_text string
	i, ok := m.list.SelectedItem().(item)
	if ok && m.width > 0 {
		wrapped_text = wrap.String(i.body, m.width)
	}
	str += m.list.View() + "\n" + wrapped_text + "\n" + m.footerView()
	v := tea.NewView(str)
	v.AltScreen = true
	return v
}

func (m mainModel) headerView() string { return "What will you do today?\n" }

func (m mainModel) footerView() string {
	return "\n(ctrl+c to quit) (ctrl+e to edit) (ctrl+d to delete) (space to mark done)"
}
