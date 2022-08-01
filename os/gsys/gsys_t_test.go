package gsys

import (
	"fmt"
	"testing"
)

func Test_GetInfo(t *testing.T) {
	GetInfo()
	fmt.Printf("%+v", SysInfo)
	if SysInfo.NumGoroutine == 0 {
		t.Fail()
	}
}
