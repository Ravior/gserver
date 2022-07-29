package cmd

import (
	"fmt"
	"strings"
)

type ICommand interface {
	Run([]string) string
	Help() string
	Name() string
}

type ICommandInterpreter interface {
	AddCommand(ICommand)
	Execute(string) string
	IsQuitCmd(string) bool
}

var (
	QUIT_CMD = [3]string{"quit", "q", "exit"}
)

type CommandInterpreter struct {
	commands map[string]ICommand
}

func NewCommandInterpreter() *CommandInterpreter {
	interpreter := &CommandInterpreter{make(map[string]ICommand)}
	return interpreter
}

func (c *CommandInterpreter) AddCommand(cmd ICommand) {
	c.commands[cmd.Name()] = cmd
	fmt.Println("add command ", cmd.Name())
}

func (c *CommandInterpreter) preExecute(rawCmdExp string) string {
	return strings.ToLower(strings.TrimSpace(rawCmdExp))
}

func (c *CommandInterpreter) IsQuitCmd(rawCmdExp string) bool {
	cmdExp := c.preExecute(rawCmdExp)
	for _, cmd := range QUIT_CMD {
		if cmd == cmdExp {
			return true
		}
	}
	return false
}

func (c *CommandInterpreter) help() string {
	helpStr := "有关某个命令的详细信息，请键入 help 命令名"
	for _, v := range c.commands {
		helpStr = fmt.Sprintf("%s\r\n%s", helpStr, v.Help())
	}
	return helpStr
}

func (c *CommandInterpreter) Execute(rawCmdExp string) string {
	defer func() string {
		if err := recover(); err != nil {
			fmt.Println("invalid rawCmdExp: ", rawCmdExp)
			return "invalid rawCmdExp: " + rawCmdExp
		}
		return "Unkown ERROR!!!"
	}()
	if rawCmdExp == "" {
		return ""
	}
	rawCmdExps := strings.Split(rawCmdExp, " ")
	if len(rawCmdExps) == 0 {
		return ""
	}
	cmdExps := make([]string, 0)
	for _, cmd := range rawCmdExps {
		cmdExps = append(cmdExps, c.preExecute(cmd))
	}

	if command, ok := c.commands[cmdExps[0]]; ok {
		return command.Run(cmdExps[1:])
	} else {
		if cmdExps[0] == "help" {
			return c.help()
		} else {
			return "command not found."
		}
	}
}
