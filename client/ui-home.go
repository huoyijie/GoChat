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

func initialHome(base base) home {
	return home{choice: CHOICE_SIGNIN, base: base}
}

func (m home) Init() tea.Cmd {
	return nil
}

func (m home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", tea.KeyEsc.String(), tea.KeyCtrlC.String():
			return m, tea.Quit
		case "j", tea.KeyDown.String():
			m.choice++
			if m.choice > len(choices)-1 {
				m.choice = len(choices) - 1
			}
		case "k", tea.KeyUp.String():
			m.choice--
			if m.choice < CHOICE_SIGNUP {
				m.choice = CHOICE_SIGNUP
			}
		case tea.KeyEnter.String():
			var next tea.Model
			if m.choice == CHOICE_SIGNIN {
				next = initialSignin(m.base)
			} else {
				next = initialSignup(m.base)
			}
			return next, next.Init()
		}
	}
	return m, nil
}

func (m home) View() string {
	tpl := "%s\n\n"
	tpl += subtle("↑/k up") + dot + subtle("↓/j down") + dot + subtle("enter select") + dot + subtle("q/esc quit")

	choices := fmt.Sprintf(
		"%s\n%s",
		checkbox(choices[CHOICE_SIGNUP], m.choice == CHOICE_SIGNUP),
		checkbox(choices[CHOICE_SIGNIN], m.choice == CHOICE_SIGNIN),
	)

	s := fmt.Sprintf(tpl, choices)
	return indent.String("\n"+s+"\n\n", 4)
}

var _ tea.Model = (*home)(nil)
