package guuid

import (
	"bytes"
	"github.com/Ravior/gserver/util/gconfig"
	"github.com/Ravior/gserver/util/gutil"
	"strconv"
	"sync/atomic"
	"time"
)

// 起始UUID
var startUuid = time.Now().Unix()

// salt
var salt = gutil.RandSeq(6)

func GetUUID() string {
	uuid := atomic.AddInt64(&startUuid, 1)
	if uuid <= 1 {
		startUuid = time.Now().Unix()
		uuid = atomic.AddInt64(&startUuid, 1)
	}
	return buildStr(uuid)
}

func buildStr(uuid int64) string {
	var bt bytes.Buffer
	bt.WriteString(gconfig.Global.ServerId)
	bt.WriteString(salt)
	bt.WriteString(strconv.FormatInt(uuid, 10))
	return bt.String()
}
