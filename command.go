package cli

type commands map[string]*Command

type commandHandler func(c Command)

type Command struct {
	shouldRun bool

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
	if f.global {
		for _, cmd := range c.Commands {
			cmd.AddFlag(f)
		}
	}

	c.Flags[f.Name] = f
}
