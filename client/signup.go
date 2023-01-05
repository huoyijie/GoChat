package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/huoyijie/GoChat/lib"
	"github.com/muesli/reflow/indent"
)

var (
	signupLabels        = []string{"用户名", "密码", "确认密码"}
	focusedSignup = focusedStyle.Copy().Render("[ 注册 ]")
	blurredSignup = fmt.Sprintf("[ %s ]", blurredStyle.Render("注册"))
)

type signup struct {
	base
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
}

func initialSignup(packChan chan<- *lib.Packet) signup {
	m := signup{
		inputs: make([]textinput.Model, 3),
		base:   base{packChan: packChan},
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 16

		switch i {
		case 0:
			t.Placeholder = "huoyijie"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 32
			t.Validate = usernameValidator
		case 1:
			t.Placeholder = "hello123"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.Validate = passwordValidator
		case 2:
			t.Placeholder = "hello123"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.Validate = passwordValidator
		}

		m.inputs[i] = t
	}

	return m
}

func (m signup) Init() tea.Cmd {
	return textinput.Blink
}

func (m signup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m, cmd := m.base.Update(msg); cmd != nil {
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				bytes, err := lib.Marshal(&lib.Signup{Auth: &lib.Auth{
					Username: m.inputs[0].Value(),
					Password: m.inputs[1].Value(),
				}})
				if err != nil {
					return m, tea.Quit
				}
				m.base.packChan <- &lib.Packet{Kind: lib.PackKind_SIGNUP, Data: bytes}
				return &users{base: base{packChan: m.base.packChan}}, nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *signup) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m signup) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(inputStyle.Width(30).Render(signupLabels[i]))
		b.WriteRune('\n')
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteString("\n\n")
		}
	}

	button := &blurredSignup
	if m.focusIndex == len(m.inputs) {
		button = &focusedSignup
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("鼠标模式: "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r 改变模式)"))

	return indent.String("\n"+b.String()+"\n\n", 4)
}

var _ tea.Model = (*signup)(nil)
