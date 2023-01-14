package main

import (
	"github.com/huoyijie/GoChat/lib"
	"google.golang.org/protobuf/proto"
)

// 处理获取用户列表请求
type users struct {
	base
}

func initialUsers(base base) *users {
	return &users{base}
}

func (u *users) do(req proto.Message, accId *uint64, accUN *string) error {
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

var _ biz = (*users)(nil)
