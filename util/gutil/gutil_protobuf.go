package gutil

import (
	"fmt"
	"github.com/Ravior/gserver/crypto/gcrc32"
	"github.com/Ravior/gserver/os/glog"
	"github.com/Ravior/gserver/util/gserialize"
	"github.com/golang/protobuf/proto"
)

// ProtoMarshal 获取MsgId/MsgData
func ProtoMarshal(msg proto.Message) (uint32, []byte) {
	defer func() {
		if err := recover(); err != nil {
			e := fmt.Sprintf("%v", err)
			glog.Errorf("ProtoMarshal has error:%v", e)
		}
	}()

	msgName := gserialize.Protobuf.GetMessageName(msg)
	if msgData, err := gserialize.Protobuf.Marshal(msg); err == nil {
		return gcrc32.EncryptString(msgName), msgData
	}
	return 0, nil
}
