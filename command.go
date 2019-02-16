package cli

import (
	"strings"
)

type commands map[string]*Command

type commandHandler func(c Command)

// A command for your program, this can either be a main command, regular command or a sub-command.
type Command struct {
	parent  *Command
	cmds    commands
	handler commandHandler

	// The name of the command that will be passed to the program for command invocation.
	Name  string

	// The arguments passed to the individual command. Starting out this will be empty, and will
	// be populated once the command's flags have been parsed from the programs input arguments.
	Args  args

	// The flags that have been added to the command. Starting out this will be empty much like the
	// command arguments. During parsing the command flags will be set accordingly, whether they be
	// value flags or boolean flags.
	Flags flags
}

// Create a new command with the given name, and given command handler.
func (c *Command) Command(name string, handler commandHandler) *Command {
	return addCommand(name, c, handler, c.cmds)
}

// Add a flag to the current command. If the given flag is a global program flag, then set it on
// each of the command's sub-commands.
func (c *Command) AddFlag(f *Flag) *Command {
	if f.global {
		for _, cmd := range c.cmds {
			cmd.AddFlag(f)
		}
	}

	c.Flags.expected[f.Name] = f

	return c
}

// Get the full name of the current command. If the current command is a sub-command, then the
// command's name, and all grandparent command names will be concatenated into a hyphenated string.
func (c Command) FullName() string {
	commands := make([]string, 0)

	next := &c

	for next != nil {
		commands = append([]string{next.Name}, commands...)
		next = next.parent
	}

	return strings.Join(commands, "-")
}

// Invoke the command's handler, and any flag handlers that are set on the flags on the command.
// If any of the flags on the command are exclusive then the command handle will not be invoked.
func (c Command) Run() {
	shouldRun := true

	for _, flags := range c.Flags.received {
		flag := flags[0]

		if flag.Handler != nil && flag.isSet {
			flag.Handler(flag, c)
			shouldRun = !flag.Exclusive
		}
	}

	if shouldRun {
		c.handler(c)
	}
}
