package main

import (
	"encoding/base64"
	"net"
	"path/filepath"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/huoyijie/GoChat/lib"
	"golang.org/x/crypto/bcrypt"
)

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

func main() {
	// 启动单独协程，监听 ctrl+c 或 kill 信号，收到信号结束进程
	go lib.SignalHandler()

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

		var (
			accId uint64
			accUN string
		)
		packChan := make(chan *lib.Packet)

		go sendTo(conn, packChan)

		// 为每个客户端启动一个协程，读取客户端发送的消息并转发
		go lib.RecvFrom(
			conn,
			func(pack *lib.Packet) error {
				switch pack.Kind {
				case lib.PackKind_SIGNUP:
					signup := &lib.Signup{}
					if err := lib.Unmarshal(pack.Data, signup); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10000}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					}

					passhashAndBcrypt, err := bcrypt.GenerateFromPassword(signup.Auth.Passhash, 14)
					if err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10001}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					}

					account := &Account{
						Username:          signup.Auth.Username,
						PasshashAndBcrypt: base64.StdEncoding.EncodeToString(passhashAndBcrypt),
					}
					if err := storage.NewAccount(account); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10002}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					}

					token, err := GenerateToken(account.Id)
					if err != nil {
						return err
					}

					if bytes, err := lib.Marshal(&lib.TokenRes{Id: account.Id, Username: account.Username, Token: token}); err != nil {
						return err
					} else {
						packChan <- &lib.Packet{
							Id:   pack.Id,
							Kind: lib.PackKind_RES,
							Data: bytes,
						}
						accId = account.Id
						accUN = account.Username
					}
				case lib.PackKind_TOKEN:
					token := &lib.Token{}
					if err := lib.Unmarshal(pack.Data, token); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10003}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					}

					if id, expired, err := ParseToken(token.Token); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10004}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					} else if expired {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10005}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					} else if account, err := storage.GetAccountById(id); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10006}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					} else if token, err := GenerateToken(account.Id); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10007}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					} else {
						if bytes, err := lib.Marshal(&lib.TokenRes{Id: account.Id, Username: account.Username, Token: token}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							accId = account.Id
							accUN = account.Username
						}
					}
				case lib.PackKind_SIGNIN:
					signin := &lib.Signin{}
					if err := lib.Unmarshal(pack.Data, signin); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10008}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					}

					if account, err := storage.GetAccountByUN(signin.Auth.Username); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10009}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					} else if passhashAndBcrypt, err := base64.StdEncoding.DecodeString(account.PasshashAndBcrypt); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10010}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					} else if err := bcrypt.CompareHashAndPassword(passhashAndBcrypt, signin.Auth.Passhash); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10011}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					} else if token, err := GenerateToken(account.Id); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10012}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					} else {
						if bytes, err := lib.Marshal(&lib.TokenRes{Id: account.Id, Username: account.Username, Token: token}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							accId = account.Id
							accUN = account.Username
						}
					}
				case lib.PackKind_USERS:
					if accId == 0 || len(accUN) == 0 {
						if bytes, err := lib.Marshal(&lib.UsersRes{Code: -10013}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					}

					if users, err := storage.GetUsers(accUN); err != nil {
						if bytes, err := lib.Marshal(&lib.UsersRes{Code: -10014}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
							return nil
						}
					} else {
						if bytes, err := lib.Marshal(&lib.UsersRes{Users: users}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_RES,
								Data: bytes,
							}
						}
					}
				case lib.PackKind_MSG:
					if accId == 0 || len(accUN) == 0 {
						if bytes, err := lib.Marshal(&lib.ErrRes{Code: -10015}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Kind: lib.PackKind_ERR,
								Data: bytes,
							}
							return nil
						}
					}

					msg := &lib.Msg{}
					if err := lib.Unmarshal(pack.Data, msg); err != nil {
						if bytes, err := lib.Marshal(&lib.ErrRes{Code: -10016}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Kind: lib.PackKind_ERR,
								Data: bytes,
							}
							return nil
						}
					}

					// 生成消息 ID
					if err := storage.NewMsg(&Message{
						Id:   int64(node.Generate()),
						Kind: uint32(lib.MsgKind_TEXT),
						From: msg.From,
						To:   msg.To,
						Data: msg.Data,
					}); err != nil {
						return err
					}
				}
				return nil
			})

		go func(accUN *string) {
			interval := time.NewTicker(100 * time.Millisecond)
			defer interval.Stop()

			for range interval.C {
				if len(*accUN) > 0 {
					msgList, _ := storage.GetMsgList(*accUN)
					for i := range msgList {
						msg := &lib.Msg{
							Id:   msgList[i].Id,
							Kind: lib.MsgKind(msgList[i].Kind),
							From: msgList[i].From,
							To:   msgList[i].To,
							Data: msgList[i].Data,
						}
						bytes, _ := lib.Marshal(msg)
						packChan <- &lib.Packet{
							Kind: lib.PackKind_MSG,
							Data: bytes,
						}
					}
				}
			}
		}(&accUN)
	}
}
