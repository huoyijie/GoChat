package main

import (
	"encoding/binary"
	"encoding/hex"
	"log"
	"os"
	"time"

	"github.com/huoyijie/GoChat/lib"
)

const (
	DEFAULT_SECRET_KEY = "3e367a60ddc0699ea2f486717d5dcd174c4dee0bcf1855065ab74c348e550b78"
)

func secretKey() (secretKey string) {
	secretKey, found := os.LookupEnv("SECRET_KEY")
	if !found {
		secretKey = DEFAULT_SECRET_KEY
	}
	return
}

func GetSecretKey() *[32]byte {
	key, err := hex.DecodeString(secretKey())
	if err != nil {
		log.Fatal(err)
	}
	return (*[32]byte)(key)
}

func GenerateToken(id uint64) (token []byte, err error) {
	gcm, err := lib.NewGCM(GetSecretKey())
	if err != nil {
		return
	}

	bytes := make([]byte, 16)
	binary.BigEndian.PutUint64(bytes, id)
	binary.BigEndian.PutUint64(bytes[8:], uint64(time.Now().Unix()))
	token = lib.Encrypt(bytes, gcm)
	return
}

func ParseToken(token []byte) (id uint64, expired bool, err error) {
	gcm, err := lib.NewGCM(GetSecretKey())
	if err != nil {
		return
	}

	bytes, err := lib.Decrypt(token, gcm)
	if err != nil {
		return
	}

	id = binary.BigEndian.Uint64(bytes)
	genTime := binary.BigEndian.Uint64(bytes[8:])
	expired = time.Since(time.Unix(int64(genTime), 0)) > 30*24*time.Hour
	return
}
