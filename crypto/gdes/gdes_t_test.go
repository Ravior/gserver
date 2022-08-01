package gdes

import (
	"fmt"
	"testing"
)

func Test_EncryptECB(t *testing.T) {
	sourceStr := "1234567812345678"
	key := []byte("12345678")
	encryptStr, err := EncryptECB([]byte(sourceStr), key, NOPADDING)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	decryptStr, _ := DecryptECB(encryptStr, key, NOPADDING)
	if sourceStr != string(decryptStr) {
		t.Fail()
	}
}
