package gipv4

import (
	"fmt"
	"testing"
)

func Test_GetIpArray(t *testing.T) {
	ips, err := GetIpArray()
	fmt.Println(ips)
	if err != nil {
		t.Fail()
	}
}

func Test_GetIntranetIp(t *testing.T) {
	ip, err := GetIntranetIp()
	fmt.Println(ip)
	if err != nil {
		t.Fail()
	}
}

func Test_GetIntranetIpArray(t *testing.T) {
	ips, err := GetIntranetIpArray()
	fmt.Println(ips)
	if err != nil {
		t.Fail()
	}
}

func Test_IsIntranet(t *testing.T) {
	if IsIntranet("127.0.0.1") == false {
		t.Fail()
	}
}
