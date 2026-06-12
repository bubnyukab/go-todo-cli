package model

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/muesli/reflow/truncate"

	"github.com/bubnyukab/go-todo-cli/internal/store"
	"github.com/bubnyukab/go-todo-cli/internal/ui"
)

type ItemDelegate struct {
	Styles ui.Styles
	State  uint
}

func (d ItemDelegate) Height() int { return 1 }

func (d ItemDelegate) Spacing() int { return 0 }

func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
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
	if index == m.Index() && d.State == listView {
		cb := d.Styles.SelectedCheckbox.Render(checkbox)
		bd := d.Styles.SelectedBody.Width(m.Width() - len(checkbox) - 1).Render(" " + body)
		str = cb + bd

	} else if i.Done {
		str = d.Styles.Done.Render(checkbox + " " + body)
	} else {
		str = d.Styles.Checkbox.Render(checkbox) + " " + body
	}

	fmt.Fprintln(w, str)
}
