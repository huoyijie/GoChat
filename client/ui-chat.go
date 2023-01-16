package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/huoyijie/GoChat/lib"
	"github.com/muesli/reflow/indent"
)

type ui_chat_t struct {
	ui_base_t
	from        string
	to          string
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func initialChat(to string, base ui_base_t) ui_chat_t {
	kv, err := base.storage.GetValue("username")
	lib.FatalNotNil(err)

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(40, 10)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return ui_chat_t{
		ui_base_t:   base,
		from:        kv.Value,
		to:          to,
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m ui_chat_t) Init() tea.Cmd {
	return tea.Batch(tick(), textarea.Blink)
}

func (m ui_chat_t) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlR:
			users := initialUsers(m.ui_base_t)
			return users, users.Init()
		case tea.KeyEnter:
			if len(strings.TrimSpace(m.textarea.Value())) == 0 {
				return m, nil
			}

			if err := m.poster.Send(&lib.Msg{Kind: lib.MsgKind_TEXT, From: m.from, To: m.to, Data: []byte(m.textarea.Value())}); err != nil {
				return m, tea.Quit
			}

			m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%s: ", m.from))+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.viewport.GotoBottom()
			m.textarea.Reset()
		}

	case tick_msg_t:
		msgList, _ := m.storage.GetMsgList(m.to)
		for i := range msgList {
			m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%s: ", msgList[i].From))+string(msgList[i].Data))
		}
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		return m, tick()

	// We handle errors just like any other message
	case err_msg_t:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m ui_chat_t) View() string {
	help := subtle("enter send") + dot + subtle("ctrl+r back") + dot + subtle("esc quit")

	s := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s",
		inputStyle.Width(32).Render("@"+m.to),
		m.viewport.View(),
		m.textarea.View(),
		help,
	) + "\n\n"

	return indent.String("\n"+s, 4)
}

var _ tea.Model = (*ui_chat_t)(nil)
