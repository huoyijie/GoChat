package main

import (
	"encoding/base64"

	"github.com/huoyijie/GoChat/lib"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
)

type signin struct {
	base
}

func initialSignin(base base) *signin {
	return &signin{base}
}

// 处理登录请求
func (s *signin) do(req proto.Message, accId *uint64, accUN *string) error {
	pack, err := s.toPacket(req)
	if err != nil {
		return err
	}

	signin := &lib.Signin{}
	if err := s.unmarshal(pack, signin); err != nil {
		return err
	}

	account, err := s.storage.GetAccountByUN(signin.Auth.Username)
	if err != nil {
		return s.poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Acc_Not_Exist.Val()})
	}

	passhashAndBcrypt, err := base64.StdEncoding.DecodeString(account.PasshashAndBcrypt)
	if err != nil {
		return s.poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Base64_Decode.Val()})
	}

	if err := bcrypt.CompareHashAndPassword(passhashAndBcrypt, signin.Auth.Passhash); err != nil {
		return s.poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Bcrypt_Compare.Val()})
	}

	return s.handleAuth(pack, account, accId, accUN)
}

var _ biz = (*signin)(nil)
