package main

import (
	"encoding/base64"
	"net"
	"path/filepath"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/huoyijie/GoChat/lib"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
)

// 把来自 packChan 的 packet 都发送到 conn
func sendTo(conn net.Conn, packChan <-chan *lib.Packet) {
	var pid uint64
	for pack := range packChan {
		if pack.Id == 0 {
			pid++
			pack.Id = pid
		}

		bytes, err := lib.MarshalPack(pack)
		if err != nil {
			return
		}

		_, err = conn.Write(bytes)
		if err != nil {
			return
		}
	}
}

// 反序列化请求对象
func unmarshal(poster lib.Post, pack *lib.Packet, req proto.Message) error {
	if err := lib.Unmarshal(pack.Data, req); err != nil {
		return poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Unmarshal.Val()})
	}
	return nil
}

// 生成 token 并向客户端发送 TokenRes packet
func handleAuth(poster lib.Post, pack *lib.Packet, account *Account, accId *uint64, accUN *string) error {
	token, err := GenerateToken(account.Id)
	if err != nil {
		return poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Gen_Token.Val()})
	}

	err = poster.Handle(pack, &lib.TokenRes{Id: account.Id, Username: account.Username, Token: token})
	if err != nil {
		return err
	}
	*accId = account.Id
	*accUN = account.Username
	return nil
}

// 处理注册请求
func signup(poster lib.Post, storage *Storage, pack *lib.Packet, accId *uint64, accUN *string) error {
	signup := &lib.Signup{}
	if err := unmarshal(poster, pack, signup); err != nil {
		return err
	}

	passhashAndBcrypt, err := bcrypt.GenerateFromPassword(signup.Auth.Passhash, 14)
	if err != nil {
		return poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Bcrypt_Gen.Val()})
	}

	account := &Account{
		Username:          signup.Auth.Username,
		PasshashAndBcrypt: base64.StdEncoding.EncodeToString(passhashAndBcrypt),
	}
	if err := storage.NewAccount(account); err != nil {
		return poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Acc_Exist.Val()})
	}

	return handleAuth(poster, pack, account, accId, accUN)
}

// 处理登录请求
func signin(poster lib.Post, storage *Storage, pack *lib.Packet, accId *uint64, accUN *string) error {
	signin := &lib.Signin{}
	if err := unmarshal(poster, pack, signin); err != nil {
		return err
	}

	account, err := storage.GetAccountByUN(signin.Auth.Username)
	if err != nil {
		return poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Acc_Not_Exist.Val()})
	}

	passhashAndBcrypt, err := base64.StdEncoding.DecodeString(account.PasshashAndBcrypt)
	if err != nil {
		return poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Base64_Decode.Val()})
	}

	if err := bcrypt.CompareHashAndPassword(passhashAndBcrypt, signin.Auth.Passhash); err != nil {
		return poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Bcrypt_Compare.Val()})
	}

	return handleAuth(poster, pack, account, accId, accUN)
}

// 处理 token 验证请求
func validateToken(poster lib.Post, storage *Storage, pack *lib.Packet, accId *uint64, accUN *string) error {
	tokenReq := &lib.Token{}
	if err := unmarshal(poster, pack, tokenReq); err != nil {
		return err
	}

	id, expired, err := ParseToken(tokenReq.Token)
	if err != nil {
		return poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Parse_Token.Val()})
	}

	if expired {
		return poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Token_Expired.Val()})
	}

	account, err := storage.GetAccountById(id)
	if err != nil {
		return poster.Handle(pack, &lib.TokenRes{Code: lib.Err_Acc_Not_Exist.Val()})
	}

	return handleAuth(poster, pack, account, accId, accUN)
}

// 处理获取用户列表请求
func users(poster lib.Post, storage *Storage, pack *lib.Packet, accUN *string) error {
	if len(*accUN) == 0 {
		return poster.Handle(pack, &lib.UsersRes{Code: lib.Err_Forbidden.Val()})
	}

	users, err := storage.GetUsers(*accUN)
	if err != nil {
		return poster.Handle(pack, &lib.UsersRes{Code: lib.Err_Get_Users.Val()})
	}

	return poster.Handle(pack, &lib.UsersRes{Users: users})
}

// 处理发送消息请求
func recvMsg(poster lib.Post, storage *Storage, pack *lib.Packet, accUN *string, node *snowflake.Node) error {
	if len(*accUN) == 0 {
		return poster.Send(&lib.ErrRes{Code: lib.Err_Forbidden.Val()})
	}

	msg := &lib.Msg{}
	if err := lib.Unmarshal(pack.Data, msg); err != nil {
		return poster.Send(&lib.ErrRes{Code: lib.Err_Unmarshal.Val()})
	}

	return storage.NewMsg(&Message{
		// 生成消息 ID
		Id:   int64(node.Generate()),
		Kind: uint32(lib.MsgKind_TEXT),
		From: msg.From,
		To:   msg.To,
		Data: msg.Data,
	})
}

// 读取并处理客户端发送的 packet
func recvFrom(conn net.Conn, poster lib.Post, storage *Storage, accId *uint64, accUN *string, node *snowflake.Node) {
	lib.RecvFrom(
		conn,
		func(pack *lib.Packet) (err error) {
			switch pack.Kind {
			case lib.PackKind_SIGNUP:
				err = signup(poster, storage, pack, accId, accUN)
			case lib.PackKind_SIGNIN:
				err = signin(poster, storage, pack, accId, accUN)
			case lib.PackKind_TOKEN:
				err = validateToken(poster, storage, pack, accId, accUN)
			case lib.PackKind_USERS:
				err = users(poster, storage, pack, accUN)
			case lib.PackKind_MSG:
				err = recvMsg(poster, storage, pack, accUN, node)
			}
			return
		})
}

// 读取并转发所有发送给 accUN 用户的未读消息
func forwardMsgs(poster lib.Post, storage *Storage, accUN *string) {
	// 间隔 100ms 检查是否有新消息
	interval := time.NewTicker(100 * time.Millisecond)
	defer interval.Stop()

	for range interval.C {
		// accUN 用户已登录客户端
		if len(*accUN) > 0 {
			// 查询所有发送给 accUN 的未读消息
			msgList, _ := storage.GetMsgList(*accUN)
			for i := range msgList {
				// 转发消息
				poster.Send(&lib.Msg{
					Id:   msgList[i].Id,
					Kind: lib.MsgKind(msgList[i].Kind),
					From: msgList[i].From,
					To:   msgList[i].To,
					Data: msgList[i].Data,
				})
			}
		}
	}
}

func main() {
	// 启动单独协程，监听 ctrl+c 或 kill 信号，收到信号结束进程
	go lib.SignalHandler()

	// 初始化存储
	storage, err := new(Storage).Init(filepath.Join(lib.WorkDir, "server.db"))
	lib.FatalNotNil(err)

	// tcp 监听地址 0.0.0.0:8888
	addr := ":8888"

	// tcp 监听
	ln, err := net.Listen("tcp", addr)
	// tcp 监听遇到错误退出进程
	lib.FatalNotNil(err)
	// 输出日志
	lib.LogMessage("Listening on", addr)

	// 创建 snowflake Node
	node, err := snowflake.NewNode(1)
	lib.FatalNotNil(err)

	// 循环接受客户端连接
	for {
		// 每当有客户端连接时，ln.Accept 会返回新的连接 conn
		conn, err := ln.Accept()
		// 如果接受的新连接遇到错误，则退出进程
		lib.FatalNotNil(err)

		// 当前连接所登录的用户
		var (
			accId uint64
			accUN string
		)

		// 通过该 channel 可向当前连接发送 packet
		packChan := make(chan *lib.Packet)
		poster := newPoster(packChan)

		// 为每个客户端启动一个协程，把来自 packChan 的 packet 都发送到 conn
		go sendTo(conn, packChan)

		// 为每个客户端启动一个协程，读取并处理客户端发送的 packet
		go recvFrom(conn, poster, storage, &accId, &accUN, node)

		// 为每个客户端启动一个协程，读取并转发所有发送给当前用户的未读消息
		go forwardMsgs(poster, storage, &accUN)
	}
}
