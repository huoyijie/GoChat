package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/huoyijie/GoChat/lib"
)

func signup_submit(m *form) (tea.Model, tea.Cmd) {
	if len(m.inputs[0].Value()) < 3 {
		m.errs[0] = "用户名至少包含3个字母或数字"
		return m, nil
	}

	if len(m.inputs[1].Value()) < 8 {
		m.errs[1] = "密码至少包含8个字母或数字"
		return m, nil
	}

	if len(m.inputs[2].Value()) < 8 {
		m.errs[2] = "密码至少包含8个字母或数字"
		return m, nil
	}

	if m.inputs[1].Value() != m.inputs[2].Value() {
		m.errs[1] = "两次密码输入不一致"
		m.errs[2] = "两次密码输入不一致"
		return m, nil
	}

	passhash := sha256.Sum256([]byte(m.inputs[1].Value()))

	tokenRes := &lib.TokenRes{}
	if err := m.poster.Handle(&lib.Signup{Auth: &lib.Auth{
		Username: m.inputs[0].Value(),
		Passhash: passhash[:],
	}}, tokenRes); err != nil {
		m.hint = fmt.Sprintf("注册帐号异常: %v", err)
		return m, nil
	} else if tokenRes.Code < 0 {
		m.hint = fmt.Sprintf("注册帐号异常: %d", tokenRes.Code)
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

type signup struct {
	form
}

func initialSignup(base base) signup {
	m := initialForm(base, 3, []string{"用户名", "密码", "确认密码"}, "注册", signup_submit)

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
		case 2:
			t.Placeholder = "hello123"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.Validate = passwordValidator
		}

		m.inputs[i] = t
	}

	return signup{form: m}
}
