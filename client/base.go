package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/huoyijie/GoChat/lib"
)

type base struct {
	packChan chan<- *lib.Packet
}

func (m base) Init() tea.Cmd {
	return nil
}

func (m base) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m base) View() string {
	return ""
}

var _ tea.Model = (*base)(nil)
