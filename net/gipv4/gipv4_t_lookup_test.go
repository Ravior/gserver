package gipv4

import (
	"fmt"
	"testing"
)

func Test_GetHostByName(t *testing.T) {
	host, err := GetHostByName("gitlib.com")
	fmt.Println(host)
	if err != nil {
		t.Fail()
	}
}

func Test_GetHostsByName(t *testing.T) {
	hosts, err := GetHostsByName("gitlib.com")
	fmt.Println(hosts)
	if err != nil {
		t.Fail()
	}
}

func Test_GetNameByAddr(t *testing.T) {
	name, err := GetNameByAddr("127.0.0.1")
	fmt.Println(name)
	if err != nil {
		t.Fail()
	}
}
