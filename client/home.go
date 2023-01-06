package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

var (
	choices = []string{"注册", "登录"}
)

const (
	CHOICE_SIGNUP int = iota
	CHOICE_SIGNIN
)

type home struct {
	base
	choice int
}

func (m home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m, cmd := m.base.Update(msg); cmd != nil {
		return m, cmd
	}
	return updateChoices(msg, m)
}

func (m home) View() string {
	s := choicesView(m)
	return indent.String("\n"+s+"\n\n", 4)
}

var _ tea.Model = (*home)(nil)

// The first view, where you're choosing a task
func choicesView(m home) string {

	tpl := "%s\n\n"
	tpl += subtle("j/k, up/down: 选择") + dot + subtle("enter: 确认") + dot + subtle("q, esc: 退出")

	choices := fmt.Sprintf(
		"%s\n%s",
		checkbox(choices[CHOICE_SIGNUP], m.choice == CHOICE_SIGNUP),
		checkbox(choices[CHOICE_SIGNIN], m.choice == CHOICE_SIGNIN),
	)

	return fmt.Sprintf(tpl, choices)
}

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m home) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.choice++
			if m.choice > len(choices)-1 {
				m.choice = len(choices) - 1
			}
		case "k", "up":
			m.choice--
			if m.choice < CHOICE_SIGNUP {
				m.choice = CHOICE_SIGNUP
			}
		case "enter":
			var next tea.Model
			if m.choice == CHOICE_SIGNIN {
				next = initialSignin(m.base)
			} else {
				next = initialSignup(m.base)
			}
			return next, nil
		}
	}
	return m, nil
}
