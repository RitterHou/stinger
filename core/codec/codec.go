package codec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

const key32Bytes = "23333333333333333333333333333333"

var key = []byte(key32Bytes)

// 更新key的值
func SetKey(k string) {
	kLen := len(k)
	if kLen < 32 {
		k += key32Bytes[:32-kLen]
	} else if kLen > 32 {
		k = k[:32]
	}
	key = []byte(k)
}

func Encrypt(source []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	cipherText := make([]byte, aes.BlockSize+len(source))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText[aes.BlockSize:], source)
	return cipherText
}

func Decrypt(source []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	if len(source) < aes.BlockSize {
		log.Fatal("cipher text too short")
	}
	iv := source[:aes.BlockSize]
	source = source[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(source, source)
	return source
}
