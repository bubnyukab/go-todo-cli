package main

import (
	"fmt"
	"io"
	"log"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wrap"
)

// sessionState is used to track which model is focused
type sessionState uint

type item struct {
	body string
	done bool
}

func (i item) FilterValue() string { return i.body }

type itemDelegate struct {
	selectedIndex int
	termWidth     int
}

func (d itemDelegate) Height() int { return 1 }

func (d itemDelegate) Spacing() int { return 0 }

func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	markedDone := " "
	i, ok := listItem.(item)
	if !ok {
		return
	}

	if i.done {
		markedDone = "X"
	} else {
		markedDone = " "
	}

	str := fmt.Sprintf("[%s] %s", markedDone, i.body)

	str = truncate.StringWithTail(str, uint(m.Width()-5), "...")
	if index == m.Index() {
		str += " <"
	}

	fmt.Fprint(w, str)
}

const (
	listHeight              = 5
	listWidth               = 20
	inputView  sessionState = iota
	listView
)

type mainModel struct {
	state sessionState
	input textinput.Model
	list  list.Model
	index int

	width int
}

func newModel() mainModel {
	m := mainModel{state: inputView}
	items := []list.Item{}
	m.input = textinput.New()
	m.list = list.New(items, itemDelegate{}, listWidth, listHeight)
	m.list.SetFilteringEnabled(false)
	m.list.SetShowHelp(false)

	return m
}

func (m mainModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == inputView {
				m.state = listView
				m.input.Blur()
			} else {
				m.state = inputView
				cmds = append(cmds, m.input.Focus())
			}
		case "enter":
			if m.state == inputView {
				m.list.InsertItem(m.index, item{body: m.input.Value(), done: false})
				m.input.Reset()
				m.Next()
			}
		}

		switch m.state {
		case inputView:
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

func (m *mainModel) Next() {
	if m.index == len(m.list.Items())-1 {
		m.index = 0
	} else {
		m.index++
	}
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m mainModel) headerView() string { return "What will you do today?\n" }

func (m mainModel) footerView() string { return "\n(esc to quit)" }
