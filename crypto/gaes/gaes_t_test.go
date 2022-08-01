package gaes

import (
	"testing"
)

func Test_Encrypt(t *testing.T) {
	sourceStr := "1234567812345678"
	key := []byte("1234567812345678")
	iv := []byte("1234567812345678")
	encryptStr, err := Encrypt([]byte(sourceStr), key, iv)
	if err != nil {
		t.Fail()
	}
	decryptStr, _ := Decrypt(encryptStr, key, iv)
	if sourceStr != string(decryptStr) {
		t.Fail()
	}
}
