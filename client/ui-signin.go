package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/huoyijie/GoChat/lib"
	"github.com/muesli/reflow/indent"
)

var (
	signinLabels  = []string{"用户名", "密码"}
	focusedSignin = focusedStyle.Copy().Render("[ 登录 ]")
	blurredSignin = fmt.Sprintf("[ %s ]", blurredStyle.Render("登录"))
)

type signin struct {
	base
	focusIndex int
	inputs     []textinput.Model
	errs       []string
	hint       string
	cursorMode textinput.CursorMode
}

func initialSignin(base base) signin {
	m := signin{
		inputs: make([]textinput.Model, 2),
		errs:   []string{"", "", ""},
		base:   base,
		hint:   "",
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
		}

		m.inputs[i] = t
	}

	return m
}

func (m signin) Init() tea.Cmd {
	return textinput.Blink
}

func (m signin) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
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
				if len(m.inputs[0].Value()) < 4 {
					m.errs[0] = "用户名至少包含4个字母或数字"
					return m, nil
				}

				if len(m.inputs[1].Value()) < 8 {
					m.errs[1] = "密码至少包含8个字母或数字"
					return m, nil
				}

				passhash := sha256.Sum256([]byte(m.inputs[1].Value()))
				bytes, err := lib.Marshal(&lib.Signin{Auth: &lib.Auth{
					Username: m.inputs[0].Value(),
					Passhash: passhash[:],
				}})
				if err != nil {
					return m, tea.Quit
				}

				req := new(request_t).init(&lib.Packet{Kind: lib.PackKind_SIGNIN, Data: bytes})
				m.reqChan <- req

				res := <-req.c
				if !res.ok() {
					m.hint = "登录超时"
					return m, nil
				}

				tokenRes := &lib.TokenRes{}
				if err := lib.Unmarshal(res.pack.Data, tokenRes); err != nil {
					return m, tea.Quit
				}

				if tokenRes.Code < 0 {
					m.hint = fmt.Sprintf("登录异常: %d", tokenRes.Code)
					return m, nil
				}

				if err := m.storage.NewKVS([]KeyValue{
					{Key: "id", Value: fmt.Sprintf("%d", tokenRes.Id)},
					{Key: "username", Value: tokenRes.Username},
					{Key: "token", Value: base64.StdEncoding.EncodeToString(tokenRes.Token)},
				}); err != nil {
					return m, tea.Quit
				}

				users := initialUsers(m.base)
				return users, users.Init()
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

func (m *signin) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		m.errs[i] = ""
	}
	m.hint = ""

	return tea.Batch(cmds...)
}

func (m signin) View() string {
	var b strings.Builder

	if len(m.hint) > 0 {
		b.WriteString(m.hint)
		b.WriteString("\n\n")
	}

	for i := range m.inputs {
		b.WriteString(inputStyle.Width(30).Render(signinLabels[i]))
		b.WriteRune('\n')
		b.WriteString(m.inputs[i].View())
		if len(m.errs[i]) > 0 {
			b.WriteRune('\n')
			b.WriteString(m.errs[i])
		}
		if i < len(m.inputs)-1 {
			b.WriteString("\n\n")
		}
	}

	button := &blurredSignin
	if m.focusIndex == len(m.inputs) {
		button = &focusedSignin
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("鼠标模式: "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r 改变模式)"))

	return indent.String("\n"+b.String()+"\n\n", 4)
}

var _ tea.Model = (*signin)(nil)
