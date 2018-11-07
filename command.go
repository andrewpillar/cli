package cli

import (
	"errors"
	"strings"
)

type commands map[string]*Command

type commandHandler func(c Command)

type Command struct {
	hasExclusive bool

	// The pointer to the command's parent if it is a sub-command. Otherwise
	// set to nil.
	Parent *Command

	// Name specifies the string the user inputs to execute the command. For
	// sub-commands this will not be the full string including the parent's
	// name, just the sub-command name itself.
	Name string

	// Args specifies the slice of arguments that are passed to the command.
	Args args

	// Flags specifies the flags that were set on the command.
	Flags flags

	// Commands specifies any sub-commands that the current command might have.
	Commands commands

	// Handler specifies the handler to call when the command is executed.
	Handler commandHandler
}

func newCommands() commands {
	return commands(make(map[string]*Command))
}

// Command creates a new command for the Command struct based on the nane, and
// handler given. Similar to the Command method on the Cli struct, only this
// creates a sub-command.
func (c *Command) Command(name string, handler commandHandler) *Command {
	return addCommand(name, c, handler, c.Commands)
}

// AddFlag takes a pointer to a Flag struct, and adds it to the Cli struct.
// Similar to the AddFlag method on the Cli struct, only the flags are
// contained to the Command struct itself, and not passed down to the
// sub-commands.
func (c *Command) AddFlag(f *Flag) {
	c.hasExclusive = f.Exclusive && f.Handler != nil

	if f.global {
		for _, cmd := range c.Commands {
			cmd.AddFlag(f)
		}
	}

	c.Flags.expected[f.Name] = f
}

// FullName returns the full name of the current command. If the command is
// a sub command then the string returned will be a concatenation of the
// command and all ancestors.
//
// For example given the command,
//
//  theme ls
//
// then the FullName would be,
//
//  theme-ls
func (c Command) FullName() string {
	commands := make([]string, 0)

	next := &c

	for next != nil {
		commands = append([]string{next.Name}, commands...)

		next = next.Parent
	}

	return strings.Join(commands, "-")
}

// Run executes the command, and takes a nilHandler. The nilHandler is only
// used as an attempted fall-back if the command does not have a handler by
// default.
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
