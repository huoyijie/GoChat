package main

import (
	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 处理 token 验证请求
type biz_val_token_t struct {
	biz_base_t
}

func initialValToken(base biz_base_t) *biz_val_token_t {
	return &biz_val_token_t{base}
}

func (vt *biz_val_token_t) do(req proto.Message, accId *uint64, accUN *string) error {
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

var _ biz_i = (*biz_val_token_t)(nil)
