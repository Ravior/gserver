package gnet

import (
	"testing"
)

func Test_NewMsg(t *testing.T) {
	msg := NewMsg(1, []byte("gserver"))
	if msg.GetMsgId() != 1 || msg.GetDataLen() != 7 {
		t.Fail()
	}
}
