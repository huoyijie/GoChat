package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func randNonce(nonceSize int) []byte {
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	return nonce
}

func NewGCM(key *[32]byte) (gcm cipher.AEAD, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return
	}
	gcm, err = cipher.NewGCM(block)
	return
}

func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		panic(err)
	}
	return &key
}

func Encrypt(plaintext []byte, gcm cipher.AEAD) []byte {
	// 随机生成字节 slice，使得每次的加密结果具有随机性
	nonce := randNonce(gcm.NonceSize())
	// Seal 方法第一个参数 nonce，会把 nonce 本身加入到加密结果
	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func Decrypt(ciphertext []byte, gcm cipher.AEAD) ([]byte, error) {
	// 首先得到加密时使用的 nonce
	nonce := ciphertext[:gcm.NonceSize()]
	// 传入 nonce 并进行数据解密
	return gcm.Open(nil, nonce, ciphertext[gcm.NonceSize():], nil)
}
