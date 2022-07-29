package command

import (
	"sync"
)

var Mgr *commandMgr

func init() {
	Mgr = NewCommandMgr()
}

type Command interface {
	Name() string
	Help() string
	Run(args []string) string
}

type commandMgr struct {
	sync.RWMutex
	commands []Command
}

func (c *commandMgr) AddCommand(command Command) {
	c.Lock()
	defer c.Unlock()

	c.commands = append(c.commands, command)
}

func (c *commandMgr) GetCommands() []Command {
	c.RLock()
	defer c.RUnlock()

	return c.commands
}

func NewCommandMgr() *commandMgr {
	mgr := &commandMgr{
		commands: make([]Command, 0),
	}
	// 系统命令
	mgr.commands = append(mgr.commands, &CommandHelp{}, &CommandSysinfo{}, &CommandProf{}, &CommandCPUProf{})
	return mgr
}
