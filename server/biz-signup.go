package main

import (
	"encoding/base64"

	"github.com/huoyijie/GoChat/lib"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
)

// 处理注册请求
type signup struct {
	base
}

func initialSignup(base base) *signup {
	return &signup{base}
}

func (s *signup) do(req proto.Message, accId *uint64, accUN *string) error {
	pack, err := s.toPacket(req)
	if err != nil {
		return err
	}

	signup := &lib.Signup{}
	if err := s.unmarshal(pack, signup); err != nil {
		return err
	}

	passhashAndBcrypt, err := bcrypt.GenerateFromPassword(signup.Auth.Passhash, 14)
	if err != nil {
		return s.poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Bcrypt_Gen.Val()})
	}

	account := &Account{
		Username:          signup.Auth.Username,
		PasshashAndBcrypt: base64.StdEncoding.EncodeToString(passhashAndBcrypt),
	}
	if err := s.storage.NewAccount(account); err != nil {
		return s.poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Acc_Exist.Val()})
	}

	return s.handleAuth(pack, account, accId, accUN)
}

var _ biz = (*signup)(nil)
