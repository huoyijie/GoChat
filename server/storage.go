package main

import (
	"log"
	"os"
	"time"

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
}

type Message struct {
	Id   int64 `gorm:"primaryKey;autoIncrement:false"`
	Kind int32
	From string
	To   string
	Data []byte
}

type Storage struct {
	db *gorm.DB
}

func (s *Storage) Init(filePath string) (*Storage, error) {
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

func (s *Storage) NewAccount(account *Account) (err error) {
	err = s.db.Create(account).Error
	return
}

func (s *Storage) GetAccountById(id uint64) (account *Account, err error) {
	account = &Account{Id: id}
	err = s.db.First(account).Error
	return
}

func (s *Storage) GetAccountByUN(username string) (account *Account, err error) {
	account = &Account{Username: username}
	err = s.db.Where(account).First(account).Error
	return
}

func (s *Storage) GetUsers(self string) (users []string, err error) {
	var accounts []Account
	err = s.db.Select("username").Find(&accounts).Order("username").Error
	if err != nil {
		return
	}
	users = make([]string, 0, len(accounts)-1)
	for i := range accounts {
		if accounts[i].Username != self {
			users = append(users, accounts[i].Username)
		}
	}
	return
}

func (s *Storage) NewMsg(msg *Message) (err error) {
	err = s.db.Create(msg).Error
	return
}

func (s *Storage) GetMsgList(to string) (msgList []Message, err error) {
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("`to` = ?", to).Find(&msgList).Order("id").Error; err != nil {
			return err
		}

		if err := tx.Where("`to` = ?", to).Delete(&Message{}).Error; err != nil {
			return err
		}

		return nil
	})
	return
}
