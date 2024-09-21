package commands

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type Command interface {
	Run(ctx context.Context, args ...string) error
	GetName() string
}

type CommandApp struct {
	commands map[string]Command
	log      *log.Helper
}

func NewCommandApp(logger log.Logger, cmds ...Command) *CommandApp {
	coms := &CommandApp{
		commands: make(map[string]Command),
		log:      log.NewHelper(logger),
	}
	for _, cmd := range cmds {
		coms.commands[cmd.GetName()] = cmd
	}
	return coms
}

func (c *CommandApp) Run(commandLine string) error {
	args := strings.Split(commandLine, ":")
	c.log.Info("Command start!")
	name := args[0]
	if cmd, ok := c.commands[name]; ok {
		return cmd.Run(context.Background(), args[1:]...)
	}
	return fmt.Errorf("command %s not found", name)
}
