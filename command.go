package cli

import (
	"strings"
)

type commands map[string]*Command

type commandHandler func(c Command)

type Command struct {
	parent  *Command
	cmds    commands
	handler commandHandler

	Name  string
	Args  args
	Flags flags
}

func (c *Command) Command(name string, handler commandHandler) *Command {
	return addCommand(name, c, handler, c.cmds)
}

func (c *Command) AddFlag(f *Flag) *Command {
	if f.global {
		for _, cmd := range c.cmds {
			cmd.AddFlag(f)
		}
	}

	c.Flags.expected[f.Name] = f

	return c
}

func (c Command) FullName() string {
	commands := make([]string, 0)

	next := &c

	for next != nil {
		commands = append([]string{next.Name}, commands...)
		next = next.parent
	}

	return strings.Join(commands, "-")
}

func (c Command) Run() {
	shouldRun := true

	for _, flags := range c.Flags.received {
		flag := flags[0]

		if flag.handler != nil && flag.isSet {
			flag.handler(flag, c)
			shouldRun = !flag.Exclusive
		}
	}

	if shouldRun {
		c.handler(c)
	}
}
