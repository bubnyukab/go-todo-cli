package main

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
)

// sessionState is used to track which model is focused
type state uint

const (
	listHeight = 5
	listWidth  = 20
)

const (
	InputView state = iota
	ListView
)

type mainModel struct {
	state    state
	input    textinput.Model
	list     list.Model
	editMode bool
	width    int
}

func newModel() mainModel {
	m := mainModel{state: InputView}
	items := []list.Item{}
	m.input = textinput.New()
	m.input.Focus()
	m.list = list.New(items, itemDelegate{}, listWidth, listHeight)
	m.list.SetFilteringEnabled(false)
	m.list.SetShowHelp(false)

	return m
}
