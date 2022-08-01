package gsha1

import (
	"fmt"
	"testing"
)

func Test_Encrypt(t *testing.T) {
	encryptStr := Encrypt("12345678")
	fmt.Println(encryptStr)
	if encryptStr == "" {
		t.Fail()
	}
}
