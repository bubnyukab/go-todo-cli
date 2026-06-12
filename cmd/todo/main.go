package main

import (
	"log"

	tea "charm.land/bubbletea/v2"

	"github.com/bubnyukab/go-todo-cli/internal/model"
	"github.com/bubnyukab/go-todo-cli/internal/store"
)

func main() {
	s := new(store.Store)
	if err := s.Init(); err != nil {
		log.Fatalf("unable to init store: %v", err)
	}

	p := tea.NewProgram(model.NewModel(s))
	if _, err := p.Run(); err != nil {
		log.Fatalf("unable to run tui: %v", err)
	}
}
