package command

type CommandHelp struct{}

func (c *CommandHelp) Name() string {
	return "help"
}

func (c *CommandHelp) Help() string {
	return "获取帮助信息"
}

func (c *CommandHelp) Run([]string) string {
	output := "当前系统支持命令如下:\r\n"
	output += "-------------------------------------------------------\r\n"
	for _, c := range Mgr.GetCommands() {
		output += "[] " + c.Name() + " - " + c.Help() + "\r\n"
	}
	output += "[] exit - 退出控制台\r\n"
	output += "-------------------------------------------------------"

	return output
}
