package main

import (
	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 处理 token 验证请求
type val_token struct {
	base
}

func initialValToken(base base) *val_token {
	return &val_token{base}
}

func (vt *val_token) do(req proto.Message, accId *uint64, accUN *string) error {
	pack, err := vt.toPacket(req)
	if err != nil {
		return err
	}

	tokenReq := &lib.Token{}
	if err := vt.unmarshal(pack, tokenReq); err != nil {
		return err
	}

	id, expired, err := ParseToken(tokenReq.Token)
	if err != nil {
		return vt.poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Parse_Token.Val()})
	}

	if expired {
		return vt.poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Token_Expired.Val()})
	}

	account, err := vt.storage.GetAccountById(id)
	if err != nil {
		return vt.poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Acc_Not_Exist.Val()})
	}

	return vt.handleAuth(pack, account, accId, accUN)
}

var _ biz = (*val_token)(nil)
