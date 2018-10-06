package cli

import "errors"

type commands map[string]*Command

type commandHandler func(c Command)

type Command struct {
	hasExclusive bool

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
	return addCommand(name, handler, c.Commands)
}

func (c *Command) AddFlag(f *Flag) {
	c.hasExclusive = f.Exclusive && f.Handler != nil

	if f.global {
		for _, cmd := range c.Commands {
			cmd.AddFlag(f)
		}
	}

	c.Flags[f.Name] = f
}

func (c Command) Run(nilHandler commandHandler) error {
	shouldRun := true

	for _,f := range c.Flags {
		if f.Handler != nil && f.isSet {
			f.Handler(*f, c)

			shouldRun = !f.Exclusive
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
