package main

import (
	"fmt"

	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/huoyijie/GoChat/lib"
	"github.com/muesli/reflow/indent"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	usersHelpStyle    = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type item struct {
	username string
	msgCount uint32
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	var str string
	if i.msgCount > 0 {
		str = fmt.Sprintf("%d. %s (%d+)", index+1, i.username, i.msgCount)
	} else {
		str = fmt.Sprintf("%d. %s", index+1, i.username)
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type users struct {
	base
	list list.Model
}

func initialUsers(base base) users {
	usersRes := &lib.UsersRes{}
	if err := base.poster.Handle(&lib.Users{}, usersRes); err != nil {
		lib.FatalNotNil(err)
	} else if usersRes.Code < 0 {
		lib.FatalNotNil(fmt.Errorf("获取用户列表异常: %d", usersRes.Code))
	}

	unReadMsgCnt, err := base.storage.UnReadMsgCount()
	lib.FatalNotNil(err)

	items := make([]list.Item, len(usersRes.Users))
	for i := range usersRes.Users {
		items[i] = item{
			username: usersRes.Users[i],
			msgCount: unReadMsgCnt[usersRes.Users[i]],
		}
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "用户列表"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = usersHelpStyle

	return users{list: l, base: base}
}

func (m users) Init() tea.Cmd {
	return tick()
}

func (m users) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", tea.KeyEsc.String(), tea.KeyCtrlC.String():
			return m, tea.Quit
		case tea.KeyCtrlX.String():
			// 删除本地存储文件
			dropDB()
			home := initialHome(m.base)
			return home, home.Init()
		case tea.KeyEnter.String():
			i, ok := m.list.SelectedItem().(item)
			if !ok {
				return m, tea.Quit
			}
			chat := initialChat(i.username, m.base)
			return chat, chat.Init()
		}

	case tickMsg:
		unReadMsgCnt, err := m.storage.UnReadMsgCount()
		if err != nil {
			return m, nil
		}

		var cmds []tea.Cmd
		for i := range m.list.Items() {
			v := m.list.Items()[i].(item)
			if count, ok := unReadMsgCnt[v.username]; ok {
				cmd := m.list.SetItem(
					i,
					item{
						username: v.username,
						msgCount: count,
					},
				)
				cmds = append(cmds, cmd)
			}
		}
		cmds = append(cmds, tick())

		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m users) View() string {
	help := subtle("↑/k up") + dot + subtle("↓/j down") + dot + subtle("q/esc quit") + dot + subtle("ctrl+x sign out") + dot + subtle("? more")

	s := fmt.Sprintf(
		"\n%s\n%s\n\n",
		m.list.View(),
		help,
	)
	return indent.String(s, 4)
}

var _ tea.Model = (*users)(nil)
