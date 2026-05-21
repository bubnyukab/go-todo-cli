package main

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/muesli/reflow/truncate"
)

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
