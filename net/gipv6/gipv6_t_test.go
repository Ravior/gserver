package gipv6

import "testing"

func Test_Validate(t *testing.T) {
	ipV6Str := "240e:45c:ce34:ef9:2f3:b70f:6888:d8c"
	if Validate(ipV6Str) == false {
		t.Fail()
	}
}
