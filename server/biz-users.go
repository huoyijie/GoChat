package main

import (
	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 处理获取用户列表请求
type biz_users_t struct {
	biz_base_t
}

func initialUsers(base biz_base_t) *biz_users_t {
	return &biz_users_t{base}
}

func (u *biz_users_t) do(req proto.Message, accId *uint64, accUN *string) error {
	pack, err := u.toPacket(req)
	if err != nil {
		return err
	}

	if len(*accUN) == 0 {
		return u.poster.Handle(pack, &lib.UsersRes{Code: lib.Err_Forbidden.Val()})
	}

	users, err := u.storage.GetUsers(*accUN)
	if err != nil {
		return u.poster.Handle(pack, &lib.UsersRes{Code: lib.Err_Get_Users.Val()})
	}

	return u.poster.Handle(pack, &lib.UsersRes{Users: users})
}

var _ biz_i = (*biz_users_t)(nil)
