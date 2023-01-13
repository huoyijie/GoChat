package main

import (
	"fmt"
	"time"

	"github.com/huoyijie/GoChat/lib"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

const (
	hotPink = lipgloss.Color("#FF06B7")
)

// General stuff for styling the view
var (
	term = termenv.EnvColorProfile()
	// keyword             = makeFgStyle("211")
	subtle              = makeFgStyle("241")
	dot                 = colorFg(" • ", "236")
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	inputStyle          = lipgloss.NewStyle().Foreground(hotPink)
)

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

type (
	errMsg  error
	tickMsg time.Time
)

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// 每个 ui 对象可嵌入 base 对象
type base struct {
	// 通过 poster 向服务器发送请求
	poster lib.Post
	// 通过 storage 读写本地存储
	storage *Storage
}

func initialBase(poster lib.Post, storage *Storage) base {
	return base{
		poster,
		storage,
	}
}

func (b *base) close() {
	b.poster.Close()
}
