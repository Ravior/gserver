package gsys

import (
	"fmt"
	"github.com/Ravior/gserver/core/os/gfile"
	"runtime"
	"time"
)

//程序运行时间
var startTime = time.Now()

var SysInfo struct {
	Uptime       int // 系统运行时长
	NumGoroutine int

	MemAllocated string // 当前内存使用量
	MemTotal     string // 所有被分配的内存
	MemSys       string // 内存占用量
	Lookups      uint64 // 指针查找次数
	MemMallocs   uint64 // 内存分配次数
	MemFrees     uint64 // 内存释放次数

	NextGC       string // 下次 GC 内存回收量
	LastGC       string // 距离上次 GC 时间
	PauseTotalNs string // 暂停时间总量
	PauseNs      string // 上次 GC 暂停时间
	NumGC        uint32 // 执行次数
}

func GetInfo() {
	SysInfo.Uptime = time.Now().Second() - startTime.Second()
	SysInfo.NumGoroutine = runtime.NumGoroutine()

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	SysInfo.MemAllocated = gfile.FormatSize(int64(m.Alloc))
	SysInfo.MemTotal = gfile.FormatSize(int64(m.TotalAlloc))
	SysInfo.MemSys = gfile.FormatSize(int64(m.Sys))
	SysInfo.Lookups = m.Lookups
	SysInfo.MemMallocs = m.Mallocs
	SysInfo.MemFrees = m.Frees

	SysInfo.NextGC = gfile.FormatSize(int64(m.NextGC))
	SysInfo.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	SysInfo.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	SysInfo.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	SysInfo.NumGC = m.NumGC
}
