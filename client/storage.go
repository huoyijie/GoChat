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

type KeyValue struct {
	Key   string `gorm:"primaryKey"`
	Value string
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
			var kv KeyValue
			if err := tx.AutoMigrate(&kv); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return nil, err
		}
		return s, nil
	}
}

func (s *Storage) NewKV(kv *KeyValue) (err error) {
	err = s.db.Create(kv).Error
	return
}

func (s *Storage) GetValue(key string) (kv *KeyValue, err error) {
	kv = &KeyValue{Key: key}
	err = s.db.First(kv).Error
	return
}

func (s *Storage) NewKVS(kvs []KeyValue) (err error) {
	err = s.db.Create(&kvs).Error
	return
}
