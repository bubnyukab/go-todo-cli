package main

import (
	"log"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

func (m mainModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.state == InputView {
				m.state = ListView
				m.input.Blur()
			} else {
				m.state = InputView
				cmds = append(cmds, m.input.Focus())
			}
		case "enter":
			if m.state == InputView && m.editMode == false {
				m.list.InsertItem(m.list.Index(), item{body: m.input.Value(), done: false})
				m.input.Reset()
			} else if m.state == InputView && m.editMode == true {
				m.list.SetItem(
					m.list.Index(),
					item{body: m.input.Value(), done: m.list.SelectedItem().(item).done},
				)
				m.input.Reset()
				m.editMode = false
			}
		case "ctrl+d":
			if m.state == ListView {
				m.list.RemoveItem(m.list.Index())
			}
		case "ctrl+e":
			if m.state == ListView {
				m.editMode = true
				i := m.list.SelectedItem().(item)
				m.input.SetValue(i.body)
				m.state = InputView
				m.input.Focus()
			}
		}

		switch m.state {
		case InputView:
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		default:
			switch msg.String() {
			case "space":
				i, ok := m.list.SelectedItem().(item)
				if !ok {
					return m, nil
				}

				i.done = !i.done
				m.list.SetItem(m.list.Index(), i)
			}
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 6)
	}

	return m, tea.Batch(cmds...)
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
