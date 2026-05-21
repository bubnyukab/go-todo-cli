package main

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/muesli/reflow/truncate"

	"github.com/bubnyukab/go-todo-cli/store"
)

type item struct {
	body string
	done bool
}

type itemDelegate struct{}

func (d itemDelegate) Height() int { return 1 }

func (d itemDelegate) Spacing() int { return 0 }

func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	markedDone := " "
	i, ok := listItem.(store.Todo)
	if !ok {
		return
	}

	if i.Done {
		markedDone = "X"
	} else {
		markedDone = " "
	}

	str := fmt.Sprintf("[%s] %s", markedDone, i.Body)

	str = truncate.StringWithTail(str, uint(m.Width()-5), "...")
	if index == m.Index() {
		str += " <"
	}

	fmt.Fprint(w, str)
}
