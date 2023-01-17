package main

import (
	"log"
	"os"
	"time"

	"github.com/huoyijie/GoChat/lib"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

type Account struct {
	Id                uint64 `gorm:"primaryKey"`
	Username          string `gorm:"uniqueIndex:uni_username"`
	PasshashAndBcrypt string
	Online            bool
}

type Message struct {
	Id       int64 `gorm:"primaryKey;autoIncrement:false"`
	Kind     int32
	From, To string
	Data     []byte
	Read     bool
}

type storage_t struct {
	db *gorm.DB
}

func (s *storage_t) Init(filePath string) (*storage_t, error) {
	if db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{Logger: gormLogger}); err != nil {
		return nil, err
	} else {
		s.db = db
		if err := s.db.Transaction(func(tx *gorm.DB) error {
			var account Account
			var msg Message
			if err := tx.AutoMigrate(&account, &msg); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return nil, err
		}
		return s, nil
	}
}

func (s *storage_t) NewAccount(account *Account) (err error) {
	err = s.db.Create(account).Error
	return
}

func (s *storage_t) GetAccountById(id uint64) (account *Account, err error) {
	account = &Account{Id: id}
	err = s.db.First(account).Error
	return
}

func (s *storage_t) GetAccountByUN(username string) (account *Account, err error) {
	account = &Account{Username: username}
	err = s.db.Where(account).First(account).Error
	return
}

func (s *storage_t) UpdateOnline(id uint64, online bool) (err error) {
	err = s.db.Model(&Account{Id: id}).Update("online", online).Error
	return
}

func (s *storage_t) GetUsers(self string) (users []*lib.User, err error) {
	var accounts []Account
	err = s.db.Select("username", "online").Order("username").Find(&accounts).Error
	if err != nil {
		return
	}

	users = make([]*lib.User, 0, len(accounts)-1)
	for i := range accounts {
		if accounts[i].Username != self {
			user := &lib.User{Username: accounts[i].Username, Online: accounts[i].Online}
			users = append(users, user)
		}
	}
	return
}

func (s *storage_t) NewMsg(msg *Message) (err error) {
	err = s.db.Create(msg).Error
	return
}

func (s *storage_t) GetMsgList(to string) (msgList []Message, err error) {
	err = s.db.Transaction(func(tx *gorm.DB) error {
		msg := &Message{To: to}

		res := tx.Model(msg).Where(msg).Update("read", true)
		if err := res.Error; err != nil {
			return err
		}

		if unReadMsgCnt := res.RowsAffected; unReadMsgCnt > 0 {
			msg.Read = true
			if err := tx.Where(msg).Order("id").Find(&msgList).Error; err != nil {
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
