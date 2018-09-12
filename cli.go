package cli

import (
	"errors"
	"strings"
)

type Cli struct {
	cmds commands

	flags flags

	main *Command
}

func New() *Cli {
	return &Cli{cmds: newCommands(), flags: newFlags()}
}

func addCommand(name string, handler commandHandler, cmds commands) *Command {
	cmd := &Command{
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
		if main == nil {
			return nil, errors.New("command '" + name + "' not found")
		}

		main.Args = args

		return main, nil
	}

	cmd.Args = args[1:]

	return findCommand(cmd.Args, cmd.Commands, cmd)
}

func (c *Cli) AddFlag(f *Flag) {
	f.global = true

	c.flags[f.Name] = f
}

func (c *Cli) Command(name string, handler commandHandler) *Command {
	return addCommand(name, handler, c.cmds)
}

func (c *Cli) Main(handler commandHandler) *Command {
	c.main = &Command{
		Args:     args([]string{}),
		Flags:    newFlags(),
		Commands: newCommands(),
		Handler:  handler,
	}

	return c.main
}

func (c *Cli) parseFlag(i int, arg string, cmd *Command) error {
	var flag *Flag

	for _, f := range cmd.Flags {
		if f.Matches(arg) {
			flag = f
			break
		}
	}

	if flag == nil {
		return errors.New("unknown option '" + arg + "'")
	}

	if strings.HasPrefix(arg, "--") {
		return c.parseLong(i, arg, cmd, flag)
	}

	return c.parseShort(i, arg, cmd, flag)
}

func (c *Cli) parseLong(i int, arg string, cmd *Command, flag *Flag) error {
	cmd.Args.set(i, "")

	if flag.Argument {
		val := ""

		if strings.Contains(arg, "=") {
			val = strings.Split(arg, "=")[1]
		} else {
			val = cmd.Args.Get(i + 1)

			cmd.Args.set(i + 1, "")
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
	cmd.Args.set(i, "")

	if flag.Argument {
		val := cmd.Args.Get(i + 1)

		cmd.Args.set(i + 1, "")

		if val == "" && flag.Default == nil {
			return errors.New("option '" + arg + "' requires an argument")
		}

		flag.Value = val

		return nil
	}

	flag.isSet = true

	return nil
}

func (c *Cli) Run(args_ []string) error {
	cmd, err := findCommand(args(args_), c.cmds, c.main)

	if err != nil {
		return err
	}

	for _, f := range c.flags {
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

	return cmd.Run()
}
