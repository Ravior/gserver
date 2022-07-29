package guuid

import (
	"fmt"
	"github.com/Ravior/gserver/util/gconfig"
	"github.com/Ravior/gserver/util/gtime"
	"github.com/Ravior/gserver/util/gutil"
	"sync/atomic"
)

// 起始UUID
var startUuid = gtime.Now()

func GetUUID() string {
	uuid := atomic.AddInt64(&startUuid, 1)
	if uuid <= 1 {
		startUuid = gtime.Now()
		uuid = atomic.AddInt64(&startUuid, 1)
	}
	salt := gutil.RandSeq(6)
	return fmt.Sprintf("%s%d%s", gconfig.Global.ServerId, uuid, salt)
}
