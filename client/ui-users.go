package main

import (
	"fmt"
	"strings"

	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/huoyijie/GoChat/lib"
	"github.com/muesli/reflow/indent"
)

const (
	listWidth  = 20
	listHeight = 14
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	usersHelpStyle    = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type item_t struct {
	username string
	online   bool
	msgCount uint32
}

func (i item_t) FilterValue() string { return i.username }

type item_proxy_t struct{}

func (d item_proxy_t) Height() int                               { return 1 }
func (d item_proxy_t) Spacing() int                              { return 0 }
func (d item_proxy_t) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d item_proxy_t) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item_t)
	if !ok {
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d. ", index+1))
	sb.WriteString(i.username)
	if i.online {
		sb.WriteRune('↑')
	}
	if i.msgCount > 0 {
		sb.WriteString(fmt.Sprintf(" (%d+)", i.msgCount))
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(sb.String()))
}

type ui_users_t struct {
	ui_base_t
	list list.Model
}

func newList(poster lib.Post, storage *storage_t) list.Model {
	usersRes := &lib.UsersRes{}
	if err := poster.Handle(&lib.Users{}, usersRes); err != nil {
		lib.FatalNotNil(err)
	} else if usersRes.Code < 0 {
		lib.FatalNotNil(fmt.Errorf("获取用户列表异常: %d", usersRes.Code))
	}

	unReadMsgCnt, err := storage.UnReadMsgCount()
	lib.FatalNotNil(err)

	items := make([]list.Item, len(usersRes.Users))
	for i := range usersRes.Users {
		items[i] = item_t{
			username: usersRes.Users[i].Username,
			online:   usersRes.Users[i].Online,
			msgCount: unReadMsgCnt[usersRes.Users[i].Username],
		}
	}

	l := list.New(items, item_proxy_t{}, listWidth, listHeight)
	l.Title = "用户列表"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = usersHelpStyle
	return l
}

func initialUsers(base ui_base_t) ui_users_t {
	return ui_users_t{list: newList(base.poster, base.storage), ui_base_t: base}
}

func (m ui_users_t) Init() tea.Cmd {
	return tick()
}

func (m ui_users_t) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", tea.KeyEsc.String(), tea.KeyCtrlC.String():
			return m, tea.Quit
		case tea.KeyCtrlX.String():
			signoutRes := &lib.SignoutRes{}
			if err := m.poster.Handle(&lib.Signout{}, signoutRes); err != nil || signoutRes.Code < 0 {
				return m, nil
			}
			// 删除本地存储数据
			m.storage.DropPrivacy()
			home := initialHome(m.ui_base_t)
			return home, home.Init()
		case tea.KeyEnter.String():
			i, ok := m.list.SelectedItem().(item_t)
			if !ok {
				return m, tea.Quit
			}
			chat := initialChat(i.username, m.ui_base_t)
			return chat, chat.Init()
		}

	case tick_msg_t:
		unReadMsgCnt, err := m.storage.UnReadMsgCount()
		if err != nil {
			return m, nil
		}

		pushes, err := m.storage.GetOnlinePushes()
		if err != nil {
			return m, nil
		}

		var cmds []tea.Cmd
	loop:
		for i := range m.list.Items() {
			v := m.list.Items()[i].(item_t)

			count, hasUnReadMsg := unReadMsgCnt[v.username]
			on, hasOnlinePush := pushes[v.username]

			if !(hasUnReadMsg || hasOnlinePush) {
				continue loop
			}

			item := item_t{
				username: v.username,
				online:   v.online,
				msgCount: v.msgCount,
			}

			if hasUnReadMsg {
				item.msgCount = count
			}

			if hasOnlinePush {
				item.online = on
			}

			cmd := m.list.SetItem(i, item)
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, tick())

		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ui_users_t) View() string {
	help := subtle("↑/k up") + dot + subtle("↓/j down") + dot + subtle("q/esc quit") + dot + subtle("ctrl+x sign out") + dot + subtle("? more")

	s := fmt.Sprintf(
		"\n%s\n%s\n\n",
		m.list.View(),
		help,
	)
	return indent.String(s, 4)
}

var _ tea.Model = (*ui_users_t)(nil)
