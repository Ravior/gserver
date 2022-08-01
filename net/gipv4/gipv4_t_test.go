package gipv4

import (
	"fmt"
	"testing"
)

func Test_Ip2long(t *testing.T) {
	fmt.Println(Ip2long("127.0.0.1"))
	if Ip2long("127.0.0.1") == 0 {
		t.Fail()
	}
}

func Test_Long2Ip(t *testing.T) {
	ipInt := Ip2long("127.0.0.1")
	if Long2ip(ipInt) != "127.0.0.1" {
		t.Fail()
	}
}

func Test_Validate(t *testing.T) {
	if Validate("127.0.0.1") == false {
		t.Fail()
	}
}

func Test_ParseAddress(t *testing.T) {
	ip, port := ParseAddress("127.0.0.1:8080")
	if ip != "127.0.0.1" || port != 8080 {
		t.Fail()
	}
}
