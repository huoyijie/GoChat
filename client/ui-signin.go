package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/huoyijie/GoChat/lib"
)

func signinSubmit(m *ui_form_t) (tea.Model, tea.Cmd) {
	passhash := sha256.Sum256([]byte(m.inputs[1].Value()))

	tokenRes := &lib.TokenRes{}
	if err := m.poster.Handle(&lib.Signin{Auth: &lib.Auth{
		Username: m.inputs[0].Value(),
		Passhash: passhash[:],
	}}, tokenRes); err != nil {
		m.hint = fmt.Sprintf("登录帐号异常: %v", err)
		return m, nil
	} else if tokenRes.Code < 0 {
		m.hint = fmt.Sprintf("登录帐号异常: %d", tokenRes.Code)
		return m, nil
	}

	if err := m.storage.StoreToken(tokenRes); err != nil {
		return m, tea.Quit
	}

	users := initialUsers(m.ui_base_t)
	return users, users.Init()
}

type ui_signin_t struct {
	ui_form_t
}

func initialSignin(base ui_base_t) ui_signin_t {
	m := initialForm(
		base,
		2,
		[]string{"用户名", "密码"},
		"登录",
		[]check_fn{usernameLenCheck, passwordLenCheck},
		signinSubmit,
	)

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

	return ui_signin_t{ui_form_t: m}
}
