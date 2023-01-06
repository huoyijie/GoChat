package main

import (
	"encoding/base64"
	"net"
	"path/filepath"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/huoyijie/GoChat/lib"
	"golang.org/x/crypto/bcrypt"
)

// 存储当前所有客户端连接
var sockets = make(map[snowflake.ID]*lib.Socket)

// 多个协程并发读写 sockets 时，需要使用读写锁
var lock sync.RWMutex

// 写锁
func wSockets(wSockets func()) {
	lock.Lock()
	defer lock.Unlock()
	wSockets()
}

// 读锁
func rSockets(rSockets func()) {
	lock.RLock()
	defer lock.RUnlock()
	rSockets()
}

func sendTo(conn net.Conn, packChan <-chan *lib.Packet) {
	var id uint64
	for pack := range packChan {
		if pack.Id == 0 {
			id++
			pack.Id = id
		}

		bytes, err := lib.MarshalPack(pack)
		if err == nil {
			_, err = conn.Write(bytes)
		}

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

		// 生成新 ID
		id := node.Generate()
		packChan := make(chan *lib.Packet)
		// 保存新连接
		wSockets(func() {
			sockets[id] = &lib.Socket{Id: id, PackChan: packChan}
		})

		go sendTo(conn, packChan)

		// 为每个客户端启动一个协程，读取客户端发送的消息并转发
		go lib.RecvFrom(
			conn,
			func(pack *lib.Packet) error {
				switch pack.Kind {
				case lib.PackKind_MSG:
					rSockets(func() {
						for k, v := range sockets {
							// 向其他所有客户端(除了自己)转发消息
							if k != id {
								v.PackChan <- &lib.Packet{
									Kind: pack.Kind,
									Data: pack.Data,
								}
							}
						}
					})
				case lib.PackKind_SIGNUP:
					signup := &lib.Signup{}
					if err := lib.Unmarshal(pack.Data, signup); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10000}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_SIGNUP,
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
								Kind: lib.PackKind_SIGNUP,
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
								Kind: lib.PackKind_SIGNUP,
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
							Kind: lib.PackKind_SIGNUP,
							Data: bytes,
						}
					}
				case lib.PackKind_TOKEN:
					token := &lib.Token{}
					if err := lib.Unmarshal(pack.Data, token); err != nil {
						if bytes, err := lib.Marshal(&lib.TokenRes{Code: -10003}); err != nil {
							return err
						} else {
							packChan <- &lib.Packet{
								Id:   pack.Id,
								Kind: lib.PackKind_TOKEN,
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
								Kind: lib.PackKind_TOKEN,
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
								Kind: lib.PackKind_TOKEN,
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
								Kind: lib.PackKind_TOKEN,
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
								Kind: lib.PackKind_TOKEN,
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
								Kind: lib.PackKind_TOKEN,
								Data: bytes,
							}
						}
					}
				}
				return nil
			},
			func() {
				// 从当前方法返回时，关闭连接
				conn.Close()
				// 删除连接
				wSockets(func() {
					delete(sockets, id)
				})
			})
	}
}
