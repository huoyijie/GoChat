package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

type users struct {
	base
}

func (m users) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m, cmd := m.base.Update(msg); cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m users) View() string {
	s := "user list"
	return indent.String("\n"+s+"\n\n", 4)
}

var _ tea.Model = (*users)(nil)
