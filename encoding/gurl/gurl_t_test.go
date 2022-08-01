package gurl

import (
	"fmt"
	"testing"
)

var (
	s  = "https://gitlib.com/?tag=gserver"
	s2 = "https%3A%2F%2Fgitlib.com%2F%3Ftag%3Dgserver"
)

func Test_Encode(t *testing.T) {
	s1 := Encode(s)
	fmt.Println(s1)
	if s1 != s2 {
		t.Fail()
	}
}

func Test_Decode(t *testing.T) {
	s1, _ := Decode(s2)
	fmt.Println(s1)
	if s1 != s {
		t.Fail()
	}
}

func Test_ParseURL(t *testing.T) {
	params, _ := ParseURL("https://gitlib.com/?tag=gserver&page=1", -1)
	if params["host"] != "gitlib.com" {
		t.Fail()
	}
}

func Test_BuildParams(t *testing.T) {
	params := map[string]string{
		"tag":  "gserver",
		"page": "1",
	}
	paramsStr := BuildParams(params)
	fmt.Println(paramsStr)
	if paramsStr == "" {
		t.Fail()
	}
}
