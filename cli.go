// Simple library for building CLI applications in Go.
package cli

import (
	"errors"
	"strings"
)

type Cli struct {
	cmds       commands
	flags      flags
	main       *Command
	nilHandler commandHandler
}

// New creates a new Cli struct for attaching commands, and flags to. There is
// no limit to how many of these you can create.
func New() *Cli {
	return &Cli{cmds: newCommands(), flags: newFlags()}
}

func addCommand(name string, parent *Command, handler commandHandler, cmds commands) *Command {
	cmd := &Command{
		Parent:   parent,
		Name:     name,
		Args:     args([]string{}),
		Flags:    newFlags(),
		Commands: newCommands(),
		Handler:  handler,
	}

	cmds[name] = cmd

	return cmd
}

func findCommand(args args, cmds commands, main *Command) (*Command, error) {
	if len(args) == 0 {
		if main == nil {
			return nil, errors.New("failed to find a command")
		}

		return main, nil
	}

	name := args[0]

	if strings.HasPrefix(name, "--") || strings.HasPrefix(name, "-") {
		main.Args = args

		return main, nil
	}

	cmd, ok := cmds[name]

	if !ok {
		if main == nil || main.Handler == nil {
			return nil, errors.New("command '" + name + "' not found")
		}

		main.Args = args

		return main, nil
	}

	cmd.Args = args[1:]

	return findCommand(cmd.Args, cmd.Commands, cmd)
}

// AddFlag takes a pointer to a Flag struct, and adds it to the Cli struct
// marking it as global. This flag will be passed down to every subsequent
// command added to the Cli struct.
func (c *Cli) AddFlag(f *Flag) {
	f.global = true

	c.flags.expected[f.Name] = f
}

// Command creates a new command for the Cli struct based on the name, and
// handler given. A pointer to the newly created Command is returned. The name
// of the command is typically what the user would type in to have the command
// run.
func (c *Cli) Command(name string, handler commandHandler) *Command {
	return addCommand(name, nil, handler, c.cmds)
}

// Main specifies the main command to run should no initial command be found
// upon the first run of the application. This only takes a command handler.
func (c *Cli) Main(handler commandHandler) *Command {
	c.main = &Command{
		Args:     args([]string{}),
		Flags:    newFlags(),
		Commands: newCommands(),
		Handler:  handler,
	}

	return c.main
}

// NilHandler specifies a handler to be used for commands which do not have
// a handler on them. This is useful if you have mulitple sub-commands that
// perform actions, but whose parent command does not.
func (c *Cli) NilHandler(handler commandHandler) {
	c.nilHandler = handler
}

func (c *Cli) parseFlag(i int, arg string, cmd *Command) error {
	var flag *Flag

	for _, f := range cmd.Flags.expected {
		if f.Matches(arg) {
			flag = f
			break
		}
	}

	if flag == nil {
		return errors.New("unknown option '" + arg + "'")
	}

	if strings.HasPrefix(arg, "--") {
		if err := c.parseLong(i, arg, cmd, flag); err != nil {
			return err
		}

		cmd.Flags.putReceived(*flag)

		return nil
	}

	if err := c.parseShort(i, arg, cmd, flag); err != nil {
		return err
	}

	cmd.Flags.putReceived(*flag)

	return nil
}

func (c *Cli) parseLong(i int, arg string, cmd *Command, flag *Flag) error {
	cmd.Args.remove(i)

	if flag.Argument {
		val := ""

		if strings.Contains(arg, "=") {
			val = strings.Split(arg, "=")[1]
		} else {
			val = cmd.Args.Get(i)
		}

		if val == "" && flag.Default == nil {
			return errors.New("option '" + arg + "' requires an argument")
		}

		flag.Value = val

		return nil
	}

	flag.isSet = true

	return nil
}

func (c *Cli) parseShort(i int, arg string, cmd *Command, flag *Flag) error {
	cmd.Args.remove(i)

	if flag.Argument {
		val := cmd.Args.Get(i)

		if val == "" && flag.Default == nil {
			return errors.New("option '" + arg + "' requires an argument")
		}

		cmd.Args.remove(i)

		flag.Value = val

		return nil
	}

	flag.isSet = true

	return nil
}

// Run takes the slice of strings, and parses them as commands and flags.
func (c *Cli) Run(args_ []string) error {
	cmd, err := findCommand(args(args_), c.cmds, c.main)

	if err != nil {
		return err
	}

	for _, f := range c.flags.expected {
		cmd.AddFlag(f)
	}

	for i, arg := range cmd.Args {
		if arg == "--" {
			break
		}

		if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
			if err = c.parseFlag(i, arg, cmd); err != nil {
				return err
			}
		}
	}

	return cmd.Run(c.nilHandler)
}
