package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

type model struct{}

func (m model) Init() tea.Cmd {
	// Start ticking.
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		// Return your Tick command again to loop.
		return m, tick()
	}
	return m, nil
}

func (m model) View() string {
	return "hello"
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
