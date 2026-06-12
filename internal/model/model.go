package model

import (
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/google/uuid"

	"github.com/bubnyukab/go-todo-cli/internal/store"
	"github.com/bubnyukab/go-todo-cli/internal/ui"
)

const (
	inputView uint = iota
	listView
)

type model struct {
	store    *store.Store
	state    uint
	styles   ui.Styles
	input    textinput.Model
	list     list.Model
	editMode bool
	width    int
}

func NewModel(s *store.Store) (model, error) {
	todos, err := s.GetTodos()
	if err != nil {
		return model{}, err
	}

	items := make([]list.Item, len(todos))
	for i, t := range todos {
		items[i] = t
	}

	m := model{
		state:  inputView,
		store:  s,
		styles: ui.NewStyles(),
	}

	m.input = newInput(m.styles)
	m.list = newList(items, m.styles)

	return m, nil
}

func newInput(styles ui.Styles) textinput.Model {
	input := textinput.New()
	input.Placeholder = "Add a new todo..."
	input.Focus()
	input.SetStyles(styles.InputPlaceholder)
	return input
}

func newList(items []list.Item, styles ui.Styles) list.Model {
	l := list.New(
		items,
		ItemDelegate{
			Styles: styles,
			State:  inputView,
		},
		0,
		0,
	)

	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)

	return l
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// global keys
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.state == inputView {
				if m.editMode {
					m.editMode = false
					m.input.Reset()
				}
				m.state = listView
				m.input.Blur()
			} else {
				m.state = inputView
				cmds = append(cmds, m.input.Focus())
			}
		}

		switch m.state {
		case inputView:
			switch msg.String() {
			case "enter":
				inputValue := strings.TrimSpace(m.input.Value())

				if inputValue == "" {
					break
				}

				// Cap todo list to 20 task max
				if !m.editMode && len(m.list.Items()) >= 20 {
					break
				}

				if m.editMode {
					selectedItem := m.list.SelectedItem().(store.Todo)
					updated := store.Todo{
						ID:   selectedItem.ID,
						Body: inputValue,
						Done: selectedItem.Done,
					}
					m.list.SetItem(
						m.list.Index(), updated,
					)
					if err := m.store.SaveTodo(updated); err != nil {
						return m, tea.Quit
					}
					m.editMode = false
				} else {
					newTodo := store.Todo{ID: uuid.New(), Body: inputValue}
					m.list.InsertItem(m.list.Index(), newTodo)
					if err := m.store.SaveTodo(newTodo); err != nil {
						return m, tea.Quit
					}
				}
				m.input.Reset()
			}
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)

		case listView:
			i, ok := m.list.SelectedItem().(store.Todo)
			listIndex := m.list.Index()
			switch msg.String() {
			case "ctrl+d":
				if ok {
					m.list.RemoveItem(listIndex)
					if err := m.store.DeleteTodo(i.ID); err != nil {
						return m, tea.Quit
					}
				}
			case "ctrl+e":
				m.editMode = true
				m.input.SetValue(i.Body)
				m.state = inputView
				cmds = append(cmds, m.input.Focus())
			case "space":
				if ok {
					i.Done = !i.Done
					m.list.SetItem(listIndex, i)
					if err := m.store.SaveTodo(i); err != nil {
						return m, tea.Quit
					}
				}
			}
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.list.SetWidth(m.width)
		m.list.SetHeight(10)
		m.input.SetWidth(m.width - 6)
	}

	return m, tea.Batch(cmds...)
}
