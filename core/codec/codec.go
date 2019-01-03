package codec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

var key = []byte("2333333333333333")

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
