package main

import (
	"encoding/base64"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

type users struct {
	base
	id       uint64
	username string
	token    []byte
}

func (m users) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m, cmd := m.base.Update(msg); cmd != nil {
		return m, cmd
	}
	return m, nil
}

func (m users) View() string {
	s := fmt.Sprintf("%d %s\n%s", m.id, m.username, base64.StdEncoding.EncodeToString(m.token))
	return indent.String("\n"+s+"\n\n", 4)
}

var _ tea.Model = (*users)(nil)
