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

type ui_home_t struct {
	ui_base_t
	choice int
}

func initialHome(base ui_base_t) ui_home_t {
	return ui_home_t{choice: CHOICE_SIGNIN, ui_base_t: base}
}

func (m ui_home_t) Init() tea.Cmd {
	return nil
}

func (m ui_home_t) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				next = initialSignin(m.ui_base_t)
			} else {
				next = initialSignup(m.ui_base_t)
			}
			return next, next.Init()
		}
	}
	return m, nil
}

func (m ui_home_t) View() string {
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

var _ tea.Model = (*ui_home_t)(nil)
