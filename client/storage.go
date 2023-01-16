package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/huoyijie/GoChat/lib"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var gormLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags),
	logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Silent,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	},
)

// 存储 key/Value 表
type KeyValue struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

// 本地聊天记录
type Message struct {
	Id   int64 `gorm:"primaryKey;autoIncrement:false"`
	Kind int32
	From string
	Data []byte
	Read bool
}

// 客户端本地存储
type storage_t struct {
	db *gorm.DB
}

func (s *storage_t) Init(filePath string) (*storage_t, error) {
	// 创建并打开数据库存储文件
	if db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{Logger: gormLogger}); err != nil {
		return nil, err
	} else {
		s.db = db
		if err := s.db.Transaction(func(tx *gorm.DB) error {
			// 自动根据模型更新表结构
			var kv KeyValue
			var message Message
			if err := tx.AutoMigrate(&kv, &message); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return nil, err
		}
		return s, nil
	}
}

// 批量插入 key/value 键值对
func (s *storage_t) NewKVS(kvs []KeyValue) (err error) {
	err = s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&kvs).Error
	return
}

// 根据 key 查询 value
func (s *storage_t) GetValue(key string) (kv *KeyValue, err error) {
	kv = &KeyValue{Key: key}
	err = s.db.First(kv).Error
	return
}

// 存储 token 对象
func (s *storage_t) StoreToken(tokenRes *lib.TokenRes) (err error) {
	err = s.NewKVS([]KeyValue{
		{Key: "id", Value: fmt.Sprintf("%d", tokenRes.Id)},
		{Key: "username", Value: tokenRes.Username},
		{Key: "token", Value: base64.StdEncoding.EncodeToString(tokenRes.Token)},
	})
	return
}

// 收到新未读消息
func (s *storage_t) NewMsg(msg *Message) (err error) {
	err = s.db.Create(msg).Error
	return
}

// 获取某个用户发给自己的未读消息列表
func (s *storage_t) GetMsgList(from string) (msgList []Message, err error) {
	err = s.db.Transaction(func(tx *gorm.DB) error {
		msg := &Message{From: from}
		res := tx.Model(msg).Where(msg).Update("read", true)
		if err := res.Error; err != nil {
			return err
		}

		if unReadMsgCnt := res.RowsAffected; unReadMsgCnt > 0 {
			msg.Read = true
			if err := tx.Where(msg).Find(&msgList).Order("id").Error; err != nil {
				return err
			}

			if err := tx.Where(msg).Delete(msg).Error; err != nil {
				msgList = nil
				return err
			}
		}

		return nil
	})
	return
}

// 获取当前登录用户的未读消息数量
func (s *storage_t) UnReadMsgCount() (msgCount map[string]uint32, err error) {
	rows, err := s.db.Model(&Message{}).Select("from", "COUNT(*) as count").Group("from").Rows()
	if err != nil {
		return
	}

	defer rows.Close()
	msgCount = make(map[string]uint32)
	for rows.Next() {
		var (
			from  string
			count uint32
		)
		if err = rows.Scan(&from, &count); err != nil {
			continue
		}
		msgCount[from] = count
	}
	return
}
