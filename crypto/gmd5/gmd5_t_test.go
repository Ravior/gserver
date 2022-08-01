package gmd5

import (
	"log"
	"testing"
)

func Test_Encrypt(t *testing.T) {
	md5Str, _ := Encrypt("123456")
	log.Println(md5Str)

	md5Str1 := MustEncrypt("123456")
	log.Println(md5Str1)

	md5Str2, _ := EncryptString("123456")
	log.Println(md5Str2)

	if md5Str != md5Str1 || md5Str1 != md5Str2 {
		t.Fail()
	}
}
