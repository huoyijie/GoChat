package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

// 用户名实时输入验证，只允许输入小写字母和数字，长度在[1,32]之间
func usernameValidator(s string) (err error) {
	usernameRegexp := "^[a-z\\d]{1,32}$"
	re, err := regexp.Compile(usernameRegexp)
	if err != nil {
		return
	}

	if !re.MatchString(s) {
		err = errors.New("username is invalid")
	}
	return
}

// 密码实时输入验证器，只允许输入小写字母和数字，长度在[1,16]之间
func passwordValidator(s string) (err error) {
	passwordRegexp := "^[a-z\\d]{1,16}$"
	re, err := regexp.Compile(passwordRegexp)
	if err != nil {
		return
	}

	if !re.MatchString(s) {
		err = errors.New("password is invalid")
	}
	return
}

type check_fn = func(string) (bool, string)

// 表单提交后检查用户名长度
func usernameLenCheck(s string) (ok bool, hint string) {
	if len(s) < 3 {
		hint = "用户名至少包含3个字母或数字"
		return
	}
	ok = true
	return
}

// 表单提交后检查密码长度
func passwordLenCheck(s string) (ok bool, hint string) {
	if len(s) < 8 {
		hint = "密码至少包含8个字母或数字"
		return
	}
	ok = true
	return
}

// 表单提交处理函数
type submit_fn = func(*ui_form_t) (tea.Model, tea.Cmd)

// 登录和注册表单基类
type ui_form_t struct {
	ui_base_t
	focusIndex int
	inputs     []textinput.Model
	errs       []string
	hint       string
	cursorMode textinput.CursorMode
	labels     []string
	button     string
	lenChecks  []check_fn
	submit     submit_fn
}

func initialForm(base ui_base_t, inputs int, labels []string, button string, lenChecks []check_fn, submit submit_fn) ui_form_t {
	return ui_form_t{
		ui_base_t: base,
		inputs:    make([]textinput.Model, inputs),
		errs:      make([]string, inputs),
		hint:      "",
		labels:    labels,
		button:    button,
		lenChecks: lenChecks,
		submit:    submit,
	}
}

func (m ui_form_t) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ui_form_t) updateInputs(msg tea.Msg) tea.Cmd {
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

func (m ui_form_t) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keyType := msg.Type; keyType {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		// Change cursor mode
		case tea.KeyCtrlR:
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
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if keyType == tea.KeyEnter && m.focusIndex == len(m.inputs) {
				for i := range m.lenChecks {
					if ok, hint := m.lenChecks[i](m.inputs[i].Value()); !ok {
						m.errs[i] = hint
						return m, nil
					}
				}
				return m.submit(&m)
			}

			// Cycle indexes
			if keyType == tea.KeyUp || keyType == tea.KeyShiftTab {
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

func (m ui_form_t) View() string {
	var b strings.Builder

	if len(m.hint) > 0 {
		b.WriteString(m.hint)
		b.WriteString("\n\n")
	}

	for i := range m.inputs {
		b.WriteString(inputStyle.Width(30).Render(m.labels[i]))
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

	button := fmt.Sprintf("[ %s ]", blurredStyle.Render(m.button))
	if m.focusIndex == len(m.inputs) {
		button = focusedStyle.Copy().Render(fmt.Sprintf("[ %s ]", m.button))
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button)

	b.WriteString(helpStyle.Render("鼠标模式: "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r 改变模式)"))
	b.WriteRune('\n')

	help := subtle("↑/shift+tab up") + dot + subtle("↓/tab down") + dot + subtle("enter select") + dot + subtle("esc quit")

	b.WriteString(help)

	return indent.String("\n"+b.String()+"\n\n", 4)
}

var _ tea.Model = (*ui_form_t)(nil)
