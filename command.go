// Simple library for building CLI applications in Go
package cli

import (
	"errors"
	"strings"
)

type commands map[string]*Command

type commandHandler func(c Command)

type Command struct {
	hasExclusive bool

	Parent *Command

	Name string

	Args args

	Flags flags

	Commands commands

	Handler commandHandler
}

func newCommands() commands {
	return commands(make(map[string]*Command))
}

func (c *Command) Command(name string, handler commandHandler) *Command {
	return addCommand(name, c, handler, c.Commands)
}

func (c *Command) AddFlag(f *Flag) {
	c.hasExclusive = f.Exclusive && f.Handler != nil

	if f.global {
		for _, cmd := range c.Commands {
			cmd.AddFlag(f)
		}
	}

	c.Flags.expected[f.Name] = f
}

func (c Command) FullName() string {
	commands := make([]string, 0)

	next := &c

	for next != nil {
		commands = append([]string{next.Name}, commands...)

		next = next.Parent
	}

	return strings.Join(commands, "-")
}

func (c Command) Run(nilHandler commandHandler) error {
	shouldRun := true

	for _, received := range c.Flags.received {
		flag := received[0]

		if flag.Handler != nil && flag.isSet {
			flag.Handler(flag, c)

			shouldRun = !flag.Exclusive
		}
	}

	if shouldRun {
		if c.Handler == nil {
			c.Handler = nilHandler
		}

		if c.Handler == nil {
			return errors.New("no command to run")
		}

		c.Handler(c)
	}

	return nil
}
