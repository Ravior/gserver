package gcrc32

import (
	"log"
	"testing"
)

func Test_Encrypt(t *testing.T) {
	crc1 := Encrypt("12345")
	log.Println(crc1)

	crc2 := EncryptString("12345")
	log.Println(crc2)

	if crc1 != crc2 {
		t.Fail()
	}
}
