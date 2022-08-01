package gipv4

import (
	"fmt"
	"testing"
)

func Test_GetMac(t *testing.T) {
	mac, _ := GetMac()
	fmt.Println(mac)
	if mac == "" {
		t.Fail()
	}
}

func Test_GetMacArray(t *testing.T) {
	macs, _ := GetMacArray()
	fmt.Println(macs)
	if len(macs) == 0 {
		t.Fail()
	}
}
