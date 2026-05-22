package main

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/muesli/reflow/truncate"

	"github.com/bubnyukab/go-todo-cli/store"
)

type itemDelegate struct {
	styles Styles
	state  uint
}

func (d itemDelegate) Height() int { return 1 }

func (d itemDelegate) Spacing() int { return 0 }

func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(store.Todo)
	if !ok {
		return
	}

	checkbox := "[ ]"
	if i.Done {
		checkbox = "[x]"
	}

	body := truncate.StringWithTail(i.Body, uint(m.Width()-5), "...")

	var str string
	if index == m.Index() && d.state == listView {
		cb := d.styles.SelectedCheckbox.Render(checkbox)
		bd := d.styles.SelectedBody.Width(m.Width() - len(checkbox) - 1).Render(" " + body)
		str = cb + bd

	} else if i.Done {
		str = d.styles.Done.Render(checkbox + " " + body)
	} else {
		str = d.styles.Checkbox.Render(checkbox) + " " + body
	}

	fmt.Fprintln(w, str)
}
