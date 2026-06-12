package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/bubnyukab/go-todo-cli/internal/model"
	"github.com/bubnyukab/go-todo-cli/internal/store"
)

func Run() error {
	s := store.New()

	if err := s.Init(); err != nil {
		return err
	}

	m, err := model.NewModel(s)
	if err != nil {
		return err
	}

	p := tea.NewProgram(m)
	_, err = p.Run()
	return err
}
