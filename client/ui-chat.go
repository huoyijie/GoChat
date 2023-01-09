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

type (
	errMsg error
)

func recvMsg(msgChan <-chan *lib.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-msgChan
	}
}

type chat struct {
	base
	from        string
	to          string
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func initialChat(to string, base base) chat {
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
	vp.SetContent(`输入消息按回车键发送`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return chat{
		base:        base,
		from:        kv.Value,
		to:          to,
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m chat) Init() tea.Cmd {
	return tea.Batch(recvMsg(m.msgChan), textarea.Blink)
}

func (m chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			if len(strings.TrimSpace(m.textarea.Value())) == 0 {
				return m, nil
			}

			chatMsg := &lib.Msg{Kind: lib.MsgKind_TEXT, From: m.from, To: m.to, Data: []byte(m.textarea.Value())}
			bytes, err := lib.Marshal(chatMsg)
			if err != nil {
				return m, tea.Quit
			}

			req := new(request_t).init(&lib.Packet{Kind: lib.PackKind_MSG, Data: bytes}, false)
			m.reqChan <- req

			m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%s: ", m.from))+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.viewport.GotoBottom()
			m.textarea.Reset()
		}

	case *lib.Msg:
		m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%s: ", msg.From))+string(msg.Data))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		return m, recvMsg(m.msgChan)

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m chat) View() string {
	s := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		inputStyle.Width(32).Render("@"+m.to),
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
	return indent.String("\n"+s, 4)
}

var _ tea.Model = (*chat)(nil)
