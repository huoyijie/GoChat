package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

type submit_fn = func(*form) (tea.Model, tea.Cmd)

type form struct {
	base
	focusIndex int
	inputs     []textinput.Model
	errs       []string
	hint       string
	cursorMode textinput.CursorMode
	labels     []string
	button     string
	submit     submit_fn
}

func initialForm(base base, inputs int, labels []string, button string, submit submit_fn) form {
	return form{
		base:   base,
		inputs: make([]textinput.Model, inputs),
		errs:   make([]string, inputs),
		hint:   "",
		labels: labels,
		button: button,
		submit: submit,
	}
}

func (m form) Init() tea.Cmd {
	return textinput.Blink
}

func (m *form) updateInputs(msg tea.Msg) tea.Cmd {
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

func (m form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m form) View() string {
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

	button := focusedStyle.Copy().Render(fmt.Sprintf("[ %s ]", m.button))
	if m.focusIndex == len(m.inputs) {
		button = fmt.Sprintf("[ %s ]", blurredStyle.Render(m.button))
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

var _ tea.Model = (*form)(nil)
