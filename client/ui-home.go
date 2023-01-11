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
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
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
			return next, next.Init()
		}
	}
	return m, nil
}

func (m home) View() string {
	tpl := "%s\n\n"
	tpl += subtle("k/j, up/down: 选择") + dot + subtle("enter: 确认") + dot + subtle("q, esc: 退出")

	choices := fmt.Sprintf(
		"%s\n%s",
		checkbox(choices[CHOICE_SIGNUP], m.choice == CHOICE_SIGNUP),
		checkbox(choices[CHOICE_SIGNIN], m.choice == CHOICE_SIGNIN),
	)

	s := fmt.Sprintf(tpl, choices)
	return indent.String("\n"+s+"\n\n", 4)
}

var _ tea.Model = (*home)(nil)
