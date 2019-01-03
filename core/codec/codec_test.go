package codec

import "testing"

func TestEncryptAndDecrypt(t *testing.T) {
	source := []byte("你好，中国人😊")
	secure := Encrypt(source)
	result := Decrypt(secure)
	t.Log(string(result))
}
