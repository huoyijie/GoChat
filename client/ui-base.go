package main

import (
	"errors"
	"fmt"
	"github.com/huoyijie/GoChat/lib"
	"regexp"

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

type base struct {
	msgChan <-chan *lib.Msg
	reqChan chan<- *request_t
	storage *Storage
}
