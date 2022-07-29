package command

import (
	"fmt"
	"github.com/Ravior/gserver/core/os/gsys"
)

type CommandSysinfo struct{}

func (c *CommandSysinfo) Name() string {
	return "sysinfo"
}

func (c *CommandSysinfo) Help() string {
	return "获取当心系统运行状态信息"
}

func (c *CommandSysinfo) Run(args []string) string {
	gsys.GetInfo()
	output := "---------------- 系统运行信息--------------------------\r\n"
	output += fmt.Sprintf("运行时间：%d\r\n", gsys.SysInfo.Uptime)
	output += fmt.Sprintf("协程数量：%d\r\n", gsys.SysInfo.NumGoroutine)
	output += fmt.Sprintf("服务器内存总量：%s\r\n", gsys.SysInfo.MemTotal)
	output += fmt.Sprintf("当前内存申请量：%s\r\n", gsys.SysInfo.MemAllocated)
	output += fmt.Sprintf("已使用内存量：%s\r\n", gsys.SysInfo.MemSys)
	output += fmt.Sprintf("指针查找次数：%d\r\n", gsys.SysInfo.Lookups)
	output += fmt.Sprintf("内存分配次数：%d\r\n", gsys.SysInfo.MemMallocs)
	output += fmt.Sprintf("下次GC内存回收量：%s\r\n", gsys.SysInfo.NextGC)
	output += fmt.Sprintf("距离上次GC时间：%s\r\n", gsys.SysInfo.LastGC)
	output += fmt.Sprintf("暂停时间总量：%s\r\n", gsys.SysInfo.PauseTotalNs)
	output += fmt.Sprintf("上次GC暂停时间：%s\r\n", gsys.SysInfo.PauseNs)
	output += fmt.Sprintf("执行次数：%d\r\n", gsys.SysInfo.NumGC)
	output += "-------------------------------------------------------"
	return output
}
